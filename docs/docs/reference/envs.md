---
sidebar_position: 3
---

# Environmental Variables

This page will go over the specifics of the environmental variables you need to define when deploying locally.

```
JURY_NAME=
JURY_ADMIN_PASSWORD=

MONGODB_URI=

EMAIL_HOST=
EMAIL_PORT=
EMAIL_FROM=
EMAIL_FROM_NAME=
EMAIL_PASSWORD=
SENDGRID_API_KEY=

PORT=
```

The `JURY_NAME` and `JURY_ADMIN_PASSWORD` are simply the name of the app and the admin password that you are using for local development. For development, the values here don't matter that much, but you should remember your admin password to log in (I personally use the classic `admin` password).

If you are using MongoDB Atlas, you should fill in the `MONGODB_URI` field using the [instructions on the MongoDB Atlas starter guide](https://www.mongodb.com/docs/atlas/getting-started/). This is the recommended way to run the application locally. If you are instead running MongoDB locally, you will need to define the `MONGODB_USER` and `MONGODB_PASS` fields -- see the [contributing page](/docs/contributing#alternative-offline-database-development) for more details.

For all email fields, refer to the information in the ["Deploying for your Hackathon"](/usage/deploy#email-hosting) page.

Finally, the definition of the `PORT` variable is optional -- specify this if you wish to connect to your app on a different port.
