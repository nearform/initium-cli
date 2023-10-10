package k8s

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"gotest.tools/v3/assert"
	"io"
	"knative.dev/pkg/apis"
	v1 "knative.dev/serving/pkg/apis/serving/v1"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
)

var root = "../../../"

var caCrt = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNsakNDQVg0Q0NRRCtnL1BXUWJDSEFUQU5CZ2txaGtpRzl3MEJBUXNGQURBTk1Rc3dDUVlEVlFRR0V3SlYKVXpBZUZ3MHlNekExTURreE9EUTNNakZhRncweU5EQTFNRGd4T0RRM01qRmFNQTB4Q3pBSkJnTlZCQVlUQWxWVApNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTMzZXhaMldSUGtGeTJiQTQ4Z3RvClVwK0JkMFd2THFROU50WDBOS014YnA1SjEwMS9VaTVsY1RIcllDMWRqYTVwdjFKVGNGUW1jcWpBajBpL3dBakMKRm5oL1JKbVFrMHE1Y2liZWNURnA5UUFvRVNmbzJxYXovUTFmbmk4OG9ONlk1b1VGd2hrdTJ1bzNUWnV6M0JDMApNRnRyNXRDSGh1UzFObFhVT05VcHpzbW1UZzdZL0R1QXpNK3VIZS9qZlo1eHFqQUx3WHV6SkRFNkNCUGdhbHh6Ck9QM0V4QWNaMmRDWnRCREpVUnpTL29qMFYxOVdsRG5FK1FkTmtGTXlaMHN1UGxPTy95V0Ercno1byt2c2dKVGgKa2ZjaVZQSXZBOUdJdDdOcDZVdzU3MW1VOVlmd2x6MlU0Qm1xTXl6M05pZFF3bVd3V0wvay8yRHlXa0JDaFRBKwpvd0lEQVFBQk1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQWQ2M01qbFRJSWlSaVdOSDZvcDJySURjY1d5aWNKClJGbHJoODllSDVWeU04Q1o0NUhPYzVjbzZvRDVFdzQ3eG9vSlI2enZEd0c2anozelpFL3ArY2I5aGJ5dFJ5cysKaWNDd1g2dnRtbVpPN0M1RHdIMFYrUzk4emowNytmZFR3dUJHTDIxSlpZVmg1bFR6cEFpdU9iSkh6OTA5d1Y3OApObVhRSHpkMmtEZnpmTWhUaXpWZERPZEs3K2k1Q1RmaENIWUc4dDY2U3pmMGU5cWJ5eUFvTndwRnpxV01lRHROCkdJTGxBNHljcm1pYzBldUpmenZjeGk5NUVwMDdaZ1dYY3pINytLTWJtVnd2RkJWblZHdE9MZC9kMWhKaXlnU20KZUZERURxam9xL1JEcndCU2thbnlJMjNVa21uNGxkNDUvbHFVNDRSMVNqZXJaQmtNUmVIeTN0ZloKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="

var endpoint = "https://127.0.0.1:6443"

