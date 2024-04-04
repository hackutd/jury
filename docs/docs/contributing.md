---
sidebar_position: 3
---

# Contributing

The entire development environment is set up using [Docker Compose](https://docs.docker.com/compose/). You simply need to have Docker and Docker Compose [installed](https://docs.docker.com/compose/install/).

## With Docker (Recommended!)

Copy `.env.template` into `.env` and fill in the environmental variables (see [environmental variables](/docs/intro#step-3-environmental-variables)). Simply run the following command to start up the dev server:

```bash
docker compose -f docker-compose.dev.yml up
```

This may take a while to install all files and build, but you should see the following messages in the log when they are done:

Frontend:
```
jury-dev-frontend  | $ vite
jury-dev-frontend  | 
jury-dev-frontend  |   VITE v5.2.2  ready in 327 ms
jury-dev-frontend  | 
jury-dev-frontend  |   ➜  Local:   http://localhost:3000/
jury-dev-frontend  |   ➜  Network: http://172.18.0.3:3000/
```

Backend (will show ALL routes):
```
jury-dev-backend   | [GIN-debug] POST   /api/admin/categories     --> server/router.SetCategories (7 handlers)
jury-dev-backend   | [GIN-debug] GET    /api/categories           --> server/router.GetCategories (7 handlers)
jury-dev-backend   | [GIN-debug] GET    /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (7 handlers)
jury-dev-backend   | [GIN-debug] HEAD   /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (7 handlers)
```

### [Alternative] Offline database development

If you don't want to create a MongoDB instance or don't want to work online (this only is viable if you've downloaded the packages before), you can use the offline development version. As with above, you will need to create your own `.env` file (again see [environmental variables](/docs/intro#step-3-environmental-variables)). Then, run the following command:

```bash
docker compose -f docker-compose-mongo.dev.yml up
```

As with above, you will need to wait for both the backend and frontend to load. Additionally, you will have to wait for the MongoDB database to intialize and create a replica set (takes more than 10 seconds).

## Manual Installation

Although the easiest development environment is through docker compose, you may manually start the frontend and backend separately. 

Requirements:

-   [yarn](https://yarnpkg.com/)
-   [go](https://go.dev/)

Copy `.env.template` into `.env` and fill in the environmental variables (same as above). Additionally, copy `client/.env.template` into `client/.env` and copy over the relavent environmental variables from `.env`. This is used to expose the correct environmental variables to the running instance of the Vite frontend as you will not have Docker compose to automatically do that for you.

Client dev server (PORT 3000):

```
cd client
yarn install
yarn start
```

Backend dev server (PORT 8000):

```
go run .
```

### Connecting to the MongoDB Container (running locally)

Use `mongodb://{MONGODB_USER}:{MONGODB_PASS}@0.0.0.0:27107/?directConnection=true` as the connection string. This works with both MongoDB Compass and through `mongosh`.

To reset the database and reload the sample data, completely remove the `data` folder that the Docker compose creates.

:::info[Note]
If you are using MongoDB Atlas, your connection string will be different. It is easier to see the data from the [MongoDB Atlas Console](https://cloud.mongodb.com/). MongoDB Atlas has connection instructions from `mongosh` [here](https://www.mongodb.com/docs/atlas/mongo-shell-connection/).
:::

## Other Resources

### Design File

Here is the Figma design file: https://www.figma.com/file/qwBWs4i7pJMpFbcjMffDZU/Jury-(Gavel-Plus)?node-id=8%3A100&t=xYwfPwRAUeJw9jNr-1
