---
sidebar_position: 2
---

# Internal API

See the [backend API page](/docs/details/backend/api) for an explanation about how the API code is structured. All routes are found at the `http://localhost:8000/api` endpoint in development and the `/api` endpoint in production, which can be accessed on the frontend through the `VITE_JURY_URL` environmental variable.

An example request looks like the following:

```js
const projRes = await getRequest<Project[]>('/project/list', 'admin');
if (projRes.status !== 200) {
    errorAlert(projRes);
    return;
}
set({ projects: projRes.data as Project[] });
```

Note that this uses the `getRequest` and `errorAlert` frontend methods, which are described in the [calling backend API page](/docs/details/frontend/api).

## All Routes

All routes are listed in `server/router/init.go` with their respective handlers.

| Path                     | Method | Auth  | Description                                     |
| ------------------------ | ------ | ----- | ----------------------------------------------- |
| /                        | GET    |       | Heartbeat route                                 |
| /judge/login             | POST   |       | Login judge                                     |
| /admin/login             | POST   |       | Log into the admin dashboard                    |
| /judge/auth              | POST   | judge | Checks to see if judge is logged in             |
| /admin/auth              | POST   | admin | Checks to see if admin is logged in             |
| /judge/new               | POST   | admin | Add a new judge                                 |
| /judge/csv               | POST   | admin | Add judges by CSV                               |
| /judge/list              | GET    | admin | Get list of all judges                          |
| /judge/:id               | DELETE | admin | Deletes a judge by ID                           |
| /judge/:id               | PUT    | admin | Edit judge info                                 |
| /admin/groups/swap       | POST   | admin | Swaps the judge groups manually                 |
| /admin/qr                | POST   | admin | Generate add judge QR code                      |
| /admin/qr/:track         | POST   | admin | Generate add judge to track QR code             |
| /qr                      | GET    |       | Gets add judge QR code                          |
| /qr/:track               | GET    |       | Gets add judge to track QR code                 |
| /qr/add                  | POST   |       | Add judge from QR code                          |
| /project/new             | POST   | admin | Add a new project                               |
| /project/devpost         | POST   | admin | Upload a Devpost CSV                            |
| /project/csv             | POST   | admin | Add projects by CSV                             |
| /project/list            | GET    | admin | get list of all projects                        |
| /project/:id             | DELETE | admin | Delete project by ID                            |
| /project/:id             | PUT    | admin | Edit project info                               |
| /admin/stats             | GET    | admin | Get all stats                                   |
| /admin/stats/:track      | GET    | admin | Get all stats for a track                       |
| /project/stats           | GET    | admin | Get the stats for projects                      |
| /judge/stats             | GET    | admin | Get the stats for judges                        |
| /admin/flags             | GET    | admin | Gets all flags                                  |
| /admin/clock             | GET    | admin | Gets the current clock state                    |
| /admin/clock/pause       | POST   | admin | Pauses the clock                                |
| /admin/clock/unpause     | POST   | admin | Resumes the clock                               |
| /admin/clock/backup      | POST   | admin | Backs up the clock to the database              |
| /admin/started           | GET    |       | Check if the clock is running                   |
| /admin/clock/reset       | POST   | admin | Resets the clock                                |
| /admin/reset             | POST   | admin | Resets the entire database                      |
| /project/reassign        | POST   | admin | Reassign all project table numbers              |
| /project/balance-groups  | POST   | admin | Balances project group numbers                  |
| /project/reassign        | POST   | admin | Reassign all project table numbers              |
| /admin/timer             | GET    | judge | Gets the judge timer length                     |
| /admin/options           | GET    | admin | Gets all config options set                     |
| /admin/options           | POST   | admin | Sets config options                             |
| /admin/num-groups        | POST   | admin | Sets num of groups and reassigns nums           |
| /admin/group-sizes       | POST   | admin | Sets the size of groups and reassigns nums      |
| /admin/block-reqs        | POST   | admin | Sets whether to block login requests            |
| /admin/max-reqs          | POST   | admin | Sets the maximum number of logins/min           |
| /admin/export/judges     | GET    | admin | Exports judges as a CSV                         |
| /admin/export/projects   | GET    | admin | Exports projects as a CSV                       |
| /admin/export/challenges | GET    | admin | Exports projects by challenge as ZIP of CSVs    |
| /admin/export/rankings   | GET    | admin | Exports a list of rankings for each judge       |
| /judge/hide/:id          | PUT    | admin | Hides a judge                                   |
| /project/hide/:id        | PUT    | admin | Hides a project                                 |
| /judge/move/:id          | PUT    | admin | Moves a judge to a different group              |
| /project/move/:id        | PUT    | admin | Moves a project to a different group            |
| /project/prioritize/:id  | PUT    | admin | Prioritizes a project                           |
| /project/prioritize      | POST   | admin | Prioritizes multiple projects                   |
| /judge/hide              | POST   | admin | Hides multiple judges                           |
| /project/hide            | POST   | admin | Hides multiple projects                         |
| /judge/move              | POST   | admin | Moves multiple judges to a different group      |
| /project/move            | POST   | admin | Moves multiple projects to a different group    |
| /admin/flag/:id          | DELETE | admin | Removes a flag                                  |
| /admin/deliberation      | POST   | admin | Toggles deliberation mode                       |
| /admin/log               | GET    | admin | Gets the audit log                              |
| /judge                   | GET    | judge | Gets judge from token cookie                    |
| /judge/welcome           | GET    | judge | Checks for `read_welcome` for a judge           |
| /judge/welcome           | PUT    | judge | Set `read_welcome` to true for a judge          |
| /judge/projects          | GET    | judge | Gets the list of projects a judge has seen      |
| /judge/next              | POST   | judge | Get next project for judge to view              |
| /judge/skip              | POST   | judge | Skips the current project with a reason         |
| /judge/finish            | POST   | judge | Finish viewing a project                        |
| /judge/rank              | POST   | judge | Update judge rankings                           |
| /judge/star/:id          | PUT    | judge | Update star ranking for a project               |
| /judge/break             | POST   | judge | Removes active project for judge (take a break) |
| /judge/notes/:id         | PUT    | judge | Update notes for a project                      |
| /project/:id             | GET    | judge | Gets a project by ID                            |
| /project/count           | GET    | judge | Gets the total number of projects               |
| /judge/project/:id       | GET    | judge | Gets a judged project by a judge                |
| /judge/deliberation      | GET    | judge | Returns if deliberation mode is on              |
| /project/list/public     | GET    |       | Gets a list of all projects for expo            |
| /challenges              | GET    |       | Gets a list of all challenges                   |
| /group-names             | GET    |       | Gets a list of all group names                  |

