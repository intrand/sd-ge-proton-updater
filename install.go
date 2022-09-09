package main

import (
	"context"
	"errors"
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

	org, repo, err = splitGithubUrl(ProtonGeUrl)
	if err != nil {
		return client, org, repo, err
	}

	return client, org, repo, err
}

func getLatestReleases(ctx context.Context) (releases []*github.RepositoryRelease, err error) {
	client := github.NewClient(nil)

	org, repo, err := splitGithubUrl(ProtonGeUrl)
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
