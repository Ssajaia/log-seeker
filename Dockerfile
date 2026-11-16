# Build the Go server and bundle the existing static frontend.
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/logseeker .

# Run as an unprivileged user; the database volume is writable by this user.
FROM alpine:3.21
RUN addgroup -S logseeker && adduser -S -G logseeker logseeker \
    && mkdir -p /app/frontend /data \
    && chown -R logseeker:logseeker /app /data
WORKDIR /app
COPY --from=build --chown=logseeker:logseeker /out/logseeker ./logseeker
COPY --chown=logseeker:logseeker frontend ./frontend
ENV LOGSEEKER_ADDR=:8080
ENV LOGSEEKER_DB=/data/logseeker.db
EXPOSE 8080
USER logseeker
ENTRYPOINT ["/app/logseeker"]
