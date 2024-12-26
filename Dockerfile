ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Install templ executable
RUN go install github.com/a-h/templ/cmd/templ@latest
# Install golang-migrate executable
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# Generate templ templates
RUN templ generate
# Build Go program to executable
RUN go build -v -o /run-app .

# Install builder image packages
RUN apt-get update && apt-get install -y \
    nodejs npm \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js dependencies
RUN npm --prefix assets install
# Run build script to generate .css file
RUN npm --prefix assets run build

# Make setup script executable
RUN chmod +x ./setup.sh
# Run setup script to download static files (HTMX and Alpine.js) from CDN
RUN ./setup.sh


FROM debian:bookworm

# Install runner image packages
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR "/app"
# Copy run-app executable to bin directory 
COPY --from=builder /run-app /usr/local/bin/
# Copy migrate executable to bin directory 
COPY --from=builder /go/bin/migrate /usr/local/bin/
# Copy migrations directory to app directory
# We don't have access to the DATABASE_URL environment variable
# inside this Dockerfile. Since we'll run the release script
# that is specified in the fly.toml file for migrations, the files need
# to be available in the runner image. 
COPY --from=builder /usr/src/app/database/migrations /app/database/migrations
# Copy static directory to app directory
COPY --from=builder /usr/src/app/static /app/static
# Copy release script to app directory
COPY --from=builder /usr/src/app/release.sh /app/release.sh
# Make release script executable
RUN chmod +x /app/release.sh
CMD ["run-app"]
