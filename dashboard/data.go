package dashboard

import (
	"github.com/MugTree/ryan_dashboard/shared"
	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/go-echarts/go-echarts/v2/types"
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

func getLiveDepthsChartSnippet(data []shared.SensorData, chartId string) render.ChartSnippet {

	times := []string{}
	depths := make([]opts.LineData, 0)

	for _, v := range data {
		times = append(times, v.Date.Format("15:04 05s"))
		depths = append(depths, opts.LineData{Value: v.Depth})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithAnimation(false),
		charts.WithInitializationOpts(opts.Initialization{Width: "600px", Height: "300px", Theme: types.ChartLine, ChartID: chartId}),
		charts.WithYAxisOpts(opts.YAxis{
			Max: 6,
			Min: 1,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Some data from a sensor",
			Subtitle: "just to illustrate",
		}))

	return line.SetXAxis(times).AddSeries("depths", depths).RenderSnippet()
}

func getLineGraphLiveDataComponent(data []shared.SensorData, id string) (templ.Component, error) {
	chart := getLiveDepthsChartSnippet(data, id)
	return LineGraphLiveData(chart.Element, chart.Script, "#"+id, "/api/charts/line"), nil
}

func getSensorData(webAddress string) ([]shared.SensorData, error) {
	data := []shared.SensorData{}
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