## Response Types

### OK Response

Will return either as a success or failure. This will be represented with the following JSON, where 0 is failure and 1 is success.

```json
{
    "ok": 1
}
```

### Error Response

All errors will be formatted in the following JSON:

```json
{
    "error": "error string"
}
```

## Default Routes

Errors will return with a non-success status (generally 4xx).

### GET /

Heartbeat route

- **Auth**: none
- **Response**: OK response

## Login Routes

### POST /judge/login

Login judge

- **Auth**: none
- **Body**: JSON

```json
{
    "code": "String | login code"
}
```

- **Response**: JSON

```json
{
    "token": "String | judge login token"
}
```

### POST /admin/login

Log into the admin dashboard                    

- **Auth**: none
- **Body**: JSON

```json
{
    "password": "String | admin password"
}
```

- **Response**: OK response

### POST /judge/auth

Checks to see if judge is logged in

- **Auth**: judge
- **Response**: OK response

### POST /admin/auth

Checks to see if admin is logged in

- **Auth**: admin
- **Response**: OK response

## Admin Panel (Judges) Routes

### POST /judge/new

Add a new judge

- **Auth**: admin
- **Body**: JSON

```json
{
    "name": "String",
    "email": "String",
    "notes": "String"
}
```

- **Response**: OK response

### POST /judge/csv

Add judges by CSV

- **Auth**: admin
- **Body**: FormData
    - `csv`: CSV file
    - `hasHeader`: Boolean, true if CSV has a header
    - `noSend`: Don't send email to judge if true
- **Response**: OK response

### GET /judge/list

Get list of all judges

