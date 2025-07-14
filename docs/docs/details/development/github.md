---
sidebar_position: 1
title: Github
description: How the Github repo works and is managed
---

# Github

The Github repo, hosted at https://github.com/hackutd/jury contains all the source code for Jury, completely open source under the MIT License.

## Issues

All active issues are located on the issues tab. When creating an issue, make sure to use the correct template (found at `.github/ISSUE_TEMPLATE` directory).

## Pull Requests

All pull requests must be approved by at least one moderator. Make sure to follow the pull request template and link any issues that your PR fixes. The PR template can be found at `.github/PULL_REQUEST_TEMPLATE.md`.

## Actions

Github will automatically run an [action](https://github.com/features/actions) to deploy the staging app when a PR is merged into the main branch. The staging app is hosted at https://jury-dev.mikz.dev. This action can be seen in the `.github/workflows/deploy.yml` file.

## Vercel

Vercel deploys the documentation site automatically when a PR is merged into main at https://jury.mikz.dev.
