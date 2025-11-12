package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/MugTree/ryan_dashboard/shared"
)

var sensor shared.Sensor

func main() {

	maxSize := 10
	sensor = shared.NewSensor(maxSize)

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			sd := shared.SensorData{
				Depth: rand.IntN(6) + 1,
				Date:  time.Now(),
			}
			sensor.AddData(sd)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/api" {
			http.Error(w, "404 page not found", 404)
			return
		}

		data := sensor.GetData()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Printf("json encode error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
		}

	})

	fmt.Println("starting sensor")

	if err := http.ListenAndServe(":8081", mux); err != nil {
		fmt.Printf("server err: %v", err)
	}

}
