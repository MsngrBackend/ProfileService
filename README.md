# Profile microservice

WIP version of profile microservice

# Installation guide

1) Clone repository

2) Make sure Go 1.26 or newer is installed

3) Install 'go-micro' framework

  ```go install go-micro.dev/v5@latest```


go mod tidy

go install github.com/pressly/goose/v3/cmd/goose@latest

export PATH=$PATH:$(go env GOPATH)/bin to ~/.bashrc

docker compose up -d



goose -dir migrations postgres \
  "postgres://profile:secret@localhost:5432/profile_db" up

# open http://localhost:9001 in browser
# login: minioadmin / minioadmin
# create bucket named "avatars"
