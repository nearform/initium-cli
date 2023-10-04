package versions

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"gotest.tools/v3/assert"
	"testing"
)

type ClientConfigFile struct {
	SchemaVersion string `yaml:"schema-version"`
	AppName       string `yaml:"app-name"`
}

const cliVersion = "v0.5.0"

func TestShouldReturnCurrentCliConfigFileSchemaVersion(t *testing.T) {
	cliConfigFileSchemaVersion := InitialConfigFileSchemaVersion
	bytes, err := getCliVersionsFileContent(cliConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	schemaVersion, err := GetCurrentCliConfigFileSchemaVersion(bytes)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	assert.Assert(t, schemaVersion == cliConfigFileSchemaVersion, fmt.Sprintf("Expected: %s, got: %s", cliConfigFileSchemaVersion, schemaVersion))
}

func TestShouldReturnCliVersion(t *testing.T) {
	cliConfigFileSchemaVersion := InitialConfigFileSchemaVersion
	cliVersionsFileContent, err := getCliVersionsFileContent(cliConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	version, err := GetCliVersion(cliVersionsFileContent)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	assert.Assert(t, version == cliVersion, fmt.Sprintf("Expected: %s, got: %s", cliVersion, version))
}

func TestShouldBeBackwardCompatibleWhenClientConfigIsMissingSchemaVersion(t *testing.T) {
	cliConfigFileSchemaVersion := InitialConfigFileSchemaVersion
	cliVersionsFileContent, err := getCliVersionsFileContent(cliConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	clientConfigFileContent, err := getClientConfigFileContent("")
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	err = ensureConfigFileSchemaVersionMatchesCli(clientConfigFileContent, cliVersionsFileContent)
	assert.NilError(t, err, "Expected error not present, got: %v", err)
}

func TestShouldAcceptTheClientConfigFileWhenSchemaVersionMatches(t *testing.T) {
	cliConfigFileSchemaVersion := InitialConfigFileSchemaVersion
	cliVersionsFileContent, err := getCliVersionsFileContent(cliConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	clientConfigFileContent, err := getClientConfigFileContent(InitialConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	err = ensureConfigFileSchemaVersionMatchesCli(clientConfigFileContent, cliVersionsFileContent)
	assert.NilError(t, err, "Expected error not present, got: %v", err)
}

func TestShouldReturnAnErrorWhenThereIsAMismatchBetweenSchemaVersions(t *testing.T) {
	cliConfigFileSchemaVersion := "v2"
	cliVersionsFileContent, err := getCliVersionsFileContent(cliConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	clientConfigFileContent, err := getClientConfigFileContent(InitialConfigFileSchemaVersion)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	err = ensureConfigFileSchemaVersionMatchesCli(clientConfigFileContent, cliVersionsFileContent)
	assert.Error(t, err, "client config file schema version does not match the one expected by the CLI")
}

func getCliVersionsFileContent(schemaVersion string) ([]byte, error) {
	versionsFile := VersionsFile{
		CliVersion:              cliVersion,
		ConfigFileSchemaVersion: schemaVersion,
	}
	return json.Marshal(versionsFile)
}

func getClientConfigFileContent(schemaVersion string) ([]byte, error) {
	if schemaVersion == "" {
		return yaml.Marshal(ClientConfigFile{
			AppName: "initium-nodejs-demo-app",
		})
	} else {
		return yaml.Marshal(ClientConfigFile{
			SchemaVersion: schemaVersion,
			AppName:       "initium-nodejs-demo-app",
		})
	}
}
