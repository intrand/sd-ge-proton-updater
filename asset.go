package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	r, err := os.Open(source)
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next() // go to next file
		if err == io.EOF {              // we ran out of files
			break // stop trying to process anything
		}
		if err != nil { // all bad conditions
			return err
		}

		filePath := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filePath, 0700); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(filePath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close() // outFile.Close error omitted as Copy error is more interesting at this point
				return err
			}

			if err := outFile.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink:
			base, _ := filepath.Split(filePath)            // get base dir for the link
			target := filepath.Join(base, header.Linkname) // start in the same directory as symlink, but the filename is the target of the symlink, not the symlink itself
			err = os.Symlink(target, filePath)
			if err != nil {
				return err
			}
		case tar.TypeLink:
			base, _ := filepath.Split(filePath)            // get base dir for the link
			target := filepath.Join(base, header.Linkname) // start in the same directory as symlink, but the filename is the target of the symlink, not the symlink itself
			err = os.Link(target, filePath)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown file type: %s in %s", string(header.Typeflag), filePath)
		}
	}

	return err
}
