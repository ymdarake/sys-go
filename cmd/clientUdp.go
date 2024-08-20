/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// clientUdpCmd represents the clientUdp command
var clientUdpCmd = &cobra.Command{
	Use:   "clientUdp",
	Short: "sample code for udp client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clientUdp called")
		startClient()
	},
}

func init() {
	rootCmd.AddCommand(clientUdpCmd)
}

func startClient() {
	conn, err := net.Dial("udp4", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Sending to server")
	_, err = conn.Write([]byte("Hello from Client"))
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1500)
	length, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received: %s\n", string(buffer[:length]))
}
