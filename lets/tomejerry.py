import argparse
import logging
import math
import sys
import traceback
import warnings
from collections import namedtuple
from typing import Iterable, Optional, Union, List, Dict, Any

import os
import threading
import time

import MySQLdb.cursors
import progressbar
from abc import abstractmethod, ABC
from enum import Enum, IntEnum
from progressbar import DynamicMessage, FormatLabel

from objects import beatmap
from objects import score
from common.db import dbConnector
from helpers import config
from objects import glob


MAX_WORKERS = 32
UNIX = os.name == "posix"
FAILED_SCORES_LOGGER = None


RecalculatorQuery = namedtuple("RecalculatorQuery", "query parameters")


class WorkerStatus(IntEnum):
    NOT_STARTED = 0
    RECALCULATING = 1
    SAVING = 2
    DONE = 3


class LwScore:
    """
    A lightweight score object, that can hold score id and pp only
    """
    __slots__ = ("score_id", "pp")

    def __init__(self, score_id: Optional[int]=None, pp: Optional[int]=None, score_: Optional[score.score]=None):
        """
        Initializes a new LwScore. Either score_id and pp OR just score must be provided.

        :param score_id: id of the score. Optional.
        :param pp: pp. Optional.
        :param score_: score object. Optional.
        """
        if score_ is not None:
            self.score_id = score_.scoreID
            self.pp = score_.pp
        elif score_id is not None and pp is not None:
            self.score_id = score_id
            self.pp = pp
        else:
            raise RuntimeError("")


class Recalculator(ABC):
    """
    Base PP Recalculator
    """
    def __init__(self, ids_query: RecalculatorQuery, count_query: RecalculatorQuery):
        """
        Instantiates a new recalculator

        :param ids_query: `RecalculatorQuery` that fetches the `id`s of the scores of which pp will be recalculated
        :param count_query: `RecalculatorQuery` that counts the _total_ number of the scoresof which pp will be
        recalculated
        """
        self.ids_query: RecalculatorQuery = ids_query
        self.count_query: RecalculatorQuery = count_query

    @abstractmethod
    def offset_ids_query(self, limit: int, offset: int) -> RecalculatorQuery:
        """
        Returns a new `RecalculatorQuery` based on `self.ids_query`, but based with LIMIT and OFFSET.
        Will be run by each worker to get their scores.

        :param limit: LIMIT value
        :param offset: OFFSET value
        :return: `RecalculatorQuery` with LIMIT and OFFSET
        """
        raise NotImplementedError()


class SimpleRecalculator(Recalculator):
    """
    A simple recalculator that can use a set of simple conditions, joined with logic ANDs
    """
    def __init__(
        self,
        conditions: Union[Iterable[str], str],
        parameters: Optional[Union[Iterable[str], Dict[str, Any]]]=None
    ):
        """
        Initializes a new SimpleRecalculator

        :param conditions: The conditions that will be joined with login ANDs.
        They can be:
        * an iterable (list, tuple, ...) of str (multiple conditions)
        * str (one condition)
        :param parameters: Iterable (list, tuple, ...) or dict that contains the query's parameters.
        These will be passed to MySQLdb to bind the query's parameters (%s and %(name)s)
        """
        if type(conditions) is list or type(conditions) is tuple:
            conditions_str = " AND ".join(conditions)
        elif type(conditions) is str:
            conditions_str = conditions
        else:
            raise TypeError("`conditions` must be either a `str`, `tuple` or `list`")
        q = "SELECT {} FROM scores JOIN beatmaps USING(beatmap_md5) WHERE {} ORDER BY scores.id DESC"
        super(SimpleRecalculator, self).__init__(
            ids_query=RecalculatorQuery(q.format("scores.id AS id", conditions_str), parameters),
            count_query=RecalculatorQuery(q.format("COUNT(*) AS c", conditions_str), parameters)
        )

    def offset_ids_query(self, limit: int, offset: int) -> str:
        return self.ids_query.query + " LIMIT {} OFFSET {}".format(limit, offset)


