# Build phase
FROM golang:1.12 AS build

ADD . /app
WORKDIR /app

# Install dependencies
RUN go mod download

# Build app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o /generated/alba .

# Copy config files
RUN cp /app/build/prod.alba.yml /generated/alba.yml

# Copy webapp
RUN mkdir /generated/web && cp -r /app/web/. /generated/web


# Final image
FROM alpine:latest
RUN apk add --no-cache \ 
  ca-certificates \
  libc6-compat
COPY --from=build /generated/ /app/
RUN chmod +x /app/alba

ENTRYPOINT cd /app && ./alba serve
EXPOSE 8888

