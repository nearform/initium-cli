package cli

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/nearform/initium-cli/src/services/project"
	"github.com/urfave/cli/v2"
	"gotest.tools/v3/assert"
)

const root = "../../"

const expectedGoAppDockerTemplate = `FROM golang:1.20.4 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v ./...
RUN go test -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
`

const expectedNodeAppDockerTemplate = `FROM node:20.2.0 AS build-env

WORKDIR /app

COPY package*.json tsconfig*.json ./

RUN npm i

COPY . /app

RUN npm run build --if-present
RUN npm test

FROM gcr.io/distroless/nodejs20-debian11
COPY --from=build-env /app /app
WORKDIR /app
USER nonroot
CMD ["index.js"]
`

var projects = map[project.ProjectType]map[string]string{
	project.NodeProject: {"directory": "example", "expectedTemplate": expectedNodeAppDockerTemplate},
	project.GoProject:   {"directory": ".", "expectedTemplate": expectedGoAppDockerTemplate},
}

func TestShouldRenderDockerTemplate(t *testing.T) {
	for projectType, props := range projects {
		var buffer bytes.Buffer

		cCtx := cli.Context{}
		instance := icli{
			project: &project.Project{
				Name:      string(projectType),
				Directory: path.Join(root, props["directory"]),
				Resources: os.DirFS(root),
			},
			Writer: &buffer,
		}

		err := instance.template(&cCtx)
		if err != nil {
			t.Fatalf(fmt.Sprintf("Error: %v", err))
		}

		expectedTemplate := props["expectedTemplate"]
		template := string(buffer.Bytes())

		assert.Assert(t, template == expectedTemplate, "Expected %s, got %s", expectedTemplate, template)
	}
}
