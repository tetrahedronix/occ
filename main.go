package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "occ",
		Short: "Opencode Configure for containerized sandboxes",
		Long:  "occ generates reproducible Docker/Podman containers for AI coding agents from a declarative orchestrator.yml",
	}
	initCmd := &cobra.Command{
		Use:   "init [docker|podman]",
		Short: "Scaffold a starter orchestrator.yml in the current project",
		// ExactArgs(1) garantisce che l'utente passi esattamente un argomento dopo init
		Args: cobra.ExactArgs(1),
		RunE: runInit,
	}

	// "force" è il nome lungo (--force)
	// "f" è lo short (-f)
	// false è il valore di default
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing file")

	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	var runtime Runtime
	// args[0] dovrebbe contenere "docker" oppure "podman"
	target := args[0]

	switch target {
	case "docker":
		runtime = RuntimeDocker
	case "podman":
		runtime = RuntimePodman
	default:
		return fmt.Errorf("runtime not supported: %s. Use 'docker' or 'podman'.", target)
	}

	// Invoca NewConfig per popolare la struttura
	cfg, err := newConfig(runtime)

	if err != nil {
		return fmt.Errorf("impossibile generare la configurazione di base: %w", err)
	}

	// Legge la flag di sovrascrittura solo dove serve, ovvero prima di scrivere il file
	// GetBool restituisce il valore booleano e un eventuale errore se la flag non esiste
	force, _ := cmd.Flags().GetBool("force")

	// Invoca la funzione di scrittura passando la config e la flag force
	err = newConfigFile(cfg, force)

	if err != nil {
		return fmt.Errorf("impossibile scrivere il file di configurazione: %w", err)
	}

	// Messaggio di successo finale (mostrato solo se tutto è andato a buon fine)
	log.Printf("Configuration for %s initialized successfully.\n", runtime)

	return nil
}
