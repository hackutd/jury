#!/bin/bash

docker compose build
docker tag jury-main ghcr.io/hackutd/jury:latest
docker push ghcr.io/hackutd/jury:latest
