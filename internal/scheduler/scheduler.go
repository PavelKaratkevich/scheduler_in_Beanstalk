package scheduler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/madflojo/tasks"
)

func StartNewScheduler() *tasks.Scheduler {
	scheduler := tasks.New()
	defer scheduler.Stop()

	return scheduler
}

// ex.: each day at 1PM

func RegisterTask(metricName string, startDate time.Time, interval time.Duration, db *sql.DB, scheduler *tasks.Scheduler, query string) {
	var metricValue interface{}

	_, err := scheduler.Add(&tasks.Task{
		Interval:   interval,
		StartAfter: startDate,
		TaskFunc: func() error {
			err := db.QueryRow(query).Scan(&metricValue)
			if err != nil {
				log.Fatal(err)
			}

			// fmt.Printf("%s: %d\n", metricName, metricValue)

			// store result in a map a
			value := map[string]interface{}{metricName: metricValue}
			jsonValue, _ := json.Marshal(value)

			// send result via HTTP request to POST endpoint
			_, err = http.Post("http://localhost:5000/result", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("metric has been sent: %s: %d\n", metricName, metricValue)
			}

			return nil
		},
	})
	if err != nil {
		log.Fatalln("Fatal error")
	}
}
