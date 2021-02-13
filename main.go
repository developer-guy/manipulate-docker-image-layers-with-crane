package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func main() {
	img, err := crane.Pull("devopps/read-file-and-write-to-sdout:latest")
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	tw := tar.NewWriter(&b)
	err = addFileToTarWriter("/Users/batuhan.apaydin/workspace/projects/personal/poc/manipulate-docker-image-layers-with-crane/layer",
		"/app",
		"/Users/batuhan.apaydin/workspace/projects/personal/poc/manipulate-docker-image-layers-with-crane/layer/hello-world.txt", tw)
	if err != nil {
		panic(err)
	}

	addLayer, err := tarball.LayerFromReader(&b)

	if err != nil {
		panic(err)
	}

	newImg, err := mutate.AppendLayers(img, addLayer)
	if err != nil {
		panic(err)
	}

	tag, err := name.NewTag("devopps/read-file-and-write-to-sdout:foo")
	if err != nil {
		panic(err)
	}

	//if s, err := daemon.Write(tag, newImg); err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println(s)
	//}

	// push to remote registry
	if err := crane.Push(newImg, tag.String()); err != nil {
		panic(err)
	}

	log.Printf("image %s pushed to the registry succesfully\n", tag.String())
}

func addFileToTarWriter(root, targetPath, filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'", filePath, err.Error()))
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get stat for file '%s', got error '%s'", filePath, err.Error()))
	}

	rel, err := filepath.Rel(root, filePath)

	header := &tar.Header{
		Name:    path.Join(targetPath, filepath.ToSlash(rel)),
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
	}

	return nil
}
