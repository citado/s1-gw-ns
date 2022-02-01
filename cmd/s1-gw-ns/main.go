package main

import (
	"time"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/lora"
	"github.com/pterm/pterm"
)

func main() {
	_ = pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("S", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("miulation", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("-1", pterm.NewStyle(pterm.FgLightGreen))).
		Render()

	a := app.New(app.Config{
		Addr:  "127.0.0.1",
		Port:  1883,
		Total: 10,
		Delay: 100 * time.Millisecond,
	})

	a.Gateway(lora.Config{
		MAC: "b827ebffff70c80a",
		Keys: lora.Keys{
			NetworkSKey:     "DB56B6C3002A4763A79E64573C629D97",
			ApplicationSKey: "94B49CD7BC621BC46571D019640804AA",
		},
		Device: lora.Device{
			Addr: "26011CF6",
		},
	})

	a.Connect()

	a.Run()

	pterm.Success.Println(a.Durations)
}
