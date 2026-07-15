# IMServer Console

English | [中文](README.zh-CN.md)

IMServer Console is a lightweight admin management console for IM services built with Go. It provides a web-based control panel and API gateway for operating IM applications, users, groups, bots, push settings, file configuration, monitoring, and message statistics.

This project is designed to help teams manage and operate an IM platform more efficiently, especially in production environments where configuration and observability matter.

## Overview

IMServer Console is built for teams that want a practical, production-friendly way to operate an IM system without building a custom admin backend from scratch. It combines management APIs, a web console, and operational features into one project that is easy to deploy and extend.

## Why this project

- Simplifies daily operations for IM deployments
- Centralizes app, account, user, and group administration
- Reduces manual database work with a structured admin interface
- Supports configuration for push, RTC, SMS, translation, and file services
- Helps monitor service health and business usage through statistics and logs

## Core capabilities

- Account and app management
- User, bot, group, and conversation administration
- Message history and recall operations
- Sensitive word and interceptor management
- Push configuration for iOS and Android
- SMS, translation, file, and RTC configuration
- Business statistics, user activity, and performance metrics
- Embedded admin web UI and API gateway

## Tech stack

- Go
- Gin framework
- GORM
- MySQL
- Vite-based web frontend

## Quick start

### Prerequisites

- Go 1.25 or newer
- MySQL server
- Access to an IM backend service that this console connects to

### Configuration

Edit the configuration file at [conf/config.yml](conf/config.yml) and set:

- the server port
- MySQL connection details
- the IM API domain and admin domain

### Run locally

```bash
go run .
```

The service will start using the configured port, and the admin interface will be available through the web console.

## Project structure

- [main.go](main.go) - application entry point
- [apis](apis) - API handlers
- [services](services) - business logic
- [dbs](dbs) - database access layer
- [routers](routers) - route definitions
- [webconsole](webconsole) - embedded admin web UI
- [conf](conf) - configuration files

## Development notes

The project uses automatic database migration/upgrade flow on startup, so the runtime environment should be prepared before launching the service.

## Contributing

Contributions, bug reports, and feature requests are welcome. If you want to improve this project, feel free to open an issue or submit a pull request.

If this project helps you, please consider giving it a star to support continued development.
