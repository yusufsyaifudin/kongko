# KONGKO

kong.ko

bentuk tidak baku: kongkow

_v Jk_ bercakap-cakap yang tidak ada artinya; mengobrol


## Background

This is sample code used in Qiscus TechTalk #97 presentation at EDS, Jl. Asem Kranji K-7, Sekip, Yogyakarta. You can see slide presentation here [z]()

## Prerequisites
- GO 1.8
- EMQTT 2.2
- [Glide package management](http://glide.sh)


## Using precompiled binary
Instead installing dependencies and build it manually, you can use precompiled binary if you are using Ubuntu (this binary compiled in Ubuntu 16.04, but may still running in higher version).

```bash
$ ./out/kongko.bin --addr localhost:9000 --ws-addr ws://localhost:8083/mqtt
```

## Install, Build and Run

Install using glide,

```bash
$ glide install
```

Then build it,

```bash
$ go build -o out/kongko.bin -i cmd/kongko/main.go
```

Finally, run it:

```bash
$ ./out/kongko.bin --addr localhost:9000 --ws-addr ws://localhost:8083/mqtt
```

