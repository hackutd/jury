# Website

This documentation site is built using [Docusaurus](https://docusaurus.io/), a modern static website generator.

### Installation

```
$ yarn
```

### Local Development

```
$ yarn start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

## Documentation Structure

- Using Jury
    - Deploying for your Hackathon
    - Physical Judging Setup
    - Judging Overview
    - Jury Admin
        - Configuration
        - Adding Projects
        - Adding Judges
        - Tracks
        - Groups
        - Admin Dashboard
    - Jury Judging
        - Judge Onboarding
        - Judging Interface
- Contributing
- Technical Details
    - Frontend
        - Routing
        - Components
        - Styling
        - Misc Conventions
    - Backend
        - API
        - Database
    - Development
        - CI/CD
        - Github
    - Rank Aggregation
    - Integration (w public API)
- Reference
    - Frontend
    - Backend (Private API)
    - Public API
    - Environmental Variables

## Image Compression

Use the `compress` script to compress PNGs if they are larger than 1000px please! Note that you will need to install [imagemagick](https://imagemagick.org/script/download.php) to use this script. Usage:

```
compress <path_to_image>
```
