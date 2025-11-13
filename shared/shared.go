package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type MemorySample struct {
	MemoryPercent float64   `json:"memory_percent"`
	Time          time.Time `json:"time"`
}

type Sensor struct {
	Data    []MemorySample
	Mu      sync.Mutex
	MaxSize int
}

func NewSensor(maxSize int) *Sensor {
	return &Sensor{
		Data:    make([]MemorySample, 0, maxSize),
		MaxSize: maxSize,
	}
}

func (s *Sensor) AddData(d MemorySample) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Data = append(s.Data, d)
	if len(s.Data) > s.MaxSize {
		excess := len(s.Data) - s.MaxSize
		s.Data = s.Data[excess:]
	}
}

func (s *Sensor) GetData() []MemorySample {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	dataCopy := make([]MemorySample, len(s.Data))
	copy(dataCopy, s.Data)
	return dataCopy
}

func CallJsonAPI(action string, endpointUrl string, apiKey string, payload io.Reader, result any) error {

	req, err := http.NewRequest(action, endpointUrl, payload)

	if err != nil {
		return fmt.Errorf("error creating getData request: %v", err)
	}

	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making getData request: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	return err

}

func MustEnv(name string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		slog.Error("Missing required environment variable", "var", name)
		os.Exit(1)
	}
	return v
}

func MustEnvGetBool(name string) bool {

	v := MustEnv(name)

	if v != "true" && v != "false" {
		slog.Error("env requires 'true'  or 'false' lowercase variable name", "var", name)
		os.Exit(1)
	}

	val, err := strconv.ParseBool(v)
	if err != nil {
		slog.Error("env can't convert value to a bool", "var", name)
		os.Exit(1)
	}

	return val
}
