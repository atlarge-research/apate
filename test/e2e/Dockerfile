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
COPY config/kind /tmp/apate/manager

ENV CP_POD_CRD_LOCATION=/app/config/crd/apate.opendc.org_podconfigurations.yaml
ENV CP_NODE_CRD_LOCATION=/app/config/crd/apate.opendc.org_nodeconfigurations.yaml
#ENV CP_LISTEN_ADDRESS=

RUN go install -v ./...

CMD ["/usr/bin/go", "test", "-v", "./test/e2e/"]
