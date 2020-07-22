package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

const (
	heads = iota
	tails
)

var (
	maxTossesPerDay int = 400
	daysToPlay      int = 768

	currentEarnings float64 = 0.0

	earningsPerDay []float64

	chartFile string = "/tmp/chart.png"
)

// Play the following game.  Every day you can toss a single coin <maxTossesPerDay>. You can stop at any time. You must stop if you reach the limit.
// Each <heads> gives you +$1, each <tails> gives you -$1.
// The goal is to maximize your sum after <daysToPlay>.

func main() {
	fmt.Println("Running...")
	rand.Seed(time.Now().UnixNano())

	earningsPerDay = make([]float64, 0)

	play1()
	// play2()
	renderGraph()
}

// play1 plays the game with the following algorithm:
// Keep tossing the coin until you are at +1, then stop.
func play1() {

	for i := 1; i < daysToPlay+1; i++ {
		dailyTotal := 0

		for j := 0; j < maxTossesPerDay; j++ {
			toss := tossCoin()

			switch toss {
			case heads:
				dailyTotal++
				currentEarnings++
			case tails:
				dailyTotal--
				currentEarnings--
			}

			if dailyTotal > 0 {
				break
			}
		}

		earningsPerDay = append(earningsPerDay, currentEarnings)
	}
}

// play2 plays the game with the following algorithm:
// Make all possible tosses every day
func play2() {

	for i := 1; i < daysToPlay+1; i++ {
		dailyTotal := 0

		for j := 0; j < maxTossesPerDay; j++ {
			toss := tossCoin()

			switch toss {
			case heads:
				dailyTotal++
				currentEarnings++
			case tails:
				dailyTotal--
				currentEarnings--
			}
		}

		earningsPerDay = append(earningsPerDay, currentEarnings)
	}
}

func renderGraph() {

	mainSeries := chart.ContinuousSeries{
		Name: "Coin Tosses",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			StrokeWidth: 4,
		},
		XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1.0).WithEnd(float64(daysToPlay))}.Values(),
		YValues: earningsPerDay,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Day",
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: float64(daysToPlay),
			},
		},
		YAxis: chart.YAxis{
			Name: "Earnings",
			GridLines: []chart.GridLine{
				{
					Value:   0,
					IsMinor: true,
				},
			},
			GridMajorStyle: chart.Style{
				Hidden:      false,
				StrokeColor: drawing.ColorBlack,
				StrokeWidth: 1.5,
			},
			GridMinorStyle: chart.Style{
				Hidden:      false,
				StrokeColor: drawing.Color{R: 0, G: 0, B: 0, A: 100},
				StrokeWidth: 1.0,
			},
		},
		Series: []chart.Series{
			mainSeries,
			// minSeries,
			// maxSeries,
		},
	}

	pngFile, err := os.Create(chartFile)
	if err != nil {
		panic(err)
	}

	if err := graph.Render(chart.PNG, pngFile); err != nil {
		panic(err)
	}

	if err := pngFile.Close(); err != nil {
		panic(err)
	}
}

// tossCoin simulates one coin toss; 0 = heads; 1 = tails
func tossCoin() int {
	return rand.Intn(2)
}
