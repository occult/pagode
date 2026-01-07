[![Go Report Card](https://goreportcard.com/badge/github.com/mikestefanello/pagoda)](https://goreportcard.com/report/github.com/mikestefanello/pagoda)
[![Test](https://github.com/mikestefanello/pagoda/actions/workflows/test.yml/badge.svg)](https://github.com/mikestefanello/pagoda/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/mikestefanello/pagoda.svg)](https://pkg.go.dev/github.com/mikestefanello/pagoda)
[![GoT](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

![pagode_banner](https://github.com/user-attachments/assets/756e2fa3-de77-4469-8f6f-6940c83696cc)

## About Pagode

Pagode is a full-stack web application starter kit with expressive, elegant architecture. We believe development must be an enjoyable and creative experience to be truly fulfilling. Pagode takes the pain out of development by combining the power of Go with modern React, providing:

- [Fast, type-safe backend](https://pagode.dev/docs/intro) with Go and Echo.
- [Modern React frontend](https://pagode.dev/docs/intro) with InertiaJS bridge.
- [Powerful ORM](https://pagode.dev/docs/database-and-orm) with Ent code generation.
- [Built-in authentication](https://pagode.dev/docs/authentication) and session management.
- [Background job processing](https://pagode.dev/docs/tasks-and-queues) with SQLite queues.
- [Admin panel](https://pagode.dev/docs/admin-panel) auto-generated for all entities.
- [Hot reload development](https://pagode.dev/docs/intro) experience.

Pagode is accessible, powerful, and provides tools required for large, robust applications.

## Learning Pagode

Pagode has comprehensive [documentation](https://pagode.dev/) and examples to help you get started quickly with the framework

## Getting Started

### Dependencies

Ensure that [Go](https://go.dev/) is installed on your system.

### Getting the Code

Start by checking out the repository. Since this repository is a _template_ and not a Go _library_, you **do not** use `go get`.

```bash
git clone git@github.com:occult/pagode.git
cd pagode
```

### Create an Admin Account

To access the admin panel, you need an admin user account. To create your first admin user, use the command-line:

```bash
make admin email=your@email.com
```

This will generate an admin account using that email address and print the randomly-generated password.

### Start the Application

Before starting, install the frontend dependencies:

```bash
npm install
```

Then, start the Vite frontend development server:

```bash
npx vite
```

From within the root of the codebase, run:

```bash
make run
```

By default, you can access the application at `localhost:8000`. Your data will be stored in the `dbs` directory.

### Live Reloading

For automatic rebuilding when code changes, install [air](https://github.com/air-verse/air) and use:

```bash
make air-install
make watch
```

## Credits

Thank you to all the following amazing projects for making this possible.

- [afero](https://github.com/spf13/afero)
- [gonertia](https://github.com/romsar/gonertia)
- [pagoda](https://github.com/mikestefanello/pagoda)
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
