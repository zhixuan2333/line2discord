# quickstart

## Requirements

- A postgresql database
- A line bot
- A discord bot
- A container

if you what do L2D whitout yout own container or database. Please see [here]()

## Configure L2D

Create an .env file with the following variables:

```
DATABASE_URL=postgresql://postgres:postgres@example.com:5998/dbname
GUILD_ID=
PARENT_ID=
LINE_CHANNEL_SECRET=
LINE_CHANNEL_TOKEN=
DISCORD_TOKEN=
```

### What is GUILD_ID?

The ID is your server id that you want to connect to.

### What is PARENT_ID?

The ID is your category id that you want to connect to.

### How to get LINE_CHANNEL_SECRET and LINE_CHANNEL_TOKEN?

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
11. Disable "Auto-reply messages" and "Greeting messages"
    Enable "Allow bot to join group chats"

### How to get DISCORD_TOKEN?

1. Go to https://discord.com/developers/applications
2. Click to "New Application"
3. Type bot name. And Create it
4. Go to Bot and add Bot
5. Click to Reveal Token
6. Visit it to add to your server 
    [https://discord.com/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=8&scope=bot](https://discord.com/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=8&scope=bot)
