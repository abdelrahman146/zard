# Service Entry Points

This directory contains the entry points for the service. The entry points are the main files that are executed when the service is started. The entry points are responsible for setting up the service and starting the server.


## Main Entry Point
the main entry point is in package `main` and is named `main.go`. This file is responsible for setting up the service and starting the server.


## Scheduled Jobs Entry Points (Cron Jobs)
The scheduled jobs entry points are in package `<job_name>` and are named `<job_name>.go`. These files are responsible for setting up the scheduled jobs and starting the cron jobs.