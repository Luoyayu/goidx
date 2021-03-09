BINARY_NAME = goidx

build:
	go mod tidy
	go build -o ${BINARY_NAME} -ldflags "-w -s -buildid="

update_deps:
	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null