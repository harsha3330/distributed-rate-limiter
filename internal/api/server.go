package api

import (
	"net/http"

	"github.com/harsha3330/distributed-rate-limiter/internal/bucket"
	jsonhelper "github.com/harsha3330/distributed-rate-limiter/internal/helper"
	"github.com/harsha3330/distributed-rate-limiter/internal/policy"
	"github.com/harsha3330/distributed-rate-limiter/internal/ratelimiter"
)

type AppServer struct {
	rateLimiter *ratelimiter.RateLimiter
	policyStore *policy.PolicyStore
	bucKetStore *bucket.BucketStore
}

func NewAppServer(bs *bucket.BucketStore, ps *policy.PolicyStore) *AppServer {
	return &AppServer{
		rateLimiter: ratelimiter.NewRateLimiter(bs, ps),
		policyStore: ps,
		bucKetStore: bs,
	}
}

func (app *AppServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /policies", app.GetPolicies)
	return mux
}

func (app *AppServer) GetPolicies(w http.ResponseWriter, r *http.Request) {
	policies := app.policyStore.ListPolicy()
	jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"policies": policies,
	}, nil)
}
