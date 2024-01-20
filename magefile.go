//go:build mage
// +build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
	return cmd.Run()
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}

// Launch local docker compose with telegram bot and sqlite database
func Dev() error {
	fmt.Println("Preparing to launch local env")
	fmt.Println(
		"launching command",
		"podman-compose",
		"up",
		"-d",
		"--build",
		"--remove-orphans",
	)
	cmd := exec.Command(
		"podman-compose",
		"up",
		"-d",
		"--build",
		"--remove-orphans",
	)
	return cmd.Run()
}

// Deploy to your fly io account, note that you must have flyctl configured for this step
func FlyDeploy() error {
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
	return cmd.Run()
}
