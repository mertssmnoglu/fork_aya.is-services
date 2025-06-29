# versions

FROM golang:1-bookworm AS upstream-builder
# FROM gcr.io/distroless/base-debian12 AS upstream-runner
FROM debian:bookworm-slim AS upstream-runner

# Create a minimal image base-debian12 or static-debian12
# (see: https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)

# ------------------
# builder-base image
# ------------------
FROM upstream-builder AS builder-base

# Configuration
ARG GH_LOGIN
ENV GH_LOGIN=$GH_USER

ARG GH_ACCESS_TOKEN
ENV GH_ACCESS_TOKEN=$GH_ACCESS_TOKEN

ARG GH_PATH
ENV GH_PATH=$GH_PATH

# Setup for private Go Modules where available
RUN echo "machine github.com login ${GH_LOGIN} password ${GH_ACCESS_TOKEN}" > ~/.netrc
ENV GOPRIVATE="${GH_PATH}/*"
RUN go env -w GOPRIVATE="${GH_PATH}/*"
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set working directory
# RUN mkdir -p /app
WORKDIR /app

# Install dependencies first to improve caching
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# Build the application
COPY . .

# ------------------------
# runner-development image
# ------------------------

FROM builder-base AS runner-development

EXPOSE 8080
ENTRYPOINT ["go", "run", "./cmd/serve/"]

# ------------------------
# builder-production image
# ------------------------

FROM builder-base AS builder-production

RUN go build ./cmd/serve/

# -----------------------
# runner-production image
# -----------------------
FROM upstream-runner AS runner-production

RUN apt-get update && \
  apt-get -y install --no-install-recommends curl && \
  groupadd -r nonroot && \
  useradd -r -g nonroot nonroot

# Copy the binary from the builder-production container
# COPY --from=builder-production /bin/sh /bin/sh
# COPY --from=builder-production /bin/cat /bin/cat
COPY --from=builder-production --chown=nonroot:nonroot /app/serve /
COPY --from=builder-production --chown=nonroot:nonroot /app/.env /
COPY --from=builder-production --chown=nonroot:nonroot /app/config.json /

# Run as a non-root user
USER nonroot

EXPOSE 8080
ENTRYPOINT ["/serve"]
