FROM golang as build

WORKDIR /build

# Force modules
ENV GO111MODULE=on

# Cache dependencies
COPY go.* ./
RUN go mod download

# Build project
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o virtual_kubelet /build/services/virtual_kubelet

# Run stage
FROM scratch
COPY --from=build /build/virtual_kubelet /virtual_kubelet
ENTRYPOINT ["/virtual_kubelet"]
