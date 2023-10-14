###############
# Build image
###############
FROM golang:latest AS builder

RUN mkdir /build
COPY go.mod go.sum /build
RUN ls -la /build
WORKDIR /build
RUN go mod download && go mod verify

COPY . /build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /wunschliste


###############
# Actual image
###############
FROM scratch

COPY --from=builder /wunschliste /bin/wunschliste
ENTRYPOINT ["/bin/wunschliste", "-file", "/data/wunschliste.bolt"]
EXPOSE 8080/tcp
