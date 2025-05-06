package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configFile string

func NewIfOnlyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ifonly",
		Short: "ifonly is hlep robot to be a winner",
		Long:  "We use ifonly to let the robot win, because ifonly is the best hand generator algorithm",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello ifonly")
			return nil
		},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
	}

	cobra.OnInitialize(onInitialize)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the miniblog configuration file.")

	return cmd
}
