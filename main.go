package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "env-injector",
		Short: "Insert env variables into SPA via meta tags",
		Run:   run,
	}
	whitelist string
)

func run(cmd *cobra.Command, args []string) {
	err := inject(whitelist, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flagSet := rootCmd.Flags()
	flagSet.StringVar(&whitelist, "whitelist", "", "Comma separated list of env vars to inject")

	cobra.MarkFlagRequired(flagSet, "whitelist")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