var token = "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklqWllaRFEwT0ZwUWVsRlJRVkp0TUMxa04wWTRNV3N6UzFwVmRXOVhSMDlvYzA5MWFESXdPVWhwVkhjaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpyZFdKbExYTjVjM1JsYlNJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKdmEyVXRhM1ZpWldOdmJtWnBaeTF6WVMxMGIydGxiaUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG01aGJXVWlPaUowWlhOMExYTmhJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpYSjJhV05sTFdGalkyOTFiblF1ZFdsa0lqb2lNelV4TkdJME9XVXRaVE15TWkwMFpHTTBMVGhsTUdVdE1qSTROemxsWm1GaE9HTXpJaXdpYzNWaUlqb2ljM2x6ZEdWdE9uTmxjblpwWTJWaFkyTnZkVzUwT210MVltVXRjM2x6ZEdWdE9uUmxjM1F0YzJFaWZRLlI5YnN3LU1MUjBubUhNOVJoUWJBZTBoZ2U0Z3JzYkJPa25RZUxVWC13SnRFY2dWejRDLVp6MFhVbUpPMEo4SVFmZXU1b3F3RnJwVHpSVEp0R24wdVdqc1RrSTZHSWNRNkpxM0FUSms3MEkwVzFqUTVJTkpJVjVmMFpfZDlIazNpZnVQaFNUUmpBZ2ljTDJCNjdMUHVBaW40T05hOHNkTE95VTZrSXFDU2Q3dURUNVMtLS1qQ0JzTnpJd1p0QVg1dVFfbDZGUHVUdlZxdGJiSUJGblVQUGExandaSWwwd3U1YV9DTG9rVFhWbDduVkZYUzBxZTJ1RWFmUlEwVGZIOWFjbWZndkFVcWFXbkNsTV9oekNzekhKOTBYaHNtOXIzQ3oya29UcDJhLWtVVWlhSWQxbGlPeWcyNFMtS29TZlNwUEpiNWcxVzZhNHB5VU1xMXNDNGFuUQ=="

const certificate = `-----BEGIN CERTIFICATE-----
MIIC/jCCAeagAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5ldGVzMB4XDTIzMTAxMDA2NTEwN1oXDTMzMTAwNzA2NTEwN1owFTETMBEGA1UE
AxMKa3ViZXJuZXRlczCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANjP
vjrW7UB3z4PFI4uw2yxlq26qOlesUY6aPq9EvHEBGfS2v74eRUXNcMt7rpvzOnjJ
nCW+E0earYAFW/seLpuSR7vvhUtnbfvUpl6GHbje2C8Q6TjE9YpqSpVj0vbSiPCy
/8AXQUUAzQgGJnC2KP62aA/vjxkJZ4SUkHD1Y1Ni6GgYS0BxCvtxfIoHUA3Tr+MX
SUnxR4IAzmQNoberifmK3DISXlXrRhc7eKnELwFdsTB/FKXDERiFJisXKDglxN/Y
ooUv++zhXAsl42ne2Y9G+7NIr/PcZBRWUSQo3wRHW6b8upS7TTHF3ScTWjAyIHPC
kMhPQ4zVtC8Zc5gDLKMCAwEAAaNZMFcwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB
/wQFMAMBAf8wHQYDVR0OBBYEFLq5/eX3KD4kKFZnFmRiauBTfRppMBUGA1UdEQQO
MAyCCmt1YmVybmV0ZXMwDQYJKoZIhvcNAQELBQADggEBALqcDaumB0+7EEsazzK5
BqnOaNW5Cq255rcz+aC+elvECuFRQ4nliG3jTADfbZ1dMZiF95fzPdkbOKDIMy4n
NkgexnqrrxcP6MmMhZeUb5o5UG9RUVg2CBd17menE2ogMmq6IFyR6l1qoecQtB3T
EQb0s5BWHqJ7HgIEgaURKCdZEtgjdXQZxvx5U7JqY7HDqeMpQlad3kYsKmZcsAGM
f1kAZcMbBJuBbCNuTJoOYeLIiBH73JjNSn4k27JsrCuAO33InyWKNDNB/WN/ghV1
tn+2QMzMAynSx2vmpWqDiomiHYZoLrv9yQhIfGoT9CwEMQn5ciOUzDvT214W/Ryp
+Mw=
-----END CERTIFICATE-----
`

type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

type knativeApiRequest struct {
	httpMethod string
	url        string
}

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

