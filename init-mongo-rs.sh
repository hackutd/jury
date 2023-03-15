#!/bin/bash

# Initiates the replica set for the MongoDB instance
# after waiting for the container to initialize completely

echo "sleeping for 10 seconds"
sleep 10

echo trying connection at `date +"%T" `
mongosh mongodb://$MONGODB_USER:$MONGODB_PASS@mongo:27017 <<EOF
    rs.initiate({_id: "rs0", "version": 1, "members": [{ "_id": 0, "host": "mongo:27017" }]});
EOF
