package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scheduler/internal/scheduler"
	"scheduler/internal/util"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DefaultPort = 5000
)

func main() {

	go func() {
		app := fiber.New()

		// Routes
		app.Get("/health", healthCheck)

		// Start server
		listenString := fmt.Sprintf(":%d", DefaultPort)
		log.Fatal(app.Listen(listenString))
	}()

	// connect to the database
	db, err := util.ConnectDB(".")
	if err != nil {
		log.Fatalln("Unable to establish connection with DB: ", err.Error())
	}
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
		time.Duration(5*time.Hour),
		db,
		newScheduler,
		// "select count(*) from registration where timestamp = '2022-12-22'",
		"select count(*) from registration",
	)

	go scheduler.RegisterTask(
		"The average weight for a registration for the most recent week",
		time.Now(),
		time.Duration(5*time.Hour),
		db,
		newScheduler,
		// "select cast(ROUND(avg(weight)) as int) from registration where timestamp BETWEEN '2022-12-16' and '2022-12-22'",
		"select cast(ROUND(avg(weight)) as int) from registration",
	)

	// Block until Ctrl+C
	c := make(chan os.Signal)
	//nolint:govet,staticcheck
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(http.StatusOK)
}
