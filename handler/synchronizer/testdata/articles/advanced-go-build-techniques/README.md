Table of contents
- [Build options](#i-build-options)
- [Which file will be included](#ii-which-file-will-be-included)
- [Build tags](#iii-build-tags)
- [Build contraints](#iv-build-contraints)

## I/ Build options

The following are some of the most commonly used options for the `go build` command:

* `-o`: Specifies the output file name. The default output file name is the name of the main package, with the `.exe` suffix added on Windows.
* `-v`: Verbose output. This option prints the names of the packages as they are being compiled.
* `-work`: Prints the name of the temporary work directory and does not delete it when exiting. This option is useful for debugging.
* `-x`: Print the commands. This option prints the commands that are being executed by the `go build` command.
* `-asmflags`: Arguments to pass to the `go tool asm` invocation.
* `-buildmode`: The build mode to use. The default build mode is `exe`. Other possible values are `shared`, `pie`, and `plugin`.
* `-buildvcs`: Whether to stamp binaries with version control information. The default value is `auto`.

For more information about the `go build` command, you can run the following command:
```bash
go help build
```

## II/ Which file will be included
When you use the `go build` command in Go, it compiles the Go source files in the current directory and its subdirectories to create an executable binary. By default, Go only compiles `.go` files and ignores other types of files in the directory. However, it's important to note that the behavior of the `go build` command can be influenced by build tags, build constraints.

The following types of files are typically ignored by `go build`:

**1. Files with Non-`.go` Extensions:** 

Any files in the directory that don't have a `.go` extension will be ignored. This includes text files, configuration files, images, etc.

**2. Files in Subdirectories:** 

The `go build` command compiles all `.go` files in the current directory and its subdirectories. Other files and directories will generally be ignored.

**3. Files with Underscores OR Periods in the Beginning:** 

Directory and file names that begin with `"."` or `"_"` are ignored by the go tool, as are directories named "testdata".

**4. Files Excluded by Build Constraints:** 

Go supports build constraints which allow you to include or exclude specific files from the build based on conditions like the target operating system or architecture. For example, files with build constraints like `//go:build windows` will be ignored when building for non-Windows platforms.

**5. Files Excluded by Build Tags:** 

Build tags are special comments in Go source files that can be used to specify which files should be included in the build based on custom conditions. Files with build tags that don't match the build context will be ignored.

**6. Files in Directories Named "testdata":** Files in directories named `testdata` are ignored by design. This directory is commonly used to include test-related data that's not meant to be compiled.

## III/ Build tags
Go's build tags provide a powerful mechanism for including or excluding specific code during the build process. By using build tags, developers can tailor their application to work with different build configurations, environments, or platform-specific requirements. This feature is particularly valuable when dealing with cross-compilation or creating binaries for specific operating systems.

A build tag is a comment placed at the beginning of a Go source file that specifies a set of conditions under which the code within that file should be included or excluded from the build. The syntax is //go:build <tag>. For instance, consider a scenario where you want to include a specific piece of code only when building for a certain version of the application:

`main.go`

```go
package main

import "fmt"

var version string

func main() {
	fmt.Println(version)
}
```

`pro.go`
```go
//go:build pro

package main

func init() {
	version = "pro"
}
```

`free.go`
```go
//go:build free

package main

func init() {
	version = "pro"
}

```

When you build with `-tags=free`, the output will be `free` because the `free.go` file was included. And when you build with `-tags=pro`, the output will be `pro`.

### Build Tags Syntax

You can combine constraints in the same way as you would any other conditional statement in programming i.e AND, OR, NOT

**NOT**

```go
//go:build !cgo
```

This will only include the file in the build if CGO is NOT enabled.

**AND**

```go
//go:build cgo && darwin
```
This will only include the file in the build if CGO is enabled and the GOOS is set to darwin.

**OR**
```go
//go:build darwin || linux
```

**Combine them all e.g**

```go
//go:build (linux || 386) && (darwin || !cgo)
```
Results to (linux AND 386) OR (darwin AND (NOT cgo)).


**Note:**
Go versions `1.16` and earlier used a different syntax for build constraints, with a `// +build` prefix. The gofmt command will add an equivalent `//go:build` constraint when encountering the older syntax.

## IV/ Build contraints
While custom build tags are set using the tags build flag, golang automatically sets a few tags based on environment variables and other factors. Here is a list of the available tags

**1. OOS and GOARCH Environment Values**
You can set constraints in your source code to only run files if a certain GOOS or GOARCH is used. e.g

```go
//go:build darwin,amd64

package utils
```

**2. GO Version constraint**

You can also constrain the inclusion of a file to the go version being used in building the entire module. EX to only build the file if the go version being used is 1.12 and above you would use `//go:build go1.18`. This would include the file if the go version is 1.18 or 1.21 (latest as of writing this).

## References
- https://kofo.dev/build-tags-in-golang
- https://pkg.go.dev/cmd/go
- https://www.digitalocean.com/community/tutorials/customizing-go-binaries-with-build-tags