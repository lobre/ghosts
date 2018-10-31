# Build stage
FROM golang:alpine3.8 AS build-env

# Install git
RUN apk add --no-cache git

# Install golang depencencies
RUN go get github.com/docker/docker/client && \
    go get github.com/lextoumbourou/goodhosts

# Add sources
ADD . /src

# Build app
RUN cd /src && go build -o ghosts

# --------

# Final stage
FROM alpine:3.8

# Install docker client
ARG DOCKER_CLI_VERSION="18.06.1-ce"
ENV DOWNLOAD_URL="https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_CLI_VERSION.tgz"
RUN apk --update add curl \
    && mkdir -p /tmp/download \
    && curl -L $DOWNLOAD_URL | tar -xz -C /tmp/download \
    && mv /tmp/download/docker/docker /usr/local/bin/ \
    && rm -rf /tmp/download \
    && apk del curl \
    && rm -rf /var/cache/apk/*

# Copy go binary, static and html template
COPY --from=build-env /src/ghosts /app/
COPY --from=build-env /src/index.html /app/
COPY --from=build-env /src/static/ /app/static/

# Define workdir
WORKDIR /app

# Create fake hosts file
RUN touch /app/hosts

ENTRYPOINT ./ghosts -hosts="/app/hosts"
