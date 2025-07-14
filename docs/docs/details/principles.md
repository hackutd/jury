---
sidebar_position: 1
title: Core Principles
description: Jury core development principles
---

# Jury Core Development Principles

## 1. Use Less Third-Party Software

Obviously Jury is built on third-party software, but our goal is to minimize the amount of third-party software we use. For example, we do not use ANY component libraries--instead, creating a concise and lightweight set of components for Jury to use based off of pure React and Tailwind. This is true for the backend as well, where we try to stick to only the core Go libraries, except for necessary libraries such as Gin and mongodb.

## 2. Minimize Package Size

The current docker image for Jury in production is **17.2 megabytes**. This is achieved through minimizing frontend assets and some nice Golang build optimizations.

## 3. Consistency is Key!

The best way to keep the code as readable as possible is to be consistent with the rest of the code. When adding new features, refer to the existing functions and match their styles as well as possible. In the backend, `go fmt` helps with small code styling while the frontend relies on `Prettier` to do that. However, the way you style components (frontend) and write routes (backend) aren't automatically formatted. That means it falls to the programmer to verify that their new code is matching the existing code.

## 4. Document Document Document

As you can see from this entire site, documentation is one of the most important parts of maintaining code for a long time. Whenever you make a change to the core functionality or add a feature to a page, make sure to document it throughly or create an issue for someone else to document. This allows for the documentation to keep up with the code--resulting in a better user experience for both the users of Jury as well as the developers.
