/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// serveUdpCmd represents the serveUdp command
var serveUdpCmd = &cobra.Command{
	Use:   "serveUdp",
	Short: "sample code for udp server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serveUdp called")
		startServe()
	},
}

func init() {
	rootCmd.AddCommand(serveUdpCmd)
}

func startServe() {
	conn, err := net.ListenPacket("udp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		// クライアントを知らない状態で開くソケット.
		length, remoteAddress, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received from %v: %v\n", remoteAddress, string(buffer[:length]))
		_, err = conn.WriteTo([]byte("Hello from Server"), remoteAddress)

		if err != nil {
			panic(err)
		}
	}

}
