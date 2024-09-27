package serve

import (
	"log"

	"osproxy/internal/osproxy"

	"github.com/spf13/cobra"
)

const (
	// FLAG NAMES

	flagNameConfig = `config`
)

type serveFlagsT struct {
	configFile string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "serve",
		DisableFlagsInUseLine: true,
		Short:                 `Execute serve process`,
		Long: `
	Run execute serve process`,

		Run: RunCommand,
	}

	cmd.Flags().String(flagNameConfig, "osproxy.yaml", "filepath to config file")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {
	flags, err := getFlags(cmd)
	if err != nil {
		log.Fatalf("unable to get flags: %s", err.Error())
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////
	osproxy, err := osproxy.NewOSProxy(flags.configFile)
	if err != nil {
		log.Fatalf("unable init proxy: %s", err.Error())
	}

	// Iniciar el servidor proxy
	err = osproxy.Run()
	if err != nil {
		log.Fatalf("unable init proxy: %s", err.Error())
	}
}
