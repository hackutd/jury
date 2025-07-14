---
sidebar_position: 2
title: Documentation
description: Docs about the docs
---

# Documentation

Writing the documentation is extremely important yet probably the hardest part of development. The docs are hosted using [Docusarus](https://docusaurus.io/). All documentation is in the `docs/docs` directory, where each markdown file represents a single page in the documentation. The configuration can be found at `docs/docusaurus.config.ts`.

## Fun thing

If you want to see how many words is in the documentation, run this command in the `docs/docs` directory:

```bash
find . -type f -name "*.md" | xargs wc -w
```
