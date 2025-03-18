# Backoffice

Backoffice - it is a platform for reporting on games and other entities, as well as managing various aspects such as games, organizations (integrators and providers), currencies, accounts and their roles.

## Project branches

The project has three main branches:

- `dev`
- `stage`
- `main`

The style of maintaining branches and commits is described [here](https://confluence.ejaw.net/display/INTERNAL/GIT).

## Documentation and access

- **https://backoffice.dev.heronbyte.com/api/swagger/index.html#/**: here you can familiarize yourself with the project's capabilities and perform various actions with entities.
- **https://backoffice.dev.heronbyte.com**: access to the project control panel.
- **https://confluence.ejaw.net/display/INTERNAL/%28TO+REWRITE%29+Backoffice**: technical documentation for the project.

## Makefile

At the root of the project there is a `Makefile` file that contains the following useful functions:

- `goose-install`: Setting `goose` for database migrations.
- `migration`: Creating a migration using `goose`.
- `migrate-up` and `migrate-down`: Apply or cancel database migrations.
- `proto` and `proto-history`: Generating code from protobuf files.
- `swag`: Initializing Swagger documentation.
- `lint`: Launching the linter for the project.

## Interaction with other services

The project interacts with two external services:

- **Bet Overlord Service** ([link to Git](https://bitbucket.org/electronicjaw/bet-overlord-service/src/master/)): provides a list of available games for a specific integrator.
- **History Service** ([link to Git](https://bitbucket.org/electronicjaw/history-service/src/main/)): serves as a repository for data about games and other entities, and also aggregates reports.

## Install and run locally
- [link to detailed instructions](https://docs.google.com/document/d/1LwFbEsnEwswbNyDdxznOLnpD2ysgp0pd7jQzv0O_7T0/edit#heading=h.q3zkiesajh1)

To start the project, you need to configure and raise the following services at default settings:

- **DB Postgre**
- **DB Redis**
- **Queue Broker RabbitMQ**


### Steps to start:

1. Set up the local configuration in the `config.example.yml` file according to your development environment. (example docker-compose file in confluence technical documentation)
2. Bring up the required services using `docker-compose` located in the root of the project.