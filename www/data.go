package www

import (
	"github.com/MugTree/ryan_dashboard/shared"
	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/goforj/godump"
)

const QueryParamsError = "query params error"
const BadDataError = "bad data error"
const SqlError = "sql error"
const JsonError = "json error"
const SensorApiError = "sensor api error"

const NoPageToEdit int64 = 0

const AssetsPathDev = "./dashboard/public/"
const AssetsPathProd = "/"
const DateLayout string = "2006-01-02 15:04:05"

func getSystemMemoryChartData(data []shared.MemorySample, chartId string) render.ChartSnippet {
	if len(data) == 0 {
		return render.ChartSnippet{}
	}

	//
	ref := data[len(data)-1].Time
	points := make([]opts.LineData, 0, len(data))

	for _, v := range data {
		x := -ref.Sub(v.Time).Seconds() // seconds ago (negative = older)
		points = append(points, opts.LineData{
			Value: [2]any{x, v.MemoryPercent}, // [x,y] pair for numeric X
		})
	}

	godump.Dump(points)

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithAnimation(false),
		charts.WithInitializationOpts(opts.Initialization{
			Width:   "900px",
			Height:  "600px",
			Theme:   types.ThemeChalk,
			ChartID: chartId,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "value",
			Name: "Seconds ago",
			Min:  -10,
			Max:  0,
			AxisLabel: &opts.AxisLabel{
				Formatter: "{value}s",
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:      "Memory usage (%)",
			AxisLabel: &opts.AxisLabel{Rotate: 90},
			Max:       100,
			Min:       0,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "System Memory Usage (last few seconds)",
			Subtitle: "from some remote system",
		}),
	)

	// Note: no SetXAxis(...) call — X values come from [2] pairs above.
	line.AddSeries("Memory %", points).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth:     opts.Bool(false), // curved line
				ShowSymbol: opts.Bool(false), // no point markers
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: opts.Float(0.4),
				Color:   "rgba(0, 255, 128, 0.5)", // semi-transparent green
			}),
		)

	return line.RenderSnippet()
}

func getSystemMemoryComponent(data []shared.MemorySample, id string) (templ.Component, error) {
	chart := getSystemMemoryChartData(data, id)
	return SystemMemory(chart.Element, chart.Script, "#"+id, "/api/charts/line"), nil
}

func getSensorData(webAddress string) ([]shared.MemorySample, error) {
	data := []shared.MemorySample{}
	if err := shared.CallJsonAPI("GET", webAddress, "", nil, &data); err != nil {
		return data, err
	}
	return data, nil
}

// func paramMustBeNumeric(w http.ResponseWriter, r *http.Request, key string) (int, bool) {
// 	v, err := strconv.Atoi(r.PathValue(key))
// 	if err != nil {
// 		logAndError(w, formatError(QueryParamsError, r, err))
// 		return 0, false
// 	}

// 	if v == 0 {
// 		logAndError(w, formatError(QueryParamsError, r, fmt.Errorf("key '%v' not numeric - %v", key, v)))
// 		return 0, false
// 	}

// 	return v, true
// }

// func paramMustBeNotEmpty(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
// 	v := r.PathValue(key)
// 	if v == "" {
// 		logAndError(w, formatError(QueryParamsError, r, fmt.Errorf("key '%v' empty string - %v", key, v)))
// 		return "", false
// 	}
// 	return v, true
// }
