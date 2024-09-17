package serve

import (
	"fmt"
	"log"
	"net/http"

	"osproxy/internal/logger"
	"osproxy/internal/osproxy"

	"github.com/spf13/cobra"
)

const (
	// FLAG NAMES

	logLevelFlagName   = `log-level`
	configFileFlagName = `config`

	// ERROR MESSAGES

	logLevelFlagErrMsg   = "unable to get flag --log-level: %s"
	configFileFlagErrMsg = "unable to get flag --config: %s"
)

type serveFlagsT struct {
	logLevel   string
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

	cmd.Flags().String(logLevelFlagName, "info", "verbosity level for logs")

	cmd.Flags().String(configFileFlagName, "osproxy.yaml", "filepath to config file")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {
	flags, err := getFlags(cmd)
	if err != nil {
		logger.Log.Fatalf("unable to parse daemon command flags")
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////
	osproxy, err := osproxy.NewOSProxy(flags.configFile)
	if err != nil {
		logger.Log.Fatalf("unable init osproxy: %s", err.Error())
	}

	// Iniciar el servidor proxy
	logger.Log.Infof("init osproxy")
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%s", osproxy.Config.Proxy.Address, osproxy.Config.Proxy.Port),
		http.HandlerFunc(osproxy.HandleFunc),
	))
}
