package cmd

import (
	"osproxy/internal/cmd/serve"
	"osproxy/internal/cmd/version"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `TODO`
	descriptionLong  = `
	TODO`
)

func NewRootCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  descriptionLong,
	}

	cmd.AddCommand(
		version.NewCommand(),
		serve.NewCommand(),
	)

	return cmd
}
