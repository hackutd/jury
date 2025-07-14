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

| Path                                                   | Method | Auth  | Description                                  |
| ------------------------------------------------------ | ------ | ----- | -------------------------------------------- |
| [/](#get-)                                             | GET    |       | Heartbeat route                              |
| [/judge/login](#post-judgelogin)                       | POST   |       | Login judge                                  |
| [/admin/login](#post-adminlogin)                       | POST   |       | Log into the admin dashboard                 |
| [/judge/auth](#post-judgeauth)                         | POST   | judge | Checks to see if judge is logged in          |
| [/admin/auth](#post-adminauth)                         | POST   | admin | Checks to see if admin is logged in          |
| [/judge/new](#post-judgenew)                           | POST   | admin | Add a new judge                              |
| [/judge/csv](#post-judgecsv)                           | POST   | admin | Add judges by CSV                            |
| [/judge/list](#get-judgelist)                          | GET    | admin | Get list of all judges                       |
| [/judge/:id](#delete-judgeid)                          | DELETE | admin | Deletes a judge by ID                        |
| [/judge/:id](#put-judgeid)                             | PUT    | admin | Edit judge info                              |
| [/admin/groups/swap](#post-admingroupsswap)            | POST   | admin | Swaps the judge groups manually              |
| [/admin/qr](#post-adminqr)                             | POST   | admin | Generate add judge QR code                   |
| [/admin/qr/:track](#post-adminqrtrack)                 | POST   | admin | Generate add judge to track QR code          |
| [/admin/qr](#get-adminqr)                              | GET    | admin | Gets add judge QR code                       |
| [/admin/qr/:track](#get-adminqrtrack)                  | GET    | admin | Gets add judge to track QR code              |
| [/qr/check/](#post-qrcheck)                            | POST   |       | Check if QR code is right                    |
| [/qr/check/:track](#post-qrchecktrack)                 | POST   |       | Check if track QR code is right              |
| [/qr/add](#post-qradd)                                 | POST   |       | Add judge from QR code                       |
| [/project/new](#post-projectnew)                       | POST   | admin | Add a new project                            |
| [/project/devpost](#post-projectdevpost)               | POST   | admin | Upload a Devpost CSV                         |
| [/project/csv](#post-projectcsv)                       | POST   | admin | Add projects by CSV                          |
| [/project/list](#get-projectlist)                      | GET    | admin | get list of all projects                     |
| [/project/:id](#delete-projectid)                      | DELETE | admin | Delete project by ID                         |
| [/project/:id](#put-projectid)                         | PUT    | admin | Edit project info                            |
| [/admin/stats](#get-adminstats)                        | GET    | admin | Get all stats                                |
| [/admin/stats/:track](#get-adminstatstrack)            | GET    | admin | Get all stats for a track                    |
| [/project/stats](#get-projectstats)                    | GET    | admin | Get the stats for projects                   |
| [/judge/stats](#get-judgestats)                        | GET    | admin | Get the stats for judges                     |
| [/admin/flags](#get-adminflags)                        | GET    | admin | Gets all flags                               |
| [/admin/clock](#get-adminclock)                        | GET    | admin | Gets the current clock state                 |
| [/admin/clock/pause](#post-adminclockpause)            | POST   | admin | Pauses the clock                             |
| [/admin/clock/unpause](#post-adminclockunpause)        | POST   | admin | Resumes the clock                            |
| [/admin/clock/backup](#post-adminclockbackup)          | POST   | admin | Backs up the clock to the database           |
| [/admin/started](#get-adminstarted)                    | GET    |       | Check if the clock is running                |
| [/admin/clock/reset](#post-adminclockreset)            | POST   | admin | Resets the clock                             |
| [/admin/reset](#post-adminreset)                       | POST   | admin | Resets the entire database                   |
| [/project/reassign](#post-projectreassign)             | POST   | admin | Reassign all project table numbers           |
| [/project/balance-groups](#post-projectbalance-groups) | POST   | admin | Balances project group numbers               |
| [/project/reassign](#post-projectreassign)             | POST   | admin | Reassign all project table numbers           |
| [/admin/timer](#get-admintimer)                        | GET    | judge | Gets the judge timer length                  |
| [/admin/options](#get-adminoptions)                    | GET    | admin | Gets all config options set                  |
| [/admin/options](#post-adminoptions)                   | POST   | admin | Sets config options                          |
| [/admin/tracks](#post-admintracks)                     | POST   | admin | Update the list of tracks                    |
| [/admin/track-views](#post-admintrack-views)           | POST   | admin | Update the min views per track               |
| [/admin/num-groups](#post-adminnum-groups)             | POST   | admin | Sets num of groups and reassigns nums        |
| [/admin/group-sizes](#post-admingroup-sizes)           | POST   | admin | Sets the size of groups and reassigns nums   |
| [/admin/block-reqs](#post-adminblock-reqs)             | POST   | admin | Sets whether to block login requests         |
| [/admin/max-reqs](#post-adminmax-reqs)                 | POST   | admin | Sets the maximum number of logins/min        |
| [/admin/export/judges](#get-adminexportjudges)         | GET    | admin | Exports judges as a CSV                      |
| [/admin/export/projects](#get-adminexportprojects)     | GET    | admin | Exports projects as a CSV                    |
| [/admin/export/challenges](#get-adminexportchallenges) | GET    | admin | Exports projects by challenge as ZIP of CSVs |
| [/admin/export/rankings](#get-adminexportrankings)     | GET    | admin | Exports a list of rankings for each judge    |
| [/judge/hide/:id](#put-judgehideid)                    | PUT    | admin | Hides a judge                                |
| [/project/hide/:id](#put-projecthideid)                | PUT    | admin | Hides a project                              |
| [/judge/move/group/:id](#put-judgemovegroupid)         | PUT    | admin | Moves a judge to a different group           |
| [/project/move/group/:id](#put-projectmovegroupid)     | PUT    | admin | Moves a project to a different group         |
| [/project/move/:id](#put-projectmoveid)                | PUT    | admin | Moves a project to a different table number  |
| [/project/prioritize/:id](#put-projectprioritizeid)    | PUT    | admin | Prioritizes a project                        |
| [/project/prioritize](#post-projectprioritize)         | POST   | admin | Prioritizes multiple projects                |
| [/judge/hide](#post-judgehide)                         | POST   | admin | Hides multiple judges                        |
| [/project/hide](#post-projecthide)                     | POST   | admin | Hides multiple projects                      |
| [/judge/move/group](#post-judgemovegroup)              | POST   | admin | Moves multiple judges to a different group   |
| [/project/move/group](#post-projectmovegroup)          | POST   | admin | Moves multiple projects to a different group |
| [/admin/flag/:id](#delete-adminflagid)                 | DELETE | admin | Removes a flag                               |
| [/admin/deliberation](#post-admindeliberation)         | POST   | admin | Toggles deliberation mode                    |
| [/admin/log](#get-adminlog)                            | GET    | admin | Gets the audit log                           |
| [/judge](#get-judge)                                   | GET    | judge | Gets judge from token cookie                 |
| [/judge/welcome](#get-judgewelcome)                    | GET    | judge | Checks for `read_welcome` for a judge        |
| [/judge/welcome](#put-judgewelcome)                    | PUT    | judge | Set `read_welcome` to true for a judge       |
| [/judge/projects](#get-judgeprojects)                  | GET    | judge | Gets the list of projects a judge has seen   |
| [/judge/next](#post-judgenext)                         | POST   | judge | Get next project for judge to view           |
| [/judge/skip](#post-judgeskip)                         | POST   | judge | Skips the current project with a reason      |
| [/judge/finish](#post-judgefinish)                     | POST   | judge | Finish viewing a project                     |
| [/judge/rank](#post-judgerank)                         | POST   | judge | Update judge rankings                        |
| [/judge/star/:id](#put-judgestarid)                    | PUT    | judge | Update star ranking for a project            |
| [/judge/notes/:id](#put-judgenotesid)                  | PUT    | judge | Update notes for a project                   |
| [/project/:id](#get-projectid)                         | GET    | judge | Gets a project by ID                         |
| [/project/count](#get-projectcount)                    | GET    | judge | Gets the total number of projects            |
| [/judge/project/:id](#get-judgeprojectid)              | GET    | judge | Gets a judged project by a judge             |
| [/judge/deliberation](#get-judgedeliberation)          | GET    | judge | Returns if deliberation mode is on           |
| [/project/list/public](#get-projectlistpublic)         | GET    |       | Gets a list of all projects for expo         |
| [/challenges](#get-challenges)                         | GET    |       | Gets a list of all challenges                |
| [/group-info](#get-group-info)                         | GET    |       | Gets a list of all group names               |

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

-   **Auth**: none
-   **Response**: OK response

## Login Routes

### POST /judge/login

Login judge

-   **Auth**: none
-   **Body**: JSON

```json
{
    "code": "String | login code"
}
```

-   **Response**: JSON

```json
{
    "token": "String | judge login token"
}
```

### POST /admin/login

Log into the admin dashboard

-   **Auth**: none
-   **Body**: JSON

```json
{
    "password": "String | admin password"
}
```

-   **Response**: OK response

### POST /judge/auth

Checks to see if judge is logged in

-   **Auth**: judge
-   **Response**: OK response

### POST /admin/auth

Checks to see if admin is logged in

-   **Auth**: admin
-   **Response**: OK response

## Admin Panel (Judges) Routes

### POST /judge/new

Add a new judge

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "name": "String",
    "email": "String",
    "notes": "String"
}
```

-   **Response**: OK response

### POST /judge/csv

Add judges by CSV

-   **Auth**: admin
-   **Body**: FormData
    -   `csv`: CSV file
    -   `hasHeader`: Boolean, true if CSV has a header
    -   `noSend`: Don't send email to judge if true
-   **Response**: OK response

### GET /judge/list

Get list of all judges

-   **Auth**: admin
-   **Response**: JSON List

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
        "rankings": ["ObjectId", "ObjectId | and so on for each ranked project"],
        "last_activity": "DateTime"
    }
]
```

### DELETE /judge/\:id

Deletes a judge by ID

-   **Auth**: admin
-   **Parameter**: ID, the ID of the judge to delete
-   **Response**: OK response

### PUT /judge/\:id

Edit judge info

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "name": "String",
    "email": "String",
    "notes": "String"
}
```

-   **Response**: OK response

### POST /admin/groups/swap

Swaps the judge groups manually

-   **Auth**: admin
-   **Response**: OK response

### POST /admin/qr

Generate add judge QR code

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### POST /admin/qr/\:track

Generate add judge to track QR code

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### GET /admin/qr

Gets add judge QR code

-   **Auth**: none
-   **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### GET /admin/qr/\:track

Gets add judge to track QR code

-   **Auth**: none
-   **Response**: JSON

```json
{
    "qr_code": "String"
}
```

### POST /qr/check

Checks if QR code is correct

-   **Auth**: none
-   **Body**: JSON

```json
{
    "code": "String"
}
```

-   **Response**: OK response

### POST /qr/check/\:track

Checks if track QR code is correct

-   **Auth**: none
-   **Body**: JSON

```json
{
    "code": "String"
}
```

-   **Response**: OK response

### POST /qr/add

Add judge from QR code

-   **Auth**: none
-   **Body**: JSON

```json
{
    "name": "String",
    "email": "String",
    "notes": "String",
    "code": "String | QR code token"
}
```

-   **Response**: OK response

## Admin Panel (Projects) Routes

### POST /project/new

Add a new project

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "name": "String",
    "description": "String",
    "url": "String",
    "try_link": "String",
    "video_link": "String",
    "challenge_list": "String | comma-separated list of challenges/tracks"
}
```

-   **Response**: OK response

### POST /project/devpost

Upload a Devpost CSV

-   **Auth**: admin
-   **Body**: FormData
    -   `csv`: CSV file
-   **Response**:

### POST /project/csv

Add projects by CSV

-   **Auth**: admin
-   **Body**: FormData
    -   `csv`: CSV file
    -   `hasHeader`: Boolean, true if CSV has a header
-   **Response**: OK response

### GET /project/list

get list of all projects

-   **Auth**: admin
-   **Response**: JSON

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
            "track2": "int"
        },
        "score": "int",
        "stars": "int",
        "track_stars": {
            "track1": "int",
            "track2": "int"
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

-   **Auth**: admin
-   **Parameter**: ID of project to delete
-   **Response**: OK response

### PUT /project/\:id

Edit project info

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "name": "String",
    "description": "String",
    "url": "String",
    "try_link": "String",
    "video_link": "String",
    "challenge_list": "String | comma-separated list of challenges/tracks"
}
```

-   **Response**: OK response

## Admin Panel (Stats/Data) Routes

### GET /admin/stats

Get all stats

-   **Auth**: admin
-   **Response**: JSON

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

-   **Auth**: admin
-   **Parameter**: Track, the track to fetch stats for
-   **Response**: JSON

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

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "num": "int",
    "num_active": "int",
    "avg_seen": "float"
}
```

### GET /judge/stats

Get the stats for judges

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "num": "int",
    "num_active": "int",
    "avg_seen": "float"
}
```

### GET /admin/flags

Gets all flags

-   **Auth**: admin
-   **Response**: JSON List

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

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "running": "bool",
    "time": "int"
}
```

### POST /admin/clock/pause

Pauses the clock

-   **Auth**: admin
-   **Response**: OK response

### POST /admin/clock/unpause

Resumes the clock

-   **Auth**: admin
-   **Response**: OK response

### POST /admin/clock/backup

Backs up the clock to the database

-   **Auth**: admin
-   **Response**: OK response

### GET /admin/started

Check if the clock is running

-   **Auth**: admin
-   **Response**: OK response (1 if true)

## Admin Options/Settings Routes

### POST /admin/clock/reset

Resets the clock

-   **Auth**: admin
-   **Response**: OK response

### POST /admin/reset

Resets the entire database

-   **Auth**: admin
-   **Response**: OK response

### POST /project/reassign

Reassign all project table numbers

-   **Auth**: admin
-   **Response**: OK response

### POST /project/balance-groups

Balances project group numbers

-   **Auth**: admin
-   **Response**: OK response

### POST /project/reassign

Reassign all project table numbers

-   **Auth**: admin
-   **Response**: OK response

### GET /admin/timer

Gets the judge timer length

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "judging_timer": "int"
}
```

### GET /admin/options

Gets all config options set

-   **Auth**: admin
-   **Response**: JSON

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

-   **Auth**: admin
-   **Body**: JSON

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
    "block_reqs": "bool"
}
```

-   **Response**: OK response

### POST /admin/tracks

Update the list of tracks

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "tracks": ["String"]
}
```

-   **Response**: OK response

### POST /admin/track-views

Update the min views per track

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "track_views": ["int"]
}
```

-   **Response**: OK response

### POST /admin/num-groups

Sets num of groups and reassigns nums

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "num_groups": "int"
}
```

-   **Response**: OK response

### POST /admin/group-sizes

Sets the size of groups and reassigns nums

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "group_sizes": ["int"]
}
```

-   **Response**: OK response

### POST /admin/block-reqs

Sets whether to block login requests

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "block_reqs": "bool"
}
```

-   **Response**: OK response

### POST /admin/max-reqs

Sets the maximum number of logins/min

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "max_req_per_min": "int"
}
```

-   **Response**:

## Admin Export Routes

### GET /admin/export/judges

Exports judges as a CSV

-   **Auth**: admin
-   **Response**: CSV Blob

### GET /admin/export/projects

Exports projects as a CSV

-   **Auth**: admin
-   **Response**: CSV Blob

### GET /admin/export/challenges

Exports projects by challenge as ZIP of CSVs

-   **Auth**: admin
-   **Response**: ZIP Blob

### GET /admin/export/rankings

Exports a list of rankings for each judge

-   **Auth**: admin
-   **Response**: CSV Blob

## Admin Table Actions Routes

### PUT /judge/hide/\:id

Hides a judge

-   **Auth**: admin
-   **Parameter**: ID | judge ID to hide
-   **Response**: OK response

### PUT /project/hide/\:id

Hides a project

-   **Auth**: admin
-   **Parameter**: ID | judge ID to hide
-   **Response**: OK response

### PUT /judge/move/group/\:id

Moves a judge to a different group

-   **Auth**: admin
-   **Parameter**: ID | judge ID to hide
-   **Body**: JSON

```json
{
    "group": "int | group number"
}
```

-   **Response**: OK response

### PUT /project/move/group/\:id

Moves a project to a different group

-   **Auth**: admin
-   **Parameter**: ID | ID of project to move
-   **Body**: JSON

```json
{
    "group": "int | group number"
}
```

-   **Response**: OK response

### PUT /project/move/\:id

Moves a project to a different table number

-   **Auth**: admin
-   **Parameter**: ID | ID of project to move
-   **Body**: JSON

```json
{
    "location": "int"
}
```

-   **Response**: OK response

### PUT /project/prioritize/\:id

Prioritizes a project

-   **Auth**: admin
-   **Parameter**: ID | ID of project to prioritize
-   **Body**: JSON

```json
{
    "prioritize": "bool | true if prioritizing"
}
```

-   **Response**: OK response

### POST /project/prioritize

Prioritizes multiple projects

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "items": ["ObjectID"],
    "prioritize": "bool | true if prioritizing"
}
```

-   **Response**: OK response

### POST /judge/hide

Hides multiple judges

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "items": ["ObjectID"],
    "hide": "bool | true if hiding"
}
```

-   **Response**: OK response

### POST /project/hide

Hides multiple projects

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "items": ["ObjectID"],
    "hide": "bool | true if hiding"
}
```

-   **Response**: OK response

### POST /judge/move/group

Moves multiple judges to a different group

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "items": ["ObjectID"],
    "group": "int | group number"
}
```

-   **Response**: OK response

### POST /project/move/group

Moves multiple projects to a different group

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "items": ["ObjectID"],
    "group": "int | group number"
}
```

-   **Response**: OK response

### DELETE /admin/flag/\:id

Removes a flag

-   **Auth**: admin
-   **Parameter**: ID, the ID of the flag to delete
-   **Response**: OK response

### POST /admin/deliberation

Toggles deliberation mode

-   **Auth**: admin
-   **Body**: JSON

```json
{
    "start": "bool | true if starting deliberation"
}
```

-   **Response**: OK response

## Admin Panel (Log) Routes

### GET /admin/log

Gets the audit log

-   **Auth**: admin
-   **Response**: JSON

```json
{
    "log": "String"
}
```

## Judging Routes

### GET /judge

Gets judge from token cookie

-   **Auth**: judge
-   **Response**: JSON

```json
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
    "rankings": ["ObjectId", "ObjectId | and so on for each ranked project"],
    "last_activity": "DateTime"
}
```

### GET /judge/welcome

Checks for `read_welcome` for a judge

-   **Auth**: judge
-   **Response**: OK response

### PUT /judge/welcome

Set `read_welcome` to true for a judge

-   **Auth**: judge
-   **Response**: OK response

### GET /judge/projects

Gets the list of projects a judge has seen

-   **Auth**: judge
-   **Response**: JSON List

```json
[
    {
        "project_id": "ObjectID",
        "name": "String",
        "location": "int",
        "description": "String",
        "notes": "String",
        "starred": "bool",
        "url": "String"
    }
]
```

### POST /judge/next

Get next project for judge to view

-   **Auth**: judge
-   **Response**: JSON

```json
{
    "project_id": "ObjectID"
}
```

### POST /judge/skip

Skips the current project with a reason

-   **Auth**: judge
-   **Body**: JSON

```json
{
    "reason": "String"
}
```

-   **Response**: OK response

### POST /judge/finish

Finish viewing a project

-   **Auth**: judge
-   **Body**: JSON

```json
{
    "notes": "String",
    "starred": "bool"
}
```

-   **Response**: OK response

### POST /judge/rank

Update judge rankings

-   **Auth**: judge
-   **Body**: JSON

```json
{
    "ranking": ["ObjectID"]
}
```

-   **Response**: OK response

### PUT /judge/star/\:id

Update star ranking for a project

-   **Auth**: judge
-   **Body**: JSON

```json
{
    "starred": "bool"
}
```

-   **Response**: OK response

### PUT /judge/notes/\:id

Update notes for a project

-   **Auth**: judge
-   **Parameter**: ID | ID for project to update notes for
-   **Body**: JSON

```json
{
    "notes": "String"
}
```

-   **Response**: OK response

### GET /project/\:id

Gets a project by ID

-   **Auth**: judge
-   **Parameter**: ID | ID for project to get
-   **Response**: JSON

```json
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
```

### GET /project/count

Gets the total number of projects

-   **Auth**: judge
-   **Response**: JSON

```json
{
    "count": "int"
}
```

### GET /judge/project/\:id

Gets a judged project by a judge

-   **Auth**: judge
-   **Parameter**: ID | Project ID to get
-   **Response**: JSON

```json
{
    "project_id": "ObjectID",
    "name": "String",
    "location": "int",
    "description": "String",
    "notes": "String",
    "starred": "bool",
    "url": "String"
}
```

### GET /judge/deliberation

Returns if deliberation mode is on

-   **Auth**: judge
-   **Response**: OK response | 1 if deliberation mode is on

## Project Expo Routes

### GET /project/list/public

Gets a list of all projects for expo

-   **Auth**: none
-   **Response**: JSON List

```json
[
    {
        "name": "String",
        "location": "u64",
        "description": "String",
        "url": "String",
        "try_link": "String",
        "video_link": "String",
        "challenge_list": ["String"],
        "group": "int"
    }
]
```

### GET /challenges

Gets a list of all challenges

-   **Auth**: none
-   **Response**: JSON List

```json
["String | challenge 1", "String | challenge 2"]
```

### GET /group-info

Gets a list of all group names and if groups are enabled

-   **Auth**: none
-   **Response**: JSON

```json
{
    "names": ["String | group 1", "String | group 2"],
    "enabled": "bool"
}
```
