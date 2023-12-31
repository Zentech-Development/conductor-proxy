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

Conductor is distributed as a standalone Docker image and as native binaries.

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
Add the following to your `docker-compose.yml` file. (Not currently working)
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

# Configuration Options
Configuration is handled using a configuration file with env file-style variables. A template configuration
file is available in `pkg/config/conductor.conf.template`.

## Environment Variables

### DATABASE_TYPE (Optional, defaults to mock)
Sets the database to use for persisting service and resource definitions and encrypted
tokens. Possible values are `sqlite`, `postgres`, `mongo`, `redis`, and `mock`.

If set to `sqlite`, Conductor will use a SQLite database in itself. A volume 
should be configured to save this data over container restarts. `CONDUCTOR_SECURE` should
probably be set to `false`, since this should only be used locally for development.

If set to `redis`, additional variables `REDIS_HOST`, and `REDIS_PASSWORD` must be supplied.

If set to `mock`, data will not be preserved when the application is stopped.

### SECURE_MODE (Optional, defaults to true)
Sets Conductor to allow non-https requests, allow unauthenticated requests, etc. This should NEVER
be `false` in production. It also does not do anything currently.

### HOST (Optional, defaults to localhost:8000)
Hostname and port that Conductor will listen for traffic on.

### DEFAULT_TOKEN_TIMEOUT_SECONDS (Optional, defaults to 1 hour)
Sets the default expiration time for account logins to use if the account does not have a 
custom value set. In seconds. A value of 0 means the generated tokens will never expire.

### ACCESS_TOKEN_SECRET_KEY (Optional, but must be provided here or in CONDUCTOR_SECRET environment variable)
Sets the secret key to be used for signing and verifying access tokens. Should be random and longer
than 36 characters. This should only be set in the config file in development. In production, set the 
secret key by setting a `CONDUCTOR_SECRET` environment variable. If one is set, the value in the config value
will be overwritten.

### DEFAULT_ADMIN_USERNAME (Optional, defaults to admin)
Username for initial admin account when Conductor first starts.

### DEFAULT_ADMIN_PASSKEY (Optional, defaults to password)
Passkey for initial admin account when Conductor first starts.

# Running Conductor Proxy
## Accounts
Accounts are representations of human or service users that has a username, password, groups, and
a token expiration. When an account logs in, an access token is generated for that user's groups
and expiration time, giving them access to other endpoints in the API. The access token must be
passed in a request header called X-CONDUCTOR-TOKEN for each subsequent request after logging in.

Accounts can be added, removed, and have their groups updated.

When Conductor starts, if there are no admin accounts created yet, it will create the default 
admin account using the CONDUCTOR_DEFAULT_ADMIN_USERNAME and CONDUCTOR_DEFAULT_ADMIN_PASSKEY
values. **The very first action to take after starting Conductor the first time is to log in as this user, create a new account in the "admins" group, logging in with the new account, and removing the default account.**

### Token Expiration
Each account has a token expiration time. This value is the expiration time, in seconds, of any
access token generated for that account. A value of 0 means no expiration. This is useful for 
controlling how often a user or service must re-authenticate.

## Groups
Groups are simply a string value. Accounts have a list of groups they are in. If a group is attempted 
to be removed while it is listed in a user, resource, or service's groups, the request will be rejected.

Groups can be added and removed.

There is a default group called "admins" in each Conductor Proxy instance, which users that are able
to manage users, services, and resources can be added to. The admin group has full permissions 
within Conductor.

## Services
A service is an application which provides access to endpoints for one or more resources. A service
can be defined with a host name, base path for all resource endpoint paths, details about the protocol
it uses (HTTP and HTTPS are supported), and authentication requirements. If there are two resources
located on the same service that use different values for any of these attributes, then there should
really be two different services created in Conductor.

Services hold a list of group names that are allowed to manage them, meaning make changes to the service
definition and its resources' definition's. The service also holds a list of group names that are 
allowed to use its resources in the proxy. Service admins may configure services and their resources,
but are not given permission to use a resource in the proxy without also being in the service's user
groups.

## Resources
Resources belong to one service and are a definition of a particular object type, its properties,
and the endpoints available for it.

### Parameters
Parameters are the primary way to pass data into requests. Parameters for HTTP requests can be set in
the request path, query parameters, body, or headers. 

Valid Data Types by Param Type for HTTP Resources:
Header: string, int, bool
Path: string, int, bool
Query: string, int, bool
Body: any
BodyFlat: any

## Proxy
If a request is made to an endpoint with a parameter value supplied that is defined as a body or
a body flat parameter, the following logic is applied. If it is a body flat parameter value is used
to replace the entirety of the incoming request body. If it is a body parameter and the supplied body
is nil or a JSON object, the parameter value will be set using the parameter name as the key and the 
supplied value as the value. If it is a body parameter and the supplied body is not nil or a JSON
object, an error will be returned.

TODO

# Conductor Proxy API Reference
There are two APIs to interact with- the proxy and the admin APIs. The admin API is used
to manage groups, services, resources, and accounts. The proxy is used to use the Conductor Proxy
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

 Admin only in secure mode.

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
> | `403` | `application/json` | `{"statusCode": 403, "message": "Not authorized", "data": {}}` |
</details>

------------------------------------------------------------------------------------------

### Remove Account (TODO)
<details>
 <summary><code>DELETE</code> <code><b>/api/accounts/:id</b></code></summary>

 Admin only in secure mode.

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

 Admin only in secure mode.

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

### Remove Group (TODO)
<details>
 <summary><code>DELETE</code> <code><b>/api/groups/:id</b></code></summary>

 Admin only in secure mode.

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

 Admin only in secure mode.

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

### Add Service
<details>
 <summary><code>POST</code> <code><b>/api/services</b></code></summary>

 Admin only in secure mode.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | name | body (required) | string | Name of service |
> | friendlyName | body (required) | string | Readable name of service |
> | host | body (required) | string | Hostname of service |
> | adminGroups | body (required) | string array | List of groups that can admin this service |
> | userGroups | body (required) | string array | List of groups that can use this service |

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

 Admin or service admin only in secure mode.

##### Parameters
> | name |  type | data type |description |
> |------|-------|-----------|------------|
> | name | body (required) | string | Name of resource |
> | friendlyName | body (required) | string | Readable name of resource |
> | serviceId | body (required) | string | ID of service this resource belongs to |
> | properties | body (required) | Property array | Properties of resource |
> | endpoints | body (required) | Endpoint array | Endpoints of resource |

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

 If running in secure mode, request must include a valid token and the accounts groups must match the service user or admin groups.

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

## Development Setup
To run the Conductor Proxy locally, use the following commands. Conductor Proxy requires
Go 1.21.5. Docker must be installed in your system to use the docker-compose commands.
```sh
git clone https://github.com/Zentech-Development/conductor-proxy.git
cd ./conductor-proxy/
go mod download

go build ./cmd/main.go # outputs executable called main
go run ./cmd/main.go # runs the application
go test ./... # runs all tests
docker-compose -f ./.docker/docker-compose.local.yml up # run Conductor Proxy with Redis
```

# Next Todos
- Add more DB adapters
- Add service and resource validation on add
- Write more test cases