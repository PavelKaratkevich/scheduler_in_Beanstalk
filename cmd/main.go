package main

import (
	"os"
	"os/signal"
	"scheduler/internal/scheduler"
	"scheduler/internal/util"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// connect to the database
	db := util.ConnectDB()
	defer db.Close()

	// Start the scheduler
	newScheduler := scheduler.StartNewScheduler()

	/* Register tasks in the following format:
		metric name,
		start time for metric collection,
		interval for execution of a task,
		db instance,
		scheduler instance,
		sql request
	*/

	go scheduler.RegisterTask(
		"The total number of registrations per day",
		time.Now(),
		time.Duration(3*time.Second),
		db,
		newScheduler,
		"select count(*) from registration where timestamp = '2022-12-22'",
	)

	go scheduler.RegisterTask(
		"The average weight for a registration for the most recent week",
		time.Now(),
		time.Duration(10*time.Second),
		db,
		newScheduler,
		"select cast(ROUND(avg(weight)) as int) from registration where timestamp BETWEEN '2022-12-16' and '2022-12-22'",
	)

	// Block until Ctrl+C
	c := make(chan os.Signal)
	//nolint:govet,staticcheck
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
