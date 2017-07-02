# aws-console

> Launch AWS Console from the command-line

## The Problem

There is no easy way to launch, or go to, the AWS console of a particular account, short of saving the page as a browser bookmark.

## This Solution

This CLI application allows users, who are typically familiar with the command-line, to launch AWS console from the command-line.

This application relies on the AWS credentials file stored under the current user's home directory (e.g. `$HOME/.aws/credentials`).

## Install

- With Go:
  ```sh
  go get -u github.com/petermbenjamin/aws-console
  ```

- [Precompiled binaries][download-link]

## Usage

```sh
$ aws-console --help
Usage: aws-console <options>
Options:
  -c, --credentials string
        Path to AWS credentials file (default "/Users/pbenjamin/.aws/credentials")
  -d, --debug
        Debug
  -h, --help
        Print Help
  -p, --profile string
        AWS Profile (default "default")
```

### Examples

- To launch AWS console using default profile credentials:
  ```sh
  aws-console
  ```
- To launch AWS console using some-other-account profile:
  ```sh
  aws-console --profile=some-other-account
  ```
- If your `.aws/credentials` file is stored in a different location:
  ```sh
  aws-console --credentials="/path/to/aws/cred/file"
  ```

## Changelog

### 0.0.1

#### Added

- Initial Release

## License

MIT &copy; [Peter Benjamin](https://github.com/petermbenjamin)

[download-link]: https://github.com/petermbenjamin/aws-console/releases/latest/
