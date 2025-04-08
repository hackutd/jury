---
sidebar_position: 1
title: Deploying for your Hackathon
description: Step-by-step instructions for how to deploy Jury for your Hackathon
---

# Welcome!

Whether you're interested in using Jury at your own hackathon or wanting to help develop and improve the software, we're glad you're here 💙. Continue reading for information to get Jury up and running in no time for your event. Or, if you're interested in setting up a development environment, go to [development setup](/docs/contributing).

:::info[How does it work?]
For a walkthrough of the judging process using Jury, check out our [Judging Walkthrough](/docs/walkthrough) page. For more technical details, read about how Jury works [here](/docs/details)!
:::

## Deploying for your Hackathon

Jury is a stand-alone application used solely for the judging of hackathon projects. Before using Jury, make sure that your hackathon is using the [recommended physical setup](/docs/judging-setup).

### Cost?

We will be deploying Jury on **Digital Ocean's [App Platform](https://www.digitalocean.com/products/app-platform)**. Jury can be built into a [Docker container](https://www.docker.com/resources/what-container/), which is perfect for serving on the App Platform. For our tier of App Platform, we are playing [$5/month](https://www.digitalocean.com/pricing/app-platform). This might sound like a lot, but Digital Ocean's pricing is **prorated hourly**. This means that (given 31 days in a month, 744 hours in a month) if you have the judging application up for 2 days (48 hours), you would be charged a whole... **32 cents**!!! And to make this deal even better, if you sign up with the referral link below, you get $200 in free Digital Ocean credits for 60 days.

The only other external resource we need is a [MongoDB Atlas Database](https://www.mongodb.com/atlas/database). Luckily, MongoDB Atlas is **completely free** at its lowest tier, which is all we need. We have used it for events that exceed 1000 participants and 200+ projects being judged simultaneously, with no issues at all. The free database gives you **512 MB** of data storage, which is plenty for text-only storage.

:::tip[TL;DR]
You are paying **16 cents** per **24 hours** of hosting, plus a database that is 100% free forever for your data.
:::

## Step 0. Set up MongoDB Atlas

We first need to create a MongoDB database cluster. Instructions can be found at: https://www.mongodb.com/docs/atlas/getting-started/. You should follow the **Atlas UI** steps 1-4.

Once you have your database set up and a user added, go to (on the sidebar) **Deployment > Database**, then click on **Connect** next to your database. You should see a popup; click on **Drivers** and on "step 3" it shows a **Connection string**, which should be in the following format:

```
mongodb+srv://<username>:<password>@<url>.mongodb.net/?retryWrites=true&w=majority&appName=<app-name>
```

The `username`, `url`, and `app-name` fields should be filled in already; you just have to put the `password` that you set in step 3 of the MongoDB setup. This connection string will be important later, so make sure you don't lose the password!

:::warning
A common issue with MongoDB is IP addresses not being whitelisted. Go to the console, click **Security > Network Access**. Click **Add IP Address** on the right and add `0.0.0.0/0` in the "Access List Entry" field.
:::

## Step 1. One-click Deploy

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/hackutd/jury/tree/master&refcode=de08cdf58df3)

Click the button above to do a one-click deploy onto Digital Ocean. For full disclosure, the URL above contains my referral code, but it will get you **$200 in free credits** if this is your first time signing up for a Digital Ocean account. (This is enough to keep Jury running for 33 months!)

Once the page loads, you should see a configuration screen with the template as **Git source** pointing to **hackutd/jury**, followed by a app settings section. Most of the settings won't need to be touched, but we will go over the settings that need to be checked or changed in the following steps.

## Step 2. Instance Size

Under "App Settings" > "Resource Settings", you will need to check (and update if not set correctly) the instance sizing. You should see **512 MB RAM / 1 Shared vCPU / 50 GB bandwidth**, with only 1 container. This is the $5.00/month tier -- make sure this is correct in the cost summary on the right.

## Step 3. Environmental Variables

