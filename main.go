package main

import (
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/term"
	"log"
	"math/rand"
	"time"

	"github.com/avoropaev/multi-armed-bandit/bandit"

	ui "github.com/gizak/termui/v3"
)

const ELEMENTS_COUNT = 300
const VIEWS_COUNT = 100000

type csvRow struct {
	group1 int
	group2 int
	group3 int
	group4 int
	group5 int
}

func main() {
	_, height, err := term.GetSize(0)
	if err != nil {
		log.Printf("error: %s", err)

		return
	}

	b := bandit.UCB1{}

	err = b.Init(ELEMENTS_COUNT)
	if err != nil {
		log.Printf("error: %s", err)

		return
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	bc := widgets.NewBarChart()
	bc.Data = []float64{}
	bc.Labels = []string{"< 1.90", "1.90 - 1.94", "1.94 - 1.98", "1.98 - 2.02", "> 2.02"}
	bc.Title = "Multi-armed bandit. X-Axis: reward, Y-Axis: Average number of views."
	bc.SetRect(5, 0, 111, height)
	bc.BarWidth = 20
	bc.LabelStyles = []ui.Style{
		ui.NewStyle(ui.ColorRed),
		ui.NewStyle(ui.ColorGreen),
		ui.NewStyle(ui.ColorYellow),
		ui.NewStyle(ui.ColorBlue),
		ui.NewStyle(ui.ColorMagenta),
	}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	bc.TitleStyle = ui.NewStyle(ui.ColorGreen)

	ui.Render(bc)
	uiEvents := ui.PollEvents()

	for i := 0; i < VIEWS_COUNT; i++ {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		default:
		}


		iteration(i, &b, bc)
	}

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}

func iteration(i int, b *bandit.UCB1, bc *widgets.BarChart)  {
	chosenArm := b.SelectArm(0.0)

	rand.Seed(time.Now().UnixNano())
	reward := float64(rand.Intn(5))

	err := b.Update(chosenArm, reward)
	if err != nil {
		log.Printf("error: %s", err)

		return
	}

	if (i + 1) % 1000 == 0 {
		rowViews := csvRow{}
		rowCounts := csvRow{}

		for index, value := range b.Counts {
			if b.Rewards[index] < 1.90 {
				rowViews.group1 += value
				rowCounts.group1++
			} else if b.Rewards[index] >= 1.90 && b.Rewards[index] < 1.94 {
				rowViews.group2 += value
				rowCounts.group2++
			} else if b.Rewards[index] >= 1.94 && b.Rewards[index] < 1.98 {
				rowViews.group3 += value
				rowCounts.group3++
			} else if b.Rewards[index] >= 1.98 && b.Rewards[index] < 2.02 {
				rowViews.group4 += value
				rowCounts.group4++
			} else {
				rowViews.group5 += value
				rowCounts.group5++
			}
		}

		if rowCounts.group1 != 0 {
			rowViews.group1 /= rowCounts.group1
		}

		if rowCounts.group2 != 0 {
			rowViews.group2 /= rowCounts.group2
		}

		if rowCounts.group3 != 0 {
			rowViews.group3 /= rowCounts.group3
		}

		if rowCounts.group4 != 0 {
			rowViews.group4 /= rowCounts.group4
		}

		if rowCounts.group5 != 0 {
			rowViews.group5 /= rowCounts.group5
		}

		bc.Data = []float64{
			float64(rowViews.group1),
			float64(rowViews.group2),
			float64(rowViews.group3),
			float64(rowViews.group4),
			float64(rowViews.group5),
		}

		ui.Render(bc)
		time.Sleep(time.Millisecond * 100)
	}
}
