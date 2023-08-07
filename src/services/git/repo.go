package git

import (
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
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

func GetGithubOrg() (string, error) {
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
			if strings.HasPrefix(url, "https://github.com/") {
				splitURL := strings.Split(url, "/")
				if len(splitURL) > 3 {
					return splitURL[3], nil
				}
			}
			// SSH (git@github.com:organization/repo.git)
			if strings.HasPrefix(url, "git@github.com:") {
				splitURL := strings.Split(url, ":")
				if len(splitURL) > 1 {
					splitPath := strings.Split(splitURL[1], "/")
					if len(splitPath) > 1 {
						return splitPath[0], nil
					}
				}
			}
		}

	}

	return "", fmt.Errorf("cannot find any remote with github.com")
}
