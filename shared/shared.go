package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type SensorData struct {
	Depth int       `db:"depth" json:"depth"`
	Date  time.Time `db:"date" json:"date"`
}

type Sensor struct {
	Data    []SensorData
	Mu      sync.Mutex
	MaxSize int
}

func NewSensor(maxSize int) Sensor {
	return Sensor{
		Data:    make([]SensorData, 0),
		Mu:      sync.Mutex{},
		MaxSize: maxSize,
	}
}

func (s *Sensor) AddData(d SensorData) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Data = append(s.Data, d)

	if len(s.Data) > s.MaxSize {
		excess := len(s.Data) - s.MaxSize
		s.Data = s.Data[excess:]
	}
}

func (s *Sensor) GetData() []SensorData {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	dataCopy := make([]SensorData, len(s.Data))
	copy(dataCopy, s.Data)
	return dataCopy
}

func CallWebsiteAPI(action string, endpointUrl string, apiKey string, payload io.Reader, result any) error {

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
