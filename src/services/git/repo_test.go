package git

import (
	"os"
	"testing"
)

func TestGetRepoName(t *testing.T) {
	_, err := GetRepoName()

	if err == nil {
		t.Error("initRepo should run on the root folder of a git repo")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	if err = os.Chdir("../../.."); err != nil {
		t.Error(err)
	}

	name, err := GetRepoName()
	if err != nil {
		t.Error(err)
	}

	if name != "initium-cli" {
		t.Error("Repo name should match git repo name")
	}

	if err = os.Chdir(cwd); err != nil {
		t.Error(err)
	}
}

func TestGetGithubOrg(t *testing.T) {
	_, err := GetGithubOrg()

	if err == nil {
		t.Error("initRepo should run on the root folder of a git repo")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	if err = os.Chdir("../../.."); err != nil {
		t.Error(err)
	}

	name, err := GetGithubOrg()
	if err != nil {
		t.Error(err)
	}

	if name != "nearform" {
		t.Error("Org doesn't match the git repo org")
	}

	if err = os.Chdir(cwd); err != nil {
		t.Error(err)
	}
}
