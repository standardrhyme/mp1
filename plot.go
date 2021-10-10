package main

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

//Plot is a function that plots the results
func Plot(mode string) {
	var title string
	var keys []int
	var values []opts.ScatterData

	for keyValue := 1; keyValue < len(desirednodesresults); keyValue++ {
		keys = append(keys, keyValue)
		values = append(values, opts.ScatterData{Value: desirednodesresults[keyValue]})
	}

	if mode == "PSH" {
		title = "Nodes vs. Convergence Time - Push Based Gossip"
	} else if mode == "PLL" {
		title = "Nodes vs. Convergence Time - Pull Based Gossip"
	} else if mode == "PPO" {
		title = "Nodes vs. Convergence Time - Push and Pull Original Based Gossip"
	} else if mode == "PPS" {
		title = "Nodes vs. Convergence Time - Push and Pull Switch Based Gossip"
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: title,
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
