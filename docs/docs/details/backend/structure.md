---
sidebar_position: 1
title: File Structure
description: How the backend code files are structured
---

# File Structure

The backend consists of multiple [modules](https://go.dev/blog/using-go-modules), each one in its own directory. Here is an overview of the modules:

### config

This small module simply contains functions for checking and getting environmental variables.

### database

This module contains function to interact with the database.

### funcs

Specialized functionality such as CSV parsing and email sending.

### judging

Complex functions used for the main judging flow. This includes aggregating judging scores, picking the next project, and maintaining the project comparisons matrix.

### logging

The logging module contains all functionality to write to the admin log.

### models

All the models (structs) that is used to represent the data in Jury. This includes judges, projects, and options.

### public

A dummy directory for the frontend. In development, this is used to serve a dummy page. In production, the statically built frontend will be copied to this directory before the Go build process.

### router

The main code for the API and route handler functions.

### util

Utility functions used across the backend app.
