package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"go.evanpurkhiser.com/dots/config"
)

var (
	sourceConfig   *config.SourceConfig
	sourceLockfile *config.SourceLockfile
)

func loadConfigs(cmd *cobra.Command, args []string) error {
	var err error

	path := config.SourceConfigPath()

	sourceConfig, err = config.LoadSourceConfig(path)
	if err != nil {
		return err
	}

	sourceLockfile, err = config.LoadLockfile(sourceConfig)
	if err != nil {
		return err
	}

	warns := config.SanitizeSourceConfig(sourceConfig)
	for _, err := range warns {
		color.New(color.FgYellow).Fprintf(os.Stderr, "warn: %s\n", err)
	}

	return nil
}

var rootCmd = cobra.Command{
	Use:   "dots",
	Short: "A portable tool for managing a single set of dotfiles",

	SilenceUsage:      true,
	SilenceErrors:     true,
	PersistentPreRunE: loadConfigs,
}

func main() {
	cobra.EnableCommandSorting = false

	rootCmd.AddCommand(&filesCmd)
	rootCmd.AddCommand(&diffCmd)
	rootCmd.AddCommand(&installCmd)
	rootCmd.AddCommand(&configCmd)

	if err := rootCmd.Execute(); err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
