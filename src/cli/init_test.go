package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/services/project"
)

func compareConfig(t *testing.T, appName string, registry string, isPrivateService bool, projectType project.ProjectType, writer io.Writer) {
	configTemplate := fmt.Sprintf(`app-name: %s
container-registry: %s
default-branch: main
dockerfile-name: null
env-var-file: .env.initium
image-pull-secrets: null
private: %t
project-type: %s
runtime-version: null
`,
		appName,
		registry,
		isPrivateService,
		projectType,
	)

	result := fmt.Sprint(writer.(*bytes.Buffer))

	if configTemplate != result {
		t.Errorf("no match between\n%sand\n%s", configTemplate, result)
	}
}

func geticliForTesting(resources fs.FS) icli {
	return NewWithOptions(
		resources,
		log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("INITIUM_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
		new(bytes.Buffer),
		Release{
			Version: "1",
			Date:    "today",
			Commit:  "unknown",
		},
	)
}

func reseticliBuffer(c *icli) {
	c.Writer = new(bytes.Buffer)
}

func icliOutput(c icli) string {
	return fmt.Sprint(c.Writer.(*bytes.Buffer))
}

func TestInitConfig(t *testing.T) {
	icli := geticliForTesting(os.DirFS("../.."))
	// Config file is read correctly

	// Generate temporary file and add app-name parameter
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Errorf("creating temporary file %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	registry := "ghcr.io/nearform"

	if _, err := f.WriteString("app-name: FromFile\ncontainer-registry: " + registry + "\nproject-type: go"); err != nil {
		t.Errorf("writing config content %v", err)
	}

	reseticliBuffer(&icli)
	if err = icli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromFile", registry, false, "go", icli.Writer)

	// Environment Variable wins over config
	os.Setenv("INITIUM_APP_NAME", "FromEnv")
	os.Setenv("INITIUM_PROJECT_TYPE", "go")
	defer os.Unsetenv("INITIUM_PROJECT_TYPE")
	defer os.Unsetenv("INITIUM_APP_NAME") // Unset the environment variable at the end
	reseticliBuffer(&icli)
	if err = icli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromEnv", registry, false, "go", icli.Writer)

	// Command line argument wins over config and Environment variable
	reseticliBuffer(&icli)
	if err = icli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam", "--project-type=go"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", registry, false, "go", icli.Writer)
}

func TestRepoNameRetrocompatibiliy(t *testing.T) {
	cli := geticliForTesting(os.DirFS("../.."))

	// Generate temporary file and add repo-name parameter
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Errorf("creating temporary file %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.WriteString("repo-name: FromFile\nproject-type: go"); err != nil {
		t.Errorf("writing config content %v", err)
	}

	reseticliBuffer(&cli)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam"}); err != nil {
		t.Error(err)
	}

	compareConfig(t, "FromParam", "FromFile", false, "go", cli.Writer)

	//Override from parameter
	reseticliBuffer(&cli)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam", "--container-registry=ghcr.io/nearform"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", "ghcr.io/nearform", false, "go", cli.Writer)
}

func TestAppName(t *testing.T) {
	cli := geticliForTesting(os.DirFS("../.."))

	err := cli.Run([]string{"initium", "build"})
	if err == nil {
		t.Errorf("CLI should ask for %s and %s if not detected", appNameFlag, repoNameFlag)
	}

	if !(strings.Contains(err.Error(), appNameFlag) && strings.Contains(err.Error(), repoNameFlag)) {
		t.Errorf("the error message should contain %s and %s", appNameFlag, repoNameFlag)
	}
}
