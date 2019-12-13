package common

import (
	"fmt"

	"github.com/thehowl/conf"
)

// Version is the git hash of the application. Do not edit. This is
// automatically set using -ldflags during build time.
var Version string

// Conf is the configuration file data for the ripple API.
// Conf uses https://github.com/thehowl/conf
type Conf struct {
	DatabaseType           string `description:"At the moment, 'mysql' is the only supported database type."`
	DSN                    string `description:"The Data Source Name for the database. More: https://github.com/go-sql-driver/mysql#dsn-data-source-name"`
	ListenTo               string `description:"The IP/Port combination from which to take connections, e.g. :8080"`
	Unix                   bool   `description:"Bool indicating whether ListenTo is a UNIX socket or an address."`
	SentryDSN              string `description:"thing for sentry whatever"`
	HanayoKey              string
	BeatmapRequestsPerUser int
	RankQueueSize          int
	OsuAPIKey              string
	RedisAddr              string
	RedisPassword          string
	RedisDB                int
}

var cachedConf *Conf

// Load creates a new Conf, using the data in the file "api.conf".
func Load() (c Conf, halt bool) {
	if cachedConf != nil {
		c = *cachedConf
		return
	}
	err := conf.Load(&c, "api.conf")
	halt = err == conf.ErrNoFile
	if halt {
		conf.MustExport(Conf{
			DatabaseType:           "mysql",
			DSN:                    "root@/ripple",
			ListenTo:               ":40001",
			Unix:                   false,
			HanayoKey:              "Potato",
			BeatmapRequestsPerUser: 2,
			RankQueueSize:          25,
			RedisAddr:              "localhost:6379",
		}, "api.conf")
		fmt.Println("Please compile the configuration file (api.conf).")
	}
	cachedConf = &c
	return
}

// GetConf returns the cachedConf.
func GetConf() *Conf {
	if cachedConf == nil {
		return nil
	}
	// so that the cachedConf cannot actually get modified
	c := *cachedConf
	return &c
}