- **Auth**: admin
- **Response**: JSON | list of judges

```json
[
    {
        "id": "ObjectId",
        "token": "String",
        "code": "String",
        "name": "String",
        "email": "String",
        "active": "bool",
        "track": "String",
        "group": "String",
        "read_welcome": "bool",
        "notes": "String",
        "current": "ObjectId",
        "last_location": "int",
        "seen": "int",
        "group_seen": "int",
        "seen_projects": [
            {
                "project_id": "ObjectID",
                "starred": "bool",
                "notes": "String",
                "name": "String",
                "location": "int",
                "description": "String"
            }
        ],
        "rankings": [
            "ObjectId",
            "ObjectId | and so on for each ranked project"
        ],
        "last_activity": "DateTime"
    }
]
```

### DELETE /judge/\:id

Deletes a judge by ID

- **Auth**: admin
- **Parameter**: ID, the ID of the judge to delete
- **Response**: OK response

### PUT /judge/\:id

Edit judge info

- **Auth**: admin
- **Body**: JSON

```json
{
    "name": "String",
    "email": "String",
    "notes": "String"
}
```

- **Response**: OK response

### POST /admin/groups/swap

Swaps the judge groups manually

- **Auth**: admin
- **Response**: OK response

### POST /admin/qr

Generate add judge QR code

- **Auth**: admin
- **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### POST /admin/qr/\:track

Generate add judge to track QR code

- **Auth**: admin
- **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### GET /qr

Gets add judge QR code

- **Auth**: none
- **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### GET /qr/\:track

Gets add judge to track QR code

- **Auth**: none
- **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### POST /qr/add

Add judge from QR code

- **Auth**: none
- **Body**: JSON

```json
{
    "name": "String",
    "email": "String",
    "notes": "String",
    "code": "String | QR code token"
}
```

- **Response**: OK response

## Admin Panel (Projects) Routes

### POST /project/new

Add a new project

- **Auth**: admin
- **Body**: JSON

```json
{
    "name": "String",
    "description": "String",
    "url": "String",
    "try_link": "String",
    "video_link": "String",
    "challenge_list": "String | comma-separated list of challenges/tracks",
}
```

- **Response**: OK response

### POST /project/devpost

Upload a Devpost CSV

- **Auth**: admin
- **Body**: FormData
    - `csv`: CSV file
- **Response**: 

### POST /project/csv

Add projects by CSV

- **Auth**: admin
- **Body**: FormData
    - `csv`: CSV file
    - `hasHeader`: Boolean, true if CSV has a header
- **Response**: OK response

### GET /project/list

get list of all projects

- **Auth**: admin
- **Response**: JSON

```json
[
    {
        "id": "ObjectId",
        "name": "String",
        "location": "u64",
        "description": "String",
        "url": "String",
        "try_link": "String",
        "video_link": "String",
        "challenge_list": ["String"],
        "seen": "int",
        "track_seen": {
            "track1": "int",
            "track2": "int",
        },
        "score": "int",
        "stars": "int",
        "track_stars": {
            "track1": "int",
            "track2": "int",
        },
        "active": "bool",
        "prioritized": "bool",
        "group": "int",
        "last_activity": "DateTime"
    }
]
```

### DELETE /project/\:id

Delete project by ID

- **Auth**: admin
- **Parameter**: ID of project to delete
- **Response**: OK response

### PUT /project/\:id

Edit project info

- **Auth**: admin
- **Body**: JSON

```json
{
    "name": "String",
    "description": "String",
    "url": "String",
    "try_link": "String",
    "video_link": "String",
    "challenge_list": "String | comma-separated list of challenges/tracks",
}
```

- **Response**: OK response

## Admin Panel (Stats/Data) Routes

### GET /admin/stats

Get all stats

- **Auth**: admin
- **Response**: JSON

```json
{
    "projects": "int",
    "avg_project_seen": "int",
    "avg_judge_seen": "int",
    "judges": "int"
}
```

### GET /admin/stats/\:track

Get all stats for a track

- **Auth**: admin
- **Parameter**: Track, the track to fetch stats for
- **Response**: JSON

