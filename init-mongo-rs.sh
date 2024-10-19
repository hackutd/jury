#!/bin/bash

# Initiates the replica set for the MongoDB instance
# after waiting for the container to initialize completely

RETRY_INTERVAL=2  # Interval between retries
MAX_ATTEMPTS=10   # Max attempts before giving up

echo "Waiting for MongoDB to be ready..."

for ((i=1; i<=MAX_ATTEMPTS; i++)); do
    echo "Attempt $i at `date +"%T"`: Trying to connect..."
    
    # Try to connect to MongoDB using mongosh and check the server status
    mongosh --quiet --eval "db.runCommand({ping: 1})" mongodb://$MONGODB_USER:$MONGODB_PASS@mongo:27017 > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "MongoDB is up! Proceeding with replica set initiation..."
        mongosh mongodb://$MONGODB_USER:$MONGODB_PASS@mongo:27017 <<EOF
            rs.initiate({_id: "rs0", "version": 1, "members": [{ "_id": 0, "host": "mongo:27017" }]});
EOF
        break
    else
        echo "MongoDB not ready, retrying in $RETRY_INTERVAL seconds..."
        sleep $RETRY_INTERVAL
    fi

    if [ $i -eq $MAX_ATTEMPTS ]; then
        echo "Failed to connect to MongoDB after $MAX_ATTEMPTS attempts. Exiting."
        exit 1
    fi
done
