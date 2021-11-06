# ksqlDB Go library

[![Go Reference](https://pkg.go.dev/badge/github.com/thmeitz/ksqldb-go.svg)](https://pkg.go.dev/github.com/thmeitz/ksqldb-go)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=thmeitz_ksqldb-go&metric=coverage)](https://sonarcloud.io/summary/new_code?id=thmeitz_ksqldb-go)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=thmeitz_ksqldb-go&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=thmeitz_ksqldb-go)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=thmeitz_ksqldb-go&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=thmeitz_ksqldb-go)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=thmeitz_ksqldb-go&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=thmeitz_ksqldb-go)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=thmeitz_ksqldb-go&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=thmeitz_ksqldb-go)

This is a unconnected fork from [Robin Moffatt](https://github.com/rmoff/ksqldb-go) and will be developed on its own.

Thank you Robin and all other contributors for their work!

## Attention - WIP

If you use this library, be warned as the client will be completely overhauled!
This client is `not production ready`!!!

⚠️ Disclaimer #1: This is a personal project and not supported or endorsed by Confluent.

## Description

This is a Go client for [ksqlDB](https://ksqldb.io/).

- [x] Execute a statement (/ksql endpoint)
- [ ] Run a query (/query endpoint)
- [x] Run push and pull queries (/query-stream endpoint)
- [ ] Terminate a cluster (/ksql/terminate endpoint)
- [ ] Introspect query status (/status endpoint)
- [x] Introspect server status (/info endpoint)
- [ ] Introspect cluster status (/clusterStatus endpoint)
- [ ] Get the validity of a property (/is_valid_property)

### KSqlParser

- parse your ksql-statements with the provided `KSqlParser`.

## Installation

Module install:

This client is a Go module, therefore you can have it simply by adding the following import to your code:

```golang
import "github.com/thmeitz/ksqldb-go"
```

Then run a build to have this client automatically added to your go.mod file as a dependency.

Manual install:

```bash
go get -u github.com/thmeitz/ksqldb-go
```

or use the client and and run

```bash
go mod tidy
```

## How to use the ksqldb-go package?

### Create a ksqlDB Client

```golang
import (
	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/logrus"
)

var (
	logger = logrus.NewStandard()
)
// than later in your code...
func main {
  options := ksqldb.Options{
    // if you need a login, do this; if not its not necessary
    Credentials: ksqldb.Credentials{Username: "myuser", Password: "mypassword"},
    // defaults to http://localhost:8088
    BaseUrl:     "http://my-super-shiny-ksqldbserver:8082",
    // this is needed, because the ksql api communicates with http2 only
    AllowHTTP:   true,
  }

  client, err := ksqldb.NewClient(options, log.Current)
  if err != nil {
     log.Fatal(err)
  }

  // then make a pull, push, execute request

  // if you finished your work, you **MUST** close the http.Transport!!!
  client.Close()
}
```

For no authentication just use blank username and password values.

### QueryBuilder (since v0.0.3)

SQL strings should be build by a QueryBuilder. Otherwise the system is open for SQL injections (see [go-webapp-scp.pdf](https://github.com/OWASP/Go-SCP/blob/master/dist/go-webapp-scp.pdf) ).

You can add multiple parameters `Bind(nil, 1, 2.5686, "string", true)`.

`nil` will be converted to `NULL`.

The number of parameters must match the parameters in the SQL statement. If not, an error is thrown.

```golang
//see file: examples/cobra-test/cmd/pull.go

k := `SELECT TIMESTAMPTOSTRING(WINDOWSTART,'yyyy-MM-dd HH:mm:ss','Europe/London') AS WINDOW_START,
TIMESTAMPTOSTRING(WINDOWEND,'HH:mm:ss','Europe/London') AS WINDOW_END,
DOG_SIZE, DOGS_CT FROM DOGS_BY_SIZE
WHERE DOG_SIZE=?;`

builder, err := ksqldb.DefaultQueryBuilder(k)
if err != nil {
	log.Fatal(err)
}

stmnt, err := builder.Bind("middle")
if err != nil {
	log.Fatal(err)
}

ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
defer cancel()
_, r, err := ksqldb.Pull(client, ctx, *stmnt, true)
if err != nil {
	log.Fatal(err)
}
```

### Pull query

```golang

// we are using the client, we are created

ctx, ctxCancel := context.WithTimeout(context.Background(), 10 \* time.Second)
defer ctxCancel()

k := "SELECT TIMESTAMPTOSTRING(WINDOWSTART,'yyyy-MM-dd HH:mm:ss','Europe/London') AS WINDOW_START, TIMESTAMPTOSTRING(WINDOWEND,'HH:mm:ss','Europe/London') AS WINDOW_END, DOG_SIZE, DOGS_CT FROM DOGS_BY_SIZE WHERE DOG_SIZE='" + s + "';"

// your select statement will be checked with integrated KSqlParser
_, r, e := ksqldb.Pull(client, ctx, k, false)
if e != nil {
  // handle the error better here, e.g. check for no rows returned
  return fmt.Errorf("error running pull request against ksqlDB:\n%v", e)
}

var dogSize string
var dogsCt float64
for *, row := range r {
  if row != nil {
    // Should do some type assertions here
    dogSize = row[2].(string)
    dogsCt = row[3].(float64)
    fmt.Printf("🐶 There are %v dogs size %v\n", dogsCt, dogSize)
  }
}
```

### Push query

```golang
rc := make(chan ksqldb.Row)
hc := make(chan ksqldb.Header, 1)

k := "SELECT ROWTIME, ID, NAME, DOGSIZE, AGE FROM DOGS EMIT CHANGES;"

// This Go routine will handle rows as and when they
// are sent to the channel
go func() {
  var name string
  var dogSize string
  for row := range rc {
    if row != nil {
      // Should do some type assertions here
      name = row[2].(string)
      dogSize = row[3].(string)
      fmt.Printf("🐾%v: %v\n",  name, dogSize)
    }
  }
}()

ctx, ctxCancel := context.WithTimeout(context.Background(), 10 \* time.Second)
defer ctxCancel()

e := ksqldb.Push(client, ctx, k, rc, hc)

if e != nil {
  // handle the error better here, e.g. check for no rows returned
  return fmt.Errorf("ksqldbPushError:\n%v", e)
}
```

### Execute a command

```golang
if err := ksqldb.Execute(client, ctx, ksqlDBServer, `CREATE STREAM DOGS (ID STRING KEY, NAME STRING, DOGSIZE STRING, AGE STRING) WITH (KAFKA_TOPIC='dogs', VALUE_FORMAT='JSON');`); err != nil {
  return fmt.Errorf("error creating the dogs stream.\n%v", err)
}
```

## Examples

You can find the examples in the [examples directory](examples).

- [cobra example](examples/cobra-test)
- [KSqlGrammar example](examples/ksqlgrammar)

See the [test environment here](examples/cobra-test/environment.adoc)

### Cobra example

The [Cobra](https://github.com/spf13/cobra) `cobra-test` example shows basic usage of the `ksqldb-go` package. To run it, you need a `Kafka` runtime environment.

The example splits the different use cases into `Cobra` commands.

Start [docker-compose](examples/cobra/docker-compose.yml).

```bash
docker-compose up -d
go run ./examples/cobra-test
```

It outputs:

```bash
ksqldb-go example with cobra

Usage:
  cobra-test [command]

Available Commands:
  completion   generate the autocompletion script for the specified shell
  help         Help about any command
  info         Displays your server infos
  pull         print the dog stats
  push         push dogs example like all-in-one example, but with ParseKSQL
  serverhealth display the server state of your servers
  setup        setup a dummy connector like in all-in-one example

Flags:
      --config string      config file (default is $HOME/.cobra-test.yaml)
  -h, --help               help for cobra-test
      --host string        set the ksqldb host (default "http://localhost:8088")
      --logformat string   set log format [text|json] (default "text")
      --loglevel string    set log level [info|debug|error|trace] (default "debug")
      --password string    set the ksqldb user password
      --username string    set the ksqldb user name

Use "cobra-test [command] --help" for more information about a command.
```

The `cobra-test setup` command sets up all needed stuff for `cobra-test pull|push` commands.

So run it first.

### KSql Grammar example

This example was written to test and fix the `Antlr4` generation problems for Golang. We changed the `Antlr4` file because there are some type issues. The `Antlr4` code generation introduced some bugs that we had to fix manually. So be careful when you use our `Makefile` to generate the `KSqlParser`. It will break the code!

We had copied the `Antlr4` file from the original sources of [confluent](https://github.com/confluentinc/ksql/blob/master/ksqldb-parser/src/main/antlr4/io/confluent/ksql/parser/SqlBase.g4).
It seems that some errors are not found by the parser because the terminal symbols are not present in the grammar.

The parser is used to check the `KSql syntax`. If there are syntax errors, we collect the errors and you get a notification about it.

The example has an error in the `Select` statement to output the errors.

Feel free to play around :)

## Docker compose

It contains the lateste versions of all products.

- zookeeper (6.2.1)
- schema-registry (6.2.1)
- ksqldb server (0.21.0)
- kafka-connect (6.2.1)
- ksqldb-cli (0.21.0)
- kafdrop (latest)

### ksqldb-cli

For testing purposes I've added `ksqldb-cli` to the `docker-compose.yml` file.

```bash
docker exec -it ksqldb-cli ksql http://ksqldb:8088
```

This starts the interctive ksqldb console.

```
OpenJDK 64-Bit Server VM warning: Option UseConcMarkSweepGC was deprecated in version 9.0 and will likely be removed in a future release.

                  ===========================================
                  =       _              _ ____  ____       =
                  =      | | _____  __ _| |  _ \| __ )      =
                  =      | |/ / __|/ _` | | | | |  _ \      =
                  =      |   <\__ \ (_| | | |_| | |_) |     =
                  =      |_|\_\___/\__, |_|____/|____/      =
                  =                   |_|                   =
                  =        The Database purpose-built       =
                  =        for stream processing apps       =
                  ===========================================

Copyright 2017-2021 Confluent Inc.

CLI v0.21.0, Server v0.21.0 located at http://ksqldb:8088
Server Status: RUNNING

Having trouble? Type 'help' (case-insensitive) for a rundown of how things work!

ksql>
```

## [Kafdrop](https://github.com/obsidiandynamics/kafdrop)

Kafdrop is a web UI for viewing Kafka topics and browsing consumer groups. The tool displays information such as brokers, topics, partitions, consumers, and lets you view messages.

![](https://raw.githubusercontent.com/obsidiandynamics/kafdrop/master/docs/images/overview.png)

## TODO

See https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-clients/contributing/

## License

[Apache License Version 2.0](LICENSE)
