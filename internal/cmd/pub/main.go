package pub

import (
	"time"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	a := app.New(cfg.App)

	for _, g := range cfg.Gateways {
		a.Gateway(g)
	}

	a.Connect()

	for i := 0; i < cfg.Tries; i++ {
		pterm.Info.Printf("iteration %d\n", i)

		a.Publish()

		// wait to cool down the system and complete the previous iteration.
		time.Sleep(time.Minute)
	}
}

// Register publish command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "pub",
			Short: "publish to network server",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
