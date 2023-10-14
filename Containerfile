###############
# Build image
###############
FROM golang:latest AS builder

# Add github to known hosts
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
# Clone private repo
RUN --mount=type=ssh git clone git@github.com:ibutra/wunschliste-go.git /wunschliste-build
# build
WORKDIR /wunschliste-build
RUN ls -la && pwd
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /wunschliste


###############
# Actual image
###############
FROM scratch

COPY --from=builder /wunschliste /bin/wunschliste
ENTRYPOINT ["/bin/wunschliste", "-file", "/data/wunschliste.bolt"]
EXPOSE 8080/tcp
