package main

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Plot() {
	var keys []int
	var values []opts.ScatterData

	for keyValue := 1; keyValue < len(desirednodesresults); keyValue++ {
		keys = append(keys, keyValue)
		values = append(values, opts.ScatterData{Value: desirednodesresults[keyValue]})
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Nodes vs. Convergence Time",
			},
		),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "# of Nodes",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Conv. Time",
		}),
	)

	// Put data into instance
	scatter.SetXAxis(keys)
	scatter.AddSeries("Category A", values)
	scatter.SetSeriesOptions(charts.WithLabelOpts(
		opts.Label{
			Show:     true,
			Position: "right",
		}),
	)
	f, _ := os.Create("nodesvsconvergencetime.html")
	err := scatter.Render(f)
	if err != nil {
		return
	}

}
