![go-build](https://github.com/ChrIgiSta/stopwatch/actions/workflows/build.yml/badge.svg) ![go-test](https://github.com/ChrIgiSta/stopwatch/actions/workflows/ci.yml/badge.svg)

# App 
This app enables stopwatch functionality for eliona.

![SmartView Picture not fould](eliona/example_stopwatch_view.png?raw=true "SmartView with stopwatch")

## Configuration

This app needs no configuration over api or database. Just use the new added asset type `Stopwatch` to start and stop an second counter as many you want.

### Registration in Eliona ###

To start and initialize an app in an Eliona environment, the app have to registered in Eliona. For this, an entry `stopwatch` in the database table `public.eliona_app` is necessary.

### Environment variables

- `APPNAME`: must be set to `stopwatch`. Some resources use this name to identify the app inside an Eliona environment.

- `CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started. (e.g. `postgres://user:pass@localhost:5432/iot`)

- `API_ENDPOINT`:  configures the endpoint to access the [Eliona API v2](https://github.com/eliona-smart-building-assistant/eliona-api). Otherwise, the app can't be initialized and started. (e.g. `http://api-v2:3000/v2`)

- `API_TOKEN`: defines the secret to authenticate the app and access the API. 

- `API_SERVER_PORT`(optional): define the port the API server listens. The default value is Port `3000`.

- `DEBUG_LEVEL`(optional): defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.

### Database tables ###

This app works without any database tables.

### App API ###

The app provides its own API to access configuration data and other functions. The full description of the API is defined in the `openapi.yaml` OpenAPI definition file.

- [API Reference](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/stopwatch/develop/openapi.yaml) shows Details of the API

**Generation**: to generate api server stub see Generation section below.


### Eliona ###

This app provides a start and a stop input (Button) to start or stop the stopwatch. In Current Value the time in seconds will provided to eliona until the timer has stopped. Then the Last Value will be actualized and the Current Value will be set to 0.

## Tools

### Testing ###

The code can be tested by running `go test`

### Generate API server stub ###

For the API server the [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) for go-server is used to generate a server stub. The easiest way to generate the server files is to use one of the predefined generation script which use the OpenAPI Generator Docker image.

```
.\generate-api-server.cmd # Windows
./generate-api-server.sh # Linux
```

### Generate Database access ###

For the database access [SQLBoiler](https://github.com/volatiletech/sqlboiler) is used. The easiest way to generate the database files is to use one of the predefined generation script which use the SQLBoiler implementation. Please note that the database connection in the `sqlboiler.toml` file have to be configured.

```
.\generate-db.cmd # Windows
./generate-db.sh # Linux
```

## Build on your own

 * Native
    * Install deps: `go install .`
    * Build: `go build .`
 * Docker
    * `docker build . --tag my-stopwatch`
