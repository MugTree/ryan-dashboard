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

		ldComp, err := getLineGraphLiveDataComponent(data, "depths_chart")

		if err != nil {
			logAndError(w, formatError(BadDataError, r, err))
		}

		HomePage(r, env.IsProd, ldComp).Render(r.Context(), w)
	})

	r.Get("/api/charts/line", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(SensorApiError, r, fmt.Errorf(sensorApiErrStr, err)))
			return
		}

		chart := getLiveDepthsChartSnippet(data, "depths_chart")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(chart.Option))
	})
}
