// @rmoff
//
package main

import (
	"fmt"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/logrus"
)

var (
	logger = logrus.NewStandard()
)

const ksqlDBServer string = "http://localhost:8088"
const ksqlDBUser string = ""
const ksqlDBPW string = ""

func main() {

	client, err := setup()
	if err != nil {
		logger.Fatalw("Failed to run setup statements. Exiting.", log.Fields{"error": err})
	}
	// Do a pull query
	fmt.Printf("\n\n" + `
	  	  🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀
		✨It's… a Golang client for ksqlDB! ✨
		  🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀
	
	Check this out, we can do pull queries, which are like K/V lookups
	against materialised views of state built from streams of events in Kafka:` + "\n\n")
	if e := getDogStats(client, "medium"); e != nil {
		logger.Errorw("error calling getDogStats", log.Fields{"error": err})
	}

	time.Sleep(3 * time.Second)
	// Do a push query
	fmt.Printf("\n\n" + `
		                      ❇️ ❇️ ❇️ ❇️ ❇️ ❇️
	
	➜ We can also do push queries, in which we subscribe to a stream of
	notifications of events. This could be every event arriving on a topic,
	or it could be events that match a given condition specified in a WHERE
	clause. Note that this is a continuous query. Here we use the cancel option
	to terminate it after 10 seconds, but by default it will run until the program
	is killed.` + "\n\n\n")
	time.Sleep(2 * time.Second)
	if e := getDogUpdates(client); e != nil {
		logger.Errorw("error calling getDogUpdates", log.Fields{"error": e})
	}
}
