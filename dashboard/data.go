package dashboard

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/MugTree/ryan_dashboard/shared"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

const QueryParamsError = "query params error"
const BadDataError = "bad data error"
const SqlError = "sql error"
const JsonError = "json error"

const NoPageToEdit int64 = 0

const AssetsPathDev = "./www/public/"
const AssetsPathProd = "/"
const DateLayout string = "2006-01-02 15:04:05"

type BaseViewModel struct {
	Request *http.Request
	IsProd  bool
}

type HomepageViewModel struct {
	BaseViewModel
	Data []shared.SensorData
}

func getSensorData(webAddress string) ([]shared.SensorData, error) {
	data := []shared.SensorData{}
	if err := shared.CallWebsiteAPI("GET", webAddress, "", nil, &data); err != nil {

		return data, err
	}

	return data, nil
}

func getChartParts(data []shared.SensorData) (chartElement string, chartScript string) {

	times := []string{}
	depths := make([]opts.LineData, 0)

	for _, v := range data {

		fmt.Println(v.Depth, v.Date.Format("3:04 PM"))
		times = append(times, v.Date.Format("3:04 PM"))
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
	return chart.Element, chart.Script + chart.Option

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
