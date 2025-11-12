package dashboard

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

var sensorApiErrStr = "error calling sensor api: %v"

func webRoutes(r chi.Router, _ *sqlx.DB, env *EnvVars) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(SensorApiError, r, fmt.Errorf(sensorApiErrStr, err)))
			return
		}

		chartId := "depths_chart"
		chart, script, _ := getLineChartParts(data, chartId)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		HomePage(
			"Homepage",
			r,
			env.IsProd,
			chart,
			script,
			"#"+chartId,
			"/api/charts/line",
		).Render(r.Context(), w)
	})

	r.Get("/api/charts/line", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(SensorApiError, r, fmt.Errorf(sensorApiErrStr, err)))
			return
		}

		chartId := "depths_chart"
		_, _, option := getLineChartParts(data, chartId)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(option))
	})
}
