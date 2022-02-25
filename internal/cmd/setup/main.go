package setup

import (
	"errors"

	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/citado/s1-gw-ns/internal/lora/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	ls := api.New(cfg.LoRaServer)

	if err := ls.CreateNetworkServer("1", "citado", "chirpstack-network-server:8000"); err != nil {
		if !errors.Is(err, api.ErrDuplicateNS) {
			pterm.Fatal.Printf("network server creation failed %s\n", err.Error())
		}
	}

	serviceProfileID, err := ls.GetOrCreateServiceProfile("fake_profile", "1", "1")
	if err != nil {
		pterm.Fatal.Printf("service profile creation failed %s\n", err.Error())
	}

	pterm.Info.Printf("service profile %s is ready for duty\n", serviceProfileID)
}

// Register pubsub command.
func Register(root *cobra.Command, cfg config.Config) {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup chirpstack server for the first time",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	}

	root.AddCommand(cmd)
}
