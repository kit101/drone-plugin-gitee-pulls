go vet .\...

go test -cover .\...

:: go build -v -ldflags \"-X main.version=${DRONE_TAG##v}\" -a -tags netgo -o release/windows/drone-plugin-gitee-pulls.exe

go build -v -a -tags netgo -o release\windows\drone-plugin-gitee-pulls.exe

.\release\windows\drone-plugin-gitee-pulls.exe --help