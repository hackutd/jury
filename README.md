# Jury

A project designed to create a new pairwise judging system using modern technologies aimed at optimizing the user experience of judges and admin users. See the inspiration for this project: [Gavel by anishathalye](https://github.com/anishathalye/gavel). Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying formulas of Gavel! The majority of our algorithm will be based off of research done in the field of [pairwise comparison](https://en.wikipedia.org/wiki/Pairwise_comparison). 

# Deployment

## Setup

Copy `.env.template` into `.env` and fill in the environmental variables

Environmental Variables:

```
JURY_NAME="Name of the jury app [Displays on the app!]"
JURY_ADMIN_PASSWORD="Password used to log into the admin dashboard"

MONGODB_URI="MongoDB connection URI string [ONLY for MongoDB Atlas]"
MONGODB_USER="Username for local mongo container [ONLY use if running local mongo instance]"
MONGODB_PASS="Password for local mongo container [ONLY use if running local mongo instance]"
```

I suggest you run the app with MongoDB Atlas! Create a free account and database [here](https://www.mongodb.com/atlas/database). It should provide you a URI string to fill into the `.env` file.

If you would rather use a local instance deployed with docker-compose, you can simply fill in the username and password you want to use with that database.

## With Docker

Run `docker compose up` after configuring the `.env` file. If you want to run mongo locally, run `docker compose -f docker-compose-mongo.yml`.

# Developing

## With Docker (Recommended!)

Requirements:

* [Docker](https://www.docker.com/)

Copy `.env.template` into `.env` and fill in the environmental variables (see above).

Simply run `docker compose -f docker-compose.dev.yml up` and open the page at `localhost:3000`.

If you want to run mongo locally, run `docker compose -f docker-compose-mongo.dev.yml`.

### Connecting to the MongoDB Container

Use `mongodb://{MONGODB_USER}:{MONGODB_PASS}@0.0.0.0:27107/?directConnection=true` as the connection string. This works with both MongoDB Compass and through `mongosh`.

To reset the database and reload the sample data, completely remove the `data` folder that the Docker compose creates.

## Manual Installation

Requirements:

* [yarn](https://yarnpkg.com/)
* [cargo](https://doc.rust-lang.org/cargo/)
* [GNU Scientific Library (GSL)](https://www.gnu.org/software/gsl/)

Copy `.env.template` into `.env` and fill in the environmental variables. You will also need a `MONGODB_URI` field.
Additionally, copy `client/.env.template` into `client/.env` and copy over the relavent environmental variables from `.env`.
This is used to expose the correct environmental variables to the running instance of CRA.

Client dev server (PORT 3000):

```
cd client
yarn install
yarn start
```

Backend dev server (PORT 8000):

```
cargo run
```

## Design File

Here is the Figma design file: https://www.figma.com/file/qwBWs4i7pJMpFbcjMffDZU/Jury-(Gavel-Plus)?node-id=8%3A100&t=xYwfPwRAUeJw9jNr-1