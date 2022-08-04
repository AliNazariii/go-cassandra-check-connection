package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"cassandra_connection_check/pkg/config"
	"cassandra_connection_check/pkg/database"
	"cassandra_connection_check/pkg/log"
)

func main() {
	conf := config.New("", "cassandra-connection-check")
	logger := log.NewLog(conf.Log.Level)
	database.NewCassandraDB(logger, &conf.Cassandra)

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		fmt.Println("Signal received: ", sig)
		done <- true
	}()
	<-done
}
