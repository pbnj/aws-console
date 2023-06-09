# aws-console

> Launch AWS Console from the command-line

## Overview

Launch AWS Console from the command-line.

## Install

```sh
go install github.com/pbnj/aws-console@latest
```

## Usage

```sh
$ aws-console -h
Usage of aws-console:
  -h Print Help
  -p string
     AWS Profile
  -v Print Version
```

Example, assuming you have a `~/.aws/config` with a profile named `my-profile`:

```sh
aws-console -p my-profile
```

## Changelog

### 0.1.0

- Upgrade to aws-sdk-go-v2
- Switch to go modules

### 0.0.1

#### Added

- Initial Release

## License

MIT
