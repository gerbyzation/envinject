package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "env-injector",
	Short: "Insert env variables into SPA via meta tags",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lets do something fun!")
		filename := "index.html"
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		fileReader := bufio.NewReader(file)
		newFile, err := os.Create("new.html")
		defer newFile.Close()
		if err != nil {
			panic(err)
		}
		writer := bufio.NewWriter(newFile)

		err = Inject("this is a test", fileReader, writer)
		if err != nil {
			panic(err)
		}
		writer.Flush()
	},
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
