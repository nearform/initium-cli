package k8s

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/joho/godotenv"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
)

const (
	root      = "../../../"
	caCrt     = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNsakNDQVg0Q0NRRCtnL1BXUWJDSEFUQU5CZ2txaGtpRzl3MEJBUXNGQURBTk1Rc3dDUVlEVlFRR0V3SlYKVXpBZUZ3MHlNekExTURreE9EUTNNakZhRncweU5EQTFNRGd4T0RRM01qRmFNQTB4Q3pBSkJnTlZCQVlUQWxWVApNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTMzZXhaMldSUGtGeTJiQTQ4Z3RvClVwK0JkMFd2THFROU50WDBOS014YnA1SjEwMS9VaTVsY1RIcllDMWRqYTVwdjFKVGNGUW1jcWpBajBpL3dBakMKRm5oL1JKbVFrMHE1Y2liZWNURnA5UUFvRVNmbzJxYXovUTFmbmk4OG9ONlk1b1VGd2hrdTJ1bzNUWnV6M0JDMApNRnRyNXRDSGh1UzFObFhVT05VcHpzbW1UZzdZL0R1QXpNK3VIZS9qZlo1eHFqQUx3WHV6SkRFNkNCUGdhbHh6Ck9QM0V4QWNaMmRDWnRCREpVUnpTL29qMFYxOVdsRG5FK1FkTmtGTXlaMHN1UGxPTy95V0Ercno1byt2c2dKVGgKa2ZjaVZQSXZBOUdJdDdOcDZVdzU3MW1VOVlmd2x6MlU0Qm1xTXl6M05pZFF3bVd3V0wvay8yRHlXa0JDaFRBKwpvd0lEQVFBQk1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQWQ2M01qbFRJSWlSaVdOSDZvcDJySURjY1d5aWNKClJGbHJoODllSDVWeU04Q1o0NUhPYzVjbzZvRDVFdzQ3eG9vSlI2enZEd0c2anozelpFL3ArY2I5aGJ5dFJ5cysKaWNDd1g2dnRtbVpPN0M1RHdIMFYrUzk4emowNytmZFR3dUJHTDIxSlpZVmg1bFR6cEFpdU9iSkh6OTA5d1Y3OApObVhRSHpkMmtEZnpmTWhUaXpWZERPZEs3K2k1Q1RmaENIWUc4dDY2U3pmMGU5cWJ5eUFvTndwRnpxV01lRHROCkdJTGxBNHljcm1pYzBldUpmenZjeGk5NUVwMDdaZ1dYY3pINytLTWJtVnd2RkJWblZHdE9MZC9kMWhKaXlnU20KZUZERURxam9xL1JEcndCU2thbnlJMjNVa21uNGxkNDUvbHFVNDRSMVNqZXJaQmtNUmVIeTN0ZloKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	endpoint  = "https://127.0.0.1:6443"
	token     = "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklqWllaRFEwT0ZwUWVsRlJRVkp0TUMxa04wWTRNV3N6UzFwVmRXOVhSMDlvYzA5MWFESXdPVWhwVkhjaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpyZFdKbExYTjVjM1JsYlNJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKdmEyVXRhM1ZpWldOdmJtWnBaeTF6WVMxMGIydGxiaUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG01aGJXVWlPaUowWlhOMExYTmhJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpYSjJhV05sTFdGalkyOTFiblF1ZFdsa0lqb2lNelV4TkdJME9XVXRaVE15TWkwMFpHTTBMVGhsTUdVdE1qSTROemxsWm1GaE9HTXpJaXdpYzNWaUlqb2ljM2x6ZEdWdE9uTmxjblpwWTJWaFkyTnZkVzUwT210MVltVXRjM2x6ZEdWdE9uUmxjM1F0YzJFaWZRLlI5YnN3LU1MUjBubUhNOVJoUWJBZTBoZ2U0Z3JzYkJPa25RZUxVWC13SnRFY2dWejRDLVp6MFhVbUpPMEo4SVFmZXU1b3F3RnJwVHpSVEp0R24wdVdqc1RrSTZHSWNRNkpxM0FUSms3MEkwVzFqUTVJTkpJVjVmMFpfZDlIazNpZnVQaFNUUmpBZ2ljTDJCNjdMUHVBaW40T05hOHNkTE95VTZrSXFDU2Q3dURUNVMtLS1qQ0JzTnpJd1p0QVg1dVFfbDZGUHVUdlZxdGJiSUJGblVQUGExandaSWwwd3U1YV9DTG9rVFhWbDduVkZYUzBxZTJ1RWFmUlEwVGZIOWFjbWZndkFVcWFXbkNsTV9oekNzekhKOTBYaHNtOXIzQ3oya29UcDJhLWtVVWlhSWQxbGlPeWcyNFMtS29TZlNwUEpiNWcxVzZhNHB5VU1xMXNDNGFuUQ=="
	namespace = "custom"
	commitSha = "93f4be93"
)

var (
	envTestFile          = path.Join(root, "assets/testdata/.env.initium")
	secretRefEnvTestFile = path.Join(root, "assets/testdata/.env.secretref.initium")
	proj                 = &project.Project{Name: "knative_test",
		Directory: path.Join(root, "example"),
		Resources: os.DirFS(root),
		IsPrivate: false,
	}
	dockerImage = docker.DockerImage{
		Registry:  "example.com",
		Directory: ".",
		Name:      "test",
		Tag:       "v1.1.0",
	}
)

