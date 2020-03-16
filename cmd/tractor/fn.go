package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/mdns"
	"github.com/manifold/qtalk/golang/mux"
	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/config"
	"github.com/spf13/cobra"
)

// `tractor fn` command
func fnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fn",
		Short: "Calls a remote qrpc function",
		Long:  "Calls a remote qrpc function",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fn := args[:2]
			args = args[2:]
			var sess mux.Session
			var err error
			if fn[0] == "agent" {
				cfg, cErr := config.OpenDefault()
				if cErr != nil {
					log.Fatal(cErr)
				}
				sess, err = mux.DialUnix(cfg.Agent.SocketPath)
			} else {
				svcs, sErr := lookupServices()
				if sErr != nil {
					log.Fatal(sErr)
				}
				svc := findRecord(svcs, fn[0])
				var addr string
				if svc != nil {
					addr = fmt.Sprintf("localhost:%d", svc.Port)
				} else {
					addr = fn[0] // if not a name, use as address
				}

				sess, err = mux.DialWebsocket(addr)
			}
			if err != nil {
				log.Fatal(err)
			}

			client := &qrpc.Client{Session: sess}
			var ret interface{}
			var resp *qrpc.Response
			if len(args) > 1 {
				resp, err = client.Call(fn[1], args, &ret)
			} else if len(args) == 1 {
				resp, err = client.Call(fn[1], args[0], &ret)
			} else {
				resp, err = client.Call(fn[1], nil, &ret)
			}
			if err != nil {
				log.Fatal(err)
			}

			if ret != nil {
				fmt.Printf("REPLY => %#v\n", ret)
			}

			if resp.Hijacked {
				go func() {
					<-sigQuit.Done()
					resp.Channel.Close()
				}()

				_, err = io.Copy(os.Stdout, resp.Channel)
				resp.Channel.Close()
				if err != nil && err != io.EOF {
					log.Fatal(err)
				}
				fmt.Println()
			}
		},
	}
	return cmd
}

func findRecord(svcs []*mdns.ServiceEntry, name string) *mdns.ServiceEntry {
	for _, svc := range svcs {
		if shortname(svc.Name) == name {
			return svc
		}
	}
	return nil
}

func lookupServices() ([]*mdns.ServiceEntry, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	var entries []*mdns.ServiceEntry
	go func() {
		for entry := range entriesCh {
			if strings.HasSuffix(entry.Name, "_tractor._tcp.local.") {
				entries = append(entries, entry)
			}
		}
	}()
	if err := mdns.Lookup("_tractor._tcp", entriesCh); err != nil {
		return nil, err
	}
	close(entriesCh)
	return entries, nil
}

func shortname(name string) string {
	return strings.Split(name, ".")[0]
}
