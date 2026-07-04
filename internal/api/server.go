package api

import (
	"fmt"
	"net/http"

	"github.com/harsha3330/distributed-rate-limiter/internal/bucket"
	jsonhelper "github.com/harsha3330/distributed-rate-limiter/internal/helper"
	"github.com/harsha3330/distributed-rate-limiter/internal/policy"
	"github.com/harsha3330/distributed-rate-limiter/internal/ratelimiter"
)

type AppServer struct {
	rateLimiter *ratelimiter.RateLimiter
	policyStore *policy.PolicyStore
	bucketStore *bucket.BucketStore
}

type CheckRequest struct {
	PolicyID string `json:"policy_id"`
	Subject  string `json:"subject"`
}

func NewAppServer(bs *bucket.BucketStore, ps *policy.PolicyStore) *AppServer {
	rl := ratelimiter.NewRateLimiter(bs, ps)
	return &AppServer{
		rateLimiter: rl,
		policyStore: ps,
		bucketStore: bs,
	}
}

func (app *AppServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /policies", app.GetPolicies)
	mux.HandleFunc("GET /policies/{id}", app.GetPolicyByID)
	mux.HandleFunc("POST /policies", app.CreatePolicy)
	mux.HandleFunc("PUT /policies/{id}", app.UpdatePolicyByID)
	mux.HandleFunc("DELETE /policies/{id}", app.DeletePolicyByID)
	mux.HandleFunc("GET /buckets", app.GetBuckets)
	mux.HandleFunc("POST /check", app.Check)
	mux.HandleFunc("DELETE /buckets/{policy}/{subject}", app.DeleteBucket)
	mux.HandleFunc("GET /health", app.Health)
	return mux
}

func (app *AppServer) Health(w http.ResponseWriter, r *http.Request) {
	if err := jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"status": "ok",
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var policy *policy.Policy = &policy.Policy{}
	err := jsonhelper.ReadJSON(w, r, policy)
	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	err = app.policyStore.CreatePolicy(policy)
	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusConflict, err.Error())
		return
	}
	if err = jsonhelper.WriteJSON(w, http.StatusCreated, jsonhelper.Envelope{
		"policy":  *policy,
		"message": "policy created",
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) GetPolicies(w http.ResponseWriter, r *http.Request) {
	policies := app.policyStore.ListPolicy()
	jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"policies": policies,
	}, nil)
}

func (app *AppServer) GetPolicyByID(w http.ResponseWriter, r *http.Request) {
	policyId := r.PathValue("id")
	if policyId == "" {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "policy id is required")
		return
	}
	policy, err := app.policyStore.GetPolicy(policyId)
	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}
	if err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"policy": *policy,
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) UpdatePolicyByID(w http.ResponseWriter, r *http.Request) {
	policyId := r.PathValue("id")
	if policyId == "" {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "policy id is required")
		return
	}

	var policy *policy.Policy = &policy.Policy{}
	err := jsonhelper.ReadJSON(w, r, policy)
	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if policy.ID != "" && policy.ID != policyId {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "policy id in URL and body do not match")
		return
	}
	policy.ID = policyId
	err = app.policyStore.UpdatePolicy(policy)
	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"policy":  *policy,
		"message": "updated the policy",
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) DeletePolicyByID(w http.ResponseWriter, r *http.Request) {
	policyId := r.PathValue("id")
	if policyId == "" {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "policy id is required")
		return
	}
	err := app.policyStore.DeletePolicy(policyId)
	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}
	if err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"message": fmt.Sprintf("Delete policy with id: %s", policyId),
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) Check(w http.ResponseWriter, r *http.Request) {
	req := &CheckRequest{}

	if err := jsonhelper.ReadJSON(w, r, req); err != nil {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.PolicyID == "" {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "policy_id is required")
		return
	}

	if req.Subject == "" {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "subject is required")
		return
	}

	resp, err := app.rateLimiter.Check(bucket.BucketKey{
		PolicyID: req.PolicyID,
		Subject:  req.Subject,
	})

	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	if err := jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"result": resp,
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) GetBuckets(w http.ResponseWriter, r *http.Request) {
	buckets := app.bucketStore.ListBucket()
	if err := jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"buckets": buckets,
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *AppServer) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	policyID := r.PathValue("policy")
	subject := r.PathValue("subject")

	if policyID == "" || subject == "" {
		jsonhelper.ErrorJSON(w, http.StatusBadRequest, "policy and subject are required")
		return
	}

	err := app.bucketStore.DeleteBucket(bucket.BucketKey{
		PolicyID: policyID,
		Subject:  subject,
	})

	if err != nil {
		jsonhelper.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	if err := jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{
		"message": "bucket deleted",
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
