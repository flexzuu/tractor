package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commitOID string

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the Tractor version",
		Long:  "Print the Tractor version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Tractor build", commitOID)
		},
	}
}
