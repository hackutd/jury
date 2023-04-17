# Jury

A project designed to create a new pairwise judging system using modern technologies aimed at optimizing the user experience of judges and admin users. See the inspiration for this project: [Gavel by anishathalye](https://github.com/anishathalye/gavel). Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying formulas of Gavel! The majority of our algorithm will be based off of research done in the field of [pairwise comparison](https://en.wikipedia.org/wiki/Pairwise_comparison). 

# Deployment

## With Docker

Run `docker compose up`

# Developing

## With Docker (Recommended!)

Requirements:

* [Docker](https://www.docker.com/)

Copy `.env.template` into `.env` and fill in the environmental variables

Simply run `docker compose -f docker-compose.dev.yml up` and open the page at `localhost:3000`

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