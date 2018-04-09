package ga

import (
	"fmt"
	"os"

	"github.com/wcharczuk/go-chart"
)

func DrawGraph(name string, MaxFitnesses, AverageFitnesses []float64) {
	generations := make([]float64, 0)
	for i := 0; i < len(MaxFitnesses); i++ {
		generations = append(generations, float64(i+1))
	}

	graph := chart.Chart{
		Title:      name,
		TitleStyle: chart.Style{Show: true},
		XAxis: chart.XAxis{
			Name: "Generations",
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
			TickPosition: chart.TickPositionUnset,
			ValueFormatter: func(v interface{}) string {
				typed := int(v.(float64))
				return fmt.Sprintf("%v", typed)
			},
		},
		YAxis: chart.YAxis{
			Name: "Fitness",
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				typed := int(v.(float64))
				return fmt.Sprintf("%v", typed)
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Max Fitness",
				XValues: generations,
				YValues: MaxFitnesses,
			},
			chart.ContinuousSeries{
				Name:    "Average Fitness",
				XValues: generations,
				YValues: AverageFitnesses,
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.LegendThin(&graph),
	}

	file, err := os.OpenFile("graphs/"+name+".png", os.O_CREATE, 777)
	Check(err)
	err = graph.Render(chart.PNG, file)
	Check(err)
}
