package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/MugTree/ryan_dashboard/shared"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
This is is assuming a rasberry pi or something
*/
func main() {

	if err := run(); err != nil {
		slog.Error("app error", "err", err)
		os.Exit(1)
	}
}

func run() error {

	// thought is here that env could come from the cmd line OR a file
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf(".env file present but can't load it %v", err)
		}
	}

	host := shared.MustEnv("SENSOR_HOST")

	rotator := &lumberjack.Logger{
		Filename:   shared.MustEnv("SENSOR_LOG"),
		MaxSize:    50, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}
	defer rotator.Close()

	const (
		updateInterval = 100 * time.Millisecond // 10th of a second
		windowSeconds  = 10
	)
	pointsPerSecond := int(time.Second / updateInterval)
	maxPoints := pointsPerSecond * windowSeconds // 100 Hz × 10 s = 1000 points

	sensor := shared.NewSensor(maxPoints)

	go func() {
		start := time.Now()
		for {
			t := time.Since(start).Seconds()

			//CGPT Smooth oscillation with small noise
			y := 60 + 10*math.Sin(t/2) + (rand.Float64()*4 - 2)

			sensor.AddData(shared.MemorySample{
				MemoryPercent: y,
				Time:          time.Now(),
			})

			time.Sleep(updateInterval)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

		data := sensor.GetData()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
		}
	})

	fmt.Println("Server running at " + host + "/api")
	if err := http.ListenAndServe(host, mux); err != nil {
		return fmt.Errorf("server error: %v", err)
	}

	return nil
}
