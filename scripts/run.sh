#!/bin/bash

docker pull ghcr.io/hackutd/jury:latest
docker stop jury-main && sleep 2 && docker run --rm -d --name jury-main --env-file ./jury.env -p 8083:8080 ghcr.io/hackutd/jury:latest
