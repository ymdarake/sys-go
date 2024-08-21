/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// sigCmd represents the sig command
var sigCmd = &cobra.Command{
	Use:   "sig",
	Short: "sample for handling signals",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sig called")
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(sigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func start() {
	ctx := context.Background()
	sigctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 5*time.Second)
	defer cancelTimeout()

	select {
	case <-sigctx.Done():
		fmt.Println("signal received")
	case <-timeoutCtx.Done():
		fmt.Println("timeout")
	}
}

func startServer() {
	ctx := context.Background()
	sigctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	server := newServer(sigctx)
	serveWaitChan := make(chan struct{})
	go server.Start(serveWaitChan)
	select {
	case <-sigctx.Done():
		fmt.Println("\nsignal received")
	}

	server.Shutdown()

	// メッセージング処理など、サーバーとしてレスポンスを返さない処理を完了してから全体を終了させる
	<-serveWaitChan
	fmt.Println("wait done")
}

type Server struct {
	Start    func(done chan struct{})
	Shutdown func()
}

func newServer(ctx context.Context) Server {
	return Server{
		Start: func(done chan struct{}) {
			fmt.Println("Starting server")
			ticker := time.NewTicker(1 * time.Second)
		loop:
			for {
				select {
				case <-ticker.C:
					fmt.Println("\ntick")
				case <-ctx.Done():
					fmt.Println("\nWaiting for existing processes to complete")
					time.Sleep(10 * time.Second)
					close(done)
					break loop
				}
			}
		},
		Shutdown: func() {
			fmt.Println("Shutting down server")
		},
	}
}
