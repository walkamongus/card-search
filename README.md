# Hearthstone Card Search

## Table of Contents

1. [Description](#description)
1. [Quickstart](#quickstart)
    * [Prerequisites](#prerequisites)
    * [Building the container image](#building-the-container-image)
    * [Starting the application](#starting-the-application)
1. [Configuration](#configuration)
    * [Available configurations](#available-configurations)
    * [Example YAML configuration](#example-yaml-configuration)
1. [Project layout](#project-layout)

## Description

This is an example Go [Gin](https://github.com/gin-gonic/gin) web application that performs a Hearthstone card search based on specific card characteristics and returns the results formatted in a table.

## Quickstart

### Prerequisites

The web app requires a valid Blizzard API client ID and client secret. Instructions for creating a Blizzard API client can be found in the Blizzard API [Getting Started Guide](https://develop.battle.net/documentation/guides/getting-started).

### Building the container image

```
$ git clone https://github.com/walkamongus/card-search.git && cd card-search
$ docker build -t card-search .
```

### Starting the application

The following will run the app listening on `localhost` on the default port `8080`:

```
$ docker run -p 8080:8080 -e CLIENT_SECRET="client-s3cr3t-string" -e CLIENT_ID="client-id-string" card-search
```

You should now be able to access the web interface at: `http://localhost:8080`

## Configuration

Application configuration may be supplied through environment variables, command-line flags, or a `config.yaml` YAML configuration file.

### Available configurations

`CLIENT_ID` (ENV) or `--client-id` (CLI): Blizzard API client ID

`CLIENT_SECRET` (ENV) or `--client-secret` (CLI): Blizzard API client secret

### Example YAML Configuration

The `config.yaml` YAML configuration file should live in the same directory at application binary

```
---
client-id: my-client-id
client-secret: my-client-secret
```

## Project Layout

```
.
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── handlers.index.go   # Route handler for main index page
├── internal            # Container for internal packages
│   ├── hsapi           # Internal package containing Blizzard API client
│   └── util            # Internal package containing helper functions
├── main.go             # Main entry point. Establishes global structures, CLI interface, and starts Gin framework
├── routes.go           # Defines application routes and their handlers
└── templates           # Container for application view templates and template helpers
    ├── error.html
    ├── footer.html
    ├── header.html
    ├── index.html
    └── menu.html
```
