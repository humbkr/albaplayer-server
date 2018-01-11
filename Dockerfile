FROM golang:1.9

# Install utils
RUN alias ll="ls -alh"

# Set go bin which doesn't appear to be set already
ENV GOBIN /go/bin

# Install dependency manager
RUN go get github.com/golang/dep/...
# Install live refresh serveer
RUN go get github.com/pilu/fresh

CMD []
