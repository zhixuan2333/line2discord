# Running on Railway

↓ は Railway の紹介文です。

[Railway](https://railway.app/) is attempting to build software development infrastructure for humans. It's founded with the core ideology that building software should be "Take what you need, leave what you don't" and that the current iteration of tools for software development is far too complicated for current generations of developers, let alone the ones that come next. As a result, Railway handles your builds, deployments, scaling, and management of infrastructure; from development to production!

## セットアップ

↓ のボタンをクリックして、Railway にデプロイしてください。

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/rvH-WB?referralCode=l7uav7)

環境変数については、[Configure L2D](../quickstart.ja.md#構成) を参照してください。:

## webhookの設定

Deploy が完了したら
1. https://railway.app/dashboard にアクセスして、line2discord を選択する
2. line2discord の下にDomain Nameがあるので、これをメモしておいてください。（実際にアクセスしたら、`Hello, world!`と表示されます。）
3. https://developers.line.biz/console/ をアクセスして、作成したチャネルを選択する
4. `Messaging Api Setting`で、`Webhook URL`を設定してください。(https://yourdomain.com/webhook と設定してください)
