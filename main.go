package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"puny-url/cache"
	"puny-url/db"
	"puny-url/routes"
)

func main() {
	flag.DurationVar(&cache.EvictInterval, "evictInterval", 30*time.Second, "Interval for eviction checks")
	flag.DurationVar(&cache.TTL, "inactiveTTL", 30*time.Second, "Time after which inactive entries are evicted")
	flag.Parse()

	log.Printf("Cache eviction interval set to %s", cache.EvictInterval.String())
	log.Printf("Inactive TTL set to %s", cache.TTL.String())

	db.InitDB()

	defer db.Close()

	go cache.StartCacheEviction()

	r := routes.SetupRoutes()

	log.Println("PunyURL now running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
