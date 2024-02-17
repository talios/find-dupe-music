// CI project for go265
package main

import (
	"context"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

const (
	minimumCoverageLevel = "5"
)

func main() {
	ctx := context.Background()
	client, _ := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	defer client.Close()

	src := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"ci", "build", ".git"},
	})

	golang := client.Container().
		From("golang:latest").
		WithExec([]string{"go", "install", "github.com/gregoryv/uncover/cmd/uncover@latest"}).
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		WithExec([]string{"go", "install", "github.com/go-critic/go-critic/cmd/gocritic@latest"}).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "get", "-d", "-v", "./..."})

	// Run tests and discard the resulting container, actual build will run separate
	// from tests.
	_, err := golang.
		WithExec([]string{"go", "test", "-coverpkg=./...", "./...", "-coverprofile", "/tmp/c.out"}).
		WithExec([]string{"uncover", "-min", minimumCoverageLevel, "/tmp/c.out"}).
		WithExec([]string{"/go/bin/golangci-lint", "run"}).
		WithExec([]string{"gocritic", "check", "-enableAll"}).
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
