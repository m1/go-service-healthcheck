package main

import (
	"log"

	"github.com/m1/go-service-healthcheck/cmd"
	"github.com/m1/go-service-healthcheck/config"
)

var (
	// GitCommit ...
	GitCommit string

	// BuildDate ...
	BuildDate string

	// Version ...
	Version string
)

func main() {
	config.GitCommit = GitCommit
	config.BuildDate = BuildDate
	config.Version = Version
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
