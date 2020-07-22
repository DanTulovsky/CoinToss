package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
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
	peoplePlaying   int = 1000 // must be at least 2 for bar chart

	chartFile string = "/tmp/chart.png"
)

// Play the following game.  Every day you can toss a single coin <maxTossesPerDay>. You can stop at any time. You must stop if you reach the limit.
// Each <heads> gives you +$1, each <tails> gives you -$1.
// The goal is to maximize your sum after <daysToPlay>.

func main() {
	fmt.Println("Running...")
	rand.Seed(time.Now().UnixNano())

	// list of earnings for each person
	finalEarnings := make([]float64, peoplePlaying)

	// stop once the daily winnings are at this number
	// play1(100)

	// stop after this many tosses
	// play2(maxTossesPerDay)

	for i := 0; i < peoplePlaying; i++ {
		finalEarning := play1(10)
		finalEarnings[i] = finalEarning
	}

	renderAllPeopleScatterGraph(finalEarnings)
	// renderAllPeopleBarGraph(finalEarnings)
}

// play1 plays the game with the following algorithm:
// Keep tossing the coin until you are at +<stopat>, then stop.
func play1(stopat int) float64 {

	totalEarnings := 0.0
	dailyRunningTotal := make([]float64, 0)

	for i := 1; i < daysToPlay+1; i++ {
		dailyTotal := 0

		for j := 0; j < maxTossesPerDay; j++ {
			toss := tossCoin()

			switch toss {
			case heads:
				dailyTotal++
				totalEarnings++
			case tails:
				dailyTotal--
				totalEarnings--
			}

			if dailyTotal >= stopat {
				break
			}
		}

		dailyRunningTotal = append(dailyRunningTotal, totalEarnings)
	}

	// total earnings after all days
	// renderSinglePersonGraph(earningsPerDay)
	return totalEarnings
}

// play2 plays the game with the following algorithm:
// Make <tosses> tosses every day
func play2(tosses int) float64 {

	totalEarnings := 0.0
	earningsPerDay := make([]float64, 0)

	for i := 1; i < daysToPlay+1; i++ {
		dailyTotal := 0

		for j := 0; j < tosses; j++ {
			toss := tossCoin()

			switch toss {
			case heads:
				dailyTotal++
				totalEarnings++
			case tails:
				dailyTotal--
				totalEarnings--
			default:
				panic("invalid value for coin toss")
			}
		}

		earningsPerDay = append(earningsPerDay, totalEarnings)
		// renderSinglePersonGraph(EarningsPerDay)
	}

	// total earnings after all days
	return totalEarnings
}

func renderAllPeopleScatterGraph(finalEarnings []float64) {
	// sort.Float64s(finalEarnings)

	numWinners := 0
	numLosers := 0
	numEven := 0

	for i := 0; i < len(finalEarnings); i++ {
		switch {
		case finalEarnings[i] > 0:
			numWinners++
		case finalEarnings[i] < 0:
			numLosers++
		default:
			numEven++
		}
	}

	fmt.Printf("Winners: %v; Losers: %v; BrokeEven: %v\n", numWinners, numLosers, numEven)

	bars := make([]chart.Value, len(finalEarnings))

	for i := 0; i < len(finalEarnings); i++ {
		bars[i] = chart.Value{Value: finalEarnings[i], Label: fmt.Sprint(i)}
	}
	viridisByY := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
		return chart.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	graph := chart.Chart{
		Title: fmt.Sprint("Eanings per Person after all days"),
		XAxis: chart.XAxis{
			Name: "People",
		},
		YAxis: chart.YAxis{
			Name: "Earnings",
			GridLines: []chart.GridLine{
				{
					Value:   0,
					IsMinor: false,
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
		Height: 1024,
		Width:  2048,
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeWidth:      chart.Disabled,
					DotWidth:         3,
					DotColorProvider: viridisByY,
				},
				XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1.0).WithEnd(float64(len(finalEarnings)))}.Values(),
				YValues: finalEarnings,
			},
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

func renderAllPeopleBarGraph(finalEarnings []float64) {
	sort.Float64s(finalEarnings)

	numWinners := 0
	numLosers := 0
	numEven := 0

	for i := 0; i < len(finalEarnings); i++ {
		switch {
		case finalEarnings[i] < 0:
			numWinners++
		case finalEarnings[i] > 0:
			numLosers++
		default:
			numEven++
		}
	}

	fmt.Printf("Winners: %v; Losers: %v; BrokeEven: %v\n", numWinners, numLosers, numEven)

	bars := make([]chart.Value, len(finalEarnings))

	for i := 0; i < len(finalEarnings); i++ {
		bars[i] = chart.Value{Value: finalEarnings[i], Label: fmt.Sprint(i)}
	}

	graph := chart.BarChart{
		Title: fmt.Sprint("Eanings per Person after all days"),
		YAxis: chart.YAxis{
			Name: "Earnings",
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 20,
			},
		},
		Height: 1024,
		// Width:        2048,
		BarWidth:     1,
		UseBaseValue: true,
		BaseValue:    0.0,
		Bars:         bars,
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

func renderSinglePersonGraph(earningsPerDay []float64) {

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
