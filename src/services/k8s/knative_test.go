package k8s

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
)

var root = "../../../"

var caCrt = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNsakNDQVg0Q0NRRCtnL1BXUWJDSEFUQU5CZ2txaGtpRzl3MEJBUXNGQURBTk1Rc3dDUVlEVlFRR0V3SlYKVXpBZUZ3MHlNekExTURreE9EUTNNakZhRncweU5EQTFNRGd4T0RRM01qRmFNQTB4Q3pBSkJnTlZCQVlUQWxWVApNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTMzZXhaMldSUGtGeTJiQTQ4Z3RvClVwK0JkMFd2THFROU50WDBOS014YnA1SjEwMS9VaTVsY1RIcllDMWRqYTVwdjFKVGNGUW1jcWpBajBpL3dBakMKRm5oL1JKbVFrMHE1Y2liZWNURnA5UUFvRVNmbzJxYXovUTFmbmk4OG9ONlk1b1VGd2hrdTJ1bzNUWnV6M0JDMApNRnRyNXRDSGh1UzFObFhVT05VcHpzbW1UZzdZL0R1QXpNK3VIZS9qZlo1eHFqQUx3WHV6SkRFNkNCUGdhbHh6Ck9QM0V4QWNaMmRDWnRCREpVUnpTL29qMFYxOVdsRG5FK1FkTmtGTXlaMHN1UGxPTy95V0Ercno1byt2c2dKVGgKa2ZjaVZQSXZBOUdJdDdOcDZVdzU3MW1VOVlmd2x6MlU0Qm1xTXl6M05pZFF3bVd3V0wvay8yRHlXa0JDaFRBKwpvd0lEQVFBQk1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQWQ2M01qbFRJSWlSaVdOSDZvcDJySURjY1d5aWNKClJGbHJoODllSDVWeU04Q1o0NUhPYzVjbzZvRDVFdzQ3eG9vSlI2enZEd0c2anozelpFL3ArY2I5aGJ5dFJ5cysKaWNDd1g2dnRtbVpPN0M1RHdIMFYrUzk4emowNytmZFR3dUJHTDIxSlpZVmg1bFR6cEFpdU9iSkh6OTA5d1Y3OApObVhRSHpkMmtEZnpmTWhUaXpWZERPZEs3K2k1Q1RmaENIWUc4dDY2U3pmMGU5cWJ5eUFvTndwRnpxV01lRHROCkdJTGxBNHljcm1pYzBldUpmenZjeGk5NUVwMDdaZ1dYY3pINytLTWJtVnd2RkJWblZHdE9MZC9kMWhKaXlnU20KZUZERURxam9xL1JEcndCU2thbnlJMjNVa21uNGxkNDUvbHFVNDRSMVNqZXJaQmtNUmVIeTN0ZloKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="

var endpoint = "https://127.0.0.1:6443"

var token = "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklqWllaRFEwT0ZwUWVsRlJRVkp0TUMxa04wWTRNV3N6UzFwVmRXOVhSMDlvYzA5MWFESXdPVWhwVkhjaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpyZFdKbExYTjVjM1JsYlNJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKdmEyVXRhM1ZpWldOdmJtWnBaeTF6WVMxMGIydGxiaUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG01aGJXVWlPaUowWlhOMExYTmhJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpYSjJhV05sTFdGalkyOTFiblF1ZFdsa0lqb2lNelV4TkdJME9XVXRaVE15TWkwMFpHTTBMVGhsTUdVdE1qSTROemxsWm1GaE9HTXpJaXdpYzNWaUlqb2ljM2x6ZEdWdE9uTmxjblpwWTJWaFkyTnZkVzUwT210MVltVXRjM2x6ZEdWdE9uUmxjM1F0YzJFaWZRLlI5YnN3LU1MUjBubUhNOVJoUWJBZTBoZ2U0Z3JzYkJPa25RZUxVWC13SnRFY2dWejRDLVp6MFhVbUpPMEo4SVFmZXU1b3F3RnJwVHpSVEp0R24wdVdqc1RrSTZHSWNRNkpxM0FUSms3MEkwVzFqUTVJTkpJVjVmMFpfZDlIazNpZnVQaFNUUmpBZ2ljTDJCNjdMUHVBaW40T05hOHNkTE95VTZrSXFDU2Q3dURUNVMtLS1qQ0JzTnpJd1p0QVg1dVFfbDZGUHVUdlZxdGJiSUJGblVQUGExandaSWwwd3U1YV9DTG9rVFhWbDduVkZYUzBxZTJ1RWFmUlEwVGZIOWFjbWZndkFVcWFXbkNsTV9oekNzekhKOTBYaHNtOXIzQ3oya29UcDJhLWtVVWlhSWQxbGlPeWcyNFMtS29TZlNwUEpiNWcxVzZhNHB5VU1xMXNDNGFuUQ=="

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

func TestLoadManifestForPrivateService(t *testing.T) {
	namespace := "custom"
	commitSha := "93f4be93"

	proj := &project.Project{Name: "knative_test",
		Directory: path.Join(root, "example"),
		Resources: os.DirFS(root),
		IsPrivate: false,
	}

	dockerImage := docker.DockerImage{
		Registry:  "example.com",
		Directory: ".",
		Name:      "test",
		Tag:       "v1.1.0",
	}

	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, path.Join(root, "example/.env.sample"))

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

func TestLoadManifestForPublicService(t *testing.T) {
	namespace := "custom"
	commitSha := "93f4be93"
	imagePullSecrets := []string{"secretPassword123"}

	proj := &project.Project{Name: "knative_test",
		Directory:        path.Join(root, "example"),
		Resources:        os.DirFS(root),
		ImagePullSecrets: imagePullSecrets,
		IsPrivate:        true,
	}

	dockerImage := docker.DockerImage{
		Registry:  "example.com",
		Directory: ".",
		Name:      "test",
		Tag:       "v1.1.0",
	}

	serviceManifest, err := LoadManifest(namespace, commitSha, proj, dockerImage, path.Join(root, "example/.env.sample"))

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
