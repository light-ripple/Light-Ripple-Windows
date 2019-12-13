from secret.achievements.utils import achievements_response


class Chart:
    """
    Chart base class
    """
    def __init__(self, id_, url, name):
        """
        Initializes a new chart.

        :param id_: chart id. Currently known values are 'beatmap' and 'overall'
        :param url: URL to open when clicking on the chart title.
        :param name: chart name displayed in the game client
        """
        self.id_ = id_
        self.url = url
        self.name = name

    def items(self):
        """
        `items()` method that allows this class to be used as a iterable dict

        :return:
        """
        return self.output_attrs.items()

    @property
    def output_attrs(self):
        """
        An unzingonified dict containing the stuff that will be sent to the game client

        :return: dict
        """
        return {
            "chartId": self.id_,
            "chartUrl": self.url,
            "chartName": self.name
        }

    @staticmethod
    def before_after_dict(name, values, none_value="0"):
        """
        Turns a tuple with two elements in a dict with two elements.

        :param name: prefix of the keys
        :param values: (value_before, value_after). value_before and value_after can be None.
        :param none_value: value to use instead of None (None, when zingonified, is not recognized by the game client)
        :return: { XXXBefore -> first element, XXXAfter -> second element }, where XXX is `name`
        """
        return {
            f"{name}{'Before' if i == 0 else 'After'}": x if x is not None else none_value for i, x in enumerate(values)
        }


class BeatmapChart(Chart):
    """
    Beatmap ranking chart
    """
    def __init__(self, old_score, new_score, beatmap_id):
        """
        Initializes a new BeatmapChart object.

        :param old_score: score object of the old score
        :param new_score: score object of the currently submitted score
        :param beatmap_id: beatmap id, for the clickable link
        """
        super(BeatmapChart, self).__init__("beatmap", f"https://ripple.moe/b/{beatmap_id}", "Beatmap Ranking")
        self.rank = (old_score.rank if old_score is not None else None, new_score.rank)
        self.max_combo = (old_score.maxCombo if old_score is not None else None, new_score.maxCombo)
        self.accuracy = (old_score.accuracy * 100 if old_score is not None else None, new_score.accuracy * 100)
        self.ranked_score = (old_score.score if old_score is not None else None, new_score.score)
        self.pp = (old_score.pp if old_score is not None else None, new_score.pp)
        self.score_id = new_score.scoreID

    @property
    def output_attrs(self):
        return {
            **super(BeatmapChart, self).output_attrs,
            **self.before_after_dict("rank", self.rank, none_value=""),
            **self.before_after_dict("maxCombo", self.max_combo),
            **self.before_after_dict("accuracy", self.accuracy),
            **self.before_after_dict("rankedScore", self.ranked_score),
            **self.before_after_dict("pp", self.pp),
            "onlineScoreId": self.score_id
        }


class OverallChart(Chart):
    """
    Overall ranking chart + achievements
    """
    def __init__(self, user_id, old_user_stats, new_user_stats, score, new_achievements, old_rank, new_rank):
        """
        Initializes a new OverallChart object.
        This constructor sucks because LETS itself sucks.

        :param user_id: id of the user
        :param old_user_stats: user stats dict before submitting the score
        :param new_user_stats: user stats dict after submitting the score
        :param score: score object of the scores that has just been submitted
        :param new_achievements: achievements unlocked list
        :param old_rank: global rank before submitting the scpre
        :param new_rank: global rank after submitting the score
        """
        super(OverallChart, self).__init__("overall", f"https://ripple.moe/u/{user_id}", "Overall Ranking")
        self.rank = (old_rank, new_rank)
        self.ranked_score = (old_user_stats["rankedScore"], new_user_stats["rankedScore"])
        self.total_score = (old_user_stats["totalScore"], new_user_stats["totalScore"])
        self.max_combo = (0, 0)     # TODO: Implement
        self.accuracy = (old_user_stats["accuracy"], new_user_stats["accuracy"])
        self.pp = (old_user_stats["pp"], new_user_stats["pp"])
        self.new_achievements = new_achievements
        self.score_id = score.scoreID

    @property
    def output_attrs(self):
        return {
            **super(OverallChart, self).output_attrs,
            **self.before_after_dict("rank", self.rank),
            **self.before_after_dict("rankedScore", self.ranked_score),
            **self.before_after_dict("totalScore", self.total_score),
            **self.before_after_dict("maxCombo", self.max_combo),
            **self.before_after_dict("accuracy", self.accuracy),
            **self.before_after_dict("pp", self.pp),
            "achievements-new": achievements_response(self.new_achievements),
            "onlineScoreId": self.score_id
        }
