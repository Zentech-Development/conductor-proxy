# Conductor Proxy
Conductor is a free, open-source solution to external API virtualization. It can integrate
with any of your existing HTTP (and more in the future) applications and the external
APIs they call to provide a common request format for your entire ecosystem. Conductor 
provides a fast, secure, and flexible solution for standardizing request formats, proxying
requests, and dynamically configuring resources available on an API without code changes.


# Overview
Imagine that you have an ecosystem of internal APIs that your company has built throughout the
years on different operating systems, different languages, different API request and response 
structures, different authentication mechanisms, and different API features. Now you need to build
an application composed of microservices that is capable of integrating with all of these systems 
with realtime requests.

One option would be to write out code which handles request and response logic, authentication
logic, and stick it in each of the new microservices. The problem is that this means any time the API's
server configuration, app paths, parameter names, or anything else changes, each microservice will
need to be rebuilt and redeployed. There is a better way.

Conductor solves this problem by becoming the central hub for API requests in your infrastructure.
The Conductor Proxy is a standalone application which acts like a smart proxy for requests sent
to it. System resources that are available to Conductor can be defined and stored. Various 
configurations must be specified about how to make requests to that resource, including defining
which endpoints are available.

Once Conductor has some defined resources, requests using a standardized format can be made to the
proxy, which will then construct the correct request to actually make to the resource, get the 
response, and send a standardized response back. This process can handle request parameters, 
authentication, and more.

Conductor is distributed as a standalone Docker image, binaries, and a Go package which can be
used to use the Conductor Proxy directly in your Go project.

In the future there will be official clients for calling the Conductor Proxy using various
programming languages, starting with Go, and a CLI.


# Features
- Supports HTTP/HTTPS incoming requests
- Supports HTTP/HTTPS external resources
- Supports HTTP endpoint parameters in the path, query params, body, and headers for external resources
- Supports static API key, simple username/password, and JWT authentication for external resources
- Supports 


# Supported Databases
- Redis (In Progress)
- Mongo (TODO)
- Postgres (TODO)
- SQLite (for testing only, TODO)


# Installation
## Standalone Docker Image
Pull the image from the public Docker registry. Check the registry page for available tags.
```sh
docker pull zentech/conductor-proxy:latest
```

## Docker Compose
Add the following to your `docker-compose.yml` file.
```yaml
services:
  conductor:
    image: zentech/conductor-proxy:latest
    container_name: conductor
    environment:
      - CONDUCTOR_DATABASE=sqlite # sqlite, postgres, mongo, redis
      - CONDUCTOR_SECURE=false
      - CONDUCTOR_HOST=localhost:7480
    ports:
      - "7480:7480"
```

## Binary - Ubuntu
TODO

## Go Package
TODO

# Configuration Options
## Environment Variables

### CONDUCTOR_DATABASE (Optional, defaults to sqlite)
Sets the database to use for persisting app and resource definitions and encrypted
tokens. Possible values are `sqlite`, `postgres`, `mongo`, `redis`.

If set to `sqlite`, Conductor will use a SQLite database in itself. A volume 
should be configured to save this data over container restarts. `CONDUCTOR_SECURE` should
probably be set to `false`, since this should only be used locally for development.

If set to `redis`, additional variables `CONDUCTOR_REDIS_HOST`, and `CONDUCTOR_REDIS_PASSWORD` must be supplied.

### CONDUCTOR_SECURE (Optional, defaults to true)
Sets Conductor to allow non-https requests, allow unauthenticated requests. This should NEVER
be `false` in production.

### CONDUCTOR_HOST (Optional, defaults to localhost:8080)
Hostname and port that Conductor will listen for traffic on.

### CONDUCTOR_DEFAULT_TOKEN_TIMEOUT (Optional, defaults to 1 hour)
Sets the default expiration time for account logins to use if the account does not have a 
custom value set. In seconds. A value of 0 means the generated tokens will never expire.

### CONDUCTOR_GIN_MODE (Optional, defaults to release)
Sets the GIN_MODE variable passed into Gin. Probably should be release unless you are developing
Conductor.

### CONDUCTOR_SECRET_KEY (Required)
Sets the secret key to be used for signing and verifying access tokens. Should be random and longer
than 36 characters.


# Concepts
## Accounts
TODO

## Groups
TODO

## Services
TODO

## Resources
### Parameters
Valid Data Types by Param Type:
Header: string, int, bool
Path: string, int, bool
Query: string, int, bool
Body: any
BodyFlat: any

TODO

## Proxy
TODO

