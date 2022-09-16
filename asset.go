package main

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v47/github"
)

func downloadAsset(ctx context.Context, asset *github.ReleaseAsset, path string) (err error) {
	client, org, repo, err := getGithubClient()
	if err != nil {
		return err
	}

	reader, _, err := client.Repositories.DownloadReleaseAsset(ctx, org, repo, *asset.ID, http.DefaultClient)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return err
}

func computeSha(tarballPath string) (sha512sum string, err error) {
	file, err := os.Open(tarballPath)
	if err != nil {
		return sha512sum, err
	}
	defer file.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return sha512sum, err
	}

	sha512sum = hex.EncodeToString(hash.Sum(nil))

	return sha512sum, err
}

func verifySha(shaPath string, tarballPath string) (err error) {
	sha512sum, err := computeSha(tarballPath)
	if err != nil {
		return err
	}

	shaBytes, err := os.ReadFile(shaPath)
	if err != nil {
		return err
	}

	shaSplit := strings.Split(string(shaBytes), " ")
	if len(shaSplit) < 1 {
		return errors.New("couldn't determine sha512sum from file")
	}

	if sha512sum != shaSplit[0] {
		return errors.New("sha512sum mismatch")
	}

	return err
}

func installTarGzAsset(source string, destination string) (err error) {
	if !strings.Contains(source, ".tar.gz") { // skip non-tar.gz
		return err
	}

	cmd := exec.Command("/usr/bin/tar", "xf", source, "-C", destination)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}
