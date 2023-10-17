# Jury

A project designed to create a new pairwise judging system using modern technologies aimed at optimizing the user experience of judges and admin users. See the inspiration for this project: [Gavel by anishathalye](https://github.com/anishathalye/gavel). Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying formulas of Gavel! The majority of our algorithm will be based off of research done in the field of [pairwise comparison](https://en.wikipedia.org/wiki/Pairwise_comparison).

# One Click Deploy!

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/acmutd/jury/tree/master&refcode=de08cdf58df3)

## 0. Set up MongoDB Atlas

MongoDB is a **free** NoSQL database that is hosted in the cloud! Their website has great instructions on how to quickly create an account and 100% lifetime free database, perfect for Jury: https://www.mongodb.com/docs/atlas/getting-started/.

Follow the instructions up to step 5: connect to your cluster. You will get a URI string that looks like `mongodb+srv://<username>:<password>@<something>.mongodb.net/?retryWrites=true&w=majority`, except \<username\>, \<password\>, and \<something\> should be filled in. You will need that for the environmental variables.

## 1. Resources

- Click Edit Plan
- Change plan to "Basic"
- Change Instance Size to "$5.00/mo -- Basic (512 MB RAM | 1 vCPU)"

Note that since Digital Ocean charges prorated, this will literally be like less than 20 cents a day.

## 2. Environmental Variables

Click on "edit" next to "jury-service" and fill in the env as follows:

- REACT_APP_JURY_NAME - Name of the app to display to the user (eg. HackUTD X Judging)
- JURY_ADMIN_PASSWORD - Password for the admin portal (suggestion: use 2 random words)
- MONGODB_URI - URI of MongoDB Atlas instance (see step 0)

### Email hosting

If using Gmail SMTP:
- EMAIL_HOST = smtp.gmail.com
- EMAIL_FROM = gmail username
- EMAIL_USERNAME = gmail username
- EMAIL_PASSWORD = google app password

If using SendGrid:
- SENDGRID_API_KEY = API key provided by SendGrid
- EMAIL_FROM = email to send from
- EMAIL_FROM_NAME = display name of sender

