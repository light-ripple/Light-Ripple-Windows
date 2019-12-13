from common.log import logUtils as log
from common.constants import gameModes
from constants import exceptions
from helpers import mapsHelper

from pp.catch_the_pp.osu_parser.beatmap import Beatmap as CalcBeatmap
from pp.catch_the_pp.osu.ctb.difficulty import Difficulty
from pp.catch_the_pp import ppCalc


class Cicciobello:
    def __init__(self, _beatmap, _score=None, accuracy=0, mods=0, combo=-1, misses=0, tillerino=False):
        # Beatmap is always present
        self.beatmap = _beatmap

        # If passed, set everything from score object
        if _score is not None:
            self.score = _score
            self.accuracy = self.score.accuracy
            self.mods = self.score.mods
            self.combo = self.score.maxCombo
            self.misses = self.score.cMiss
        else:
            # Otherwise, set acc and mods from params (tillerino)
            self.accuracy = accuracy
            self.mods = mods
            self.combo = combo
            if self.combo < 0:
                self.combo = self.beatmap.maxCombo
            self.misses = misses

        # Multiple acc values computation
        self.tillerino = tillerino

        # Result
        self.pp = 0
        self.calculate_pp()

    def calculate_pp(self):
        try:
            # Cache beatmap
            mapFile = mapsHelper.cachedMapPath(self.beatmap.beatmapID)
            mapsHelper.cacheMap(mapFile, self.beatmap)

            # TODO: Sanizite mods

            # Gamemode check
            if self.score and self.score.gameMode != gameModes.CTB:
                raise exceptions.unsupportedGameModeException()

            # Accuracy check
            if self.accuracy > 1:
                raise ValueError("Accuracy must be between 0 and 1")

            # Calculate difficulty
            calcBeatmap = CalcBeatmap(mapFile)
            difficulty = Difficulty(beatmap=calcBeatmap, mods=self.mods)

            # Calculate pp
            if self.tillerino:
                results = []
                for acc in [1, 0.99, 0.98, 0.95]:
                    results.append(ppCalc.calculate_pp(
                        diff=difficulty,
						accuracy=acc,
						combo=self.combo if self.combo >= 0 else calcBeatmap.max_combo,
						miss=self.misses
                    ))
                self.pp = results
            else:
                self.pp = ppCalc.calculate_pp(
                    diff=difficulty,
					accuracy=self.accuracy,
					combo=self.combo if self.combo >= 0 else calcBeatmap.max_combo,
					miss=self.misses
                )
        except exceptions.osuApiFailException:
            log.error("cicciobello ~> osu!api error!")
            self.pp = 0
        except exceptions.unsupportedGameModeException:
            log.error("cicciobello ~> Unsupported gamemode")
            self.pp = 0
        except Exception as e:
            log.error("cicciobello ~> Unhandled exception: {}".format(str(e)))
            self.pp = 0
            raise
        finally:
            log.debug("cicciobello ~> Shutting down, pp = {}".format(self.pp))