func TestConfig(t *testing.T) {
	decodedCert, err := base64.StdEncoding.DecodeString(caCrt)

	if err != nil {
		t.Fatalf("Not possible to decode base64 cert into a string: %v", err)
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		t.Fatalf("Not possible to decode base64 token into a string: %v", err)
	}

	_, err = Config(endpoint, fmt.Sprintf("%x", decodedToken), []byte(decodedCert))

	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	mockDummyValues := map[string]string{
		"caCrt":    "certificatestring",
		"endpoint": "endpoint",
		"token":    "tokenstring",
	}

	_, err = Config(mockDummyValues["endpoint"], mockDummyValues["token"], []byte(mockDummyValues["caCrt"]))

	if err == nil {
		t.Fatalf("Error: strings shouldn't be not supported")
	}

}

func TestLoadManifestForPublicService(t *testing.T) {
	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, envTestFile, secretRefEnvTestFile)

	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	annotations := serviceManifest.Spec.Template.ObjectMeta.Annotations
	assert.Assert(t, annotations[UpdateTimestampAnnotationName] != "", "Missing %s annotation", UpdateTimestampAnnotationName)
	assert.Assert(t, annotations[UpdateShaAnnotationName] == commitSha, "Expected %s SHA, got %s", commitSha, annotations[UpdateShaAnnotationName])

	labels := serviceManifest.GetLabels()
	_, ok := labels[visibilityLabel]
	assert.Assert(t, !ok, "Visibility label should not be set for public services")
}

func TestLoadManifestForPrivateService(t *testing.T) {
	imagePullSecrets := []string{"secretPassword123"}

	privateProj := proj
	privateProj.IsPrivate = true
	privateProj.ImagePullSecrets = imagePullSecrets

	dockerImage := docker.DockerImage{
		Registry:  "example.com",
		Directory: ".",
		Name:      "test",
		Tag:       "v1.1.0",
	}

	serviceManifest, err := LoadManifest(namespace, commitSha, privateProj, dockerImage, envTestFile, secretRefEnvTestFile)

	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	annotations := serviceManifest.Spec.Template.ObjectMeta.Annotations
	pullSecret := serviceManifest.Spec.Template.Spec.ImagePullSecrets[0].Name
	assert.Assert(t, annotations[UpdateTimestampAnnotationName] != "", "Missing %s annotation", UpdateTimestampAnnotationName)
	assert.Assert(t, annotations[UpdateShaAnnotationName] == commitSha, "Expected %s SHA, got %s", commitSha, annotations[UpdateShaAnnotationName])
	assert.Assert(t, pullSecret == imagePullSecrets[0], "Expected secret value to be %s, got %s", imagePullSecrets, pullSecret)

	labels := serviceManifest.GetLabels()
	assert.Assert(t, labels[visibilityLabel] == visibilityLabelPrivateValue)
}

func TestLoadManifestEnvironmentVariables(t *testing.T) {
	envVariablesFromFile, err := godotenv.Read(envTestFile)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %v", err))
	}
	secretRefEnvVariablesFromFile, err := godotenv.Read(secretRefEnvTestFile)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, envTestFile, secretRefEnvTestFile)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	for _, envVar := range serviceManifest.Spec.Template.Spec.Containers[0].Env {
		delete(envVariablesFromFile, envVar.Name)
		delete(secretRefEnvVariablesFromFile, envVar.Name)
	}
	assert.Assert(t, len(envVariablesFromFile) == 0, "Missing environment variables: %s", envVariablesFromFile )
	assert.Assert(t, len(secretRefEnvVariablesFromFile) == 0, "Missing secret environment variables: %s", secretRefEnvVariablesFromFile)
}

func TestLoadManifestEnvironmentVariablesInvalidFormat(t *testing.T) {
	invalidEnvTestFile := path.Join(root, "assets/testdata/.env.initium.invalid")
	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, invalidEnvTestFile, secretRefEnvTestFile)
	assert.Assert(t, err != nil && strings.Contains(err.Error(), "Error loading .env file"), "There should be a validation error when missing a mandatory character" )
	assert.Assert(t, serviceManifest == nil, "Expected nil manifest, got %v", serviceManifest)
}

func TestLoadManifestSecretRefEnvironmentMandatoryChars(t *testing.T) {
	invalidSecretRefEnvTestFile := path.Join(root, "assets/testdata/.env.secretref.initium.invalid2")
	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, envTestFile, invalidSecretRefEnvTestFile)
	assert.Assert(t, err != nil && strings.Contains(err.Error(), "Value must be in the format <secret-name>/<secret-key>"), "There should be a validation error when missing a mandatory character" )
	assert.Assert(t, serviceManifest == nil, "Expected nil manifest, got %v", serviceManifest)
}

func TestLoadManifestSecretRefEnvironmentConflict(t *testing.T) {
	invalidSecretRefEnvTestFile := path.Join(root, "assets/testdata/.env.secretref.conflictingvar")
	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, envTestFile, invalidSecretRefEnvTestFile)
	assert.Assert(t, err != nil && strings.Contains(err.Error(), "Conflicting environment variable"), "There should be a validation error when missing a mandatory character" )
	assert.Assert(t, serviceManifest == nil, "Expected nil manifest, got %v", serviceManifest)
}

func TestLoadManifestSecretEnvironmentVariablesFileDoesNotExist(t *testing.T) {
	nonExistingEnvTestFile := "idontexist"
	nonExistingSecretRefEnvTestFile := "idontexist"
	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, nonExistingEnvTestFile, nonExistingSecretRefEnvTestFile)
	assert.Assert(t, serviceManifest != nil, "Expected maifest to be created without issues. Dotenv files are optional. Err: %s", err)
}

func TestLoadManifestSecretEnvironmentVariablesEmptyFile(t *testing.T) {
	emtpyFile := path.Join(root, "assets/testdata/.env.initium.empty")
	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, emtpyFile, emtpyFile)
	assert.Assert(t, serviceManifest != nil, "Expected maifest to be created without issues. Dotenv files are optional. Err: %s", err)
}
