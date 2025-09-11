#!/bin/bash

# This script will run the test using docker compose!
# It will poll for the jury-testing container status until it has exited.
# Then it will read the container logs to see if any tests have failed.
# Logs will be printed to the console.
# Finally, it will exit and kill the docker compose containers

# Start docker containers and wait a second (just for fun)
# docker compose -f docker-compose.test.yml up -d
# sleep 1

echo "Running tests..."

timeout=0   # Max timeout set to 60 seconds, may need to change if longer tests!
while [ "$timeout" -lt 60 ]; do
    if [ "$(docker inspect --format '{{.State.Status}}' jury-testing)" = "exited" ]; then
        break
    fi
    sleep 1
done

# Get logs and get rid of containers
logs=$(docker logs jury-testing)
# docker compose -f docker-compose.test.yml down

# If no failed lines, exit with success
failed=$(echo "$logs" | grep "failed")
if [[ -z "$failed" ]]; then
    echo "Success :)"
    exit 0
fi

# Otherwise, exit with failure
printf "$logs"
printf "\n\n###################################\n##### THERE ARE TEST FAILURES #####\n###################################\n"
exit 1
