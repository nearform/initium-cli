package local

import (
    "os"
    "k8s-kurated-addons.cli/src/utils/defaults"
    "k8s-kurated-addons.cli/src/utils/logger"
    "k8s-kurated-addons.cli/src/utils/bindata"
)
type LocalService struct {
    HasDockerfile bool // Flag will tell the application if it has a Dockerfile
}


// Create the Dockerfile
func (ls LocalService) CreateDockerfile(dockerFilePath string) {
    logger.PrintInfo("Creating Dockerfile under " + dockerFilePath + "/" + defaults.DefaultDockerfileName)
    data, err := bindata.Asset(dockerFilePath + "/" + defaults.DefaultDockerfileName)
    if err != nil {
        logger.PrintError("Dockerfile asset not found", err)
    }

    err = os.WriteFile(dockerFilePath + "/" + defaults.DefaultDockerfileName, []byte(data), 0644)

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