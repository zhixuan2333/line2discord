# quickstart

## 必要条件

-   `postgresql` のデータベース
-   Internet に接続かつ逆接続できるのサーバー
-   LINE のボット（作れます）
-   Discord のボット（作れます）

**_サーバーがない？データベースがない？_** 心配ありません[こちら](./starter/railway.ja.md)を参照してください。

## インストール

[こちら](https://github.com/zhixuan2333/line2discord/releases/latest)から最新のバージョンを取得してください。
※サポートされていない OS は、Issue で報告してください、または Repo を Clone して`go build`で作れます。

## 構成

`.env` ファイルを作成するか、または環境変数を設定してください。

```sh
# .envを作成する場合
$ cp .env.example .env
$ vim .env
```

```ini
DATABASE_URL=postgresql://postgres:postgres@example.com:5998/dbname
GUILD_ID=
PARENT_ID=
LINE_CHANNEL_SECRET=
LINE_CHANNEL_TOKEN=
DISCORD_TOKEN=
```

-   `GUILD_ID`: L2D を使用するサーバーの ID
-   `PARENT_ID`: サーバーにカテゴリーを作成して、その ID を設定してください。
-   `LINE_CHANNEL_SECRET` と `LINE_CHANNEL_TOKEN` について

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

-   `DISCORD_TOKEN` について

    1. Go to https://discord.com/developers/applications
    2. Click to "New Application"
    3. Type bot name. And Create it
    4. Go to Bot and add Bot
    5. Click to Reveal Token
    6. Visit it to add to your server
       [https://discord.com/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=8&scope=bot](https://discord.com/oauth2/authorize?client_id={{CLIENT_ID}}&permissions=8&scope=bot)

## アプリケーションを起動する

```sh
$ ./line2discord
```

## webhook の設定

Deploy が完了したら

1. https://railway.app/dashboard にアクセスして、line2discord を選択する
2. line2discord の下に Domain Name があるので、これをメモしておいてください。（実際にアクセスしたら、`Hello, world!`と表示されます。）
3. https://developers.line.biz/console/ をアクセスして、作成したチャネルを選択する
4. `Messaging Api Setting`で、`Webhook URL`を設定してください。(https://yourdomain.com/webhook と設定してください)