class ScoreIdsPool:
    """
    Pool of score ids that needs to be recalculated.
    """
    logger = logging.getLogger("score_ids_pool")

    def __init__(self):
        """
        Initializes a new pool
        """
        self._lock = threading.RLock()
        self.scores = []

    def load(self, recalculator: Recalculator):
        """
        Loads score ids in the pool from a Recalculator instance

        :param recalculator: The recalculator instance that will be used to fetch the score ids
        :return:
        """
        with self._lock:
            query_result = glob.db.fetchAll(recalculator.ids_query.query, recalculator.ids_query.parameters)
            self.scores += [LwScore(x["id"], 0) for x in query_result]
        self.logger.debug("Loaded {} scores".format(len(self.scores)))

    def chunk(self, chunk_size: int) -> List[int]:
        """
        Returns a chunk of score ids of the specified size, and removes the chunk from the pool.

        :param chunk_size: size of the chunk
        :return: score ids list
        """
        with self._lock:
            chunked_scores = self.scores[:chunk_size]
            self.scores = self.scores[chunk_size:]
        self.logger.debug("Chunked {} scores. Current scores in pool: {}".format(chunk_size, len(self.scores)))
        return chunked_scores

    @property
    def is_empty(self):
        """
        Whether the pool is empty or not

        :return: `True` if the pool is empty else `False`
        """
        return not bool(self.scores)


