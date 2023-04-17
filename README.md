# Jury

Table of Contents
-----------------
- [Introduction](#introduction)
- [Deployment](#deployment)
  - [With Docker](#with-docker)
- [Development](#development)
  - [With Docker (Recommended!)](#with-docker-recommended)
  - [Manual Installation](#manual-installation)
- [License](#license)

## Introduction

Jury is a project aimed at creating a new pairwise judging system using modern technologies to optimize the user experience for judges and admin users. The inspiration for this project is [Gavel by anishathalye](https://github.com/anishathalye/gavel), and [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) provides more details on the underlying formulas of Gavel. The majority of our algorithm is based on research done in the field of [pairwise comparison](https://en.wikipedia.org/wiki/Pairwise_comparison).

To see the designs for this project, check out the [Figma page](https://www.figma.com/file/qwBWs4i7pJMpFbcjMffDZU/Jury-(Gavel-Plus)?node-id=8%3A100&t=egLV7iVmwvNRXef5-1).

## Deployment

### With Docker

To deploy the application using Docker, run `docker compose up`.

## Development

### With Docker (Recommended!)

Requirements:

* [Docker](https://www.docker.com/)

Copy `.env.template` into `.env` and fill in the environmental variables.

To start the development environment, simply run `docker compose -f docker-compose.dev.yml up` and open the page at `localhost:3000`.

### Manual Installation

Requirements:

* [yarn](https://yarnpkg.com/)
* [cargo](https://doc.rust-lang.org/cargo/)
* [GNU Scientific Library (GSL)](https://www.gnu.org/software/gsl/)

Copy `.env.template` into `.env` and fill in the environmental variables. You will also need a `MONGODB_URI` field.
Additionally, copy `client/.env.template` into `client/.env` and copy over the relevant environmental variables from `.env`.
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

Note that manual installation is not recommended, and the Docker setup is the preferred development environment.

## License

This project is licensed under the [MIT License](LICENSE).
