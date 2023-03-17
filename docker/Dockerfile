FROM node:lts-alpine

# Connect the built container to a repository when push to the Github Package registry, automatically
# See https://docs.github.com/en/packages/learn-github-packages/connecting-a-repository-to-a-package
ARG LABEL_IMAGE_SOURCE
LABEL org.opencontainers.image.source=${LABEL_IMAGE_SOURCE}

ADD . /app

WORKDIR /app

ENTRYPOINT [ "npm", "start" ]