If using [AWS SES](https://docs.aws.amazon.com/ses/latest/dg/smtp-credentials.html):
- EMAIL_HOST = email-smtp.<region>.amazonaws.com
- EMAIL_FROM = verified email to send from
- EMAIL_USERNAME = SES SMTP username
- EMAIL_PASSWORD = SES SMTP password

**NOTE: Only set `SENDGRID_API_KEY` if using Sendgrid!!**

## 3. Info

You can change the start of the URL if you want:
- App Info > edit
- Change app name (URL will be <app_name><some_characters>.ondigitalocean.app)
- Click Next

# Deploy Manually

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

> If you would rather use a local instance deployed with docker-compose, you can simply fill in the username and password you want to use with that database. You will use the files `docker-compose-mongo.yml` for production and `docker-compose-mongo.dev.yml` for development. This is **not recommended** as atlas provides a much easier interface!!

## With Docker

Run `docker compose up` after configuring the `.env` file.

# Contributing

To run the dev environment:

## With Docker (Recommended!)

Copy `.env.template` into `.env` and fill in the environmental variables (see above).

Simply run `docker compose -f docker-compose.dev.yml up` and open the page at `localhost:3000`.

### Connecting to the MongoDB Container

Use `mongodb://{MONGODB_USER}:{MONGODB_PASS}@0.0.0.0:27107/?directConnection=true` as the connection string. This works with both MongoDB Compass and through `mongosh`.

To reset the database and reload the sample data, completely remove the `data` folder that the Docker compose creates.

## Manual Installation

Requirements:

-   [yarn](https://yarnpkg.com/)
-   [cargo](https://doc.rust-lang.org/cargo/)
-   [GNU Scientific Library (GSL)](https://www.gnu.org/software/gsl/)

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

## Backend Routes

All routes are found at the `http://localhost:8000/api` endpoint, which can be accessed on the frontend through the `REACT_APP_JURY_URL` environmental variable.

An example fetch request (The `credentials: include` allow for authentication cookies to be sent over):

```js
const fetchedProjects = await fetch(`${process.env.REACT_APP_JURY_URL}/project/list`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
}).then((data) => data.json());
```

| Path                 | Method | Auth  | Description                             |
| -------------------- | ------ | ----- | --------------------------------------- |
| /judge/list          | GET    | admin | get list of all judges                  |
| /judge/new           | POST   | admin | Add a new judge                         |
| /judge/csv           | POST   | admin | Preview a judges csv                    |
| /judge/csv/upload    | POST   | admin | Upload judges using csv                 |
| /judge/stats         | GET    | admin | Get the stats for the add judges page   |
| /judge/login         | POST   |       | Login request for a judge               |
| /judge/auth          | POST   | token | Checks to see if judge is logged in     |
| /judge/welcome       | GET    | token | Checks for `read_welcome` for a judge   |
| /judge/welcome       | PUT    | token | Set `read_welcome` to true for a judge  |
| /project/list        | GET    | admin | get list of all projects                |
| /project/new         | POST   | admin | Add a new project                       |
| /project/csv         | POST   | admin | Preview a projects csv                  |
| /project/csv/upload  | POST   | admin | Upload projects using csv               |
| /project/stats       | GET    | admin | Get the stats for the add projects page |
| /project/devpost     | POST   | admin | Upload a Devpost CSV                    |
| /project/:id         | GET    | token | Get project by ID                       |
| /admin/login         | POST   |       | Log into the admin dashboard            |
| /admin/stats         | GET    | admin | Get all stats                           |
| /admin/sync          | GET    | admin | Establish event stream with server      |
| /admin/clock         | GET    | admin | Gets the current clock state            |
| /admin/clock/pause   | POST   | admin | Pauses the clock                        |
| /admin/clock/unpause | POST   | admin | Resumes the clock                       |
| /admin/clock/reset   | POST   | admin | Resets the clock                        |

### GET /judge/list

Get a list of all judges

**Response**

```json
[
    {
        "id": "ObjectId",
        "token": "String",
        "code": "String",
        "name": "String",
        "email": "String",
        "active": "bool",
        "last_activity": "DateTime",
        "read_welcome": "bool",
        "notes": "String",
        "next": "ObjectId",
        "prev": "ObjectId",
        "alpha": "f64",
        "beta": "f64"
    }
    // ...
]
```

### POST /judge/new

_Requires admin password_ | Add a singular new judge

**Request Body**

```json
{
    "name": "String",
    "email": "String",
    "notes": "String"
}
```

**Response**

-   `200 Ok`
-   Error + String

### POST /judge/csv

_Requires admin password_ | PREVIEW judge CSV file (does not upload)

**Request Body**

```
Content-Type: multipart/form-data

Fields:
- csv: CSV file upload
- hasHeader: Bool that is true if csv contains header
```

**Response**

```json
[
    {
        "name": "String",
        "email": "String",
        "notes": "String"
    }
    // ...
]
```

or Error + String

### POST /judge/csv/upload

_Requires admin password_ | **Upload** judge CSV file

**Request Body**

```
Content-Type: multipart/form-data

Fields:
- csv: CSV file upload
- hasHeader: Bool that is true if csv contains header
```

**Response**

-   `200 Ok` + Number of added judges as String
-   Error + String

### GET /judge/stats

_Requires admin password_ | Get the stats to display on the add judges page

**Response**

```json
{
    "num": 0, // u64 - Total # of judges
    "alpha": 0.0, // f64 - Average alpha of judges
    "beta": 0.0, // f64 - Average beta of judges
},
```

or Error + String

### POST /judge/login

Login request for a judge

**Request Body**

```json
{
    "code": "String (6-digit login code for judge)"
}
```

**Response**

-   `200 Ok` + Token as String
-   Error + String

### POST /judge/auth

_Requires judge token_ | Checks to see if judge is logged in

**Response**

-   `200 Ok` if token is valid
-   `401 Unauthorized` if token is invalid
-   `500 Internal Server Error`

### GET /judge/welcome

_Requires judge token_ | Checks to see if a judge has read the welcome page

**Response**

```json
{
    "ok": true // Boolean, true or false
}
```

or Error

### POST /judge/welcome

_Requires judge token_ | Set the `read_welcome` field for a judge to true

**Response**

-   `202 Accepted`
-   Error + String

### GET /project/list

Get a list of all projects

**Response**

```json
[
    {
        "id": "ObjectId",
        "name": "String",
        "location": "u64",
        "description": "String",
        "try_link": "String",
        "video_link": "String",
        "challenge_list": ["String"],
        "seen": "u64",
        "votes": "u64",
        "mu": "f64",
        "sigma_sq": "f64",
        "active": "bool",
        "prioritized": "bool",
        "last_activity": "DateTime"
    }
    // ...
]
```

### POST /project/new

_Requires admin password_ | Add a singular new project

**Request Body**

```json
{
    "name": "String",
    "description": "String",
    "try_link": "String (Optional)",
    "video_link": "String (Optional)",
    "challenge_list": "String (Optional)"
}
```

**Response**

-   `200 Ok`
-   Error + String

### POST /project/csv

_Requires admin password_ | PREVIEW project CSV file (does not upload)

**Request Body**

```
Content-Type: multipart/form-data

Fields:
- csv: CSV file upload
- hasHeader: Bool that is true if csv contains header
```

**Response**

```json
[
    {
        "name": "String",
        "location": "String",
        "description": "String"
    }
    // ...
]
```

or Error + String

### POST /project/csv/upload

_Requires admin password_ | **Upload** project CSV file

**Request Body**

```
Content-Type: multipart/form-data

Fields:
- csv: CSV file upload
- hasHeader: Bool that is true if csv contains header
```

**Response**

-   `200 Ok` + Number of added judges as String
-   Error + String

### POST /project/devpost

_Requires admin password_ | Upload CSV from devpost

**Request Body**

```
Content-Type: multipart/form-data

Fields:
- csv: CSV file upload
```

**Response**

-   Ok 200
-   Error + String

### GET /project/:id

_Requires judge token_ | Gets a project by ID

* id - ID of the project

**Response**

-   Ok 200 with project data as JSON
-   Error

### POST /admin/login

Login request for admins

**Request Body**

```json
{
    "password": "String (admin password)"
}
```

**Response**

-   `200 Ok` + Password as String
-   Error + String

### GET /admin/stats

_Requires admin password_ | Get the stats to display on the admin dashboard

**Response**

```json
{
    "projects": 0, // u64 - Total # of projects
    "seen": 0.0, // f64 - Average seen # of projects
    "votes": 0.0, // f64 - Average votes of projects
    "time": 0, // u64 - Current judging time in milliseconds
    "avg_mu": 0.0, // f64 - Average mu value of projects
    "avg_sigma": 0.0, // f64 - Average sigma^2 value of projects
    "judges": 0 // u64 - Total # of judges
},
```

or Error

### GET /admin/sync

_Requires admin password_ | Establish event stream with server

**Request**

-   [EventStream](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)

**Response**

```
Content-Type: text/event-stream
```

### GET /admin/clock

_Requires admin password_ | Get the current clock state

**Response**

```json
{
    "start": 0, // Start time in millis
    "prev": 0, // Previously elapsed time, relavent if paused/previously paused
    "paused": true // Bool, true if paused
}
```

or Error

### POST /admin/clock/pause

_Requires admin password_ | Pauses the clock

**Response**

-   Ok 200
-   InternalServerError 500 + String

### POST /admin/clock/unpause

_Requires admin password_ | Resumes the clock

**Response**

-   Ok 200
-   InternalServerError 500 + String

### POST /admin/clock/reset

_Requires admin password_ | Resets the clock

**Response**

-   Ok 200
-   InternalServerError 500 + String
