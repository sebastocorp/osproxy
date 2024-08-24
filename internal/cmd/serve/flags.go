package serve

import (
	"context"
	"log"

	"osproxy/internal/logger"

	"github.com/spf13/cobra"
)

const (

)


func getFlags(cmd *cobra.Command) (flags serveFlagsT, err error) {

	// Get root command flags
	flags.logLevel, err = cmd.Flags().GetString(logLevelFlagName)
	if err != nil {
		log.Fatalf(logLevelFlagErrMsg, err.Error())
	}

	level, err := logger.GetLevel(flags.logLevel)
	if err != nil {
		log.Fatalf(logLevelFlagErrMsg, err.Error())
	}

	logger.InitLogger(context.Background(), level)

	// Get server command flags

	flags.configFile, err = cmd.Flags().GetString(configFileFlagName)
	if err != nil {
		logger.Log.Fatalf(configFileFlagErrMsg, err.Error())
	}

	return flags, err
}
