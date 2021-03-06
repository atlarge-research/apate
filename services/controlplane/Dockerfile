FROM golang as build

WORKDIR /build

# Force modules
ENV GO111MODULE=on

# Cache dependencies
COPY go.* ./
RUN go mod download

# Build project
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o controlplane /build/services/controlplane

# Run stage
FROM docker
WORKDIR /
RUN apk add curl --no-cache && \
    curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl

RUN curl -L -o /helm.tar.gz https://get.helm.sh/helm-v3.2.1-linux-amd64.tar.gz && \
    tar -zxvf /helm.tar.gz && \
    mv /linux-amd64/helm /usr/bin/helm && \
    rm -rf /helm.tar.gz /linux-amd64

RUN apk del curl

WORKDIR /app
COPY --from=build /build/controlplane /app/controlplane
COPY --from=build /build/config/kind.yml /app/config/kind.yml
COPY --from=build /build/config/prometheus.yml /app/config/prometheus.yml
COPY --from=build /build/config/crd /app/config/crd

CMD ["/app/controlplane"]
