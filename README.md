# manipulate-docker-image-layers-with-crane
Totally inspired by @ahmetb's [latest post about building container images in Go.](ahmet.im/blog/building-container-images-in-go), please read this post before you move on to the hands on section.

# What is crane tool ?

<img src="https://github.com/google/go-containerregistry/raw/main/images/crane.png" height="150"/>

Google has a repository called ["go-containerregistry"](https://github.com/google/go-containerregistry) which provides Go library and CLIs for working with container registries, [crane](https://github.com/google/go-containerregistry/blob/main/cmd/crane/README.md) is one of them. More techically, crane is a tool for interacting with remote images and registries.


