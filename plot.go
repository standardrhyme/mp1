package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

//IsWritable is a function used to check whether the user's current directory can be written to.
func IsWritable() (isWritable bool) {
	ex, errone := os.Executable()
	if errone != nil {
		fmt.Println("Error identifying the path of the current directory.")
		fmt.Println("This is used to write the results file to the current directory. Please adjust settings and try again.")
	}
	exPath := filepath.Dir(ex)

	info, err := os.Stat(exPath)
	if err != nil {
		fmt.Println("Path to current directory doesn't exist.")
		return false
	}

	err = nil
	if !info.IsDir() {
		fmt.Println("Path isn't a directory.")
		return false
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		fmt.Println("Write permission bit is not set on this file for user. Please adjust settings and try again.")
		return false
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(exPath, &stat); err != nil {
		fmt.Println("Unable to get stat.")
		return false
	}

	err = nil
	if uint32(os.Geteuid()) != stat.Uid {
		isWritable = false
		fmt.Println("User doesn't have permission to write to this directory. Please adjust permissions and try again.")
		return false
	}
	return true
}

//Plot is a function that plots the results
func Plot(mode string) {
	var title string
	var keys []int
	var values []opts.ScatterData

	for keyValue := 1; keyValue < len(desiredNodesResults); keyValue++ {
		keys = append(keys, keyValue)
		values = append(values, opts.ScatterData{Value: desiredNodesResults[keyValue]})
	}

	if mode == "1" {
		title = "Nodes vs. Convergence Time - Push Based Gossip"
	} else if mode == "2" {
		title = "Nodes vs. Convergence Time - Pull Based Gossip"
	} else if mode == "3" {
		title = "Nodes vs. Convergence Time - Push and Pull Original Based Gossip"
	} else if mode == "4" {
		title = "Nodes vs. Convergence Time - Push and Pull Switch Based Gossip"
	} else {
		title = "Nodes vs. Convergence Time"
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
			Name: "# of Rounds",
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

	if IsWritable() {
		f, _ := os.Create("nodesvsconvergencetime.html")
		err := scatter.Render(f)
		if err != nil {
			fmt.Println("The results file was not able to be created in the current directory.")
			fmt.Println("Please try adjusting reading and writing permissions.")
		}
		fmt.Println("\nTo see the number of nodes vs number of rounds results, " +
			"open 'nodesvsconvergencetime.html' from the current directory.")
		fmt.Println("(If the file does not exist in the current directory, check directory read and write permissions and try again.)")
	} else {
		fmt.Println("The current directory is not writeable. Please adjust permissions and try again.")
	}

}
