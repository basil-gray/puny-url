package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"puny-url/config"
	"puny-url/db"
	"puny-url/internal/cache"
	"puny-url/internal/logger"
	"puny-url/routes"
)

func main() {

	flag.DurationVar(&config.Args.EvictInterval, "e", 1*time.Second, "Interval for eviction checks")
	flag.DurationVar(&config.Args.Ttl, "ttl", 30*time.Second, "Time after which inactive entries are evicted")
	flag.IntVar(&config.Args.Port, "port", 8080, "Port to run the server on")
	silent := flag.Bool("s", false, "Enable silent logging")
	flag.Parse()

	logger.Init(*silent)
	logger.Warn("Interval for eviction job set to " + config.Args.EvictInterval.String())
	logger.Warn("TTL set to " + config.Args.Ttl.String())

	db.InitDB()

	defer db.Close()

	c := cache.New(config.Args.EvictInterval, config.Args.Ttl)
	r := routes.SetupRoutes(c)

	logger.Warn("PunyURL now running on :" + strconv.Itoa(config.Args.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Args.Port), r))
}
