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

	var evictInterval time.Duration
	var ttl time.Duration

	flag.DurationVar(&evictInterval, "e", 30*time.Second, "Interval for eviction checks")
	flag.DurationVar(&ttl, "ttl", 30*time.Second, "Time after which inactive entries are evicted")
	flag.Parse()

	log.Printf("Cache eviction interval set to %s", evictInterval.String())
	log.Printf("Inactive TTL set to %s", ttl)

	db.InitDB()

	defer db.Close()

	c := cache.New(evictInterval, ttl)

	r := routes.SetupRoutes(c)

	log.Println("PunyURL now running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
