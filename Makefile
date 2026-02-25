binary_name = res
build_dir = ./bin

.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n] ' && read ans && [ $${ans:-n} = y ]

## build: build the application
.PHONY: build
build:
	@mkdir -p ${build_dir}
	go build -o ${build_dir}/${binary_name} main.go

## export: export the application
.PHONY: export
export: build
	@mkdir -p ${build_dir}
	GOARCH=amd64 GOOS=darwin go build -o ${build_dir}/${binary_name}_darwin_amd64 main.go
	GOARCH=arm64 GOOS=darwin go build -o ${build_dir}/${binary_name}_darwin_arm64 main.go
	GOARCH=amd64 GOOS=linux go build -o ${build_dir}/${binary_name}_linux_amd64 main.go
	GOARCH=amd64 GOOS=windows go build -o ${build_dir}/${binary_name}_windows_amd64.exe main.go

## run: run the application
.PHONY: run
run:
	@cd ${build_dir} && ./${binary_name}

## clean: clean up the build binaries
.PHONY: clean
clean: confirm
	@echo "Cleaning up..."
	@rm -rf ${build_dir}
