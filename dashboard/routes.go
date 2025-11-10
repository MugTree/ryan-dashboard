package dashboard

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/MugTree/ryan_dashboard/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/jmoiron/sqlx"
	"github.com/starfederation/datastar-go/datastar"
)

func webRoutes(r chi.Router, _ *sqlx.DB, env *EnvVars, sensorData []shared.SensorData) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		startPoint := 0
		endPoint := 10

		fullChart := getChartData(sensorData, startPoint, endPoint)
		index := startPoint + 1

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		HomePage("Homepage", r, env.IsProd, fullChart, index).Render(r.Context(), w)
	})

	r.Patch("/api/charts/line/{index}", func(w http.ResponseWriter, r *http.Request) {

		index, ok := paramMustBeNumeric(w, r, "index")
		if !ok {
			return
		}

		startPoint := index * 10
		endPoint := startPoint + 10

		fmt.Println("-------------------------------")
		fmt.Printf("startPoint: %v", startPoint)
		fmt.Printf("end: %v", endPoint)
		fmt.Println("-------------------------------")

		chart := getChartData(sensorData, startPoint, endPoint)

		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(LineGraph(chart, index+1))

	})
}

// TODO needs to return two values (chartElement, chartScript)
func getChartData(data []shared.SensorData, start int, end int) string {

	times := []string{}
	depths := make([]opts.LineData, 0)

	for _, v := range data[start:end] {

		fmt.Println(v.Depth, v.DataSensed.Format("3:04 PM"))
		times = append(times, v.DataSensed.Format("3:04 PM"))
		depths = append(depths, opts.LineData{Value: v.Depth})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line title",
			Subtitle: "Sub",
		}))

	chart := line.SetXAxis(times).AddSeries("depths", depths).RenderSnippet()
	return chart.Element + chart.Script + chart.Option

}

// Create some dummy sensor data - give some sort of graph idea
func getSensorData() []shared.SensorData {
	now := time.Now().UTC()

	data := []shared.SensorData{}
	var depth int

	for i := 480; i > 0; i-- {
		t := now.Add(-time.Duration(i) * time.Minute)

		r := rand.IntN(6)

		if r%3 == 0 || depth == 0 {
			depth = depth + 1
		} else {
			depth = depth - 1
		}

		data = append(data, shared.SensorData{Depth: depth, DataSensed: t})
	}

	return data
}
