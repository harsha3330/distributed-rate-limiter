type Policy struct {
    ID        string
    Limit     uint32
    Window    time.Duration
}

type BucketKey struct {
    PolicyID string
    Subject  string
}

type Bucket struct {
    Count uint32
    Start int64
}

type PolicyStore struct {
    policies map[string]*Policy
}

type BucketStore struct {
    buckets map[BucketKey]*Bucket
}

type StateMachine struct {
    policies map[string]*Policy
    buckets  map[BucketKey]*Bucket
}



POST   /v1/policies
GET    /v1/policies
GET    /v1/policies/{id}
PUT    /v1/policies/{id}
DELETE /v1/policies/{id}

POST   /v1/check

GET    /v1/buckets/{policy}/{key}
DELETE /v1/buckets/{policy}/{key}

GET    /healthz
GET    /metrics