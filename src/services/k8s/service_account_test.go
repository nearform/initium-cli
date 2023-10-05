package k8s

import (
	"fmt"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

const rootDirectory = "../../../"

func TestShouldReturnSingleFileContent(t *testing.T) {
	expectedFileContent := `apiVersion: v1
kind: ServiceAccount
metadata:
  name: initium-cli-sa`

	bytes, err := getFileContent("service-account", os.DirFS(rootDirectory))
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}
	fileContent := string(bytes)

	assert.Assert(t, expectedFileContent == fileContent, "Expected %s, got %s", expectedFileContent, fileContent)
}

func TestShouldReturnServiceAccountManifest(t *testing.T) {
	expectedFileContent := `---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: initium-cli
rules:
  - apiGroups:
      - ''
      - serving.knative.dev
      - apps
      - networking.k8s.io
    resources:
      - namespaces
      - deployments
      - replicasets
      - ingresses
      - services
      - secrets
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: initium-cli
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: initium-cli
subjects:
  - kind: ServiceAccount
    name: initium-cli-sa
    namespace: default
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: initium-cli-sa
---
apiVersion: v1
kind: Secret
type: kubernetes.io/service-account-token
metadata:
  name: initium-cli-token
  annotations:
    kubernetes.io/service-account.name: "initium-cli-sa"
`

	serviceAccountManifest, err := getServiceAccountManifest(os.DirFS(rootDirectory))
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	assert.Assert(t, expectedFileContent == serviceAccountManifest, "Expected %s, got %s", expectedFileContent, serviceAccountManifest)
}
