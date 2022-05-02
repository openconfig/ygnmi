package main

import (
	"log"

	"github.com/openconfig/ygnmi/cmd/generator"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use: "ygnmi",
	}
	root.AddCommand(generator.New())

	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