```json
{
    "projects": "int",
    "avg_project_seen": "int",
    "avg_judge_seen": "int",
    "judges": "int"
}
```

### GET /project/stats

Get the stats for projects

- **Auth**: admin
- **Response**: JSON

```json
{
    "num": "int",
    "num_active": "int",
    "avg_seen": "float"
}
```

### GET /judge/stats

Get the stats for judges

- **Auth**: admin
- **Response**: JSON

```json
{
    "num": "int",
    "num_active": "int",
    "avg_seen": "float"
}
```

### GET /admin/flags

Gets all flags

- **Auth**: admin
- **Response**: JSON List

```json
[
    {
        "id": "ObjectID",
        "project_id": "ObjectID",
        "judge_id": "ObjectID",
        "time": "DateTime",
        "project_name": "String",
        "project_location": "int",
        "judge_name": "String",
        "reason": "String"
    }
]
```

## Admin Panel (Clock) Routes

### GET /admin/clock

Gets the current clock state

- **Auth**: admin
- **Response**: JSON

```json
{
    "running": "bool",
    "time": "int"
}
```

### POST /admin/clock/pause

Pauses the clock

- **Auth**: admin
- **Response**: OK response

### POST /admin/clock/unpause

Resumes the clock

- **Auth**: admin
- **Response**: OK response

### POST /admin/clock/backup

Backs up the clock to the database

- **Auth**: admin
- **Response**: OK response

### GET /admin/started

Check if the clock is running

- **Auth**: admin
- **Response**: OK response (1 if true)

## Admin Options/Settings Routes

### POST /admin/clock/reset

Resets the clock

- **Auth**: admin
- **Response**: OK response

### POST /admin/reset

Resets the entire database

- **Auth**: admin
- **Response**: OK response

### POST /project/reassign

Reassign all project table numbers

- **Auth**: admin
- **Response**: OK response

### POST /project/balance-groups

Balances project group numbers

- **Auth**: admin
- **Response**: OK response

### POST /project/reassign

Reassign all project table numbers

- **Auth**: admin
- **Response**: OK response

### GET /admin/timer

Gets the judge timer length

- **Auth**: admin
- **Response**: JSON

```json
{
    "judging_timer": "int"
}
```

### GET /admin/options

Gets all config options set

- **Auth**: admin
- **Response**: JSON

```json
{
    "id": "ObjectID",
    "ref": "int",
    "clock": {
        "start_time": "int",
        "pause_time": "int",
        "running": "bool"
    },
    "judging_timer": "int",
    "min_views": "int",
    "clock_sync": "bool",
    "deliberation": "bool",
    "judge_tracks": "bool",
    "tracks": ["String"],
    "multi_group": "bool",
    "num_groups": "int",
    "group_sizes": ["int"],
    "switching_mode": "String",
    "auto_switch_prop": "float",
    "manual_switches": "int",
    "qr_code": "String",
    "track_qr_codes": {
        "track1": "String",
        "track2": "String"
    },
    "group_names": ["String"],
    "ignore_tracks": ["String"],
    "max_req_per_min": "int",
    "block_reqs": "bool"
}
```

### POST /admin/options

Sets config options

- **Auth**: admin
- **Body**: JSON

```json
{
    "judging_timer": "int",
    "min_views": "int",
    "clock_sync": "bool",
    "judge_tracks": "bool",
    "tracks": ["String"],
    "multi_group": "bool",
    "num_groups": "int",
    "group_sizes": ["int"],
    "switching_mode": "String",
    "auto_switch_prop": "float",
    "group_names": ["String"],
    "ignore_tracks": ["String"],
    "max_req_per_min": "int",
    "block_reqs": "bool",
}
```

- **Response**: OK response

### POST /admin/num-groups

Sets num of groups and reassigns nums

- **Auth**: admin
- **Body**: JSON

```json
{
    "num_groups": "int"
}
```

- **Response**: OK response

### POST /admin/group-sizes

Sets the size of groups and reassigns nums

- **Auth**: admin
- **Body**: JSON

```json
{
    "group_sizes": ["int"]
}
```

- **Response**: OK response

### POST /admin/block-reqs

Sets whether to block login requests

- **Auth**: admin
- **Body**: JSON

