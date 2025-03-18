# ComHel

Docker `Com`pose `Hel`per is just a command helper to help run docker compose up and down

# Installation

```bash
go install github.com/FaisalBudiono/comhel@latest
```

# Logs

You can see the logs on `[home-dir]/.comhel/logs.log`

# Config

You can change your config ENV on `[home-dir]/.comhel/.env`

| Name             | Value                                | Description                                                                                                  |
| ---------------- | ------------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| COMHEL_DEV_MODE  | `true` or `false` (default: `false`) | When run on devmode, the program will use the .env and logging in the directory where the program is called. |
| COMHEL_LOG_LEVEL | any string (default: `debug`)        | Will log in debug mode if set to `debug`. The default log level is `warning`.                                |
