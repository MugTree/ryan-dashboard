package dashboard

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/starfederation/datastar-go/datastar"
)

func webRoutes(r chi.Router, _ *sqlx.DB, env *EnvVars) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(JsonError, r, fmt.Errorf("error calling sensor api: %v", err)))
			return
		}

		chart, script := getChartParts(data)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		HomePage("Homepage", r, env.IsProd, chart, script).Render(r.Context(), w)
	})

	r.Patch("/api/charts/line", func(w http.ResponseWriter, r *http.Request) {

		data, err := getSensorData(env.SensorAddress)
		if err != nil {
			logAndError(w, formatError(JsonError, r, fmt.Errorf("error calling sensor api: %v", err)))
			return
		}

		script, element := getChartParts(data)

		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(LineGraph(script, element))
	})
}
