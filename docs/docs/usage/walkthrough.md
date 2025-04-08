---
sidebar_position: 3
title: Judging Walkthrough
description: A quick walkthrough of the judging process through an admin and judge POV.
---

# Judging Walkthrough

It's the end of your hackathon! You've made it through the event and just have one last thing left to do: judge the projects. This isn't as easy of a task as it may seem. The list can get quite long: organizing and orientating judges, setting up judging booths for participants, distributing projects to judges, tallying up scores, and figuring out the final result, all while dealing with the chaos of the ongoing problems. It's a LOT, even for a seasoned hackathon organizer. That's why we wanted to design an all-in-one solution to handle the hassles of judging hackathons: **✨ Jury ✨**.

:::warning[Note]
If you haven't read the [Physical Judging Setup](/docs/judging-setup) doc, you should read that first!
:::

This page will walk you through a typical end-of-hackathon experience using Jury, from checking in judges to a final deliberation!

## Pre-event Setup

To ensure Jury works correctly, you should set it up a day or two in advance. Create a MongoDB instance and use the one-click deploy to have Jury up in less than 15 minutes (Details are on the [Deploying for your Hackathon](/docs/intro) page). Once deployed, share your admin password with all organizers that need to access the admin side of the judging system.

## Judges Enter

Fast forward to Sunday. Project are getting finished up, and judges are coming through the door! Once judges are checked in, you can put them into the Jury system -- in one of two ways:

1. If you already have a list of judges that you would like to add, simply upload a CSV file of all judges to Jury.
2. Otherwise, you can use the "Add Judges" form on the admin dashboard to add individual judges -- you only need their name and email.

When judges are added, they will receive an email, with information about how to access their own judging portal and login. Each judge will get a unique code to login. Once they log in, they will be able to start judging whenever an admin starts the judging session.

## Projects Submitted

As the project deadline approaches, hackers will be hard at work submitting their projects. We've typically used [Devpost](https://devpost.com) as the platform for submitting projects, so Jury has an integration that allows you to take the CSV export from Devpost and import that directly into Jury. Once the deadline passes, you can export the data from devpost and import it into Jury, all within a minute. If you have any projects that are submitted late, you can manually add them to Jury through the "Add Projects" form. As projects are added to the system, they will automatically be assigned table numbers (sequentially).

At this point, your team would be setting up the judging room(s), with tables laid out for hackers to stand at (judging in the [science fair style](/docs/judging-setup)). You can then go to the live judging page ([example here](https://jury-dev.mikz.dev/expo)), which will show a list of all projects and their assigned table numbers. Once you have tables set up, this live judging page can be safely sent out to hackers without exposing any of the other parts of the judging application. Hackers should go to their tables and set up there before judging starts.

In the admin settings, you are also able to export projects grouped by their opt-in prize (specified on Devpost) as a CSV file. This may be beneficial for any sponsors that want lists of projects that submitted specifically to their challenge. There are other output options such as a list of all projects or judges on that page as well. Some other settings include limiting the project presentation times and enabling track judging.

## Judging Time!

To start judging, any admin can start the judging session using the "Start Judging" button at the top of the admin dashboard. Once started, judges can refresh the page to start judging!

### As a Judge

Once you click "Start Judging", you will be presented with a randomly assigned project. On the page, you will see the project name, table number, and a description of the project. This description is just for your reference; you do not need to read through it. Go to the table number that is displayed. Once the group starts presenting, click the "Start Judging" problem at the top of the screen to start judging the project. It will start a countdown timer (default 5 minutes), which will make a beeping sound once completed.

If the group is not there or is currently in the process of being judged, you can click the "Skip" button, which will allow you to skip the project. This will have no net effect on the ranking/status of that given project. If you see a project that may violate our rules or code of conduct, you can flag it for organizers to look at.

Once a group is finished presenting, click the "Finish" button to open the rating popup. This menu will ask you to rate the judges on a 1-10 scale in a multitude of categories. These numbers will be considered in the final deliberation but doesn't directly impact the project's ranking. Once finished, click "Done". You will be taken back to the inital menu, which should now show 1 project. Click the "Next Project" button to move to the second project. You will once again be presented with another random project. Go to that table, listen to their project, then rate their project. Once you get back to the initial page again, you now should have 2 projects. You can now drag and drop the projects into the ranking, where the top project is your highest-ranked project.

You will continue this process until you see all projects or the judging session ends. If at any time you would like to take a break, you can click the "I want to take a break!" button, which allows you to take a break for as long as you need. To resume judging, simply click the "Next Project" button.

### As an Admin (Organizer)

As judging happens, the admin page will continuously be updated with the information regarding scores of projects and the stats of each judge. The table presented on the dashboard can be sorted by any column, from the name of the project to their score. 

Projects and judge can also be dynamically be added, and they will be simply circulated into the system. The key metric you are looking for is the average number of times each project is seen. As a general metric (and especially if you have a high project to judges ratio), you should aim for every project to be seen at least 5 times. This ensures that the algorithm has enough data to fairly aggregate project rankings.

At any time, you may pause judging for everyone. This will prevent new projects from being assigned to judges. This feature is good if you ever want to have a break during judging or if you want to end judging.

## Final Deliberation

Once judging has ended, organizers will get together and decide on the winners. The admin dashboard will show you a list of scores for each project, which is the aggregated Borda count score for each project. For more details on this algorithm, refer to the [technical details](/docs/details) page.

:::warning[WIP]
🚧 The judging walkthrough page is almost done but still under construction (need to add pictures!) 🚧
:::
