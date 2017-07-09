package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"flag"
	"path/filepath"
	"strings"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

var dockerComposeFile string

func waitForTermination(p project.APIProject) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	err := p.Delete(context.Background(), options.Delete{
		RemoveRunning: true,
		RemoveVolume:  true,
	})
	if err != nil {
		log.Fatalf("Failed to delete project: %s", err)
	}
}

func main() {
	flag.StringVar(&dockerComposeFile, "f", "", "path to the docker compose file")
	flag.Parse()

	if dockerComposeFile == "" {
		log.Fatalln("docker compose file not given")
	}

	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{dockerComposeFile},
			ProjectName:  fmt.Sprintf("docker-compose-test-%s", strings.TrimSuffix(dockerComposeFile, filepath.Ext(dockerComposeFile))),
		},
	}, nil)

	if err != nil {
		log.Fatalf("Failed to create project: %s", err)
	}

	err = project.Up(context.Background(), options.Up{})
	if err != nil {
		log.Fatalf("Failed to up project: %s", err)
	}

	log.Println("Hit CTRL+C to stop project...")
	waitForTermination(project)
}