# Conductor Proxy API Reference
There are two APIs to interact with- the proxy and the admin APIs. The admin API is used
to manage apps, resources, and accounts. The proxy is used to use the Conductor Proxy
virtual API.

## Admin API

### Login
<details>
 <summary><code>POST</code> <code><b>/api/login</b></code></summary>

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | username | body (required) | string | Username of user or service account |
> | passkey | body (required) | string | Passkey of user or service account |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200` | `application/json` | `{"statusCode": 200, "message": "Login successful", "data": {"token": "ey...."}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request", "data": {}}` |
> | `401` | `application/json` | `{"statusCode": 403, "message": "Failed to login", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Add Account
<details>
 <summary><code>POST</code> <code><b>/api/accounts</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | username | body (required) | string | Username of user or service account |
> | passkey | body (required) | string | Passkey of user or service account |
> | groups | body (required) | string array | List of groups to add the user to |
> | tokenExpiration | body (required) | int | Access token expiration time for account in seconds. 0 means never expires. |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201` | `application/json` | `{"statusCode": 201, "message": "Added account successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request, account name might already exist or group does not", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Remove Account
<details>
 <summary><code>DELETE</code> <code><b>/api/accounts/:id</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type | description |
> |------|-------|-----------|-------------|
> | id | path (required) | string | ID of account to remove |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200` | `application/json` | `{"statusCode": 200, "message": "Removed account successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request, make sure account ID is in request path", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Add Group
<details>
 <summary><code>POST</code> <code><b>/api/groups</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type | description |
> |------|-------|-----------|-------------|
> | name | body (required) | string | Name of group |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201` | `application/json` | `{"statusCode": 201, "message": "Added group successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request, group name might already exist", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Remove Group
<details>
 <summary><code>DELETE</code> <code><b>/api/groups/:id</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type | description |
> |------|-------|-----------|-------------|
> | id | path (required) | string | ID of group to remove |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200` | `application/json` | `{"statusCode": 200, "message": "Removed group successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request, make sure group ID is in request path", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Change Account Groups
<details>
 <summary><code>PUT</code> <code><b>/api/accounts/:id</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | id | path (required) | string | ID of account to update |
> | groupsToAdd | body (optional) | string array | List of groups to add the user to |
> | groupsToRemove | body (optional) | string array | List of groups to remove the user from |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200` | `application/json` | `{"statusCode": 201, "message": "Account groups updated successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Add App
<details>
 <summary><code>POST</code> <code><b>/api/apps</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | name | body (required) | string | Name of app |
> | friendlyName | body (required) | string | Readble name of app |
> | host | body (required) | string | Hostname of server |
> | adminGroups | body (required) | string array | List of groups that can admin this app |
> | userGroups | body (required) | string array | List of groups that can use this app |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201` | `application/json` | `{"statusCode": 201, "message": "Added app successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Add Resource
<details>
 <summary><code>POST</code> <code><b>/api/resources</b></code></summary>

 If running in secure mode, only admins will be able to do this.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | name | body (required) | string | Name of resource |
> | friendlyName | body (required) | string | Readble name of resource |
> | adminGroups | body (required) | string array | List of groups that can admin this resource |
> | userGroups | body (required) | string array | List of groups that can use this resource |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201` | `application/json` | `{"statusCode": 201, "message": "Added resource successfully", "data": {}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

## Proxy API
### Proxy Request
<details>
 <summary><code>POST</code> <code><b>/proxy</b></code></summary>

 If running in secure mode, request must include a valid token and the accounts groups must match the request.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | resourceId | body (required) | string | ID of requested resource |
> | endpoint | body (required) | string | Name of requested endpoint |
> | method | body (required) | string | HTTP method of requested endpoint |
> | params | body (required) | object | Object of param keys and values |
> | data | body (required) | any | Optional custom body to send to override body and bodyFlat params |
> | X_CONDUCTOR_KEY | header (required if secure mode on) | string | Valid access token |

##### Responses
> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200` | `application/json` | `{"statusCode": (proxied response status code), "message": "Success", "data": {(proxied response data)}}` |
> | `400` | `application/json` | `{"statusCode": 400, "message": "Bad request", "data": {}}` |
> | `401` | `application/json` | `{"statusCode": 401, "message": "Unauthenticated", "data": {}}` |
> | `403` | `application/json` | `{"statusCode": 403, "message": "Forbidden", "data": {}}` |
</details>

------------------------------------------------------------------------------------------


# Contributing
If you'd like to contribute to the project, awesome! Documentation for Contributing is
coming soon.