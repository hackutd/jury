---
sidebar_position: 2
---

# Internal API

All routes are found at the `http://localhost:8000/api` endpoint, which can be accessed on the frontend through the `VITE_JURY_URL` environmental variable.

An example fetch request (The `credentials: include` allow for authentication cookies to be sent over):

```js
const fetchedProjects = await fetch(`${import.meta.env.VITE_JURY_URL}/project/list`, {
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