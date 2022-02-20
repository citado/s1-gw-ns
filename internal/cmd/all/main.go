package all

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/citado/s1-gw-ns/internal/lora"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	ls := lora.NewAPI(cfg.LoRaServer)

	a := app.New(cfg.App)

	f, err := os.Create("result.csv")
	if err != nil {
		pterm.Fatal.Printf("cannot create result.cvs %s", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for _, g := range cfg.Gateways {
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
		a.PublishSubscribe()

		pterm.Success.Println(a.Durations)

		r := make([]string, 0)
		for _, d := range a.Durations {
			r = append(r, fmt.Sprintf("%g", d.Seconds()))
		}

		if err := w.Write(r); err != nil {
			pterm.Fatal.Printf("cannot write to result.cvs %s", err)
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
