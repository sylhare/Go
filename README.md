# Go

Experiments on Go

## Go-lang

Go is an open source programming language.
To install it, you can either download it from [golang.org](https://golang.org/) or use brew:

```bash
brew install golang
...
go version
>> go version go1.11.2 darwin/amd64
```

## Introduction

### Run your script

Run you script `hello.go` with:

```bash
go run src/examples/hello.go
>> Hello, World!
```

### Install a Go script

Install it in your `$GOPATH/bin` with:

```bash
go install hello.go
```

You need to set your own `$GOPATH` for it to work. 
Check `setup.sh` to see which path to export, Add the change to your `.bashrc` to make it permanent.

### Create a Go module

Run the following command to create a new module:

```bash
go mod init <module-name>
```

This will create a `go.mod` file in the current directory.

### Dependency management

To add a dependency, run:

```bash
go get <package-name>
```

This will add the package to the `go.mod` file.
To install all dependencies, and clean up the `go.mod` file, run:

```bash
go mod tidy
```

You can also use `go mod download` which will download all dependencies in the cache.

### Run a unit test

Run unit tests:

```bash
go test ./...
```

This should run all the tests in the current directory and subdirectories.

## Sources

Here are a couple of useful links:

- [How to write Go code](https://golang.org/doc/code.html#Workspaces)
- [Interactive Go tutorial](https://tour.golang.org/welcome/1)