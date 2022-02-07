package main

import (
	"encoding/csv"
	"os"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/pterm/pterm"
)

func main() {
	_ = pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("S", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("miulation", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("-1", pterm.NewStyle(pterm.FgLightGreen))).
		Render()

	cfg := config.New()

	a := app.New(cfg.App)

	f, err := os.Create("result.csv")
	if err != nil {
		pterm.Fatal.Printf("cannot create result.cvs %s", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for _, g := range cfg.Gateways {
		a.Gateway(g)
	}

	a.Connect()

	for i := 0; i < cfg.Tries; i++ {
		a.Run()

		pterm.Success.Println(a.Durations)

		r := make([]string, 0)
		for _, d := range a.Durations {
			r = append(r, d.String())
		}

		if err := w.Write(r); err != nil {
			pterm.Fatal.Printf("cannot write to result.cvs %s", err)
		}

		w.Flush()
	}
}
