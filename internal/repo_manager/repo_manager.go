package repo_manager

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Checkout(repoPath string, targetBranch string) error {

	gitPath, err := exec.LookPath("git")

	if err != nil {
		log.Fatal("Git not found in PATH")
	}

	cmd := exec.Command(gitPath, "-C", repoPath, "checkout", targetBranch)

	err = cmd.Run()

	return err
}

func GitStatus(repoPath string) (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Git not found in PATH")
	}

	cmd := exec.Command(gitPath, "-C", repoPath, "status")

	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	err = cmd.Run()

	return outBuffer.String(), err
}

func GitPull(repoPath string) error {
	remoteURL, err := getRemoteURL(repoPath)
	if err != nil {
		return err
	}

	gitPath, err := exec.LookPath("git")

	if err != nil {
		log.Fatal("Git not found in PATH")
	}

	repoUrl := addCredentialsToURL(remoteURL, os.Getenv("GIT_EMAIL"), os.Getenv("GIT_PASSWORD"))

	fmt.Println("Remote Repository URL:", repoUrl)

	cmd := exec.Command(gitPath, "-C", repoPath, "pull", repoUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return err
}

func addCredentialsToURL(url, username, password string) string {
	parts := strings.SplitN(url, "://", 2)
	if len(parts) != 2 {
		log.Fatal("Invalid URL format")
	}

	credentials := fmt.Sprintf("%s:%s@", username, password)
	urlWithCredentials := fmt.Sprintf("%s://%s%s", parts[0], credentials, parts[1])

	return urlWithCredentials
}

func getRemoteURL(repoPath string) (string, error) {
	gitPath, err := exec.LookPath("git")
	cmd := exec.Command(gitPath, "config", "--get", "remote.origin.url")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	remoteURL := strings.TrimSpace(string(output))
	return remoteURL, nil
}
