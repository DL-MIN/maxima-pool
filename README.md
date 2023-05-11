# Moodle Maxima Pool

A simple serverside MaximaPool implementation to deal with requests from [moodle-qtype_stack](https://github.com/maths/moodle-qtype_stack).


## Features

- Automatically fetch maxima scripts from [moodle-qtype_stack](https://github.com/maths/moodle-qtype_stack)
- Supports multiple plugin versions
- Prebuild maxima snapshots
- Supports *HTTP Basic Auth* and API token via HTTP header


## Requirements

- golang >= 1.19
- maxima


## Build

```shell
go mod tidy
go build -a -buildmode=exe -trimpath
```


## Usage

Adjust the given configuration file `config.yaml` to your desired values. The search paths for the configuration file are `./config.yaml` and `/etc/maxima-pool/config.yaml`, otherwise define the path with the `-config` flag.

```shell
./Moodle_Maxima_Pool -config /path/to/config.yaml
```

**Important:** Build snapshots before the first run:

```shell
./Moodle_Maxima_Pool -config /path/to/config.yaml -create-snapshots
```
