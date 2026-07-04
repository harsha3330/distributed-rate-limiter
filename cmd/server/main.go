package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/harsha3330/distributed-rate-limiter/internal/api"
	"github.com/harsha3330/distributed-rate-limiter/internal/bucket"
	"github.com/harsha3330/distributed-rate-limiter/internal/policy"
	"github.com/harsha3330/distributed-rate-limiter/internal/ratelimiter"
)

func main() {

	addr := flag.String("addr", ":8080", "the host and port where server is gonna run")
	flag.Parse()
	bs := bucket.NewBucketStore()
	ps := policy.NewPolicyStore()
	rl := ratelimiter.NewRateLimiter(bs, ps)
	app := api.NewAppServer(bs, ps, rl)

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on %s", *addr)

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
