package all

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/citado/s1-gw-ns/internal/lora/api"
	"github.com/citado/s1-gw-ns/internal/sim"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	ls := api.New(cfg.LoRaServer)

	a := app.New(cfg.App)

	sim := sim.Read()

	for _, g := range sim.Gateways {
		a.Gateway(g)

		for _, d := range g.Devices {
			if err := ls.Activate(d.DevEUI, d.Addr, g.Keys.ApplicationSKey, g.Keys.NetworkSKey); err != nil {
				pterm.Fatal.Printf("device ativation failed %+v %s\n", d, err)
			}

			pterm.Info.Printf("device activated %s\n", d.DevEUI)
		}
	}

	a.Connect()

	for i := 0; i < cfg.Tries; i++ {
		f, err := os.Create(fmt.Sprintf("result_%d.csv", i+1))
		if err != nil {
			pterm.Fatal.Printf("cannot create result.cvs %s", err)
		}
		defer f.Close()

		w := csv.NewWriter(f)

		a.PublishSubscribe()

		pterm.Success.Println(a.Durations)

		for dev, d := range a.Durations {
			r := make([]string, 0)

			r = append(r, fmt.Sprintf("%d", dev))

			for _, d := range d {
				r = append(r, fmt.Sprintf("%g", d.Seconds()))
			}

			if err := w.Write(r); err != nil {
				pterm.Fatal.Printf("cannot write to result.cvs %s", err)
			}
		}

		w.Flush()
	}
}

// Register pubsub command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "all",
			Short: "publish to network server and consume from application server",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
