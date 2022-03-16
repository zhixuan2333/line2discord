# quickstart

## Requirements

-   A postgresql database
-   A server
-   A line bot
-   A discord bot
-   A container

No server? No database? [Look here](./starter/railway.md)

## Install

Get lastest version from [Release](https://github.com/zhixuan2333/line2discord/releases/latest)

## Configure L2D

Create an `.env` file.

```sh
cp .env.example .env
```

You will get this.

```ini
DATABASE_URL=postgresql://postgres:postgres@example.com:5998/dbname
GUILD_ID=
PARENT_ID=
LINE_CHANNEL_SECRET=
LINE_CHANNEL_TOKEN=
DISCORD_TOKEN=
```

-   `GUILD_ID`: Your discord server id.
-   `PARENT_ID`: Your discord category channel id.
-   `LINE_CHANNEL_SECRET` and `LINE_CHANNEL_TOKEN`

    1. Go to https://developers.line.biz/
    2. Click to "Console"
    3. Login as your LINE account
    4. Select your name at Providers
    5. Click to "Create a new channel"
    6. Type of "Message API"
    7. Type some info for your Channel
    8. Click to "Create"
    9. Copy the Channel Secret at Basic setting
    10. Copy the Channel Token at Message API
    11. Disable "Auto-reply messages" and "Greeting messages" Enable "Allow bot to join group chats"

-   `DISCORD_TOKEN`

    1. Go to https://discord.com/developers/applications
    2. Click to "New Application"
    3. Type bot name. And Create it
    4. Go to Bot and add Bot
    5. Click to Reveal Token
    6. Visit it to add to your server
       [https://discord.com/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=8&scope=bot](https://discord.com/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=8&scope=bot)

## Start the application

```
./line2discord
```

## Setting Webhook

When Deploy is completed,

1. Access to https://railway.app/dashboard and select line2discord
2. Below line2discord, there will display a domain, please memorize it. (When you access, you will see `Hello, world!`)
3. Access to https://developers.line.biz/console/ and select the channel you created
4. In the `Messaging Api Setting`, set the `Webhook URL` to https://yourdomain.com/webhook

