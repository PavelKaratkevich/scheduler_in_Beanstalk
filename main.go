package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/madflojo/tasks"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DefaultDatabaseName = "db.sqlite"
)

var (
	DataDir = getEnvDefault("DATA_DIR", "../data")
)

func main() {

	// connect to the database
	db := connectDB("sqlite3", filepath.Join(DataDir, DefaultDatabaseName))
	defer db.Close()

	// Start the scheduler
	scheduler := tasks.New()
	defer scheduler.Stop()

	// Register tasks in the following format: metric name, start date, interval, db instance, scheduler instance, sql request
	go registerTask("The total number of registrations per day", time.Now(), time.Duration(3*time.Second), db, scheduler, "select count(*) from registration where timestamp = '2022-12-22'")
	go registerTask("The average weight for a registration for the most recent week", time.Now(), time.Duration(10*time.Second), db, scheduler, "select cast(ROUND(avg(weight)) as int) from registration where timestamp BETWEEN '2022-12-16' and '2022-12-22'")

	// Block until Ctrl+C
	c := make(chan os.Signal)
	//nolint:govet,staticcheck
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func connectDB(driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath.Join(DataDir, DefaultDatabaseName))

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// getEnvDefault return the value of the environment variable specified by name, or the defaultValue if not set
func getEnvDefault(name string, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultValue
}


func registerTask(metricName string, startDate time.Time, interval time.Duration, db *sql.DB, scheduler *tasks.Scheduler, query string) {
	var metricValue interface{}

	_, err := scheduler.Add(&tasks.Task{

		Interval:   interval,
		StartAfter: startDate,
		TaskFunc: func() error {
			err := db.QueryRow(query).Scan(&metricValue)
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Printf("%s: %d\n", metricName, metricValue)
			value := map[string]interface{}{metricName: metricValue}
			jsonValue, _ := json.Marshal(value)
			http.Post("http://localhost:5000/result", "application/json", bytes.NewBuffer(jsonValue))
			return nil
		},
	})
	if err != nil {
		log.Fatalln("Fatal error")
	}
}
