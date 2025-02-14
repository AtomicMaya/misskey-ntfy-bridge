# misskey->ntfy webhook bridge

Because i don't like my notifications to be JSON blobs, here is a utility to convert your Misskey/Sharkey webhooks into ntfy compatible blurbs! 

This is for people like me that do not like having an `up_*` type channel on their ntfy instance.

Contributions welcome, especially if you want to do things with attachments or other shenanigans!

## Compatible software

- [Misskey](https://github.com/misskey-dev/misskey-hub-next)
- [Sharkey](https://activitypub.software/TransFem-org/Sharkey)

## Setup

Rename `default.env` to `.env` to store your environment variables.

Use one of the provided binaries on the [releases page](https://github.com/AtomicMaya/misskey-ntfy-bridge/releases) (or build from source)

### Run from source

`go run app/main.go`

### Build from source

```shell
GCO_ENABLED=0 go build -tags netgo -a -ldflags "-w" -o build/misskey-ntfy-bridge-latest ./app
```

## Deploy

### docker

See the provided template docker file in [./Dockerfile](./Dockerfile)

### docker compose

See the provided template compose file in [./compose.yaml](./compose.yaml)

### systemd (why would you do this to yourself?)

See the provided template systemd service file in [./misskey-ntfy-bridge.service](./misskey-ntfy-bridge.service)

## License 

Licensed under [EUPL-1.2](https://interoperable-europe.ec.europa.eu/sites/default/files/custom-page/attachment/2020-03/EUPL-1.2%20EN.txt)

See [./LICENSE.md](./LICENSE.md)

## Notes

built with programming socks and e.s. motion 2020 gear