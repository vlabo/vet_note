VERSION 0.8

ARG --global go_version = 1.22.5
ARG --global node_version = 20.5.1

ARG --global outputDir = "./dist"

ionic-build:
  FROM "node:${node_version}"

  WORKDIR /web

  # Install Ionic and Angular CLI globally
  RUN npm install -g @ionic/cli @angular/cli
  SAVE IMAGE --cache-hint

  COPY web .
  RUN npm install
  RUN npm install --save-dev fuse.js # Dont know why I need this
  RUN ionic build

  SAVE ARTIFACT --keep-ts "./www" "${outputDir}/www"

    # we need to do some magic here because tauri expects the binaries to include the rust target tripple.
  
build:
  FROM "golang:${go_version}-alpine"
  WORKDIR /go-server
  COPY . .
  COPY (+ionic-build/${outputDir}/www) ./web/www

  RUN --no-cache go build .

  SAVE ARTIFACT "vet_note" AS LOCAL ${outputDir}/vet_note