func TestLoadManifest(t *testing.T) {
	namespace := "custom"
	commitSha := "93f4be93"

	proj := &project.Project{Name: "knative_test",
		Directory: path.Join(root, "example"),
		Resources: os.DirFS(root),
	}

	dockerImage := docker.DockerImage{
		Registry:  "example.com",
		Directory: ".",
		Name:      "test",
		Tag:       "v1.1.0",
	}

	serviceManifest, err := loadManifest(namespace, commitSha, proj, dockerImage, path.Join(root, "example/.env.sample"))

	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	annotations := serviceManifest.Spec.Template.ObjectMeta.Annotations
	assert.Assert(t, annotations[UpdateTimestampAnnotationName] != "", "Missing %s annotation", UpdateTimestampAnnotationName)
	assert.Assert(t, annotations[UpdateShaAnnotationName] == commitSha, "Expected %s SHA, got %s", commitSha, annotations[UpdateShaAnnotationName])
}

func TestApply(t *testing.T) {
	config, err := Config(endpoint, token, []byte(certificate))
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	var apiRequests []knativeApiRequest
	getServiceRequestCount := 0
	handler := func(request *http.Request) (*http.Response, error) {
		apiRequests = append(apiRequests, knativeApiRequest{
			httpMethod: request.Method,
			url:        request.URL.Path,
		})

		header := http.Header{}
		header.Set("Content-Type", "application/json")

		// simulate new service creation
		if request.Method == http.MethodGet && request.URL.Path == "/apis/serving.knative.dev/v1/namespaces/default/services/initium-cli" {
			if getServiceRequestCount == 0 {
				getServiceRequestCount += 1
				return &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
					Header:     header,
				}, nil
			} else {
				statusUrl, handlerErr := url.Parse("http://initium-nodejs-demo-app.initium-chore-another-pr.example.com")
				if handlerErr != nil {
					t.Fatalf(fmt.Sprintf("Error: %s", err))
				}
				service := v1.Service{
					Status: v1.ServiceStatus{
						RouteStatusFields: v1.RouteStatusFields{
							URL: (*apis.URL)(statusUrl),
						},
					},
				}
				marshaledJson, handlerErr := json.Marshal(service)
				if handlerErr != nil {
					t.Fatalf(fmt.Sprintf("Error: %s", err))
				}

				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader(marshaledJson)),
					Header:     header,
				}, nil
			}
		} else {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     header,
			}, nil
		}
	}

	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return transportFunc(handler)
	})

	proj := project.Project{
		Name:      "initium-cli",
		Directory: path.Join("../../../", "."),
		Resources: os.DirFS("../../../"),
	}

	dockerImage := docker.DockerImage{
		Registry:  "example.org",
		Directory: defaults.ProjectDirectory,
		Name:      "test",
		Tag:       "v1.1.0",
	}

	err = Apply("default", "6ac63179", config, &proj, dockerImage, defaults.EnvVarFile)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	assert.Check(t, len(apiRequests) == 4, "Expected 4 requests to be sent to Knative API")

	assert.Assert(t, apiRequests[0].httpMethod == http.MethodPost, "Expected POST method to be called")
	assert.Assert(t, apiRequests[0].url == "/api/v1/namespaces", "Expected URL suffix to be /api/v1/namespaces")

	assert.Assert(t, apiRequests[1].httpMethod == http.MethodGet, "Expected GET method to be called")
	assert.Assert(t, apiRequests[1].url == "/apis/serving.knative.dev/v1/namespaces/default/services/initium-cli", "Expected URL suffix to be /apis/serving.knative.dev/v1/namespaces/default/services/initium-cli")

	assert.Assert(t, apiRequests[2].httpMethod == http.MethodPost, "Expected POST method to be called")
	assert.Assert(t, apiRequests[2].url == "/apis/serving.knative.dev/v1/namespaces/default/services", "Expected URL suffix to be /apis/serving.knative.dev/v1/namespaces/default/services")

	assert.Assert(t, apiRequests[3].httpMethod == http.MethodGet, "Expected GET method to be called")
	assert.Assert(t, apiRequests[3].url == "/apis/serving.knative.dev/v1/namespaces/default/services/initium-cli", "Expected URL suffix to be /apis/serving.knative.dev/v1/namespaces/default/services/initium-cli")
}
