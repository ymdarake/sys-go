/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")
		StartChunkClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartClient() {
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	current := 0

	var conn net.Conn = nil

	for {
		var err error
		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:8888")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}

		request, err := http.NewRequest(
			"POST",
			"http://localhost:8888",
			strings.NewReader(sendMessages[current]),
		)
		if err != nil {
			panic(err)
		}
		request.Header.Set("Accept-Encoding", "gzip")

		err = request.Write(conn)
		if err != nil {
			panic(err)
		}

		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}

		dump, err := httputil.DumpResponse(response, false)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))
		defer response.Body.Close()

		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				panic(err)
			}
			io.Copy(os.Stdout, reader)
			reader.Close()
		} else {
			io.Copy(os.Stdout, response.Body)
		}

		current++
		if current == len(sendMessages) {
			break
		}
	}
	conn.Close()
}

func StartChunkClient() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request, err := http.NewRequest(
		"GET",
		"http://localhost:8888",
		nil,
	)
	if err != nil {
		panic(err)
	}

	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}

	dump, err := httputil.DumpResponse(response, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))

	if len(response.TransferEncoding) < 1 || response.TransferEncoding[0] != "chunked" {
		panic("Wrong TransaferEncoding. Not 'chunked'.")
	}

	for {
		// serve.go で指定している形式
		// Readは新しい書き込みをタイムアウト付きで待ち受ける。
		// 下記、connのReadについて抜粋。
		// Read reads data from the connection.
		// Read can be made to time out and return an error after a fixed
		// time limit; see SetDeadline and SetReadDeadline.
		// Read(b []byte) (n int, err error)
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		size, err := strconv.ParseInt(
			string(sizeStr[:len(sizeStr)-2]),
			16,
			64,
		)

		if size == 0 {
			break
		}

		if err != nil {
			panic(err)
		}

		line := make([]byte, int(size))
		io.ReadFull(reader, line)
		reader.Discard(2)
		fmt.Printf("	%d bytes: %s\n", size, string(line))
	}
}
