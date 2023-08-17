// CI project for go265
package main

import (
	"context"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	client, _ := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	defer client.Close()

	src := client.Host().Directory(".")

	golang := client.Container().
		From("golang:1.21").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "get"})

	path := "build/target"
	outpath := filepath.Join(".", path)
	_ = os.MkdirAll(outpath, os.ModePerm)

	//	version := "go1.21.0"
	//	gobin := "/go/bin/" + version
	//	gobin := "go"

	//	build := golang.
	//		WithExec([]string{"go", "install", "golang.org/dl/" + version + "@latest"}).
	//		WithExec([]string{gobin, "download"}).
	//		WithExec([]string{gobin, "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
	//		WithExec([]string{"/go/bin/golangci-lint", "run"}).
	//		WithExec([]string{gobin, "get", "-u"}).
	//		WithExec([]string{gobin, "mod", "tidy"})

	_, _ = buildIt(golang, "darwin", "amd64", path).Directory(path).Export(ctx, path)
	_, _ = buildIt(golang, "linux", "amd64", path).Directory(path).Export(ctx, path)
}

func buildIt(container *dagger.Container, os string, arch string, path string) *dagger.Container {
	return container.
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{"go", "build", "-o", path + "/" + os})
}
