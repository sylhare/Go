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

## Get started

### Run your script

Run you script `hello.go` with:

```bash
go run hello.go
```

### Install a Go script

Install it in your `$GOPATH/bin` with:

```bash
go install hello.go
```

You need to set your own `$GOPATH` for it to work. 
Check `setup.sh` to see which path to export, Add the change to your `.bashrc` to make it permanent.


## Hugo

Hugo is a go framework for building websites, a bit like Jekyll.

Install [Hugo](https://gohugo.io/getting-started/quick-start/) with brew on mac:
```bash
brew install hugo
...
hugo version
>> Hugo Static Site Generator v0.51
```

And then create your first website with:
```bash
hugo new site quickstart
```

## Sources

Here are a couple of useful links:

- [How to write Go code](https://golang.org/doc/code.html#Workspaces)
- [Interactive Go tutorial](https://tour.golang.org/welcome/1)
- [Hugo a static site generator](https://gohugo.io/)