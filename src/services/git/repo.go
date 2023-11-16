package git

import (
	"fmt"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
)

const (
	httpgithubprefix = "https://github.com/"
	gitgithubprefix  = "git@github.com:"
)

func initRepo() (*git.Repository, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(wd) // path to the repository
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetHash() (string, error) {
	repo, err := initRepo()

	if err != nil {
		return "", fmt.Errorf("error getting git repo %v", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("cannot get HEAD reference: %v", err)
	}

	return headRef.Hash().String(), nil
}

func GetBranchName() (string, error) {
	repo, err := initRepo()

	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	branchName := head.Name().Short()
	return branchName, nil
}

func getGithubRemote() (string, error) {
	repo, err := initRepo()

	if err != nil {
		return "", fmt.Errorf("error getting git repo %v", err)
	}

	c, err := repo.Config()
	if err != nil {
		return "", fmt.Errorf("error getting the repo config: %v", err)
	}

	for _, remote := range c.Remotes {
		for _, url := range remote.URLs {
			// HTTPS (https://github.com/organization/repo.git)
			if strings.HasPrefix(url, httpgithubprefix) {
				return strings.Replace(url, httpgithubprefix, "", 1), nil
			}
			// SSH (git@github.com:organization/repo.git)
			if strings.HasPrefix(url, gitgithubprefix) {
				return strings.Replace(url, gitgithubprefix, "", 1), nil
			}
		}
	}

	return "", fmt.Errorf("no github remote found")
}

func GetRepoName() (string, error) {
	remote, err := getGithubRemote()

	if err != nil {
		return "", err
	}

	splitRemote := strings.Split(remote, "/")
	return strings.Replace(splitRemote[1], ".git", "", 1), nil
}

func GetGithubOrg() (string, error) {
	remote, err := getGithubRemote()

	if err != nil {
		return "", err
	}

	splitRemote := strings.Split(remote, "/")
	return splitRemote[0], nil
}

func PublishCommentPRGithub (url string) {
	fmt.Printf("Application URL: %s", url)
	commitSha, err := GetHash()
	fmt.Printf("Commit hash: %s", commitSha)
	fmt.Printf("Timestamp: %v", time.Now())
}
