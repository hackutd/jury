# Jury

A project designed to create a new pairwise judging system using modern technologies aimed at optimizing the user experience of judges and admin users. See the inspiration for this project: [Gavel by anishathalye](https://github.com/anishathalye/gavel). Refer to [this excellent article](https://www.anishathalye.com/2015/03/07/designing-a-better-judging-system/) for more details on the underlying formulas of Gavel! The majority of our algorithm will be based off of research done in the field of [pairwise comparison](https://en.wikipedia.org/wiki/Pairwise_comparison).

# One Click Deploy!

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/acmutd/jury/tree/master&refcode=de08cdf58df3)

## 0. Set up MongoDB Atlas

MongoDB is a **free** NoSQL database that is hosted in the cloud! Their website has great instructions on how to quickly create an account and 100% lifetime free database, perfect for Jury: https://www.mongodb.com/docs/atlas/getting-started/.

Follow the instructions up to step 5: connect to your cluster. You will get a URI string that looks like `mongodb+srv://<username>:<password>@<something>.mongodb.net/?retryWrites=true&w=majority`, except \<username\>, \<password\>, and \<something\> should be filled in. You will need that for the environmental variables.

Note: You may need to whitelist all IPs in MongoDB for Jury to deploy successfully. To do this, add `0.0.0.0/0` as an IP address under "Network Access".

## 1. Resources

- Click Edit Plan
- Change plan to "Basic"
- Change Instance Size to "$5.00/mo -- Basic (512 MB RAM | 1 vCPU)"

Note that since Digital Ocean charges prorated, this will literally be like less than 20 cents a day.

## 2. Environmental Variables

Click on "edit" next to "jury-service" and fill in the env as follows:

- VITE_JURY_NAME - Name of the app to display to the user (eg. HackUTD X Judging)
- JURY_ADMIN_PASSWORD - Password for the admin portal (suggestion: use 2 random words)
- MONGODB_URI - URI of MongoDB Atlas instance (see step 0)

### Email hosting

If using Gmail SMTP:
- EMAIL_HOST = smtp.gmail.com
- EMAIL_FROM = gmail username
- EMAIL_USERNAME = gmail username
- EMAIL_PASSWORD = [google app password](https://support.google.com/accounts/answer/185833?hl=en#app-passwords)

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

Check out [our contributing docs](/CONTRIBUTING.md).