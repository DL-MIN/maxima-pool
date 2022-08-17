# Maxima Pool

A simple serverside Maxima pool implementation to deal with requests from [moodle-qtype_stack](https://github.com/maths/moodle-qtype_stack).


## Requirements

- golang
- maxima


## Build

```shell
go build maxima-pool.go
```


## Usage

All listed options are optional and have a reasonable default value, usage information is available via `-h`.

```shell
./maxima-pool -host "0.0.0.0" -port "8000" -maxima "/usr/bin/maxima" -timeout "10s" -tmp "/tmp"
```
