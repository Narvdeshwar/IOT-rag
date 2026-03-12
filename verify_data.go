package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	url := os.Getenv("POSTGRES_URL")
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM sensor_events 
			WHERE device_id = '1234' 
			AND value = 32.5
		)`).Scan(&exists)
	
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		fmt.Println("Record FOUND in database. This is REAL data from your sensor events.")
	} else {
		fmt.Println("Record NOT FOUND. This might be a hallucination or dummy data from the model.")
		
		// Let's see what IS there
		fmt.Println("\nSample records from database:")
		rows, _ := pool.Query(context.Background(), "SELECT device_id, event_time, metric, value FROM sensor_events LIMIT 5")
		defer rows.Close()
		for rows.Next() {
			var d, m string
			var t interface{}
			var v float64
			rows.Scan(&d, &t, &m, &v)
			fmt.Printf("Device: %s | Time: %v | Metric: %s | Value: %f\n", d, t, m, v)
		}
	}
}
