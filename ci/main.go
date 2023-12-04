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

	src := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"ci", "build", ".git"},
	})

	golang := client.Container().
		From("golang:1.21").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "get"}).
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"})

	// Run tests and discard the resulting container, actual build will run separate
	// from tests.
	_, err := golang.
		WithExec([]string{"go", "get", "-t"}).
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"/go/bin/golangci-lint", "run"}).
		WithExec([]string{"go", "test", "-coverpkg=./...", "./..."}).
		Sync(ctx)

	if err != nil {
		panic(err)
	}

	path := "build/target"
	outpath := filepath.Join(".", path)
	_ = os.MkdirAll(outpath, os.ModePerm)

	_, _ = buildIt(golang, "darwin", "amd64", path).Directory(path).Export(ctx, path)
	_, _ = buildIt(golang, "linux", "amd64", path).Directory(path).Export(ctx, path)
}

func buildIt(container *dagger.Container, os string, arch string, path string) *dagger.Container {
	return container.
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", path + "/" + os + "/find-dupe-music"})
}
