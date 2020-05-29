FROM golang as build

WORKDIR /build

# Force modules
ENV GO111MODULE=on

# Cache dependencies
COPY go.* ./
RUN go mod download
RUN mkdir /new_tmp

# Build project
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o apatelet /build/services/apatelet

# Run stage
FROM scratch
WORKDIR /app

COPY --from=build /build/apatelet /app/apatelet
COPY --from=build /build/config/crd /app/config/crd
COPY --from=build /new_tmp /tmp/apate

ENTRYPOINT ["/app/apatelet"]
