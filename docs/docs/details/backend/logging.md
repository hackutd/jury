---
sidebar_position: 4
title: Logging
description: How to log requests and changes to the global audit log
---

# Logging

The logger will log all actions that happen. To do this, use the `state.Logger` object in a route. The `Logger` object has a couple of functions that you can use. The most important are the following:

- `SystemLogf`: System logging
- `JudgeLogf`: Log something the judge does
- `AdminLogf`: Log something the admin does

The logger will write to the backend's memory, as well as the database. When the server restarts, Jury will load all logs from the database.

## Example log output

```
[2025-07-09T21:16:45Z] SYSTEM | Server restarted
[2025-07-09T21:16:45Z] SYSTEM | API created and running
[2025-07-14T14:51:59Z] ADMIN  | Log in
[2025-07-14T14:52:09Z] JUDGE  | Log in, assigned token dFsh5ZcRttD4CHmz
[2025-07-14T14:52:13Z] JUDGE  | Judge Judge Name (68715dc620fd42313b223755) | Read welcome message
[2025-07-14T14:52:45Z] ADMIN  | Unpaused clock
[2025-07-14T14:52:49Z] JUDGE  | Judge Judge Name (68715dc620fd42313b223755) | Picked new project TestProject-MO (6842e19239e3c24050d215ac)
```
