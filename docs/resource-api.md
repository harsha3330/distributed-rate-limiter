Resources: 

Policies
Buckets (runtime state)
Cluster (Raft)


i) Create Policy

POST /v1/policies

req :
{
    "name": "login",
    "limit": 5,
    "window": "1m"
}
res :
{
    "message" : "succesful / error"
}

ii) List Policies 

GET /v1/policies
res :
[
    {
        "id": "policy-1",
        "name": "login",
        "limit": 5,
        "window": "1m"
    },
    {
        "id": "policy-2",
        "name": "search",
        "limit": 100,
        "window": "1m"
    }
]

iii) Get Policy 

GET /v1/policies/policy-1

iv) Update Policy 

PUT /v1/policies/policy-1
{
    "limit": 10,
    "window": "1m"
}

v) DELETE /v1/policies/policy-1


vi) Check Ratelimiter status

POST /v1/check

req:
{
    "key": "user123",
    "policy": "login"
}

res: 
{
    "allowed": true,
    "remaining": 4,
    "limit": 5,
    "reset_after_seconds": 42
}

vii) Get Bucket

GET /v1/buckets/login/user123

{
    "key": "user123",
    "policy": "login",
    "count": 3,
    "window_start": "2026-07-01T12:00:00Z",
    "window_end": "2026-07-01T12:01:00Z"
}


viii) Delete Bucket

DELETE /v1/buckets/login/user123

{
    "success": true
}

ix) List Buckets

GET /v1/buckets
GET /v1/buckets?policy=login

[
    {
        "key": "user1",
        "policy": "login",
        "count": 4
    },
    {
        "key": "user2",
        "policy": "login",
        "count": 1
    }
]