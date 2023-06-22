package cli

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

func TestInitConfig(t *testing.T) {
	configTemplate := fmt.Sprintf(`app-name: %%s
app-version: %s
default-branch: null
dockerfile-name: %s
endpoint: null
registry-user: null
repo-name: %s
runtime-version: null
`,
		defaults.AppVersion,
		defaults.DockerfileName,
		defaults.RepoName,
	)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory")
	}

	cli := CLI{
		CWD:    cwd,
		Writer: new(bytes.Buffer),
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("KKA_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
	}

	if err = cli.Run([]string{"kka", "init", "config"}); err != nil {
		t.Error(err)
	}
	expected := fmt.Sprintf(configTemplate, "cli")
	result := fmt.Sprint(cli.Writer.(*bytes.Buffer))
	if expected != result {
		t.Errorf("no match between %s and %s", expected, result)
	}

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

	expected = fmt.Sprintf(configTemplate, "FromFile")
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"kka", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	result = fmt.Sprint(cli.Writer.(*bytes.Buffer))
	if expected != result {
		t.Errorf("no match between %s\n and\n %s", expected, result)
	}

	// Environment Variable wins over config
	os.Setenv("KKA_APP_NAME", "FromEnv")
	expected = fmt.Sprintf(configTemplate, "FromEnv")
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"kka", fmt.Sprintf("--config-file=%s", f.Name()), "init", "config"}); err != nil {
		t.Error(err)
	}
	result = fmt.Sprint(cli.Writer.(*bytes.Buffer))
	if expected != result {
		t.Errorf("no match between %s\n and\n %s", expected, result)
	}

	// Command line argument wins over config and Environment variable
	expected = fmt.Sprintf(configTemplate, "FromParam")
	cli.Writer = new(bytes.Buffer)
	if err = cli.Run([]string{"kka", fmt.Sprintf("--config-file=%s", f.Name()), "--app-name=FromParam", "init", "config"}); err != nil {
		t.Error(err)
	}

	result = fmt.Sprint(cli.Writer.(*bytes.Buffer))

	if expected != result {
		t.Errorf("no match between %s\n and\n %s", expected, result)
	}
}
