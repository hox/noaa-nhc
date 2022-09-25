# 🌴 NHC Atlantic Outlook Discord Webhook

When the National Hurricane Center posts a new [hurricane / tropical weather outlook](https://www.nhc.noaa.gov/), a notification will be sent to a discord channel via a webhook. Last outlook publish-date is stored in redis to prevent sending duplicate messages of the same outlook.

# Docker

You can run this project with docker by using the image `ghcr.io/hox/noaa-nhc`

# Build/Run

To build/run this project use the following commands.

```
go run .
```

or

```
go build
./noaa-nhc
```

# Environment Variables

- Use `WEBHOOK_TOKEN` for your Discord Webhook token. ex: `https://discord.com/api/webhooks/xxxxxxxx/xxxxxxxxxxxxxxxxxxxx`
- Use `REDIS_DSN` for the DSN of your running Redis server. ex: `redis://127.0.0.1:6379`

# In use

![](https://cdn.eli.tf/mRmnVjKZdDa7.png)

# Side-note

This is my first project built in Go, I just wanted to have fun with it so please don't judge lol :)
