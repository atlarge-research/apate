FROM golang as build

WORKDIR /build

# Force modules
ENV GO111MODULE=on

# Cache dependencies
COPY go.* ./
RUN go mod download

# Build project
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o control_plane /build/services/control_plane

# Run stage
FROM docker
COPY --from=build /build/control_plane /control_plane
CMD ["/control_plane"]
