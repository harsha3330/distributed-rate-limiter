i) Cluster status

GET /v1/cluster

{
    "leader": "node1",
    "term": 8,
    "commit_index": 145,
    "applied_index": 145,
    "members": [
        "node1",
        "node2",
        "node3"
    ]
}

ii) Add node

POST /v1/cluster/members

{
    "id": "node4",
    "address": "10.0.0.4:9000"
}

ConfChangeAddNode


iii) Remove Node 

DELETE /v1/cluster/members/node4

ConfChangeRemoveNode

iv) GET /healthz

