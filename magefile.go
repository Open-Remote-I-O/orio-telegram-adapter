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

	"github.com/magefile/mage/mg"
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
	scanner := bufio.NewScanner(os.Stdin)
	if scanned := scanner.Scan(); scanned {
		return strings.TrimSpace(scanner.Text()), nil
	}

	return "", scanner.Err()
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(Deps)
	fmt.Println("Building...")
	cmd := exec.Command(
		"RUN CGO_ENABLED=0",
		"GOOS=linux",
		"GOARCH=amd64",
		`-ldflags="-s -w`,
		"go",
		"build",
		"./src/cmd/main.go",
		"-o",
		"./orio-telegram-adapter",
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%e", err)
		return err
	}
	return nil
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

// Manage your deps, or running package managers.
func Deps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	if err := cmd.Run(); err != nil {
		fmt.Printf("%e", err)
	}
	return nil
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
	cmd := exec.Command(
		command,
		"--env-file",
		"./docker-compose.env",
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

func dockerLoginFlyRegistry() error {
	fmt.Println("")
	dockerAuthCmd := exec.Command(
		"flyctl",
		"auth",
		"docker",
	)
	out, err := dockerAuthCmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func podmanLoginFlyRegistry() error {
	fmt.Println("Retrieving fly auth token")
	fmt.Println(
		"flyctl",
		"auth",
		"token",
	)
	authTokenCmd := exec.Command(
		"flyctl",
		"auth",
		"token",
	)
	rawAuthToken, err := authTokenCmd.Output()
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
	cmd := exec.Command(
		"podman",
		"login",
		"-u",
		"x",
		"-p",
		authToken,
		registryUrl,
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Something went wring while completing auth step on %s registry url \n", registryUrl)
		fmt.Printf("%e \n", err.Error())
	}
	return nil
}

func LoginImageRegistry() error {
	if hasBinary(podmanCommand) {
		return podmanLoginFlyRegistry()
	}
	return dockerLoginFlyRegistry()
}

func podmanBuildContainerImage(imageTag string) error {
	cmd := exec.Command(
		"podman",
		"build",
		"-t",
		fmt.Sprintf("localhost/%s", imageTag),
		"--build-arg-file",
		"./argfile.conf",
		".",
	)
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

func podmanPushContainerImage(imageTag string) error {
	cmd := exec.Command(
		"podman",
		"push",
		"--format",
		"v2s2",
		fmt.Sprintf("localhost/%s", imageTag),
		fmt.Sprintf("docker://registry.fly.io/%s:latest", imageTag),
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("error occured during image push to registry %e", err)
		return err
	}
	return nil
}

func BuildPushImageToRegistry() error {
	if !hasBinary(podmanCommand) {
		return fmt.Errorf("Currently docker build and push step compatibility is not available")
	}
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

// Deploy to your fly io account, note that you must have flyctl configured for this step
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
	cmd := exec.Command(
		"flyctl",
		"deploy",
		"-i",
		fmt.Sprintf("registry.fly.io/%s:latest", containerImageTag),
		"--ha=false",
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%e", err)
	}
	return nil
}
