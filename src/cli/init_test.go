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
)

func compareConfig(t *testing.T, appName string, registry string, writer io.Writer) {
	configTemplate := fmt.Sprintf(`app-name: %s
cluster-endpoint: null
container-registry: %s
default-branch: main
dockerfile-name: null
env-var-file: .env.initium
registry-user: null
runtime-version: null
`,
		appName,
		registry,
	)

	result := fmt.Sprint(writer.(*bytes.Buffer))

	if configTemplate != result {
		t.Errorf("no match between\n%sand\n%s", configTemplate, result)
	}
}

func GeticliForTesting(resources fs.FS) icli {
	return NewWithOptions(
		resources,
		log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("INITIUM_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
		new(bytes.Buffer),
	)
}

func TestInitConfig(t *testing.T) {
	icli := GeticliForTesting(os.DirFS("../.."))
	// Config file is read correctly

	// Generate temporary file and add app-name parameter
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Errorf("creating temporary file %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	registry := "ghcr.io/nearform"

	if _, err := f.WriteString("app-name: FromFile\ncontainer-registry: " + registry); err != nil {
		t.Errorf("writing config content %v", err)
	}

	icli.Writer = new(bytes.Buffer)
	if err = icli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromFile", registry, icli.Writer)

	// Environment Variable wins over config
	os.Setenv("INITIUM_APP_NAME", "FromEnv")
	defer os.Unsetenv("INITIUM_APP_NAME") // Unset the environment variable at the end
	icli.Writer = new(bytes.Buffer)
	if err = icli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromEnv", registry, icli.Writer)

	// Command line argument wins over config and Environment variable
	icli.Writer = new(bytes.Buffer)
	if err = icli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", registry, icli.Writer)

}

func TestRepoNameRetrocompatibiliy(t *testing.T) {
	cli := GeticliForTesting(os.DirFS("../.."))

	// Generate temporary file and add repo-name parameter
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Errorf("creating temporary file %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.WriteString("repo-name: FromFile"); err != nil {
		t.Errorf("writing config content %v", err)
	}

	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", "FromFile", cli.Writer)

	//Override from parameter
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam", "--container-registry=ghcr.io/nearform"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", "ghcr.io/nearform", cli.Writer)
}

func TestAppName(t *testing.T) {
	cli := GeticliForTesting(os.DirFS("../.."))

	err := cli.Run([]string{"initium", "build"})
	if err == nil {
		t.Errorf("CLI should ask for %s and %s if not detected", appNameFlag, repoNameFlag)
	}

	if !(strings.Contains(err.Error(), appNameFlag) && strings.Contains(err.Error(), repoNameFlag)) {
		t.Errorf("the error message should contain %s and %s", appNameFlag, repoNameFlag)
	}
}
