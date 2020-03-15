package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/manifold/tractor/pkg/agent"
	"github.com/manifold/tractor/pkg/agent/console"
	"github.com/manifold/tractor/pkg/agent/systray/subprocess"
	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tractor-agent",
		Short: "Tractor Agent",
		Long:  "Tractor Agent",
		Run:   runAgent,
	}

	tractorUserPath string
	devMode         bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&devMode, "dev", "d", false, "run in debug mode")
	rootCmd.PersistentFlags().StringVarP(&tractorUserPath, "path", "p", "", "path to the user tractor directory (default is ~/.tractor)")
}

func main() {
	rootCmd.Execute()
}

func runAgent(cmd *cobra.Command, args []string) {
	if os.Getenv("SYSTRAY_SUBPROCESS") != "" {
		subprocess.Run()
		return
	}

	ag, err := agent.New(tractorUserPath, console.New(), devMode)
	fatal(err)

	if ag.Config.Agent.CheckSocket() && devMode {
		fmt.Println("Agent will not run in dev mode if agent socket exists.")
		return
	}

	dm := daemon.New(ag.DaemonServices()...)
	ctx := context.Background()
	fatal(dm.Run(ctx))
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
