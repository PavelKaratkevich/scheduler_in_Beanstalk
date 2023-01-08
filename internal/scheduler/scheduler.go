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

func RegisterTask(metricName string, startDate time.Time, interval time.Duration, db *sql.DB, scheduler *tasks.Scheduler, query string) {
	var metricValue interface{}

	_, err := scheduler.Add(&tasks.Task{
		Interval:   interval,
		StartAfter: startDate,
		TaskFunc: func() error {
			err := db.QueryRow(query).Scan(&metricValue)
			if err != nil {
				log.Fatalln("Error while making SQL request: ", err)
			}

			// store result in a map a
			value := map[string]interface{}{metricName: metricValue}
			jsonValue, _ := json.Marshal(value)

			// send result via HTTP request to POST endpoint
			_, err = http.Post("http://apiforscheduledtasks-env.eba-pr7nv2hg.eu-central-1.elasticbeanstalk.com/result", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err.Error())
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
