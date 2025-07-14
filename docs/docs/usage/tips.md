---
sidebar_position: 7
title: Tips and Contingencies
description: Tips for best ways to use Jury and possible Contingencies that can occur related to Jury
---

# Tips and Contingencies

## Tips

### Hackathon Size

Jury is designed to be used for **relatively larger hackathons**. If your hackathon has less than 20 projects, it might be a good idea to use a spreadsheet and manually assign judges to view projects. There just comes a point where the technical complexity of Jury outweights the benefits it provides. Any hackathons with more than 20 projects should be able to use Jury.

The largest event that Jury has been used at is at [HackUTD 2024](https://ripple.hackutd.co), where we had 281 projects being judged using Jury. This worked decently well, but we did run into a couple of issues regarding size. The key metric when considering how many judges you need is the **judge to project ratio**. We've found that having a ratio of **at least 1 judge for every 2 projects** is generally requied for using Jury. Any less and there simply isn't enough data for judging. The more data the better, and events where that ratio was closer to 1:1 had a lot more data points and more accurate results.

### Overseeing Judging

We generally assign one or two organizers to be in charge of the overall judging process. They will sit at the front and monitor the judging app. There are a couple of items that the manager should be looking at and delegating:

1. Keep an eye on the **average views per project**; we want to get that number to at least 3-5
2. Look for any projects that get flagged--if not flagged as absent, send an organizer to see if the flag is legit
3. Look for the **"hidden bc absent"** flags; send an organizer to make sure these are actually absent before dismissing the flag
4. Check the admin log occassionally to make sure nothing bad is happening (AHEM someone trying to hack into Jury)
5. Add judges and projects if needed
6. Monitor logins and disable if needed

## Contingencies

### Brute Force Attack

At our event, we had an instance where someone tried to brute force login to judges. They weren't able to do any damage luckily, but we needed to create a method to fight this obviously. Our solution is a way to rate limit logins, as well as a way to disable logins completely. Any judges that are already logged in can keep judging with their login token. This can be done through the [admin settings](/docs/usage/admin/configuration#judge-login).

You can detect these types of attacks by keeping an eye on the **audit log**. A brute force attack generally is easy to tell as it will show a lot of failed login attempts in a short or continuous amount of time.
