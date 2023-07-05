package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

func compareConfig(t *testing.T, appName string, writer io.Writer) {
	configTemplate := fmt.Sprintf(`app-name: %s
app-version: %s
cluster-endpoint: null
default-branch: null
dockerfile-name: %s
registry-user: null
repo-name: %s
runtime-version: null
`,
		appName,
		defaults.AppVersion,
		defaults.DockerfileName,
		defaults.RepoName,
	)

	result := fmt.Sprint(writer.(*bytes.Buffer))
	if configTemplate != result {
		t.Errorf("no match between\n%sand\n%s", configTemplate, result)
	}

}

func TestInitConfig(t *testing.T) {

	cli := CLI{
		Writer: new(bytes.Buffer),
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("KKA_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
	}

	if err := cli.Run([]string{"kka", "init", "config"}); err != nil {
		t.Error(err)
	}

	compareConfig(t, "cli", cli.Writer)
	// Config file is read correctly

	// Generate temporary file and add app-name parameter
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Errorf("creating temporary file %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.Write([]byte("app-name: FromFile")); err != nil {
		t.Errorf("writing config content %v", err)
	}

	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"kka", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromFile", cli.Writer)

	// Environment Variable wins over config
	os.Setenv("KKA_APP_NAME", "FromEnv")
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"kka", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromEnv", cli.Writer)

	// Command line argument wins over config and Environment variable
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"kka", fmt.Sprintf("--config-file=%s", f.Name()), "--app-name=FromParam", "init", "config"}); err != nil {
		t.Error(err)
	}
	compareConfig(t, "FromParam", cli.Writer)
}
