//go:build mage
// +build mage

package main

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(Deps)
	fmt.Println("Building...")
	cmd := exec.Command(
		"RUN CGO_ENABLED=0",
		"GOOS=linux",
		"GOARCH=amd64",
		"go",
		"build",
		"./src/cmd/main.go",
		"-o",
		"./orio-telegram-adapter",
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%e", err)
	}
	return nil
}

// Manage your deps, or running package managers.
func Deps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	if err := cmd.Run(); err != nil {
		fmt.Printf("%e", err)
	}
	return nil
}

const (
	podmanComposeCommand = "podman-compose"
	dockerComposeCommand = "docker compose"
)

func hasBinary(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

func launchDockerOrPodman() (string, error) {
	if hasBinary(podmanComposeCommand) {
		return podmanComposeCommand, nil
	}

	if hasBinary(dockerComposeCommand) {
		return dockerComposeCommand, nil
	}

	return "", errors.New("missing command to launch local env, please consult manual for local env required tools")
}

// Launch local docker compose with telegram bot and sqlite database
func Dev() error {
	fmt.Println("Preparing to launch local env")

	command, err := launchDockerOrPodman()
	if err != nil {
		return err
	}

	fmt.Println(
		"launching command",
		command,
		"up",
		"-d",
		"--build",
		"--remove-orphans",
		"--force-recreate",
	)
	cmd := exec.Command(
		command,
		"up",
		"-d",
		"--build",
		"--remove-orphans",
		"--force-recreate",
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Deploy to your fly io account, note that you must have flyctl configured for this step
func Deploy_fly() error {
	fmt.Println("Preparing to deploy to fly io")
	fmt.Println(
		"flyctl",
		"deploy",
		"--ha=false",
	)
	cmd := exec.Command(
		"flyctl",
		"deploy",
		"--ha=false",
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%e", err)
	}
	return nil
}
