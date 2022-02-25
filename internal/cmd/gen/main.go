package gen

import (
	"errors"
	"fmt"
	"time"

	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/citado/s1-gw-ns/internal/lora"
	"github.com/citado/s1-gw-ns/internal/lora/api"
	"github.com/citado/s1-gw-ns/internal/sim"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(cfg config.Config, devCount, gwCount int) {
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

	gateways := make([]lora.Config, 0)

	for i := 0; i < gwCount; i++ {
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

		gateway := lora.Config{
			MAC: mac,
			Keys: lora.Keys{
				NetworkSKey:     "DB56B6C3002A4763A79E64573C629D97",
				ApplicationSKey: "94B49CD7BC621BC46571D019640804AA",
			},
			Devices: make([]lora.Device, 0),
		}

		for j := 0; j < devCount; j++ {
			devEUI := api.GenerateDevEUI()
			if err := ls.CreateDevice(
				devEUI,
				fmt.Sprintf("generated-device-%d-%d", i, j),
				1,
				fmt.Sprintf("generated on %s", time.Now()),
				deviceProfileID,
				true,
				0.0,
				false,
			); err != nil {
				pterm.Error.Printf("device generation failed %s %s\n", devEUI, err)
			} else {
				pterm.Info.Printf("device %s generated\n", devEUI)
			}

			gateway.Devices = append(gateway.Devices, lora.Device{
				DevEUI: devEUI,
				Addr:   api.GenerateDevAddr(),
			})
		}

		gateways = append(gateways, gateway)
	}

	sim.Write(sim.Config{
		Gateways: gateways,
	})
}

// Register gen command.
func Register(root *cobra.Command, cfg config.Config) {
	var (
		devCount     int
		gatewayCount int
	)

	// nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate lora devices/gatways in lorawan server and writes their configuration",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, devCount, gatewayCount)
		},
	}

	cmd.Flags().IntVar(&devCount, "dev-count", 10, "number of generated devices")          // nolint: gomnd
	cmd.Flags().IntVar(&gatewayCount, "gateway-count", 10, "number of generated gateways") // nolint: gomnd

	root.AddCommand(cmd)
}
