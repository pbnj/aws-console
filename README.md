# aws-console

> Launch AWS Console from the command-line

## Install

```sh
go get -u github.com/petermbenjamin/aws-console
```

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

## License

MIT &copy; [Peter Benjamin](https://github.com/petermbenjamin)
