## Pagode: Modern Go + React starter kit

[![Go Report Card](https://goreportcard.com/badge/github.com/mikestefanello/pagoda)](https://goreportcard.com/report/github.com/mikestefanello/pagoda)
[![Test](https://github.com/mikestefanello/pagoda/actions/workflows/test.yml/badge.svg)](https://github.com/mikestefanello/pagoda/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/mikestefanello/pagoda.svg)](https://pkg.go.dev/github.com/mikestefanello/pagoda)
[![GoT](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

![pagode_banner](https://github.com/user-attachments/assets/756e2fa3-de77-4469-8f6f-6940c83696cc)

## Table of Contents

- [Introduction](#introduction)
  - [Overview](#overview)
  - [Foundation](#foundation)
    - [Backend](#backend)
    - [Frontend](#frontend)
    - [Storage](#storage)
  - [Screenshots](#screenshots)
- [Getting started](#getting-started)
  - [Dependencies](#dependencies)
  - [Getting the code](#getting-the-code)
  - [Create an admin account](#create-an-admin-account)
  - [Start the application](#start-the-application)
  - [Live reloading](#live-reloading)
- [Service container](#service-container)
  - [Dependency injection](#dependency-injection)
  - [Test dependencies](#test-dependencies)
- [Configuration](#configuration)
  - [Environment overrides](#environment-overrides)
  - [Environments](#environments)
- [Database](#database)
  - [Auto-migrations](#auto-migrations)
  - [Separate test database](#separate-test-database)
- [ORM](#orm)
  - [Entity types](#entity-types)
  - [New entity type](#new-entity-type)
- [Sessions](#sessions)
  - [Encryption](#encryption)
- [Authentication](#authentication)
  - [Login / Logout](#login--logout)
  - [Forgot password](#forgot-password)
  - [Registration](#registration)
  - [Admins](#admins)
  - [Authenticated user](#authenticated-user)
    - [Middleware](#middleware)
  - [Email verification](#email-verification)
- [Admin panel](#admin-panel)
  - [Code generation](#code-generation)
  - [Access](#access)
  - [Considerations](#considerations)
  - [Roadmap](#roadmap)
- [Routes](#routes)
  - [Custom middleware](#custom-middleware)
  - [Handlers](#handlers)
  - [Errors](#errors)
  - [Redirects](#redirects)
  - [Testing](#testing)
    - [HTTP server](#http-server)
    - [Request / Request helpers](#request--response-helpers)
    - [Goquery](#goquery)
- [User interface](#user-interface)
  - [Why Gomponents?](#why-gomponents)
  - [HTMX support](#htmx-support)
    - [Header management](#header-management)
    - [Conditional and partial rendering](#conditional-and-partial-rendering)
    - [CSRF token](#csrf-token)
  - [Request](#request)
    - [Title and metatags](#title-and-metatags)
    - [URL generation](#url-generation)
  - [Components](#components)
  - [Layouts](#layouts)
  - [Pages](#pages)
    - [Rendering](#rendering)
  - [Forms](#forms)
    - [Submission processing](#submission-processing)
    - [Inline validation](#inline-validation)
    - [CSRF](#csrf)
  - [Models](#models)
  - [Node caching](#node-caching)
  - [Flash messaging](#flash-messaging)
- [Pager](#pager)
- [Cache](#cache)
  - [Set data](#set-data)
  - [Get data](#get-data)
  - [Flush data](#flush-data)
  - [Flush tags](#flush-tags)
- [Tasks](#tasks)
  - [Queues](#queues)
  - [Dispatcher](#dispatcher)
  - [Monitoring tasks and queues](#monitoring-tasks-and-queues)
- [Cron](#cron)
- [Files](#files)
- [Static files](#static-files)
  - [Cache control headers](#cache-control-headers)
  - [Cache-buster](#cache-buster)
- [Email](#email)
- [HTTPS](#https)
- [Logging](#logging)
- [Credits](#credits)

## Introduction

### Overview

**Pagode** is not a framework — it’s a modern starter kit for building full-stack web applications using **Go**, **InertiaJS**, and **React**, powered by **Tailwind CSS** for styling.

Pagode provides the structure and tooling you need to hit the ground running, without locking you into rigid conventions or heavyweight abstractions. It balances flexibility and productivity by letting you swap out pieces as needed while still offering a solid foundation of battle-tested technologies.

While JavaScript frontends dominate the landscape, Pagode embraces a hybrid approach: it combines server-side rendering and client-side interactivity to deliver fast, modern user experiences — without sacrificing simplicity. Thanks to tools like InertiaJS and Tailwind, you can build reactive, beautiful interfaces with minimal boilerplate and zero custom Webpack/Vite configuration.

Pagode proves that Go is not just for APIs — it's a powerful full-stack solution when paired with the right tools. And yes, you still get the control, speed, and simplicity you love about Go.

### Foundation

While many great projects were used to build this, all of which are listed in the [credits](#credits) section, the following provide the foundation of the back and frontend. It's important to note that you are **<ins>not required to use any of these</ins>**. Swapping any of them out will be relatively easy.

#### Backend

- [Echo](https://echo.labstack.com/): High performance, extensible, minimalist Go web framework.
- [Ent](https://entgo.io/): Simple, yet powerful ORM for modeling and querying data.

#### Frontend

With **server-side rendered HTML** powered by **Go** and **InertiaJS**, you get the best of both worlds — a modern SPA-like experience with server-driven logic. Combined with the tools below, you can build beautiful, dynamic UIs without the usual frontend overhead.

- [InertiaJS](https://inertiajs.com/): Bridges your Go backend with modern JavaScript frameworks like React, enabling full-page SPA experiences without building a separate API.
- [React](https://reactjs.org/): A declarative library for building interactive UIs, perfectly paired with Inertia for dynamic frontend behavior.
- [Tailwind CSS v4](https://tailwindcss.com/): A utility-first CSS framework for rapidly building custom designs directly in your markup — no context switching or naming class conflicts.
- [shadcn/ui](https://ui.shadcn.com/): A beautifully designed, accessible component library built on top of Tailwind CSS and Radix UI — perfect for rapidly building consistent interfaces.

#### Storage

- [SQLite](https://sqlite.org/): A small, fast, self-contained, high-reliability, full-featured, SQL database engine and the most used database engine in the world.

Originally, Postgres and Redis were chosen as defaults but since the aim of this project is rapid, simple development, it was changed to SQLite which now provides the primary data storage as well as persistent, background [task queues](#tasks). For [caching](#cache), a simple in-memory solution is provided. If you need to use something like Postgres or Redis, swapping those in can be done quickly and easily. For reference, [this branch](https://github.com/mikestefanello/pagoda/tree/postgres-redis) contains the code that included those (but is no longer maintained).

### Screenshots

#### Inline form validation

<img src="https://raw.githubusercontent.com/mikestefanello/readmeimages/main/pagoda/inline-validation.png" alt="Inline validation"/>

#### Switch layout templates, user registration

<img src="https://raw.githubusercontent.com/mikestefanello/readmeimages/main/pagoda/register.png" alt="Registration"/>

#### Alpine.js modal, HTMX AJAX request

<img src="https://raw.githubusercontent.com/mikestefanello/readmeimages/main/pagoda/modal.png" alt="Alpine and HTMX"/>

#### User entity list (admin panel)

<img src="https://raw.githubusercontent.com/mikestefanello/readmeimages/main/pagoda/admin-user_list.png" alt="User entity list"/>

#### User entity edit (admin panel)

<img src="https://raw.githubusercontent.com/mikestefanello/readmeimages/main/pagoda/admin-user_edit.png" alt="User entity edit"/>

#### Monitor task queues (provided by Backlite via the admin panel)

<img src="https://raw.githubusercontent.com/mikestefanello/readmeimages/main/backlite/failed.png" alt="Manage task queues"/>

## Getting started

### Dependencies

Ensure that [Go](https://go.dev/) is installed on your system.

### Getting the code

Start by checking out the repository. Since this repository is a _template_ and not a Go _library_, you **do not** use `go get`.

```
git clone git@github.com:occult/pagode.git
cd pagode
```

### Create an admin account

In order to access the [admin panel](#admin-panel), you must log in with an admin user and in order to create your first admin user account, you must use the command-line. Execute `make admin email=your@email.com` from the root of the codebase, and an admin account will be generated using that email address. The console will print the randomly-generated password for the account.

Once you have one admin account, you can use that account to manage other users and admins from within the UI.

### Start the application

Before starting, install the frontend dependencies:

`npm install`

Then, start the Vite frontend development server:

`npx vite`

From within the root of the codebase, simply run `make run`.

By default, you should be able to access the application in your browser at `localhost:8000`. Your data will be stored within the `dbs` directory. If you ever want to quickly delete all data, just remove this directory.

These settings, and many others, can be changed via the [configuration](#configuration).

### Live reloading

Rather than using `make run`, if you prefer live reloading so your app automatically rebuilds and runs whenever you save code changes, start by installing [air](https://github.com/air-verse/air) by running `make air-install`, then use `make watch` to start the application with automatic live reloading.

## Service container

The container is located at `pkg/services/container.go` and is meant to house all of your application's services and/or dependencies. It is easily extensible and can be created and initialized in a single call. The services currently included in the container are:

- Authentication
- Cache
- Configuration
- Database
- Files
- Graph
- Mail
- ORM
- Tasks
- Validator
- Web

A new container can be created and initialized via `services.NewContainer()`. It can be later shutdown via `Shutdown()`, which will attempt to gracefully shutdown all services.

### Dependency injection

The container exists to facilitate easy dependency-injection both for services within the container and areas of your application that require any of these dependencies. For example, the container is automatically passed to the `Init()` method of your route [handlers](#handlers) so that the handlers have full, easy access to all services.

### Test dependencies

It is common that your tests will require access to dependencies, like the database, or any of the other services available within the container. Keeping all services in a container makes it especially easy to initialize everything within your tests. You can see an example pattern for doing this [here](#environments).

## Configuration

The `config` package provides a flexible, extensible way to store all configuration for the application. Configuration is added to the `Container` as a _Service_, making it accessible across most of the application.

Be sure to review and adjust all the default configuration values provided in `config/config.yaml`.

### Environment overrides

Leveraging the functionality of [viper](https://github.com/spf13/viper) to manage configuration, all configuration values can be overridden by environment variables. The name of the variable is determined by the set prefix and the name of the configuration field in `config/config.yaml`.

In `config/config.go`, the prefix is set as `pagoda` via `viper.SetEnvPrefix("pagoda")`. Nested fields require an underscore between levels. For example:

```yaml
http:
  port: 1234
```

can be overridden by setting an environment variable with the name `PAGODA_HTTP_PORT`.

### Environments

The configuration value for the current _environment_ (`Config.App.Environment`) is an important one as it can influence some behavior significantly (will be explained in later sections).

A helper function (`config.SwitchEnvironment`) is available to make switching the environment easy, but this must be executed prior to loading the configuration. The common use-case for this is to switch the environment to `Test` before tests are executed:

```go
func TestMain(m *testing.M) {
    // Set the environment to test
    config.SwitchEnvironment(config.EnvTest)

    // Start a new container
    c = services.NewContainer()

    // Run tests
    exitVal := m.Run()

    // Shutdown the container
    if err := c.Shutdown(); err != nil {
        panic(err)
    }

    os.Exit(exitVal)
}
```

## Database

The database currently used is [SQLite](https://sqlite.org/) but you are free to use whatever you prefer. If you plan to continue using [Ent](https://entgo.io/), the incredible ORM, you can check their supported databases [here](https://entgo.io/docs/dialects). The database driver is provided by [go-sqlite3](https://github.com/mattn/go-sqlite3). A reference to the database is included in the `Container` if direct access is required.

Database configuration can be found and managed within the `config` package.

### Auto-migrations

[Ent](https://entgo.io/) provides automatic migrations which are executed on the database whenever the `Container` is created, which means they will run when the application starts.

### Separate test database

Since many tests can require a database, this application supports a separate database specifically for tests. Within the `config`, the test database can be specified at `Config.Database.TestConnection`, which is the database connection string that will be used. By default, this will be an in-memory SQLite database.

When a `Container` is created, if the [environment](#environments) is set to `config.EnvTest`, the database client will connect to the test database instead and run migrations so your tests start with a clean, ready-to-go database.

When this project was using Postgres, it would automatically drop and recreate the test database. Since the current default is in-memory, that is no longer needed. If you decide to use a test database not in-memory, you can alter the `Container` initialization code to do this for you.

## ORM

As previously mentioned, [Ent](https://entgo.io/) is the supplied ORM. It can be swapped out, but I highly recommend it. I don't think there is anything comparable for Go, at the current time. If you decide to remove Ent, you will lose the dynamic [admin panel](#admin-panel) which allows you to administer all entity types from within the UI. If you're not familiar with Ent, take a look through their top-notch [documentation](https://entgo.io/docs/getting-started).

An Ent client is included in the `Container` to provide easy access to the ORM throughout the application.

Ent relies on code-generation for the entities you create to provide robust, type-safe data operations. Everything within the `ent` directory in this repository is generated code for the two entity types listed below except the [schema declaration](#https://github.com/mikestefanello/pagoda/tree/main/ent/schema) and [custom extension](https://github.com/mikestefanello/pagoda/tree/main/ent/admin) to generate code for the [admin panel](#admin-panel).

### Entity types

The two included entity types are:

- User
- PasswordToken

### New entity type

While you should refer to their [documentation](https://entgo.io/docs/getting-started) for detailed usage, it's helpful to understand how to create an entity type and generate code. To make this easier, the `Makefile` contains some helpers.

1. Ensure all Ent code is downloaded by executing `make ent-install`.
2. Create the new entity type by executing `make ent-new name=User` where `User` is the name of the entity type. This will generate a file like you can see in `ent/schema/user.go` though the `Fields()` and `Edges()` will be left empty.
3. Populate the `Fields()` and optionally the `Edges()` (which are the relationships to other entity types).
4. When done, generate all code by executing `make ent-gen`.

The generated code is extremely flexible and impressive. An example to highlight this is one used within this application:

```go
entity, err := c.ORM.PasswordToken.
    Query().
    Where(passwordtoken.ID(tokenID)).
    Where(passwordtoken.HasUserWith(user.ID(userID))).
    Where(passwordtoken.CreatedAtGTE(expiration)).
    Only(ctx.Request().Context())
```

This executes a database query to return the _password token_ entity with a given ID that belong to a user with a given ID and has a _created at_ timestamp field that is greater than or equal to a given time.

## Sessions

Sessions are provided and handled via [Gorilla sessions](https://github.com/gorilla/sessions) and configured as middleware in the router located at `pkg/handlers/router.go`. Session data is currently stored in cookies but there are many [options](https://github.com/gorilla/sessions#store-implementations) available if you wish to use something else.

Here's a simple example of loading data from a session and saving new values:

```go
func SomeFunction(ctx echo.Context) error {
    sess, err := session.Get(ctx, "some-session-key")
    if err != nil {
        return err
    }
    sess.Values["hello"] = "world"
    sess.Values["isSomething"] = true
    return sess.Save(ctx.Request(), ctx.Response())
}
```

### Encryption

Session data is encrypted for security purposes. The encryption key is stored in [configuration](#configuration) at `Config.App.EncryptionKey`. While the default is fine for local development, it is **imperative** that you change this value for any live environment otherwise session data can be compromised.

## Authentication

Included are standard authentication features you expect in any web application. Authentication functionality is bundled as a _Service_ within `services/AuthClient` and added to the `Container`. If you wish to handle authentication in a different manner, you could swap this client out or modify it as needed.

Authentication currently requires [sessions](#sessions) and the session middleware.

### Login / Logout

The `AuthClient` has methods `Login()` and `Logout()` to log a user in or out. To track a user's authentication state, data is stored in the session including the user ID and authentication status.

Prior to logging a user in, the method `CheckPassword()` can be used to determine if a user's password matches the hash stored in the database and on the `User` entity.

Routes are provided for the user to login and logout at `user/login` and `user/logout`.

### Forgot password

Users can reset their password in a secure manner by issuing a new password token via the method `GeneratePasswordResetToken()`. This creates a new `PasswordToken` entity in the database belonging to the user. The actual token itself, however, is not stored in the database for security purposes. It is only returned via the method so it can be used to build the reset URL for the email. Rather, a hash of the token is stored, using `bcrypt` the same package used to hash user passwords. The reason for doing this is the same as passwords. You do not want to store a plain-text value in the database that can be used to access an account.

Tokens have a configurable expiration. By default, they expire within 1 hour. This can be controlled in the `config` package. The expiration of the token is not stored in the database, but rather is used only when tokens are loaded for potential usage. This allows you to change the expiration duration and affect existing tokens.

Since the actual tokens are not stored in the database, the reset URL must contain the user and password token ID. Using that, `GetValidPasswordToken()` will load a matching, non-expired _password token_ entity belonging to the user, and use `bcrypt` to determine if the token in the URL matches stored hash of the password token entity.

Once a user claims a valid password token, all tokens for that user should be deleted using `DeletePasswordTokens()`.

Routes are provided to request a password reset email at `user/password` and to reset your password at `user/password/reset/token/:user/:password_token/:token`.

### Registration

The actual registration of a user is not handled within the `AuthClient` but rather just by creating a `User` entity. When creating a user, use `HashPassword()` to create a hash of the user's password, which is what will be stored in the database.

A route is provided for the user to register at `user/register`.

### Admins

A checkbox field has been added to the `User` entity type to indicate if the user has admin access. If your app requires a more robust authorization system, such as roles and permissions, you could easily replace this field and adjust all usage of it accordingly. If a user has this field checked, they will be able to access the [admin panel](#admin-panel). [Middleware](#middleware) is provided to easily restrict access to routes based on admin status.

### Authenticated user

The `AuthClient` has two methods available to get either the `User` entity or the ID of the user currently logged in for a given request. Those methods are `GetAuthenticatedUser()` and `GetAuthenticatedUserID()`.

#### Middleware

Registered for all routes is middleware that will load the currently logged in user entity and store it within the request context. The middleware is located at `middleware.LoadAuthenticatedUser()` and, if authenticated, the `User` entity is stored within the context using the key `context.AuthenticatedUserKey`.

If you wish to require either authentication or non-authentication for a given route, you can use either `middleware.RequireAuthentication()` or `middleware.RequireNoAuthentication()`.

If you wish to restrict a route to admins only, you can use `middleware.RequireAdmin`.

### Email verification

Most web applications require the user to verify their email address (or other form of contact information). The `User` entity has a field `Verified` to indicate if they have verified themself. When a user successfully registers, an email is sent to them containing a link with a token that will verify their account when visited. This route is currently accessible at `/email/verify/:token` and handled by `pkg/handlers/auth.go`.

There is currently no enforcement that a `User` must be verified in order to access the application. If that is something you desire, it will have to be added in yourself. It was not included because you may want partial access of certain features until the user verifies; or no access at all.

Verification tokens are [JSON Web Tokens](https://jwt.io/) generated and processed by the [jwt](https://github.com/golang-jwt/jwt) module. The tokens are _signed_ using the encryption key stored in [configuration](#configuration) (`Config.App.EncryptionKey`). **It is imperative** that you override this value from the default in any live environments otherwise the data can be comprimised. JWT was chosen because they are secure tokens that do not have to be stored in the database, since the tokens contain all of the data required, including built-in expirations. These were not chosen for password reset tokens because JWT cannot be withdrawn once they are issued which poses a security risk. Since these tokens do not grant access to an account, the ability to withdraw the tokens is not needed.

By default, verification tokens expire 12 hours after they are issued. This can be changed in configuration at `Config.App.EmailVerificationTokenExpiration`. There is currently not a route or form provided to request a new link.

Be sure to review the [email](#email) section since actual email sending is not fully implemented.

To generate a new verification token, the `AuthClient` has a method `GenerateEmailVerificationToken()` which creates a token for a given email address. To verify the token, pass it in to `ValidateEmailVerificationToken()` which will return the email address associated with the token and an error if the token is invalid.

## Admin panel

The admin panel functionality is considered to be in _beta_ and remains under active development, though all features described here are expected to be fully-functional. Please use caution when using these features and be sure to report any issues you encounter.

The _admin panel_ currently includes:

- A completely dynamic UI to manage all entities defined by _Ent_.
- A section to monitor all [background tasks and queues](#tasks).

There are no separate templates or interfaces for the admin section (see [screenshots](#screenshots)).

Users with admin [access](#access) will see additional links on the default sidebar at the bottom. As with all default UI components, you can easily move these pages and links to a dedicated section, layout, etc. Clicking on the link for any given entity type will provide a pageable table of entities and the ability to add/edit/delete.

### Code generation

In order to automatically and dynamically provide admin functionality for entities, code generation is used by means of leveraging Ent's [extension API](https://entgo.io/docs/extensions) which makes generating code using the Ent graph schema very easy. A [custom extension](https://github.com/mikestefanello/pagoda/blob/master/ent/admin/extension.go) is provided to generate code that provides flat entity type structs and handler code that work directly with Echo. So, both of those are required in order for any of this to work. Whenever you modify one of your entity types or generate a new one, the admin code will also automatically generate.

Without going in to too much detail here, the generated code provides a [handler](https://github.com/mikestefanello/pagoda/blob/master/ent/admin/handler.go) that is then used by a provided [web handler](https://github.com/mikestefanello/pagoda/blob/master/pkg/handlers/admin.go) to power all the routes used in the admin UI. While the rest of the related code should be simple enough to follow, it's worth calling attention to the highly-dynamic [entity form](https://github.com/mikestefanello/pagoda/blob/master/pkg/ui/forms/admin_entity.go) that is constructed using the _Ent_ graph data structure.

### Access

Only admin users can access the admin panel. The details are outlined in the [admins](#admins) and [middleware](#middleware) sections. If you haven't yet generated an admin user, follow [these instructions](#create-an-admin-account).

### Considerations

Since the generated code is completely dynamic, all entity functionality related to creating and editing must be defined within your _Ent_ schema. Refer to the [User](https://github.com/mikestefanello/pagoda/blob/master/ent/schema/user.go) entity schema as an example.

- Field validation must be defined within each entity field (ie, validating an email address in a _string_ field).
- Pre-processing must be defined within entity hooks (ie, hashing the user's password).
- _Sensitive_ fields will be omitted from the UI, and only modified if a value is provided during creation or editing.
- _Edges_ must be bound to an [edge field](https://entgo.io/docs/schema-edges#edge-field) if you want them visible and editable.

### Roadmap

- Determine which tests should be included and provide them.
- Inline validation.
- Either exposed sorting, or allow the _handler_ to be configured with sort criteria for each type.
- Exposed filters.
- Support all field types (types such as _JSON_ as currently not supported).
- Control which fields appear in the entity list table.

## Routes

The router functionality is provided by [Echo](https://echo.labstack.com/guide/routing/) and constructed within via the `BuildRouter()` function inside `pkg/handlers/router.go`. Since the _Echo_ instance is a _Service_ on the `Container` which is passed in to `BuildRouter()`, middleware and routes can be added directly to it.

### Custom middleware

By default, a middleware stack is included in the router that makes sense for most web applications. Be sure to review what has been included and what else is available within _Echo_ and the other projects mentioned.

A `middleware` package is included which you can easily add to along with the custom middleware provided.

### Handlers

A `Handler` is a simple type that handles one or more of your routes and allows you to group related routes together (ie, authentication). All provided handlers are located in `pkg/handlers`. _Handlers_ also handle self-registering their routes with the router.

#### Example

The provided patterns are not required, but were designed to make development as easy as possible.

For this example, we'll create a new handler which includes a GET and POST route and uses the ORM. Start by creating a file at `pkg/handlers/example.go`.

1. Define the handler type:

```go
type Example struct {
    orm *ent.Client
}
```

2. Register the handler so the router automatically includes it

```go
func init() {
    Register(new(Example))
}
```

3. Initialize the handler (and inject any required dependencies from the _Container_). This will be called automatically.

```go
func (e *Example) Init(c *services.Container) error {
    e.orm = c.ORM
    return nil
}
```

4. Declare the routes

**It is highly recommended** that you provide a `Name` for your routes. Most methods on the back and frontend leverage the route name and parameters in order to generate URLs. All route names are currently stored as consts in the `routenames` package so they are accessible from within the `ui` layer.

```go
func (e *Example) Routes(g *echo.Group) {
    g.GET("/example", e.Page).Name = routenames.Example
    g.POST("/example", c.PageSubmit).Name = routenames.ExampleSubmit
}
```

5. Implement your routes

```go
func (e *Example) Page(ctx echo.Context) error {
    // add your code here
}

func (e *Example) PageSubmit(ctx echo.Context) error {
    // add your code here
}
```

### Errors

Routes can return errors to indicate that something wrong happened and an error page should be rendered for the request. Ideally, the error is of type `*echo.HTTPError` to indicate the intended HTTP response code, and optionally a message that will be logged. You can use `return echo.NewHTTPError(http.StatusInternalServerError, "optional message")`, for example. If an error of a different type is returned, an _Internal Server Error_ is assumed.

The [error handler](https://echo.labstack.com/guide/error-handling/) is set to the provided `Handler` in `pkg/handlers/error.go` in the `BuildRouter()` function. That means that if any middleware or route return an error, the request gets routed there. This route passes the status code to the `pages.Error` UI component page, allowing you to easily adjust the markup depending on the error type.

### Redirects

The `pkg/redirect` package makes it easy to perform redirects, especially if you provide names for your routes. The `Redirect` type provides the ability to chain redirect options and also supports automatically handling HTMX redirects for boosted requests.

For example, if your route name is `user_profile` with a URL pattern of `/user/profile/:id`, you can perform a redirect by doing:

```go
return redirect.New(ctx).
    Route("user_profile").
    Params(userID).
    Query(queryParams).
    Go()
```

### Testing

Since most of your web application logic will live in your routes, being able to easily test them is important. The following aims to help facilitate that.

The test setup and helpers reside in `pkg/handlers/router_test.go`.

Only a brief example of route tests were provided in order to highlight what is available. Adding full tests did not seem logical since these routes will most likely be changed or removed in your project.

#### HTTP server

When the route tests initialize, a new `Container` is created which provides full access to all the _Services_ that will be available during normal application execution. Also provided is a test HTTP server with the router added. This means your tests can make requests and expect responses exactly as the application would behave outside of tests. You do not need to mock the requests and responses.

#### Request / Response helpers

With the test HTTP server setup, test helpers for making HTTP requests and evaluating responses are made available to reduce the amount of code you need to write. See `httpRequest` and `httpResponse` within `pkg/handlers/router_test.go`.

Here is an example how to easily make a request and evaluate the response:

```go
func TestAbout_Get(t *testing.T) {
    doc := request(t).
        setRoute("about").
        get().
        assertStatusCode(http.StatusOK).
        toDoc()
}
```

#### Goquery

A helpful, included package to test HTML markup from HTTP responses is [goquery](https://github.com/PuerkitoBio/goquery). This allows you to use jQuery-style selectors to parse and extract HTML values, attributes, and so on.

In the example above, `toDoc()` will return a `*goquery.Document` created from the HTML response of the test HTTP server.

Here is a simple example of how to use it, along with [testify](https://github.com/stretchr/testify) for making assertions:

```go
h1 := doc.Find("h1.title")
assert.Len(t, h1.Nodes, 1)
assert.Equal(t, "About", h1.Text())
```

## User interface

### Why React + InertiaJS?

Modern single-page interactions, rich component ecosystems, and type-safety are now table-stakes for productive web development. By pairing **React 18** with **InertiaJS**, Pagode keeps the simplicity of server-side routing while delivering a fully interactive SPA experience.

- **No separate API layer** – Controllers still live in Go; they just return JSON “page props” for React.
- **Reuse the entire npm ecosystem** – Charts, editors, drag-and-drop, and every other React package drop right in.
- **Zero client-side routing boilerplate** – Inertia intercepts links and form submissions automatically.
- **Typed front-end** – Ship confident UIs with TypeScript and shadcn/ui primitives styled by Tailwind v4.
- **Faster iteration** – Hot-reload for both Go and React via Vite; no template compilation steps.

#### CSRF token

If [CSRF](#csrf) protection is enabled, the token value will automatically be passed to HTMX to be included in all non-GET requests. This is done in the `JS()` [component](#components) by leveraging HTMX [events](https://htmx.org/reference/#events).

### Request

The `Request` type in the `ui` package is a foundational helper that provides useful data from the current request as well as resources and methods that make rendering UI components much easier. Using the `echo.Context`, a `Request` can be generated by executing `ui.NewRequest(ctx)`. As you develop and expand your application, you will likely want to expand this type to include additional data and methods that your frontend requires.

`NewRequest()` will automatically populate the following fields using the `echo.Context` from the current request:

- `Context`: The provided _echo.Context_
- `CurrentPath`: The requested URL path
- `IsHome`: If the request was for the homepage
- `IsAuth`: If the user is authenticated
- `AuthUser`: The logged-in user entity, if one
- `CSRF`: The CSRF token, if the middleware is being used
- `Htmx`: Data from the HTMX headers, if HTMX made the request
- `Config`: The application configuration, if the middleware is being used

#### Title and metatags

The `Request` type has additional fields to make it easy to set static values within components being rendered on a given page. While the _title_ is always important, the others are provided as an example:

- `Title`: The page title
- `Metatags`:
  - `Description`: The description of the page
  - `Tags`: A slice of keyword tags

#### URL generation

As mentioned in the [Routes](#routes) section, it is recommended, though not required, to provide names for each of your routes. These are currently defined as consts in the `routenames` package. If you use names for your routes, you can leverage the URL generation methods on the `Request`. This allows you to prevent hard-coding your route paths and parameters in multiple places.

The methods both take a route name and optional variadic route parameters:

- `Path()`: Generates a relative path for a given route.
- `Url()`: Generates an absolute URL for a given route. This uses the `App.Host` field in your [configuration](#configuration) to determine the host of the URL.

**Example:**

```go
g.GET("/user/:uid", profilePage).Name = routenames.Profile
```

```go
func ProfileLink(r *ui.Request, userName string, userID int64) gomponents.Node {
    return A(
        Class("profile"),
        Href(r.Path(routenames.Profile, userID)),
        Text(userName),
    )
}
```

### Components

The [components package](https://github.com/mikestefanello/pagoda/tree/templates/pkg/ui/components) is meant to be your library of reusable _gomponent_ components. Having this makes building your [layouts](#layouts), [pages](#pages), [forms](#forms), [models](#models) and the rest of your user interface much easier. Some of the examples provided include components for [flash messages](#flash-messaging), navigation menus, tabs, metatags, and form elements used to automatically provide [inline validation](#inline-form-validation).

Your components can also make using utility-based CSS libraries, such as [Tailwind CSS](https://tailwindcss.com/), much easier by avoiding excessive duplication of classes across elements.

### Layouts

_Layouts_ are full HTML templates that are used by [pages](#pages) to inject themselves in to, allowing you to easily have multiple pages that all use the same layout, and to easily switch layouts between different pages. [Included](https://github.com/mikestefanello/pagoda/tree/templates/pkg/ui/layouts) is a _primary_ and _auth_ layout as an example, which you can see in action by navigating between the links on the _General_ and _Account_ sidebar menus.

### Pages

_Pages_ are what [route handlers](#handlers) ultimately assemble and render. They may accept primitives, [models](#models), [forms](#forms), or nothing at all, and they embed themselves in a [layout](#layouts) of their choice. Each _page_ represents a different page of your web application and many [examples](https://github.com/mikestefanello/pagoda/tree/templates/pkg/ui/pages) are provided for reference. See below for a minimal example.

#### Rendering

The `Request` type contains a `Render()` method which makes rendering your page within a given layout simple. It automatically handles partial rendering, omitting the [layout](#layouts) and only rendering the [page](#pages) if the request is made by HTMX and is not boosted. Using HTMX is completely optional. This is accomplished by passing in your layout and _page_ separately, for example:

```go
func MyPage(ctx echo.Context, username string) error {
    r := ui.NewRequest(ctx)
    r.Title = "My page"
    node := Div(Textf("Hello, %s!", username))
    return r.Render(layouts.Primary, node)
}
```

Using `Render()`, in this example, only `node` will render if HTMX made the request in a non-boosted fashion, otherwise `node` will render within `layouts.Primary`.

And from within your [route handler](#handlers), simply:

```go
func (e *ExampleHandler) Page(ctx echo.Context) error {
    return pages.MyPage(ctx, "abcd")
}
```

### Forms

Building, rendering, validating and processing forms is made extremely easy with [Echo binding](https://echo.labstack.com/guide/binding/), [validator](https://github.com/go-playground/validator), [form.Submission](https://github.com/mikestefanello/pagoda/blob/templates/pkg/form/submission.go), and the provided _gomponent_ [components](#components).

Start by declaring the form within the [forms](https://github.com/mikestefanello/pagoda/tree/templates/pkg/ui/forms) package:

```go
type Guestbook struct {
    Message    string `form:"message" validate:"required"`
    form.Submission
}
```

Embedding `form.Submission` satisfies the `form.Form` interface and handles submissions and validation for you.

Next, provide a method that renders the form:

```go
func (f *Guestbook) Render(r *ui.Request) Node {
    return Form(
        ID("guestbook"),
        Method(http.MethodPost),
        Attr("hx-post", r.Path(routenames.GuestbookSubmit)),
        TextareaField(TextareaFieldParams{
            Form:      f,
            FormField: "Message",
            Name:      "message",
            Label:     "Message",
            Value:     f.Message,
        }),
        ControlGroup(
            FormButton("is-link", "Submit"),
        ),
        CSRF(r),
    )
}
```

Then, create a _page_ that includes your form:

```go
func UserGuestbook(ctx echo.Context, form *forms.Guestbook) error {
    r := ui.NewRequest(ctx)
    r.Title = "User page"

    content := Div(
        Class("guestbook"),
        H2(Text("My guestbook")),
        P(Text("Hi, please sign my guestbook!")),
        form.Render(r)
    )

    return r.Render(layouts.Primary, content)
}
```

And last, have your handler render the _page_ in a route, and provide a route for the submission.

```go
func (e *Example) Routes(g *echo.Group) {
    g.GET("/guestbook", e.Page).Name = routenames.Guestbook
    g.POST("/guestbook", c.PageSubmit).Name = routenames.GuestbookSubmit
}

func (e *Example) Page(ctx echo.Context) error {
    return pages.UserGuestbook(ctx, form.Get[forms.Guestbook](ctx))
}
```

`form.Get` will either initialize a new form, or load one previously stored in the context (ie, if it was already submitted).

#### Submission processing

Using the example form above, this is all you would have to do within the _POST_ callback for your route:

Start by submitting the form via `form.Submit()`, along with the request context. This will:

1. Store a pointer to the form in the _context_ so that your _GET_ callback can access the form values (shown previously). That allows the form to easily be re-rendered with any validation errors it may have as well as the values that were provided.
2. Parse the input in the _POST_ data to map to the struct so the fields becomes populated. This uses the `form` struct tags to map form input values to the struct fields.
3. Validate the values in the struct fields according to the rules provided in the optional `validate` struct tags.

Then, evaluate the error returned, if one, and process the form values however you need to:

```go
func (e *Example) Submit(ctx echo.Context) error {
    var input forms.Guestbook

    // Submit the form.
    err := form.Submit(ctx, &input)

    // Check the error returned, and act accordingly.
    switch err.(type) {
    case nil:
        // All good!
    case validator.ValidationErrors:
        // The form input was not valid, so re-render the form with the errors included.
        return e.Page(ctx)
    default:
        // Request failed, show the error page.
        return err
    }

    msg.Success(fmt.Sprintf("Your message was: %s", input.Message))

    return redirect.New(ctx).
        Route(routenames.Home).
        Go()
}
```

#### Inline validation

The `Submission` makes inline validation easier because it will store all validation errors in a map, keyed by the form struct field name. It also contains helper methods that the provided form [components](#components), such as `TextareaField` shown in the example above, use to automatically provide classes and error messages. The example form above will have inline validation without requiring anything other than what is shown above.

While [validator](https://github.com/go-playground/validator) is a great package that is used to validate based on struct tags, the downside is that the messaging, by default, is not very human-readable or easy to override. Within `Submission.setErrorMessages()` the validation errors are converted to more readable messages based on the tag that failed validation. Only a few tags are provided as an example, so be sure to expand on that as needed.

#### CSRF

By default, all non `GET` requests will require a CSRF token be provided as a form value. This is provided by middleware and can be adjusted or removed in the router.

The `Request` automatically extracts the CSRF token from the context, but you must include it in your forms by using the provided `CSRF()` [component](#components) as shown in the example above.

### Models

Models are objects built and provided by your _routes_ that can be rendered by your _ui_. Though not required, they reside in the [models package](https://github.com/mikestefanello/pagoda/tree/main/pkg/ui/models) and each has a `Render()` method, making them easy to render within your [pages](#pages). Please see example routes such as the homepage and search for an example.

### Node caching

While most likely unnecessary for most applications, but because optimizing software is fun, a simple `gomponents.Node` cache is provided. This is not because _gomponents_ is inefficient, in fact my basic benchmarks put it as either similar or slightly better than Go templates, but rather because there are _some_ performance gains to be seen by caching static nodes and it may seem wasteful to build and render static HTML on every single page load. It is important to note, you can only cache nodes that are static and will never change.

A good example of this, and one included, is the entire upper navigation bar, search form, and search modal in the _Primary_ layout. It contains a large amount of nested _gomponent_ function calls and a lot of rendering is required. There is no reason to do this more than once.

The cache functions are available in `pkg/ui/cache` and can most easily used like this:

```go
func SearchModal() gomponents.Node {
    return cache.SetIfNotExists("searchModal", func() gomponents.Node {
        return Div(...your entire nested node...)
    })
}
```

`cache.SetIfNotExists()`is a helper function that uses `cache.Get()` to check if the `Node` is already cached under the provided _key_, and if not, executes the _func_ to generate the `Node`, and caches that via `cache.Set()`.

`cache.Set()` does more than just cache the `Node` in-memory. It renders the entire `Node` into a `bytes.Buffer`, then stores a `Raw()` `Node` using the rendered content. This means that everytime the `Node` is taken from the cache and rendered, the pre-rendered `string` is used rather than having to iterate through the nested component, executing all of the element functions and rendering and building the entire HTML output.

It's worth noting that my benchmarking was very limited and cannot be considered anything definitive. In my tests, gomponents was faster, allocated less overall, but had more allocations in total. If you're able to cache static nodes, gomponents can perform significantly better. Reiterating, for most applications, these differences in nanoseconds and bytes will most likely be completely insignificant and unnoticed; but it's worth being aware of.

### Flash messaging

Flash messaging functionality is provided within the `msg` package. It is used to provide one-time status messages to users.

Flash messaging requires that [sessions](#sessions) and the session middleware are in place since that is where the messages are stored.

#### Creating messages

There are four types of messages, and each can be created as follows:

- Success: `msg.Success(ctx echo.Context, message string)`
- Info: `msg.Info(ctx echo.Context, message string)`
- Warning: `msg.Warning(ctx echo.Context, message string)`
- Danger: `msg.Danger(ctx echo.Context, message string)`

#### Rendering messages

When a flash message is retrieved from storage in order to be rendered, it is deleted from storage so that it cannot be rendered again.

A [component](#components), `FlashMessages()`, is provided to render flash messages within your UI.

## Pager

A very basic mechanism is provided to handle and facilitate paging located in `pkg/pager` and can be initialized via `pager.NewPager()`. If the requested URL contains a `page` query parameter with a numeric value, that will be set as the page number in the pager. This query key can be controlled via the `QueryKey` constant.

Methods include:

- `SetItems(items int)`: Set the total amount of items in the entire result-set
- `IsBeginning()`: Determine if the pager is at the beginning of the pages
- `IsEnd()`: Determine if the pager is at the end of the pages
- `GetOffset()`: Get the offset which can be useful in constructing a paged database query

There is currently no generic component to easily render a pager, but the homepage does have an example.

## Cache

As previously mentioned, the default cache implementation is a simple in-memory store, backed by [otter](https://github.com/maypok86/otter), a lockless cache that uses [S3-FIFO](https://s3fifo.com/) eviction. The `Container` houses a `CacheClient` which is a useful wrapper to interact with the cache (see examples below). Within the `CacheClient` is the underlying store interface `CacheStore`. If you wish to use a different store, such as Redis, and want to keep using the `CacheClient`, simply implement the `CacheStore` interface with a Redis library and adjust the `Container` initialization to use that.

The built-in usage of the cache is currently only used for a simple example route located at `/cache` where you can set and view the value of a given cache entry.

Since the current cache is in-memory, there's no need to adjust the `Container` during tests. When this project used Redis, the configuration had a separate database that would be used strictly for tests to avoid writing to your primary database. If you need that functionality, it is easy to add back in.

### Set data

**Set data with just a key:**

```go
err := c.Cache.
    Set().
    Key("my-key").
    Data(myData).
    Expiration(time.Hour * 2).
    Save(ctx)
```

**Set data within a group:**

```go
err := c.Cache.
    Set().
    Group("my-group").
    Key("my-key").
    Expiration(time.Hour * 2).
    Data(myData).
    Save(ctx)
```

**Include cache tags:**

```go
err := c.Cache.
    Set().
    Key("my-key").
    Tags("tag1", "tag2").
    Expiration(time.Hour * 2).
    Data(myData).
    Save(ctx)
```

### Get data

```go
data, err := c.Cache.
    Get().
    Group("my-group").
    Key("my-key").
    Fetch(ctx)
```

### Flush data

```go
err := c.Cache.
    Flush().
    Group("my-group").
    Key("my-key").
    Execute(ctx)
```

### Flush tags

This will flush all cache entries that were tagged with the given tags.

```go
err := c.Cache.
    Flush().
    Tags("tag1", "tag2").
    Execute(ctx)
```

### Tagging

As shown in the previous examples, cache tags were provided because they can be convenient. However, maintaining them comes at a cost and it may not be a good fit for your application depending on your needs. When including tags, the `CacheClient` must lock in order to keep the tag index in sync. And since the tag index cannot support eviction, since that could result in a flush call not actually flushing the tag's keys, the maps that provide the index do not have a size limit. See the code for more details.

## Tasks

Tasks are queued operations to be executed in the background, either immediately, at a specific time, or after a given amount of time has passed. Some examples of tasks could be long-running operations, bulk processing, cleanup, notifications, etc.

Since we're already using [SQLite](https://sqlite.org/) for our database, it's available to act as a persistent store for queued tasks so that tasks are never lost, can be retried until successful, and their concurrent execution can be managed. [Backlite](https://github.com/mikestefanello/backlite) is the library chosen to interface with [SQLite](https://sqlite.org/) and handle queueing tasks and processing them asynchronously. I wrote that specifically to address the requirements I wanted to satisfy for this project.

To make things easy, the _Backlite_ client is provided as a _Service_ on the `Container` which allows you to register queues and add tasks.

Configuration for the _Backlite_ client is exposed through the app's yaml configuration. The required database schema will be automatically installed when the app starts.

### Queues

A full example of a queue implementation can be found in `pkg/tasks` with an interactive form to create a task and add to the queue at `/task` (see `pkg/handlers/task.go`). Also refer to the [Backlite](https://github.com/mikestefanello/backlite) documentation for reference and examples.

See `pkg/tasks/register.go` for a simple way to register all of your queues and to easily pass the `Container` to them so the queue processor callbacks have access to all of your app's dependencies.

### Dispatcher

The _task dispatcher_ is what manages the worker pool used for executing tasks in the background. It monitors incoming and scheduled tasks and handles sending them to the pool for execution by the queue's processor callback. This must be started in order for this to happen. In `cmd/web/main.go`, the _task dispatcher_ is automatically started when the app starts via:

```go
c.Tasks.Start(ctx)
```

The app [configuration](#configuration) contains values to configure the client and dispatcher including how many goroutines to use, when to release stuck tasks back into the queue, and how often to cleanup retained tasks in the database.

When the app is shutdown, the dispatcher is given 10 seconds to wait for any in-progress tasks to finish execution. This can be changed in `cmd/web/main.go`.

### Monitoring tasks and queues

The [admin panel](#admin-panel) contains the UI provided by [Backlite](https://github.com/mikestefanello/backlite) in order to fully monitor all tasks and queues from within your browser.

## Cron

By default, no cron solution is provided because it's very easy to add yourself if you need this. You can either use a [ticker](https://pkg.go.dev/time#Ticker) or a [library](https://github.com/robfig/cron).

## Files

To handle file management functionality such as file uploads, an abstracted file system interface is provided as a _Service_ on the `Container` powered by [afero](https://github.com/spf13/afero). This allows you to easily change the file system backend (ie, local, GCS, SFTP, in-memory) without having to change any of the application code other than the initialization on the `Container`. By default, the local OS is used with a directory specified in the application configuration (which defaults to `uploads`). When running tests, an in-memory file system backend is automatically used.

A simple file upload form example is provided at `/files` which also dynamically lists all files previously uploaded. No database entities or entries are created or provided for files and uploaded files are not available to be served. You will have to implement whatever functionality your application needs.

## Static files

Static files are currently configured in the router (`pkg/handler/router.go`) to be served from the `static` directory. If you wish to change the directory, alter the constant `config.StaticDir`. The URL prefix for static files is `/files` which is controlled via the `config.StaticPrefix` constant.

### Cache control headers

Static files are grouped separately so you can apply middleware only to them. Included is a custom middleware to set cache control headers (`middleware.CacheControl`) which has been added to the static files router group.

The cache max-life is controlled by the configuration at `Config.Cache.Expiration.StaticFile` and defaults to 6 months.

### Cache-buster

While it's ideal to use cache control headers on your static files so browsers cache the files, you need a way to bust the cache in case the files are changed. In order to do this, a function, `File()`, is provided in the `ui` package to generate a static file URL for a given file that appends a cache-buster query. This query string is generated using the timestamp of when the app started and persists until the application restarts.

For example, to render a file located in `static/picture.png`, you would use:

```go
return Img(Src(ui.File("picture.png")))
```

Which would result in:

```html
<img src="/files/picture.png?v=1741053493" />
```

Where `1741053493` is the cache-buster.

## Email

An email client was added as a _Service_ to the `Container` but it is just a skeleton without any actual email-sending functionality. The reason is that there are a lot of ways to send email and most prefer using a SaaS solution for that. That makes it difficult to provide a generic solution that will work for most applications.

The structure in the client (`MailClient`) makes composing emails very easy, and you have the option to construct the body using either a simple string or with a renderable _gomponent_, as explained in the [user interface](#user-interface), in order to produce HTML emails. A simple example is provided in `pkg/ui/emails`.

The standard library can be used if you wish to send email via SMTP and most SaaS providers have a Go package that can be used if you choose to go that direction. **You must** finish the implementation of `MailClient.send`.

The _from_ address will default to the configuration value at `Config.Mail.FromAddress`. This can be overridden per-email by calling `From()` on the email and passing in the desired address.

See below for examples on how to use the client to compose emails.

**Sending with a string body**:

```go
err = c.Mail.
    Compose().
    To("hello@example.com").
    Subject("Welcome!").
    Body("Thank you for registering.").
    Send(ctx)
```

**Sending an HTML body using a gomponent**:

```go
err = c.Mail.
    Compose().
    To("hello@example.com").
    Subject("Confirm your email address").
    Component(emails.ConfirmEmailAddress(ctx, username, token)).
    Send(ctx)
```

This will use the HTML provided when rendering the _gomponent_ as the email body.

## HTTPS

By default, the application will not use HTTPS but it can be enabled easily. Just alter the following configuration:

- `Config.HTTP.TLS.Enabled`: `true`
- `Config.HTTP.TLS.Certificate`: Full path to the certificate file
- `Config.HTTP.TLS.Key`: Full path to the key file

To use _Let's Encrypt_ follow [this guide](https://echo.labstack.com/cookbook/auto-tls/#server).

## Logging

By default, the [Echo logger](https://echo.labstack.com/guide/customization/#logging) is not used for the following reasons:

1. It does not support structured logging, which makes it difficult to deal with variables, especially if you intend to store a logger in context with pre-populated attributes.
2. The upcoming v5 release of Echo will not contain a logger.
3. It should be easy to use whatever logger you prefer (even if that is Echo's logger).

The provided implementation uses the relatively new [log/slog](https://go.dev/blog/slog) library which was added to Go in version 1.21 but swapping that out for whichever you prefer is very easy.

### Context

The simple `pkg/log` package provides the ability to set and get a logger from the Echo context. This is especially useful when you use the provided logger middleware (see below). If you intend to use a different logger, modify these methods to receive and return the logger of your choice.

### Usage

Adding a logger to the context:

```go
logger := slog.New(logHandler).With("id", requestId)
log.Set(ctx, logger)
```

Access and use the logger:

```go
func (h *handler) Page(ctx echo.Context) error {
    log.Ctx(ctx).Info("send a message to the log",
      "value_one", valueOne,
      "value_two", valueTwo,
    )
}
```

### Log level

When the _Container_ configuration is initialized (`initConfig()`), the `slog` default log level is set based on the environment. `INFO` is used for production and `DEBUG` for everything else.

### Middleware

The `SetLogger()` middleware has been added to the router which sets an initialized logger on the request context. It's recommended that this remains after Echo's `RequestID()` middleware because it will add the request ID to the logger which means that all logs produced for that given request will contain the same ID, so they can be linked together. If you want to include more attributes on all request logs, set those fields here.

The `LogRequest()` middleware is a replacement for Echo's `Logger()` middleware which produces a log of every request made, but uses our logger rather than Echo's.

```
2024/06/15 09:07:11 INFO GET /contact request_id=gNblvugTKcyLnBYPMPTwMPEqDOioVLKp ip=::1 host=localhost:8000 referer="" status=200 bytes_in=0 bytes_out=5925 latency=107.527803ms
```

## Credits

Thank you to all the following amazing projects for making this possible.

- [afero](https://github.com/spf13/afero)
- [gonertia](https://github.com/romsar/gonertia)
- [inertiajs](https://inertiajs.com/)
- [laravel](https://github.com/laravel)
- [tailwindcss](https://github.com/tailwindlabs/tailwindcss)
- [shadcn](https://github.com/shadcn-ui/ui)
- [air](https://github.com/air-verse/air)
- [backlite](https://github.com/mikestefanello/backlite)
- [echo](https://github.com/labstack/echo)
- [ent](https://github.com/ent/ent)
- [go](https://go.dev/)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)
- [goquery](https://github.com/PuerkitoBio/goquery)
- [jwt](https://github.com/golang-jwt/jwt)
- [otter](https://github.com/maypok86/otter)
- [sessions](https://github.com/gorilla/sessions)
- [sqlite](https://sqlite.org/)
- [testify](https://github.com/stretchr/testify)
- [validator](https://github.com/go-playground/validator)
- [viper](https://github.com/spf13/viper)
