package cmd

import (
    "fmt"

    "vault-cli/internal/server"

    "github.com/spf13/cobra"
)

var (
    listenAddr string
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Start the web UI server",
    RunE: func(cmd *cobra.Command, args []string) error {
        srv := server.New(cfg, database)
        fmt.Printf("Starting Vault UI server on %s...\n", listenAddr)
        return srv.Start(listenAddr)
    },
}

func init() {
    serverCmd.Flags().StringVar(&listenAddr, "addr", "127.0.0.1:8080", "address to bind the web server")
    rootCmd.AddCommand(serverCmd)
}