class Worker:
    """
    A tomejerry worker. Recalculates pp for a set of scores.
    """
    score_ids_pool = ScoreIdsPool()

    def __init__(self, chunk_size: int, worker_id: int=-1, start: bool=True):
        """
        Initializes a new worker.

        :param chunk_size: Number of scores to process
        :param worker_id: This worker's id. Optional. Default: -1.
        :param start: Whether to start the worker immediately or not
        :param
        """
        self.worker_id: int = worker_id
        self.thread: threading.Thread = None
        self.logger: logging.Logger = logging.getLogger("w{}".format(worker_id))
        self.recalculated_scores_count: int = 0
        self.saved_scores_count: int = 0
        self.chunk_size: int = chunk_size
        self.scores: List[LwScore] = self.score_ids_pool.chunk(self.chunk_size)
        self.status: WorkerStatus = WorkerStatus.NOT_STARTED
        self.failed_scores: int = 0
        if start:
            self.threaded_work()

    def recycle(self, start: bool=True):
        """
        Recycles this worker with a new chunk of scores

        :param start: Whether to start the worker immediately or not
        :return:
        """
        if self.thread.is_alive():
            raise RuntimeError("The thread is still alive")
        del self.thread
        self.thread = None
        self.status = WorkerStatus.NOT_STARTED
        self.scores = self.score_ids_pool.chunk(self.chunk_size)
        self.logger.debug("Recycled with {} new scores".format(self.chunk_size))
        if start:
            self.threaded_work()

    def recalc_score(self, score_data: Dict) -> score:
        """
        Recalculates pp for a score

        :param score_data: dict containing score and beatmap information about a score.
        :return: new `score` object, with `pp` attribute set to the new value
        """
        # Create score object and set its data
        s: score.score = score.score()
        s.setDataFromDict(score_data)
        s.passed = True

        # Create beatmap object and set its data
        b: beatmap.beatmap = beatmap.beatmap()
        b.setDataFromDict(score_data)

        # Calculate score pp
        s.calculatePP(b)
        del b
        return s

    def _work(self):
        """
        Run worker's work. Fetches scores, recalculates pp and saves the results in the database.

        :return:
        """
        # Make sure the worker hasn't been disposed
        if self.status == WorkerStatus.DONE:
            raise RuntimeError("This worker has been disposed")

        self.logger.info("Started worker. Assigned {} scores".format(self.chunk_size))
        try:
            # Recalculate all pp and save results in memory using LwScore objects
            self.recalculate_pp()

            # Store the new pp values permanently in the database
            self.save_recalculations()
        finally:
            # Mark the worker as disposed at the end
            self.logger.debug("Disposing worker")
            self.status = WorkerStatus.DONE

    def recalculate_pp(self):
        """
        Recalculates the pp and saves results in memory

        :return:
        """
        # We cannot use a SSDictCursor directly, because the connection will time out
        # if the cursor doesn't consume every result before the `wait_timeout`, which is
        # 600 seconds in MariaDB's default configuration. This means that we have to recalculate
        # PPs for all scores in no more than 600 seconds, or we'll get a 'MySQL server has
        # gone away error'. Fetching every score (joined with the respective beatmap)
        # directly would take up too much RAM, so we fetch all the score_ids at the
        # beginning with one query, store them in memory and fetch the data for
        # each score, one by one, using the same connection (to avoid pool overhead)
        self.status = WorkerStatus.RECALCULATING
        # self.recalculated_scores_count = 0

        # Fetch all score_ids
        # self.scores = [LwScore(x["id"], 0) for x in glob.db.fetchAll(self.ids_query.query, self.ids_query.parameters)]

        # Get a db worker
        cursor = None
        db_worker = glob.db.pool.getWorker()
        if db_worker is None:
            self.logger.warning("Cannot fetch scores. No database worker available!!")
            return

        try:
            # Get a cursor (normal DictCursor)
            cursor = db_worker.connection.cursor(MySQLdb.cursors.DictCursor)
            for i, lw_score in enumerate(self.scores):
                if i % self.log_every == 0:
                    self.logger.debug("Processed {}/{} scores".format(i, self.chunk_size))

                # Fetch score and beatmap data for this id
                cursor.execute(
                    "SELECT * FROM scores JOIN beatmaps USING(beatmap_md5) WHERE scores.id = %s LIMIT 1",
                    (lw_score.score_id,)
                )
                score_ = cursor.fetchone()
                try:
                    # Recalculate pp
                    recalculated_score = self.recalc_score(score_)

                    if recalculated_score is not None:
                        # New score returned, store new pp in memory
                        self.scores[i].pp = recalculated_score.pp
                        if recalculated_score.pp == 0:
                            # PP calculator error
                            self.log_failed_score(score_, "0 pp")

                    # Mark for garbage collection
                    del score_
                    del recalculated_score
                except Exception as e:
                    self.log_failed_score(score_, str(e), traceback_=True)
                finally:
                    self.recalculated_scores_count += 1
        finally:
            # Close cursor and connection
            if cursor is not None:
                cursor.close()
            if db_worker is not None:
                glob.db.pool.putWorker(db_worker)
            self.logger.debug("PP Recalculated")

    def save_recalculations(self):
        """
        Saves the recalculated performance points in the database

        :return:
        """
        self.status = WorkerStatus.SAVING
        # self.saved_scores_count = 0

        # Make sure we've at least fetched the scores
        if not self.scores:
            self.logger.warning("No scores to update.")
            return

        # Update db
        self.logger.debug("Updating scores in database")
        for i, lw_score in enumerate(self.scores):
            if i % self.log_every == 0:
                self.logger.debug("Updated {}/{} scores".format(i, self.chunk_size))
            glob.db.execute("UPDATE scores SET pp = %s WHERE id = %s LIMIT 1", (lw_score.pp, lw_score.score_id))
            self.saved_scores_count += 1

        self.logger.debug("Scores updated")

    @property
    def log_every(self) -> int:
        """
        Number of scores that have to be processed before logging the worker's status

        :return:
        """
        return max(min((self.chunk_size // 3), 1000), 1)

    def threaded_work(self):
        """
        Starts this worker's work in a new thread

        :return:
        """
        self.thread = threading.Thread(target=self._work)
        self.thread.start()

    def log_failed_score(self, score_: Dict[str, Any], additional_information: str="", traceback_: bool=False):
        """
        Logs a failed score.

        :param score_: score dict (from db) that triggered the error
        :param additional_information: additional information (type of error)
        :param traceback_: Whether the traceback should be logged or not.
        It should be `True` if the logging was triggered by an unhandled exception
        :return:
        """
        msg = ""
        if traceback_:
            msg = "\n\n\nUnhandled exception: {}\n{}".format(sys.exc_info(), traceback.format_exc())
        msg += "score_id:{} ({})".format(score_["id"], additional_information).strip()
        FAILED_SCORES_LOGGER.error(msg)
        self.failed_scores += 1


def mass_recalc(recalculator: Recalculator, workers_number: int=MAX_WORKERS, chunk_size: Optional[int]=None):
    """
    Recalculate performance points for a set of scores, using multiple workers

    :param recalculator: the recalculator that will be used
    :param workers_number: the number of workers to spawn
    :return:
    """
    start_time = time.time()
    global FAILED_SCORES_LOGGER
    workers = []

    logging.info("Query: {} ({})".format(recalculator.ids_query.query, recalculator.ids_query.parameters))

    # Fetch the total number of scores
    total_scores = glob.db.fetch(recalculator.count_query.query, recalculator.count_query.parameters)
    if total_scores is None:
        logging.warning("No scores to recalc.")
        return

    # Set up failed scores logger (creates file too)
    FAILED_SCORES_LOGGER = logging.getLogger("failed_scores")
    FAILED_SCORES_LOGGER.addHandler(
        logging.FileHandler("tomejerry_failed_scores_{}.log".format(time.strftime("%d-%m-%Y--%H-%M-%S")))
    )

    # Get the number of total scores from the result dict
    total_scores = total_scores[next(iter(total_scores))]
    logging.info("Total scores: {}".format(total_scores))
    if total_scores == 0:
        return

    # for some reason `typing` believes that `math.ceil` returns a `float`, so we need an extra cast here...
    scores_per_worker = int(math.ceil(total_scores / workers_number))
    logging.info("Using {} workers and {} scores per worker".format(workers_number, scores_per_worker))

    # Load score ids in the pool
    logging.info("Filling score ids pool")
    Worker.score_ids_pool.load(recalculator)

    # Spawn the workers and start them
    for i in range(workers_number):
        workers.append(
            Worker(
                chunk_size=chunk_size
                if chunk_size is not None
                else len(Worker.score_ids_pool.scores) // workers_number // 3,
                worker_id=i,
                start=True
            )
        )

    # Progress bar loop
    steps_text = {
        WorkerStatus.NOT_STARTED: "Starting workers",
        WorkerStatus.RECALCULATING: "Recalculating pp",
        WorkerStatus.SAVING: "Updating db"
    }
    recycles = 0
    widgets = [
        "[ ", "Starting", " ]",
        "w_pp:<>", "w_db:<>", "w_done:<>", "rec:0",
        progressbar.FormatLabel(" %(value)s/%(max)s "),
        progressbar.Bar(marker="#", left="[", right="]", fill="."),
        progressbar.Percentage(),
        " (", progressbar.ETA(), ") "
    ]
    with progressbar.ProgressBar(
        widgets=widgets,
        max_value=total_scores,
        redirect_stdout=True,
        redirect_stderr=True
    ) as bar:
        while True:
            lowest_status = min([x.status for x in workers])

            # Loop through all workers to get progress value
            total_progress_value = sum(
                [
                    x.recalculated_scores_count if lowest_status != WorkerStatus.SAVING else x.saved_scores_count
                    for x in workers
                ]
            )

            # Recycle the workers if needed
            workers_done = [x for x in workers if x.status == WorkerStatus.DONE]
            if workers_done and not Worker.score_ids_pool.is_empty:
                logging.info("Recycling workers")
                recycles += 1
                for worker in workers_done:
                    worker.recycle(start=True)

            # Output total status information
            widgets[1] = steps_text.get(lowest_status, "...")
            widgets[3] = " w_pp:<{}/{}>".format(
                len([x for x in workers if x.status == WorkerStatus.RECALCULATING]), len(workers)
            )
            widgets[4] = " w_db:<{}/{}>".format(
                len([x for x in workers if x.status == WorkerStatus.SAVING]), len(workers)
            )
            widgets[5] = " w_done:<{}/{}>".format(len(workers_done), len(workers))
            widgets[6] = " rec:{}".format(recycles)
            bar.update(total_progress_value)

            # Exit from the loop if every worker has finished its work
            if len(workers_done) == len(workers):
                break

            # Wait 0.5 s and update the progress bar again
            time.sleep(0.5)

    # Recalc done. Print some stats
    end_time = time.time()
    failed_scores = sum([x.failed_scores for x in workers])
    logging.info(
        "\n\nDone!\n"
        ":: Recalculated\t{} scores\n"
        ":: Failed\t{} scores\n"
        ":: Total\t{} scores\n\n"
        ":: Took\t{:.2f} seconds".format(
            total_scores - failed_scores,
            failed_scores,
            total_scores,
            end_time - start_time
        )
    )


def main():
    # CLI stuff
    parser = argparse.ArgumentParser(description="pp recalc tool for ripple, new version.")
    recalc_group = parser.add_mutually_exclusive_group(required=False)
    recalc_group.add_argument(
        "-r", "--recalc", help="calculates pp for all high scores", required=False, action="store_true"
    )
    recalc_group.add_argument(
        "-z", "--zero", help="calculates pp for 0 pp high scores", required=False, action="store_true"
    )
    recalc_group.add_argument("-i", "--id", help="calculates pp for the score with this score_id", required=False)
    recalc_group.add_argument(
        "-m", "--mods", help="calculates pp for high scores with these mods (flags)", required=False
    )
    recalc_group.add_argument(
        "-g", "--gamemode", help="calculates pp for scores played on this game mode (std:0, taiko:1, ctb:2, mania:3)",
        required=False
    )
    recalc_group.add_argument(
        "-u", "--userid", help="calculates pp for high scores set by a specific user (user_id)", required=False
    )
    recalc_group.add_argument(
        "-b", "--beatmapid", help="calculates pp for high scores played on a specific beatmap (beatmap_id)", required=False
    )
    recalc_group.add_argument(
        "-fhd", "--fixstdhd", help="calculates pp for std hd high scores (14/05/2018 pp algorithm changes)",
        required=False, action="store_true"
    )
    parser.add_argument("-w", "--workers", help="number of workers. {} by default. Max {}".format(
        MAX_WORKERS // 2, MAX_WORKERS
    ), required=False)
    parser.add_argument("-cs", "--chunksize", help="score chunks size", required=False)
    parser.add_argument("-v", "--verbose", help="verbose/debug mode", required=False, action="store_true")
    args = parser.parse_args()

    # Logging
    progressbar.streams.wrap_stderr()
    logging.basicConfig(level=logging.DEBUG if args.verbose else logging.INFO)
    logging.info("Running under {}".format("UNIX" if UNIX else "WIN32"))

    # Load config
    logging.info("Reading config file")
    glob.conf = config.config("config.ini")

    # Get workers from arguments if set
    workers_number = MAX_WORKERS // 2
    if args.workers is not None:
        workers_number = int(args.workers)

    # Get chunk size from arguments if set
    chunk_size = None
    if args.chunksize is not None:
        chunk_size = int(args.chunksize)

    # Disable MySQL db warnings (it spams 'Unsafe statement written to the binary log using statement...'
    # because we use UPDATE with LIMIT 1 when updating performance points after recalculation
    warnings.filterwarnings("ignore", category=MySQLdb.Warning)

    # Connect to MySQL
    logging.info("Connecting to MySQL db")
    glob.db = dbConnector.db(
        glob.conf.config["db"]["host"],
        glob.conf.config["db"]["username"],
        glob.conf.config["db"]["password"],
        glob.conf.config["db"]["database"],
        max(workers_number, MAX_WORKERS)
    )

    # Set verbose
    glob.debug = args.verbose

    # Get recalculator
    recalculators_gen = {
        "zero": lambda: SimpleRecalculator(("scores.completed = 3", "pp = 0")),
        "recalc": lambda: SimpleRecalculator(("scores.completed = 3",)),
        "mods": lambda: SimpleRecalculator(("scores.completed = 3", "mods & %s > 0"), (args.mods,)),
        "id": lambda: SimpleRecalculator(("scores.id = %s",), (args.id,)),
        "gamemode": lambda: SimpleRecalculator(("scores.completed = 3", "scores.play_mode = %s",), (args.gamemode,)),
        "userid": lambda: SimpleRecalculator(("scores.completed = 3", "scores.userid = %s",), (args.userid,)),
        "beatmapid": lambda: SimpleRecalculator(("scores.completed = 3", "beatmaps.beatmap_id = %s",), (args.beatmapid,)),
        "fixstdhd": lambda: SimpleRecalculator(("scores.completed = 3", "scores.play_mode = 0", "scores.mods & 8 > 0"))
    }
    recalculator = None
    for k, v in vars(args).items():
        if v is not None and ((type(v) is bool and v) or type(v) is not bool):
            if k in recalculators_gen:
                recalculator = recalculators_gen[k]()
                break

    # Execute mass recalc
    if recalculator is not None:
        mass_recalc(recalculator, workers_number, chunk_size)
    else:
        logging.warning("No recalc option specified")
        parser.print_help()


if __name__ == "__main__":
    main()
