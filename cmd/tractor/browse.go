package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/mdns"
	"github.com/spf13/cobra"
)

func browseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "browse",
		Short: "Browse discoverable workspaces",
		Long:  "Browse discoverable workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			entriesCh := make(chan *mdns.ServiceEntry, 4)
			go func() {
				for entry := range entriesCh {
					if strings.HasSuffix(entry.Name, "_tractor._tcp.local.") {
						parts := strings.Split(entry.Name, ".")
						fmt.Printf("%s (%d) %s\n", parts[0], entry.Port, entry.Info)
					}
				}
			}()
			err := mdns.Lookup("_tractor._tcp", entriesCh)
			if err != nil {
				log.Fatal(err)
			}
			close(entriesCh)
		},
	}
	return cmd
}
