package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/MugTree/ryan_dashboard/shared"
)

var sensor shared.Sensor

func main() {

	sensor = shared.NewSensor(10)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			sd := shared.SensorData{
				Depth:      rand.IntN(4) + 1,
				DataSensed: time.Now(),
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

		fmt.Printf("items: %v", len(data))
		fmt.Println()

		for _, v := range data {
			fmt.Printf("depth: %v time: %s\n", v.Depth, v.DataSensed.Format("2006-01-02 15:04:05"))
		}

	})

	fmt.Println("starting sensor")

	if err := http.ListenAndServe(":8081", mux); err != nil {
		fmt.Printf("server err: %v", err)
	}

}
