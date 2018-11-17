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

## Get Started with Hugo

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

## Configure your site

Once you have created the project follow the instructions:

1. Download a theme into the same-named folder.
   Choose a [theme](https://themes.gohugo.io/), or
   create your own with the `hugo new theme <THEMENAME>` command.
2. You can add content by creating a single files
   with `hugo new <SECTIONNAME>/<FILENAME>.<FORMAT>`.
3. Start the built-in live server via `hugo server`.

## Sources

Here are a couple of useful links:

- [How to write Go code](https://golang.org/doc/code.html#Workspaces)
- [Interactive Go tutorial](https://tour.golang.org/welcome/1)
- [Hugo a static site generator](https://gohugo.io/)