```json
{
    "block_reqs": "bool"
}
```

- **Response**: OK response

### POST /admin/max-reqs

Sets the maximum number of logins/min

- **Auth**: admin
- **Body**: JSON

```json
{
    "max_req_per_min": "int"
}
```

- **Response**:

## Admin Export Routes

### GET /admin/export/judges

Exports judges as a CSV

- **Auth**: admin
- **Response**: CSV Blob

### GET /admin/export/projects

Exports projects as a CSV

- **Auth**: admin
- **Response**: CSV Blob

### GET /admin/export/challenges

Exports projects by challenge as ZIP of CSVs

- **Auth**: admin
- **Response**: ZIP Blob

### GET /admin/export/rankings

Exports a list of rankings for each judge

- **Auth**: admin
- **Response**: CSV Blob

## Admin Table Actions Routes

### PUT /judge/hide/\:id

Hides a judge

- **Auth**: admin
- **Parameter**: ID | judge ID to hide
- **Response**: OK response

### PUT /project/hide/\:id

Hides a project

- **Auth**: admin
- **Parameter**: ID | judge ID to hide
- **Response**: OK response

### PUT /judge/move/\:id

Moves a judge to a different group

- **Auth**: admin
- **Parameter**: ID | judge ID to hide
- **Response**: OK response

### PUT /project/move/\:id

Moves a project to a different group

- **Auth**: 
- **Response**:

### PUT /project/prioritize/\:id

Prioritizes a project

- **Auth**: 
- **Response**:

### POST /project/prioritize

Prioritizes multiple projects

- **Auth**: 
- **Response**:

### POST /judge/hide

Hides multiple judges

- **Auth**: 
- **Response**:

### POST /project/hide

Hides multiple projects

- **Auth**: 
- **Response**:

### POST /judge/move

Moves multiple judges to a different group

- **Auth**: 
- **Response**:

### POST /project/move

Moves multiple projects to a different group

- **Auth**: 
- **Response**:

### DELETE /admin/flag/\:id

Removes a flag

- **Auth**: 
- **Response**:

### POST /admin/deliberation

Toggles deliberation mode

- **Auth**: 
- **Response**:

### GET /admin/log

Gets the audit log

- **Auth**: 
- **Response**:

### GET /judge

Gets judge from token cookie

- **Auth**: 
- **Response**:

### GET /judge/welcome

Checks for `read_welcome` for a judge

- **Auth**: 
- **Response**:

### PUT /judge/welcome

Set `read_welcome` to true for a judge

- **Auth**: 
- **Response**:

### GET /judge/projects

Gets the list of projects a judge has seen

- **Auth**: 
- **Response**:

### POST /judge/next

Get next project for judge to view

- **Auth**: 
- **Response**:

### POST /judge/skip

Skips the current project with a reason

- **Auth**: 
- **Response**:

### POST /judge/finish

Finish viewing a project

- **Auth**: 
- **Response**:

### POST /judge/rank

Update judge rankings

- **Auth**: 
- **Response**:

### PUT /judge/star/\:id

Update star ranking for a project

- **Auth**: 
- **Response**:

### POST /judge/break

Removes active project for judge (take a break)

- **Auth**: 
- **Response**:

### PUT /judge/notes/\:id

Update notes for a project

- **Auth**: 
- **Response**:

### GET /project/\:id

Gets a project by ID

- **Auth**: 
- **Response**:

### GET /project/count

Gets the total number of projects

- **Auth**: 
- **Response**:

### GET /judge/project/\:id

Gets a judged project by a judge

- **Auth**: 
- **Response**:

### GET /judge/deliberation

Returns if deliberation mode is on

- **Auth**: 
- **Response**:

### GET /project/list/public

Gets a list of all projects for expo

- **Auth**: 
- **Response**:

### GET /challenges

Gets a list of all challenges

- **Auth**: 
- **Response**:

### GET /group-names

Gets a list of all group names

- **Auth**: 
- **Response**:

### GET /judge/list

Get a list of all judges

**Response**

## OUTDATED DOCS!!!

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

### GET /project/\:id

_Requires judge token_ | Gets a project by ID

-   id - ID of the project

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
