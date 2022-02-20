package cmd

import (
	"os"

	"github.com/citado/s1-gw-ns/internal/cmd/all"
	"github.com/citado/s1-gw-ns/internal/cmd/gen"
	"github.com/citado/s1-gw-ns/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("S", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("miulation", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("-1", pterm.NewStyle(pterm.FgLightGreen))).
		Render()

	cfg := config.New()

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "s1-gw-ns",
		Short: "Simulation 1: between gateway and network server lorawan",
	}

	all.Register(root, cfg)
	gen.Register(root, cfg)

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
