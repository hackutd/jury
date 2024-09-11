#!/bin/bash

docker pull michaelzhao21/jury
docker stop jury-main
docker run --rm -d --name jury-main --env-file ./jury.env -p 8080:8080 michaelzhao21/jury

