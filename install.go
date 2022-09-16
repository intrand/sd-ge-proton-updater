package main

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v47/github"
)

// split [https://|http://|git://]github.com/owner/repo[.git]/... into owner, repo
func splitGithubUrl(url string) (owner string, repo string, err error) {
	str := url
	str = strings.TrimPrefix(str, "https://")
	str = strings.TrimPrefix(str, "http://")
	str = strings.TrimPrefix(str, "git://")

	if strings.HasPrefix(str, "github.com/") {
		str = strings.TrimPrefix(str, "github.com/")
	} else {
		err = errors.New("not a github url")
		return
	}

	fields := strings.Split(str, "/")
	if len(fields) < 2 {
		err = errors.New("invalid format for github")
		return
	}

	owner = fields[0]
	repo = strings.TrimRight(fields[1], ".git")

	return
}

func getGithubClient() (client *github.Client, org string, repo string, err error) {
	client = github.NewClient(nil)

	org, repo, err = splitGithubUrl(protonGeUrl)
	if err != nil {
		return client, org, repo, err
	}

	return client, org, repo, err
}

func getLatestReleases(ctx context.Context) (releases []*github.RepositoryRelease, err error) {
	client := github.NewClient(nil)

	org, repo, err := splitGithubUrl(protonGeUrl)
	if err != nil {
		return releases, err
	}

	allReleases, _, err := client.Repositories.ListReleases(ctx, org, repo, nil)
	if err != nil {
		return releases, err
	}

	if len(allReleases) < 1 {
		return releases, errors.New("no releases found")
	}

	var stableReleases []*github.RepositoryRelease
	for _, unfilteredRelease := range allReleases {
		if *unfilteredRelease.Prerelease { // omit pre-releases
			continue
		}
		if *unfilteredRelease.Draft { // omit drafts
			continue
		}
		stableReleases = append(stableReleases, unfilteredRelease)
	}

	if len(stableReleases) < 2 {
		return releases, errors.New("not enough stable releases found (possibly 0)")
	}

	releases = []*github.RepositoryRelease{
		stableReleases[0],
		stableReleases[1],
	}

	return releases, err
}

func install() (err error) {
	ctx := context.Background()             // create context
	releases, err := getLatestReleases(ctx) // get latest stable releases
	if err != nil {
		return err
	}

	latestRelease := releases[0]
	// latestMinusOneRelease := releases[1]

	latestPath := filepath.Join(protonPath, *latestRelease.TagName)

	exist, err := exists(latestPath)
	if err != nil {
		return err
	}

	if exist {
		log.Println("Release " + *latestRelease.TagName + " already exists on this console. Nothing to do. Exiting.")
		return err
	} // end exists

	var shaAsset *github.ReleaseAsset
	var tarballAsset *github.ReleaseAsset
	for _, asset := range latestRelease.Assets {
		if strings.Contains(*asset.Name, ".sha512sum") {
			shaAsset = asset
		}
		if strings.Contains(*asset.Name, ".tar.gz") {
			tarballAsset = asset
		}
	} // end list of assets

	var naked *github.ReleaseAsset
	if shaAsset == naked || tarballAsset == naked { // check we got data for both
		return errors.New("couldn't get enough info about releases. Did something change?")
	}

	dir, err := mkTempDir(*latestRelease.TagName)
	if err != nil {
		return err
	}

	var shaPath string = filepath.Join(dir, *shaAsset.Name)
	var tarballPath string = filepath.Join(dir, *tarballAsset.Name)

	// download SHA-512 checksum
	err = downloadAsset(ctx, shaAsset, shaPath)
	if err != nil {
		return err
	}

	// download GE-Proton tarball distribution
	err = downloadAsset(ctx, tarballAsset, tarballPath)
	if err != nil {
		return err
	}

	// verify SHA of tarball against SHA-512 checksum file (which is not signed!)
	err = verifySha(shaPath, tarballPath)
	if err != nil {
		return err
	}

	err = installTarGzAsset(tarballPath, protonPath)
	if err != nil {
		return err
	}

	log.Println("Successfully installed: " + *latestRelease.TagName)
	return err
}
