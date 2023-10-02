package versions

import (
	"encoding/json"
	"errors"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"path"
)

const (
	ConfigFileSchemaVersionFlagName = "schema-version"
	InitialConfigFileSchemaVersion  = "v1"
)

type VersionsFile struct {
	CliVersion              string `json:"cliVersion"`
	ConfigFileSchemaVersion string `json:"configFileSchemaVersion"`
}

func LoadCliVersionsFileContent(resources fs.FS) ([]byte, error) {
	bytes, err := fs.ReadFile(resources, path.Join("assets", "versions", "versions.json"))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GetCurrentCliConfigFileSchemaVersion(cliVersionsFileContent []byte) (string, error) {
	var versionsFile VersionsFile
	err := json.Unmarshal(cliVersionsFileContent, &versionsFile)
	if err != nil {
		return "", err
	}

	return versionsFile.ConfigFileSchemaVersion, nil
}

func CheckClientConfigFileSchemaMatchesCli(clientConfigFile string, resources fs.FS) error {
	//TODO: copied from flags.go, it would be nice to extract config file loading
	//if the default config file doesn't exist we can ignore the rest and return nil
	_, err := os.Stat(clientConfigFile)
	if err != nil && errors.Is(err, os.ErrNotExist) && clientConfigFile == defaults.ConfigFile {
		return nil
	}

	clientConfigFileContent, err := os.ReadFile(clientConfigFile)
	if err != nil {
		return err
	}

	cliVersionsFileContent, err := LoadCliVersionsFileContent(resources)
	if err != nil {
		return err
	}

	return ensureConfigFileSchemaVersionMatchesCli(clientConfigFileContent, cliVersionsFileContent)
}

func ensureConfigFileSchemaVersionMatchesCli(clientConfigFileContent []byte, cliVersionsFileContent []byte) error {
	currentCliConfigFileVersion, err := GetCurrentCliConfigFileSchemaVersion(cliVersionsFileContent)
	if err != nil {
		return err
	}

	clientConfigFile := make(map[string]string)
	err = yaml.Unmarshal(clientConfigFileContent, &clientConfigFile)
	if err != nil {
		return err
	}

	clientConfigFileVersion := ""
	if clientConfigFile[ConfigFileSchemaVersionFlagName] == "" && currentCliConfigFileVersion == InitialConfigFileSchemaVersion {
		clientConfigFileVersion = InitialConfigFileSchemaVersion
	} else {
		clientConfigFileVersion = clientConfigFile[ConfigFileSchemaVersionFlagName]
	}

	if currentCliConfigFileVersion != clientConfigFileVersion {
		return errors.New("client config file schema version does not match the one expected by the CLI")
	}

	return nil
}
