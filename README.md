# Jury

[WIP] A project designed to create a new pairwise judging system using modern technologies aimed at optimizing the user experience of judges and admin users. See the inspiration for this project: [Gavel by anishathalye](https://github.com/anishathalye/gavel). Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying formulas of Gavel!

# External System Requirements

> Note: This might be mitigated with the use of a dockerized development environment (WIP)

* [GNU Scientific Library (GSL)](https://www.gnu.org/software/gsl/)

# Deployment

## With Docker

Run `docker compose up`

# Developing

## With Docker

<!-- **RECOMMENDED** -->
**WIP**

Requirements:

* [Docker](https://www.docker.com/)

Copy `.env.template` into `.env` and fill in the environmental variables

Simply run `docker compose -f docker-compose.dev.yml up` and open the page at `localhost:<PORT>`

## Manual Installation

Requirements:

* [yarn](https://yarnpkg.com/)
* [cargo](https://doc.rust-lang.org/cargo/)

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
