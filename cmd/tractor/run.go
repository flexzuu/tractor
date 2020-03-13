package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/manifold/tractor/pkg/misc/logging/zap"
	"github.com/manifold/tractor/pkg/workspace/supervisor"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run workspace daemon",
		Long:  "Run workspace daemon",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat("tractor.go"); os.IsNotExist(err) {
				fmt.Println("not a valid tractor workspace")
				os.Exit(1)
			}
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			svr := supervisor.New(wd, "", os.Stdout)
			svr.Log = zap.NewLogger(os.Stdout)
			dm := daemon.New(svr)
			if err := dm.Run(context.Background()); err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}
