VERSION 0.8

ARG --global go_version = 1.22.5
ARG --global node_version = 20.5.1

ARG --global outputDir = "./dist"

build-base:
  FROM debian:bullseye
  WORKDIR /root

  # Install necessary dependencies
  RUN apt-get update && apt-get install -y \
      git \
      curl \
      gnupg \
      build-essential \
      git \
      wget \
      gcc \
      tar \
      libc6-dev


  RUN wget https://go.dev/dl/go${go_version}.linux-amd64.tar.gz
  RUN tar -C /usr/local -xzf go${go_version}.linux-amd64.tar.gz
  ENV PATH=$PATH:/usr/local/go/bin
  RUN go version

  ENV GOCACHE = "/.go-cache"
  ENV GOMODCACHE = "/.go-mod-cache"
  ENV CGO_ENABLED = 1
  
  RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.0/install.sh | bash
  ENV NVM_DIR=/root/.nvm
  RUN . "$NVM_DIR/nvm.sh" && nvm install ${node_version}
  RUN . "$NVM_DIR/nvm.sh" && nvm use v${node_version}
  RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${node_version}
  ENV PATH="/root/.nvm/versions/node/v${node_version}/bin/:${PATH}"
  RUN node --version
  RUN npm --version

  # Install Ionic and Angular CLI globally
  RUN npm install -g @ionic/cli @angular/cli

  SAVE IMAGE --cache-hint
  
build:
  FROM +build-base
  WORKDIR /go-workdir
  COPY . .
  RUN cd web && npm install
  RUN cd web && npm install --save-dev fuse.js # Dont know why I need this

  RUN --no-cache ./build.sh

  SAVE ARTIFACT "vet_note" AS LOCAL ${outputDir}/vet_note

