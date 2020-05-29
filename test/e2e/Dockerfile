# This Dockerfile can be used to run e2e tests locally if you do not want to set up an environment.
FROM docker

RUN apk add --no-cache ca-certificates go git make protobuf

# Cache dependencies
COPY go.* ./
RUN go mod download

RUN apk add curl --no-cache && \
    curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl

RUN curl -L -o /helm.tar.gz https://get.helm.sh/helm-v3.2.1-linux-amd64.tar.gz && \
    tar -zxvf /helm.tar.gz && \
    mv /linux-amd64/helm /usr/bin/helm && \
    rm -rf /helm.tar.gz /linux-amd64

WORKDIR /app
# Copy project
COPY . .

RUN go install -v ./...

ENV CI_PROJECT_DIR=/app

CMD ["/usr/bin/go", "test", "-v", "./test/e2e/"]
