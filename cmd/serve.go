/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartServer() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is now running on localhost:8888")

	for {
		conn, err := listener.Accept()
		fmt.Println("listener.Accept")
		if err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			fmt.Printf("go func() Accept %v\n", conn.RemoteAddr())
			for {
				fmt.Println("inner for")
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))

				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					netErr, ok := err.(net.Error)
					if ok && netErr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						fmt.Println("EOF")
						break
					}
					panic(err)
				}
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))
				content := `Hello World\n`

				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					Body:          io.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)
			}
		}()
	}
}
