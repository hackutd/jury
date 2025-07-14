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

For the deployment action, it requires the SSH hostname/user/key to deploy the app. I ([Michael Zhao](https://github.com/MichaelZhao21)) am currently using my personal server to host the staging app, and the main Jury repo is configured to my own server.

### Packages

Build container artifacts are stored in the [Github Container Registry](https://github.com/hackutd/jury/pkgs/container/jury). Currently, I ([Michael Zhao](https://github.com/MichaelZhao21)) have this configured with my own PAT (personal access token). If you would like to configure your own, you can follow the [official instructions by Github](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-with-a-personal-access-token-classic).

## Vercel

Vercel deploys the documentation site automatically when a PR is merged into main at https://jury.mikz.dev.
