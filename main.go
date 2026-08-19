package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "occ",
		Short: "Opencode Configure for containerized sandboxes",
		Long:  "occ generates reproducible Docker/Podman containers for AI coding agents from a declarative container.yml",
	}
	initCmd := &cobra.Command{
		Use:   "init [docker|podman]",
		Short: "Scaffold a starter container.yml in the current project",
		// ExactArgs(1) garantisce che l'utente passi esattamente un argomento dopo init
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] conterrà "docker" oppure "podman"
			target := args[0]

			switch target {
			case "docker":
				fmt.Println("Inizializzazione per container Docker")
			case "podman":
				fmt.Println("Inizializzazione per container Podman")
			default:
				fmt.Println("Tipo container non supportato")
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
