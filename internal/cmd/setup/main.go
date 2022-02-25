package setup

import (
	"errors"
	"fmt"
	"time"

	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/citado/s1-gw-ns/internal/lora/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	ls := api.New(cfg.LoRaServer)

	pterm.Info.Printf("these helpers are only for simulation purpose, don't use them in production!")

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

	deviceProfileID, err := ls.GetOrCreateDeviceProfile("fake_dp", "1", "1")
	if err != nil {
		pterm.Fatal.Printf("device profile creation failed %s\n", err.Error())
	}

	pterm.Info.Printf("device profile %s is ready for duty\n", deviceProfileID)

	if err := ls.CreateApplication("citado", "application for load testing", "1", serviceProfileID); err != nil {
		if !errors.Is(err, api.ErrDuplicateApp) {
			pterm.Fatal.Printf("application creation failed %s\n", err.Error())
		}
	}

	for i := 0; i < 10; i++ {
		mac := api.GenerateGWID()

		if err := ls.CreateGateway(
			mac,
			fmt.Sprintf("generated-gateway-%d", i),
			fmt.Sprintf("generated on %s", time.Now()),
			"1",
			"1",
			serviceProfileID,
		); err != nil {
			if !errors.Is(err, api.ErrDuplicateGateway) {
				pterm.Error.Printf("gateway generation failed")
			}
		}

		pterm.Info.Printf("gateway %s created without problem!\n", mac)
	}
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