Click on "edit" next to "jury-service" and fill in the environmental variables as follows:

-   `VITE_JURY_NAME` - Name of the app to display to the user (eg. HackUTD X Judging)
-   `JURY_ADMIN_PASSWORD` - Password for the admin portal (suggestion: use 2 random words)
-   `MONGODB_URI` - URI of MongoDB Atlas instance (see step 0)

### Email Hosting

Jury needs to send emails to judges with their judging code. Each judge will receive an email when you add them, so if you have 100 judges, then you will need to send out 100 emails (but you should account for at least double to handle situations where you have to re-send emails or try different emails for judges). There are 3 ways we recommend hosting emails:

1. The first method is through [**Gmail SMTP**](https://support.google.com/a/answer/176600?hl=en#gmail-smtp-option), which means that you simply send emails through your personal `xxx@gmail.com` email. This is completely free but may face issues with sending limits (2000/day), spam filtering, and email sending performance.
2. The second method is through [**Sendgrid**](https://sendgrid.com/en-us). Many organizations already use Sendgrid, so this may be convenient for you. Note that you need to enable billing and [upgrade your account](https://sendgrid.com/en-us/marketing/sendgrid-services-cro) if you want to send more than 100 emails in a day. Sendgrid is nice because it's a fixed price up to a high email limit and has a separate API for sending emails outside of SMTP.
3. The final option that has been tested is [**AWS SES**](https://aws.amazon.com/ses/). This is a paid service also used by many organizations. The pricing for AWS SES is per 1000 emails, so it may be more expensive depending on how many emails you send. We use the SMTP protocol to call AWS SES's service too.

You will have to fill out specific environmental variables for each service:

If using **Gmail SMTP**:

-   `EMAIL_HOST` = smtp.gmail.com
-   `EMAIL_FROM` = gmail username
-   `EMAIL_USERNAME` = gmail username
-   `EMAIL_PASSWORD` = [google app password](https://support.google.com/accounts/answer/185833?hl=en#app-passwords)

If using **SendGrid**:

-   `SENDGRID_API_KEY` = [API key](https://docs.sendgrid.com/ui/account-and-settings/api-keys) provided by SendGrid
-   `EMAIL_FROM` = email to send from
-   `EMAIL_FROM_NAME` = display name of sender

If using [**AWS SES**](https://docs.aws.amazon.com/ses/latest/dg/smtp-credentials.html):

-   `EMAIL_HOST` = email-smtp.\<region\>.amazonaws.com
-   `EMAIL_FROM` = verified email to send from
-   `EMAIL_USERNAME` = SES SMTP username
-   `EMAIL_PASSWORD` = SES SMTP password

Make sure you **do not change the pre-filled fields**. You may leave any unused fields blank (ie. if you are using Sendgrid, you do not need to fill out the `EMAIL_PASSWORD` field).

## Step 4. App Name

Optionally, you can adjust the name of the app. If you scroll down to the "Finalize" section, there is a place to name your project. The URL of your deployed application will be `<app_name><some_characters>.ondigitalocean.app`. You may also want to set up a custom domain name for your application (see [below](#custom-domain-name)).

## Step 5. Done!

Review all of the settings to make sure it looks correct (especially the instance size). Once you make sure all settings are correct, click "Create Resources". After the app has finished building (which may take a couple of minutes), you should be able to see the app online at the provided URL!

Note that if you need to change any environmental variables, you can do so in the online Digital Ocean console (Note the app will take a few minutes to re-deploy before reflecting changes).

## Next Steps

### Troubleshooting

If you have any issues with setting up Jury, please feel free to reach out to me ([Michael Zhao](mailto:michaelzhao314@gmail.com)).

### Custom Domain Name

Follow the provided DigitalOcean guide to deploy to a [custom domain name](https://docs.digitalocean.com/products/app-platform/how-to/manage-domains/). Make sure you use **Option 2: Using a CNAME Pointer** if have your own domain hosting set up.
