# golog-tools
[![Build Status](https://travis-ci.org/weldpua2008/golog-tools.svg?branch=master)](https://travis-ci.org/weldpua2008/golog-tools) Log tools is the repo with the following tools:
* trace calls in logs

### Build

```
go build cmd/logtracer/logtracer.go
```

### Run

```
asomeprogram | logtracer
```

### Demo
![Demo](demo.gif)

##### Gif was generated with:

```
docker run --rm -v $PWD:/data asciinema/asciicast2gif -s 2 -t solarized-dark demo.cast demo.gif
```
