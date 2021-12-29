go vet ./...

go test -cover ./...

# go build -v -ldflags \"-X main.version=${DRONE_TAG##v}\" -a -tags netgo -o release/linux/amd64/drone-plugin-gitee-pulls

go build -v -a -tags netgo -o release/linux/amd64/drone-plugin-gitee-pulls

./release/linux/amd64/drone-plugin-gitee-pulls --help

docker build -f docker/Dockerfile.linux.amd64 -t kit101z/drone-plugin-gitee-pulls