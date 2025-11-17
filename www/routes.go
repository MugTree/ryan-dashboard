package www

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var sensorApiErrStr = "error calling sensor api: %v"

func webRoutes(r chi.Router, env *EnvVars) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(SensorApiError, r, fmt.Errorf(sensorApiErrStr, err)))
			return
		}

		chartComponent, err := getSystemMemoryComponent(data, "depths_chart")
		if err != nil {
			logAndError(w, formatError(BadDataError, r, err))
		}

		HomePage(r, env.IsProd, chartComponent).Render(r.Context(), w)
	})

	r.Get("/api/charts/linepart", func(w http.ResponseWriter, r *http.Request) {

		// templ.RenderFragments(r.Context(), w, DummyComp(), "blah")
	})

	r.Get("/api/charts/line", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(SensorApiError, r, fmt.Errorf(sensorApiErrStr, err)))
			return
		}

		chartData := getSystemMemoryChartData(data, "depths_chart")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(chartData.Option))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<!DOCTYPE html><html><head><title>health</title></head><body></body></html>"))
	})
}
