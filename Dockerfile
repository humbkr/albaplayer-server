FROM dockercore/golang-cross:latest


RUN apt-get update

# Vim
RUN apt-get install -y vim

# Create aliases
RUN echo 'alias ll="ls -lah"' >> ~/.bashrc
RUN echo 'alias vi="vim"' >> ~/.bashrc

# Set go bin which doesn't appear to be set already
ENV GOBIN /go/bin

# Install dependency manager
RUN go get github.com/golang/dep/...
# Install live refresh serveer
RUN go get github.com/pilu/fresh

CMD []
