# misskey->ntfy webhook bridge

Because i don't like my notifications to be JSON blobs, here is a utility to convert your Misskey/Sharkey webhooks into ntfy compatible blurby! 

Contributions welcome, especially if you want to do things with attachments or other shenanigans!

## Setup

Use the `default.env` to give you some indication on how to build your environment variables.

Use one of the provided binaries in `/build` (or build from source)

Setup a systemd unit to keep it running (like below)

```
[Unit]
Description=Bridge to ntfy for Misskey notifications
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=<FILL IN YOUR USER>
ExecStart=<PATH TO EXEC>
WorkingDirectory=<PATH TO WORKING DIR>

Environment="HOST="
Environment="PORT="
Environment="ORIGIN_URL="
Environment="NTFY_URL="
Environment="NTFY_CHANNEL="
Environment="NTFY_TOKEN="

[Install]
WantedBy=multi-user.target
```