# go-env

Runs a command with environment gotten from [etcd](https://github.com/coreos/etcd).

## Usage

To start a new bash with environment from etcd use:

```bash
ETCD_URL=http://127.0.0.1:4001 ETCD_ENV=appname go-env bash
```

It'll read environment variables from both `"/environment"` and `"/appname/environment"` directories.
`"/appname/environment"` takes precedense.


## License

See [LICENSE](LICENSE) file.
