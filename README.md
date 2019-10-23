# icheck

A Go program to batch-check images downloaded if completed.

## Installation

```
$ go get github.com/yashiro1899/icheck
```

## Usage

```
NAME:
   icheck - find out incomplete images in paths

USAGE:
   icheck [paths...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --quiet, -q     print incomplete images only
   --sniffing, -s  determine the image type of the first 32 bytes of data
   --help, -h      show help
   --version, -v   print the version
```

## Examples

![example](https://raw.githubusercontent.com/yashiro1899/icheck/master/example.png)
