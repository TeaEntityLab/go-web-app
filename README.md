
go-web-app
==

A web app template and Cli with structures inspired by RoR/Laravel/AdonisJS.

# Features

* cli (*`urfave/cli`*)
	* Generate Model/Migration/Seeder
	* seed
	* migrate
* DB (*`gorm`*)
* Env/Logger (*`caarlos0/env`* *`logrus`*)
* Basic Auth middleware
* APIDocs serve/gen (*`swaggo`*)

# Why

It's nice to write something in RoR ways.

I hope you'll enjoy this template

# Example

Run:
```bash
./run.sh app

# Run with Docs:
ENABLE_DOCUMENT=true ./run.sh app
```

Generate Docs (*`swaggo`*):
```bash
cd app; ./generateDocs.sh; cd ..
```

CLI:
```bash
# Help
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh help

# Make
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh make:model
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh make:migration
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh make:seeder

# Migrate
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh migrate
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh migrate:rollback
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh migrate:refresh

# Seed
DB_ENDPOINTS=mydb.sqlite3 DB_TYPE=sqlite ./cli.sh seed
```

# Mod Dependencies/Suggestions

## WSGI/Routing/Auth

* Auth
	* JWT (Auth Token) https://github.com/dgrijalva/jwt-go
	* Crypto cypto "golang.org/x/crypto"
* Env (as Config structs) env "github.com/caarlos0/env/v6"
* ORM (GORM) https://gorm.io/gorm
  	* Migration https://github.com/go-gormigrate/gormigrate/v2
	* driver-sqlite "gorm.io/driver/sqlite"
	* driver-mysql "gorm.io/driver/mysql"
	* driver-postgres "gorm.io/driver/postgres"
* Gin

## Log/Error

* errors "github.com/pkg/errors"
* logrus "github.com/sirupsen/logrus"
* slack-go "github.com/johntdyer/slack-go"

## Cache/Data/Serialize

* Cache
	* LRUCache github.com/hashicorp/golang-lru
	* Redis gopkg.in/redis.v5
* Data/Serialize
	* jsoniter github.com/json-iterator/go
	* jsondiff github.com/nsf/jsondiff
	* diffmatchpatch github.com/sergi/go-diff/diffmatchpatch
	* hashring github.com/serialx/hashring
* XID github.com/rs/xid
* Net golang.org/x/net

## Cli/Etc

* Cli https://github.com/urfave/cli/blob/master/docs/v2/manual.md
