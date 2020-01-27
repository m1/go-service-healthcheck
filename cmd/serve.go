package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/m1/go-service-healthcheck/api"
	"github.com/m1/go-service-healthcheck/runner"
)

var (
	serveCmd = cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}
)

func serve() {
	go func() {
		r := runner.New()
		if err := r.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	a := api.New()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
