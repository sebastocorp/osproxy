package serve

import (
	"github.com/spf13/cobra"
)

func getFlags(cmd *cobra.Command) (flags serveFlagsT, err error) {
	flags.configFile, err = cmd.Flags().GetString(flagNameConfig)

	return flags, err
}
