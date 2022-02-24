package gen

import (
	"fmt"
	"time"

	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/citado/s1-gw-ns/internal/lora"
	"github.com/citado/s1-gw-ns/internal/lora/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func main(cfg config.Config, count int) {
	ls := api.New(cfg.LoRaServer)
	devices := make([]lora.Device, 0)

	for i := 0; i < count; i++ {
		devEUI := api.GenerateDevEUI()
		if err := ls.CreateDevice(
			devEUI,
			fmt.Sprintf("generated-device-%d", i),
			cfg.LoRaServer.ApplicationID,
			fmt.Sprintf("generated on %s", time.Now()),
			cfg.LoRaServer.DeviceProfileID,
			true,
			0.0,
			false,
		); err != nil {
			pterm.Error.Printf("device generation failed %s %s\n", devEUI, err)
		} else {
			pterm.Info.Printf("device %s generated\n", devEUI)
		}

		devices = append(devices, lora.Device{
			DevEUI: devEUI,
			Addr:   api.GenerateDevAddr(),
		})
	}

	b, err := yaml.Marshal(devices)
	if err != nil {
		pterm.Fatal.Printf("cannot create yaml from generated devices %s\n", err)
	}

	// nolint: forbidigo
	fmt.Println(string(b))
}

// Register pubsub command.
func Register(root *cobra.Command, cfg config.Config) {
	var count int

	// nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate lora devices in lorawan server and writes its configuration",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, count)
		},
	}

	cmd.Flags().IntVar(&count, "count", 10, "number of generated devices") // nolint: gomnd

	root.AddCommand(cmd)
}
