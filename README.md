> Credit: Inspired by @ahmetb's [latest blog post about building container images in Go.](ahmet.im/blog/building-container-images-in-go), please read this post before you move on to the hands-on section.

# What is a crane tool?

<img src="https://github.com/google/go-containerregistry/raw/main/images/crane.png" height="300"/>

Google has a repository called ["go-containerregistry"](https://github.com/google/go-containerregistry) which provides Go library and CLIs for working with container registries, [crane](https://github.com/google/go-containerregistry/blob/main/cmd/crane/README.md) is one of them. More technically, the crane is a tool for interacting with remote images and registries.

# Hands On

Let's start with explaining the demo, first, we have a directory that includes a basic Go application that prints the content of the file to stdout, we'll start with building the container image, then with crane, we'll add a new layer to it using hello-world.txt that is available in the [layer/](./layer) directory, by doing so, we'll update the content of the file that is available within the container image.

Let's take a look at the Dockerfile of the project.
```Dockerfile
FROM golang:1.15.7-alpine

WORKDIR /app

COPY ./ ./

ENTRYPOINT ["go", "run", "main.go"]
```

it's very straightforward, then take a look at the go code.
```golang
package main

import (
	"io/ioutil"
	"log"
)

func main() {
	content, err := ioutil.ReadFile("hello-world.txt")

	if err != nil {
		log.Fatal("could not read file, error:", err)
	}

	log.Println("Content of the file is : ", string(content))
}
```

it's very straightforward too.

Let's build the container image, in that case, we need docker to build the container image.
```bash
$ docker image build -t devopps/read-file-and-write-to-sdout:latest .
[+] Building 7.3s (9/9) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                                                      0.0s
 => => transferring dockerfile: 131B                                                                                                                                                                      0.0s
 => [internal] load .dockerignore                                                                                                                                                                         0.0s
 => => transferring context: 2B                                                                                                                                                                           0.0s
 => [internal] load metadata for docker.io/library/golang:1.15.7-alpine                                                                                                                                   1.7s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                                                                                             0.0s
 => [internal] load build context                                                                                                                                                                         0.0s
 => => transferring context: 463B                                                                                                                                                                         0.0s
 => [1/3] FROM docker.io/library/golang:1.15.7-alpine@sha256:dbda4e47937a3abb515c386d955002be5116d060c90d936127cc24ac439c815c                                                                             4.9s
 => => resolve docker.io/library/golang:1.15.7-alpine@sha256:dbda4e47937a3abb515c386d955002be5116d060c90d936127cc24ac439c815c                                                                             0.0s
 => => extracting sha256:4c0d98bf9879488e0407f897d9dd4bf758555a78e39675e72b5124ccf12c2580                                                                                                                 0.2s
 => => sha256:8b36f00a8e74ce31a867744519cc5db8c4aaeb181cffcda1b4d8269b1cc7f336 106.77MB / 106.77MB                                                                                                        0.0s
 => => sha256:5e5ebcc3e85238e4fbf5ab2428f9ed61dcede6c59b605d56b2f02fb991c70850 126B / 126B                                                                                                                0.0s
 => => sha256:dbda4e47937a3abb515c386d955002be5116d060c90d936127cc24ac439c815c 1.65kB / 1.65kB                                                                                                            0.0s
 => => sha256:18100456495c42bcdccab3411d8cfd56f3fdaa8527dd2b5a83800f96c7074a41 1.36kB / 1.36kB                                                                                                            0.0s
 => => sha256:54d042506068c9699d4236315fa76ea8789415c1079bcaff35fb3730ea649547 4.61kB / 4.61kB                                                                                                            0.0s
 => => sha256:4c0d98bf9879488e0407f897d9dd4bf758555a78e39675e72b5124ccf12c2580 2.81MB / 2.81MB                                                                                                            0.0s
 => => sha256:9e181322f1e7b3ebee5deeef0af7d13619801172e91d2d73dcf79b5d53d82d91 281.20kB / 281.20kB                                                                                                        0.0s
 => => sha256:6422294da7d35128e72551ecf15f3a4d9577e5cfa516b6d62fe8b841a9470cb3 154B / 154B                                                                                                                0.0s
 => => extracting sha256:9e181322f1e7b3ebee5deeef0af7d13619801172e91d2d73dcf79b5d53d82d91                                                                                                                 0.1s
 => => extracting sha256:6422294da7d35128e72551ecf15f3a4d9577e5cfa516b6d62fe8b841a9470cb3                                                                                                                 0.0s
 => => extracting sha256:8b36f00a8e74ce31a867744519cc5db8c4aaeb181cffcda1b4d8269b1cc7f336                                                                                                                 4.2s
 => => extracting sha256:5e5ebcc3e85238e4fbf5ab2428f9ed61dcede6c59b605d56b2f02fb991c70850                                                                                                                 0.0s
 => [2/3] WORKDIR /app                                                                                                                                                                                    0.5s
 => [3/3] COPY ./ ./                                                                                                                                                                                      0.0s
 => exporting to image                                                                                                                                                                                    0.0s
 => => exporting layers                                                                                                                                                                                   0.0s
 => => writing image sha256:ac1ec869614296ba300d64189d6706865396fd9c45caefbb3a9a614dfa1cdd81                                                                                                              0.0s
 => => naming to docker.io/devopps/read-file-and-write-to-sdout:latest                                                                                                                                    0.0s
```

Run it and verify the output because it should match with the [hello-world.txt](./read-file-and-write-to-sdout/hello-world.txt).
```bash
$ docker container run devopps/read-file-and-write-to-sdout:latest
2021/02/13 15:12:42 Content of the file is:  hello world
```

Let's edit this image with the crane by adding a new layer to it, the layer that we are going to add is the same file but with different content. So, if we add the file to the workdir of the image by crane, this code will be going to start to use the file that we add with the layer.

This is the following [content](./layer/hello-world.txt) that we'll add as a final layer of the image.
```text
hello world made by crane
```

Let's take a look at the code.
```golang
package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
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

    // uncomment these lines if you want to save this image to the local image cache - "NEEDS Dockerin this case"
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

```

> IMPORTANT: Before running this code, please shutdown the Docker.

You should notice that, first, we pull the image then we create a layer as tar format that includes my hello-world.txt then we append the new layer to the image.

Lets run this code.
```bash
$ go run -v ./main.go
2021/02/13 18:42:28 image devopps/read-file-and-write-to-sdout:foo pushed to the registry succesfully
```

> IMPORTANT: start the Docker again.

Lets verify the output of the container from the edited image.
```bash
$ docker container run devopps/read-file-and-write-to-sdout:foo
2021/02/13 16:29:09 Content of the file is :  hello world made by crane
```

Tada, it worked.🎉🎉🎉🎉.


> BONUS: crane also can be installed as a cli, go to the [installation page](https://github.com/google/go-containerregistry/blob/main/cmd/crane/README.md#installation) and download it.

