#!/bin/bash

docker compose build
docker tag jury-main michaelzhao21/jury:latest
docker push michaelzhao21/jury:latest
