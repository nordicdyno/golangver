# golangver

Yet Another Go versions manager with its own mojo! 

## FAQ

> Q: Why not [gvm](https://github.com/moovweb/gvm) or any of available similiar solutions?

Why not "both"? 😁

Actually `golangver` does version management in its own opinionated style:

1. uses symlink (can be redefined with flag `--go-bin`) for local switch to desired Go version
2. uses `go install https://go.dev/dl/go<version>` distros, so no compilation step after distro fetch 
4. IDE-aware (for now supports IDEA/Goland only) – i.e. suggests patching Go SDK version in project settings if project files are detected
5. go.mod aware – i.e. suggests patching Go's version in `go.mod` if detected
6. doesn't store any local artifacts like its own dir and not require patched local `.<shell>rc` or `.profile` files
7. is written in Go not bash


## Requirements

* Already installed Go distributive:
https://go.dev/doc/install (or just use package manager like `brew` or `apt-get`)
* git binary

## Installation

    go install github.com/nordicdyno/golangver@latest

## How to

show Go distributives: all available locally and all latest Y-minor version with latest Z path (1.Y.Z) available remotely (`-r` flag):

    
    golangver list -r
    
    # Output:
    downloaded by `go install`:
      1.18beta2   /Users/user/sdk/go1.18beta2/bin/go
    * 1.17.6      /Users/user/sdk/go1.17.6/bin/go
      1.17.3      /Users/user/sdk/go1.17.3/bin/go
    
    managed by other tools (IDEA):
      1.17.6      /Users/user/go/go1.17.6/bin/go
      1.17.5      /Users/user/go/go1.17.5/bin/go
    
    # remote Go versions:
      1.18beta2	https://go.dev/blog/go1.18beta2
      1.17.6	https://golang.org/doc/devel/release#go1.17
      1.16.13	https://golang.org/doc/devel/release#go1.16
      1.15.15	https://golang.org/doc/devel/release#go1.15
      1.14.15	https://golang.org/doc/devel/release#go1.14
      1.13.15	https://golang.org/doc/devel/release#go1.13

show all available Go distros locally:

    golangver list

show all Go distributives available locally and all available (`-a`) remotely with outdated stuff Go<1.13 (`-o`):

    golangver list -r -a -o

install Go by version number (available on https://go.dev/dl/):

    golangver get 1.17.6

switch symlink by Go version number (supports distros are installed by `go install` command only):

    golangver use 1.17.6

switch symlink to provided full path to Go binary:

    golangver use /Users/user/sdk/go1.17.6/bin/go
