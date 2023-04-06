package local

import (
    "os"
    "k8s-kurated-addons.cli/src/utils/defaults"
    "k8s-kurated-addons.cli/src/utils/logger"
    "k8s-kurated-addons.cli/src/utils/bindata"

    "github.com/a8m/envsubst"
)


type LocalService struct {
    HasDockerfile bool // Flag will tell the application if it has a Dockerfile
    HasManifest bool
}

// Create the Dockerfile
func (ls LocalService) CreateDockerfile(dockerFilePath string) {
    logger.PrintInfo("Creating Dockerfile under " + dockerFilePath + "/" + defaults.DefaultDockerfileName)

    data := ls.readAsset(dockerFilePath + "/" + defaults.DefaultDockerfileName)

    err := os.WriteFile(dockerFilePath + "/" + defaults.DefaultDockerfileName, []byte(data), 0644)

    if err != nil {
        logger.PrintError("Error writing to file", err)
    }
}

// Create manifest
func (ls LocalService) CreateManifest(appPort string, appName string, repoName string) {
    logger.PrintInfo("Creating manifest under " + defaults.DefaultManifestPath)

    data := ls.readAsset(defaults.DefaultManifestPath)

    os.Setenv("KKA_APP_NAME", appName)
    os.Setenv("KKA_APP_REGISTRY", repoName)
    os.Setenv("KKA_APP_PORT", appPort)

    buf, err := envsubst.Bytes(data)

    if err != nil {
        logger.PrintError("Error passing environment variables", err)
    }

    logger.PrintInfo(string(buf))

    err = os.WriteFile(defaults.DefaultManifestPath, []byte(buf), 0644)

    if err != nil {
        logger.PrintError("Error writing to file", err)
    }
}

// Remove the Dockerfile
func (ls LocalService) RemoveDockerfile(dockerFilePath string) {
    logger.PrintInfo("Removing Dockerfile")
    err := os.Remove(dockerFilePath + "/" + defaults.DefaultDockerfileName)

    if err != nil {
        logger.PrintError("Error removing Dockerfile", err)
    }
}

func (ls LocalService) readAsset(asset string) []byte {
    data, err := bindata.Asset(asset)
    if err != nil {
        logger.PrintError(asset + " not found", err)
    }

    return data
}