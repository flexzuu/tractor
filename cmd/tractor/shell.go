package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/manifold/qtalk/golang/mux"
	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/agent"
	"github.com/spf13/cobra"
)

func shellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "An interactive repl",
		Long:  "An interactive repl",
		Run:   shellRun,
	}
	cmd.PersistentFlags().StringVarP(&tractorUserPath, "path", "p", "", "path to the user tractor directory (default is ~/.tractor)")
	return cmd
}

func workspaceAddr(dir string) string {
	ag, err := agent.New(tractorUserPath, nil, false)
	if err != nil {
		fatal(err)
	}
	sess, err := mux.DialUnix(ag.SocketPath)
	if err != nil {
		log.Fatal(err)
	}
	client := &qrpc.Client{Session: sess}

	var addr string
	_, err = client.Call("connect", dir, &addr)
	if err != nil {
		log.Fatal(err)
	}
	sess.Close()
	return addr
}

func shellRun(cmd *cobra.Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	sess, err := mux.DialUnix(workspaceAddr(dir))
	if err != nil {
		log.Fatal(err)
	}
	client := &qrpc.Client{Session: sess}
	resp, err := client.Call("repl", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Hijacked {
		scanner := bufio.NewScanner(resp.Channel)
		scanner.Split(ScanLines)
		executor := func(in string) {
			fmt.Fprintf(resp.Channel, "%s\n", in)
			if scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "remote:", err)
			}
		}
		p := prompt.New(
			executor,
			func(in prompt.Document) (out []prompt.Suggest) { return },
			prompt.OptionPrefix(">>> "),
		)
		p.Run()
	}

}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
