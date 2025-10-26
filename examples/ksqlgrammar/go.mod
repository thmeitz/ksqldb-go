module github.com/thmeitz/ksqldb-go/examples/ksqlgrammar

go 1.24.0

require (
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230321174746-8dcc6526cfb1
	github.com/thmeitz/ksqldb-go v0.2.0
)

require golang.org/x/exp v0.0.0-20251023183803-a4bb9ffd2546 // indirect

// Use local ksqldb-go instead of remote
replace github.com/thmeitz/ksqldb-go => ../..
