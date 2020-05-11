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

# Build 
FROM alpine/helm as prom
WORKDIR /
RUN helm repo add google https://kubernetes-charts.storage.googleapis.com/
RUN helm template prometheus google/prometheus > /prom

# Run stage
FROM docker
COPY --from=build /build/controlplane /controlplane
COPY --from=build /build/config/kind /tmp/apate/manager
COPY --from=prom  /prom /tmp/prometheus

RUN apk add curl --no-cache && \
    curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.8.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl && \
    apk del curl

CMD ["/controlplane"]
