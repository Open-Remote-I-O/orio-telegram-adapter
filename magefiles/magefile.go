//go:build mage
// +build mage

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/sh"
)

const (
	podmanComposeCommand = "podman-compose"
	dockerComposeCommand = "docker compose"
	podmanCommand        = "podman"
	dockerCommand        = "docker"
)

func hasBinary(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

func inputConsolePrompt(label string) (string, error) {
	fmt.Println(label)
	scanner := bufio.NewScanner(os.Stdin)
	if scanned := scanner.Scan(); scanned {
		return strings.TrimSpace(scanner.Text()), nil
	}

	return "", scanner.Err()
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
		"--env-file",
		"./docker-compose.env",
		"up",
		"-d",
		"--build",
		"--remove-orphans",
		"--force-recreate",
	)

	if err := sh.RunV(
		command,
		"--env-file",
		"./docker-compose.env",
		"up",
		"-d",
		"--build",
		"--remove-orphans",
		"--force-recreate",
	); err != nil {
		return err
	}
	return nil
}

func dockerLoginFlyRegistry() error {
	err := sh.RunV(
		"flyctl",
		"auth",
		"docker",
	)
	if err != nil {
		return err
	}
	return nil
}

func podmanLoginFlyRegistry() error {
	fmt.Println("Retrieving fly auth token")
	fmt.Println(
		"flyctl",
		"auth",
		"token",
	)

	// TODO The 'fly auth token' command is deprecated. Use 'fly tokens create' instead.
	rawAuthToken, err := sh.Output(
		"flyctl",
		"auth",
		"token",
	)
	if err != nil {
		return err
	}

	registryUrlInput, err := inputConsolePrompt("provide your registry to push image to:")
	if err != nil {
		return err
	}

	registryUrl := strings.TrimSuffix(registryUrlInput, "\n")
	authToken := strings.TrimSuffix(string(rawAuthToken), "\n")

	fmt.Println("Executing registry login operation")
	if err := sh.RunV(
		"podman",
		"login",
		"-u",
		"x",
		"-p",
		authToken,
		registryUrl,
	); err != nil {
		fmt.Printf("Something went wring while completing auth step on %s registry url \n", registryUrl)
		fmt.Printf("%e \n", err.Error())
	}
	return nil
}

// Login to fly image registry
func LoginImageRegistry() error {
	if hasBinary(podmanCommand) {
		return podmanLoginFlyRegistry()
	}
	return dockerLoginFlyRegistry()
}

func podmanBuildContainerImage(imageTag string) error {
	if err := sh.RunV(
		"podman",
		"build",
		"-t",
		fmt.Sprintf("localhost/%s", imageTag),
		"--build-arg-file",
		"./argfile.conf",
		".",
	); err != nil {
		fmt.Printf("%e", err)
		return err
	}
	return nil
}

func podmanPushContainerImage(imageTag string) error {
	if err := sh.RunV(
		"podman",
		"push",
		"--format",
		"v2s2",
		fmt.Sprintf("localhost/%s", imageTag),
		fmt.Sprintf("docker://registry.fly.io/%s:latest", imageTag),
	); err != nil {
		fmt.Printf("%e", err)
		return err
	}
	return nil
}

// Build image locally and pushes into fly.io image registry
func BuildPushImageToRegistry() error {
	if !hasBinary(podmanCommand) {
		return fmt.Errorf("Currently docker build and push step compatibility is not available")
	}
	// TODO default as latest
	rawContainerImageTag, err := inputConsolePrompt("provide container image tag:")
	if err != nil {
		return err
	}
	containerImageTag := strings.TrimSuffix(rawContainerImageTag, "\n")
	err = podmanBuildContainerImage(containerImageTag)
	if err != nil {
		return err
	}
	err = podmanPushContainerImage(containerImageTag)
	if err != nil {
		return err
	}
	fmt.Println("image built and pushed successfully")
	return nil
}

// Deploy to your fly.io account from fly.io registry, note that you must have flyctl configured for this step
func DeployImageFromRegistry() error {
	fmt.Println("Preparing to deploy")
	rawContainerImageTag, err := inputConsolePrompt("provide container image tag:")
	if err != nil {
		return err
	}
	containerImageTag := strings.TrimSuffix(rawContainerImageTag, "\n")
	fmt.Println(
		"flyctl",
		"deploy",
		"-i",
		fmt.Sprintf("registry.fly.io/%s:latest", containerImageTag),
		"--ha=false",
	)
	if err := sh.RunV(
		"flyctl",
		"deploy",
		"-i",
		fmt.Sprintf("registry.fly.io/%s:latest", containerImageTag),
		"--ha=false",
	); err != nil {
		fmt.Printf("%e", err)
		return err
	}
	return nil
}
