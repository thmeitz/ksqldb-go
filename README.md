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

This client is `not production ready` and the interfaces can be changed without notification!!!

⚠️ Disclaimer #1: This is a personal project and not supported or endorsed by Confluent.

## Migration?

Checkout [ksqldb-migrate](https://github.com/thmeitz/ksqldb-migrate), a tool to run your ksqlDB migrations.

## Description

This is a Go client for [ksqlDB](https://ksqldb.io/).

- [x] Execute a statement (/ksql endpoint)
- [x] Run push and pull queries (/query-stream endpoint)
- [x] Close push query (/close-query endpoint)
- [x] Terminate a cluster (/ksql/terminate endpoint)
- [x] Introspect query status (/status endpoint)
- [x] Introspect server status (/info endpoint)
- [x] Introspect cluster status (/clusterStatus endpoint)
- [x] Get the validity of a property (/is_valid_property)

> Deprecation:
> the [Run a query](https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/query-endpoint/) endpoint is deprecated and willl not be implemented.

### KSqlParser

- the lexer works like the Confluent Java lexer case insensitive (ex `SELECT * FROM BLA` is identical to `select * from bla`). (since v0.0.4)
- parse your ksql-statements with the provided `parser.ParseSql` method.
- `Push`, `Pull`, `Execute` queries parsed by default with `parser.ParseSQL`.
- `<client-instance>.EnableParseSQL(false)` enables / disables the parser

#### Supported ksqlDB versions

It seems that mdrogalis-voluble is no longer provided on confluent-hub.

For that reason I built it locally and put it in the docker compose volume for the Kafka Connectors.

- tested with ksqldb: v0.28.2, v0.22.0, v0.21.0

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

> #### Breaking Change v0.0.4
>
> The HTTP client has now its own package

```golang
import (
  "github.com/Masterminds/log-go"
  "github.com/Masterminds/log-go/impl/logrus"
  "github.com/thmeitz/ksqldb-go"
  "github.com/thmeitz/ksqldb-go/net"
)

var (
	logger = logrus.NewStandard()
)
// than later in your code...
func main {
  options := net.Options{
    // if you need a login, do this; if not its not necessary
    Credentials: net.Credentials{Username: "myuser", Password: "mypassword"},
    // defaults to http://localhost:8088
    BaseUrl:     "http://my-super-shiny-ksqldbserver:8082",
    // this is needed, because the ksql api communicates with http2 only
    // default value in v0.0.4
    AllowHTTP:   true,
  }

  // only log.Logger is allowed or nil (since v0.0.4)
  // logrus is in maintenance mode, so I'll using zap in the future
  client, err := net.NewHTTPClient(options, nil)
  if err != nil {
     log.Fatal(err)
  }
  defer client.Close()

  // then make a pull, push, execute request
}
```

For no authentication remove `Credentials` from options.

### QueryBuilder

> #### Breaking Change v0.0.4
>
> The QueryBuilder was to complicated, so I've refactored it

SQL strings should be build by a QueryBuilder. Otherwise the system is open for SQL injections (see [go-webapp-scp.pdf](https://github.com/OWASP/Go-SCP/blob/master/dist/go-webapp-scp.pdf) ).

You can add multiple parameters `QueryBuilder("insert into bla values(?,?,?,?,?)", nil, 1, 2.5686, "string", true)`.

`nil` will be converted to `NULL`.

The number of parameters must match the parameters in the SQL statement. If not, an error is thrown.

```golang
//see file: examples/cobra-test/cmd/pull.go

k := `SELECT TIMESTAMPTOSTRING(WINDOWSTART,'yyyy-MM-dd HH:mm:ss','Europe/London') AS WINDOW_START,
TIMESTAMPTOSTRING(WINDOWEND,'HH:mm:ss','Europe/London') AS WINDOW_END,
DOG_SIZE, DOGS_CT FROM DOGS_BY_SIZE
WHERE DOG_SIZE=?;`

stmnt, err := ksqldb.QueryBuilder(k, "middle")
if err != nil {
	log.Fatal(err)
}

fmt.Println(*stmnt)
```

### Pull query

```golang

options := net.Options{
	Credentials: net.Credentials{Username: "user", Password: "password"},
	BaseUrl:     "http://localhost:8088",
	AllowHTTP:   true,
}

kcl, err := ksqldb.NewClientWithOptions(options)
if err != nil {
	log.Fatal(err)
}
defer kcl.Close()

query := `select timestamptostring(windowstart,'yyyy-MM-dd HH:mm:ss','Europe/London') as window_start,
timestamptostring(windowend,'HH:mm:ss','Europe/London') as window_end, dog_size, dogs_ct
from dogs_by_size where dog_size=?;`

stmnt, err := ksqldb.QueryBuilder(query, dogsize)
if err != nil {
	log.Fatal(err)
}

ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
defer cancel()

qOpts := (&ksqldb.QueryOptions{Sql: *stmnt}).EnablePullQueryTableScan(false)

_, r, err := kcl.Pull(ctx, *qOpts)
if err != nil {
	log.Fatal(err)
}

var windowStart string
var windowEnd string
var dogSize string
var dogsCt float64
for _, row := range r {

	if row != nil {
		// Should do some type assertions here
		windowStart = row[0].(string)
		windowEnd = row[1].(string)
		dogSize = row[2].(string)
		dogsCt = row[3].(float64)
		log.Infof("🐶 There are %v dogs size %v between %v and %v", dogsCt, dogSize, windowStart, windowEnd)
	}
}
```

### Push query

```golang
options := net.Options{
	Credentials: net.Credentials{Username: "user", Password: "password"},
	BaseUrl:     "http://localhost:8088",
	AllowHTTP:   true,
}

kcl, err := ksqldb.NewClientWithOptions(options)
if err != nil {
	log.Fatal(err)
}
defer kcl.Close()

// you can disable parsing with `kcl.EnableParseSQL(false)`
query := "select rowtime, id, name, dogsize, age from dogs emit changes;"

rowChannel := make(chan ksqldb.Row)
headerChannel := make(chan ksqldb.Header, 1)

// This Go routine will handle rows as and when they
// are sent to the channel
go func() {
	var dataTs float64
	var id string
	var name string
	var dogSize string
	var age string
	for row := range rowChannel {
		if row != nil {
			// Should do some type assertions here
			dataTs = row[0].(float64)
			id = row[1].(string)
			name = row[2].(string)
			dogSize = row[3].(string)
			age = row[4].(string)

			// Handle the timestamp
			t := int64(dataTs)
			ts := time.Unix(t/1000, 0).Format(time.RFC822)

			log.Infof("🐾 New dog at %v: '%v' is %v and %v (id %v)\n", ts, name, dogSize, age, id)
		}
	}
}()

ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
defer cancel()

e := kcl.Push(ctx, ksqldb.QueryOptions{Sql: query}, rowChannel, headerChannel)
if e != nil {
	log.Fatal(e)
}
```

### Execute a command

```golang
  options := net.Options{
		Credentials: net.Credentials{Username: "user", Password: "password"},
		BaseUrl:     "http://localhost:8088",
		AllowHTTP:   true,
	}

	kcl, err := ksqldb.NewClientWithOptions(options)
	if err != nil {
		log.Fatal(err)
	}
	defer kcl.Close()

	resp, err := kcl.Execute(ksqldb.ExecOptions{KSql: `
		CREATE SOURCE CONNECTOR DOGS WITH (
		'connector.class'                = 'io.mdrogalis.voluble.VolubleSourceConnector',
		'key.converter'                  = 'org.apache.kafka.connect.storage.StringConverter',
		'value.converter'                = 'org.apache.kafka.connect.json.JsonConverter',
		'value.converter.schemas.enable' = 'false',
		'genkp.dogs.with'                = '#{Internet.uuid}',
		'genv.dogs.name.with'            = '#{Dog.name}',
		'genv.dogs.dogsize.with'         = '#{Dog.size}',
		'genv.dogs.age.with'             = '#{Dog.age}',
		'topic.dogs.throttle.ms'         = 1000
		);
		`})
	if err != nil {
		log.Fatalf("create source connector dogs failed %w", err)
		os.Exit(-1)
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
  check             check a <example>.ksql file with the integrated parser
  cluster-status    get cluster status
  completion        generate the autocompletion script for the specified shell
  health            display the server state of your servers
  help              Help about any command
  info              Displays your server infos
  pull              print the dog stats
  push              push dogs example
  setup             setup a dummy connector
  terminate-cluster terminates your cluster
  validate          validates a property

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

This example was written to test and fix the `Antlr4` generation problems for Golang. We changed the `Antlr4` file because there are some type issues (type is a reserved word in golang). The `Antlr4` code generation introduced some bugs that we had to fix manually (no Antlr4 output for needed package names). So be careful when you use our `Makefile` to generate the `KSqlParser`. It will break the code!

We had copied the `Antlr4` file from the original sources of [confluent](https://github.com/confluentinc/ksql/blob/master/ksqldb-parser/src/main/antlr4/io/confluent/ksql/parser/SqlBase.g4).

The parser is used to check the `KSql syntax`. If there are syntax errors, the errors will be collected and you get a notification about it.

## Docker compose

It contains the latest versions of all products.

- zookeeper (6.2.1)
- schema-registry (6.2.1)
- ksqldb server (0.21.0)
- kafka-connect (6.2.1)
- ksqldb-cli (0.21.0)
- kafdrop (latest)

### ksqldb

I've added following options to `docker-compose` to get the `ClusterStatus`.

```yaml
KSQL_OPTS: "-Dksql.heartbeat.enable=true -Dksql.lag.reporting.enable=true"
```

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

Kafdrop runs on port 9000 on your localhost.

```
http://localhost:9000
```

![](https://raw.githubusercontent.com/obsidiandynamics/kafdrop/master/docs/images/overview.png)

## TODO

See https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-clients/contributing/

## License

[Apache License Version 2.0](LICENSE)
