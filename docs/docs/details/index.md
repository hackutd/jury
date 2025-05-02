---
sidebar_position: 3
title: Technical Details
description: The nitty-gritty details behind how Jury and its algorithm works!
---

# Technical Details

## Technology Overview

### Frontend

The frontend itself is a simple single-page, statically built [Vite](https://vite.dev/) application using Typescript and React. Here are some libraries that are used:
- **Tailwind** - used for styling, preferred over CSS
- [**React Router**](https://reactrouter.com/en/main) - simple browser-level routing
- [**ESLint**](https://eslint.org/) - static analysis and linting
- [**Zustand**](https://github.com/pmndrs/zustand) - simple global store/state management solution

Some other minor libraries were also used to do specific things:
- [**dnd kit**](https://dndkit.com/) - the main drag and drop functionality
- **React Hook Form** - form state management library
- **React Range** - used to create range sliders
- **Tailwind Merge** - simple library that will merge different tailwind class strings together

### Backend

The backend of the application is a [Golang](https://go.dev/) application, built primarily off of the [Gin](https://github.com/gin-gonic/gin) HTTP framework. This backend communicates with a [MongoDB](https://mongodb.com/) cluster, which serves as the primary database. Here are some other libraries that were used:
- [**godotenv**](https://github.com/joho/godotenv) - used to load environmental variables in a development environment
- [**MongoDB Go Driver**](https://www.mongodb.com/docs/drivers/go/current/) - SDK for connecting to mongodb
- [**Air**](https://github.com/air-verse/air) - hot-reloading of the app

### Development

The development was done through [Docker Compose](https://docs.docker.com/compose/), which allowed us to run both the frontend and backend simultaneously from any computer that had Docker installed. These instances of the app would also be able to communicate and had hot-reloading built in. There is also a separate version of the Docker Compose file that instantiates a local mongodb container (instead of using a MongoDB Atlas cluster).

### Deployment

The main deployment Dockerfile has 3 main stages:
1. **Static client builder** - This will take our Vite app and build it statically, creating HTML/JS files to serve to the user
2. **Backend compilation** - Due to the nature of Golang, we can actually compile the entire backend into a single executable, with all libraries included! We copy over the public folder from the first stage to ensure no compilation errors and build a standalone go app
3. **Final scratch container** - We start with the empty [scratch](https://hub.docker.com/_/scratch) docker container to minimize the final executable size. We only need to copy over the output files from each stage and SSL certificates (bc scratch container literally has nothing). We set the entrypoint to be our jury executable, and we're done building!

This built container can now be served anywhere that can host a single docker image. This is extremely convenient as the [Digital Ocean App Platform](https://www.digitalocean.com/products/app-platform) allows us to do just that -- run a single docker image. 

## Rank Aggregation Algorithm

[Borda Count](https://en.wikipedia.org/wiki/Borda_count)

:::warning[WIP]
🚧 The technical details page is still under construction and will be split into sub-pages 🚧
:::
