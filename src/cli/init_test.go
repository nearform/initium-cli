package cli

import (
	"bytes"
	"fmt"
	"io"
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

func getCLI() CLI {
	return CLI{
		Writer: new(bytes.Buffer),
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("INITIUM_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
	}
}

func TestInitConfig(t *testing.T) {
	cli := getCLI()
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

	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromFile", registry, cli.Writer)

	// Environment Variable wins over config
	os.Setenv("INITIUM_APP_NAME", "FromEnv")
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromEnv", registry, cli.Writer)

	// Command line argument wins over config and Environment variable
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"initium", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config", "--app-name=FromParam"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", registry, cli.Writer)

}

func TestRepoNameRetrocompatibiliy(t *testing.T) {
	cli := getCLI()

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
