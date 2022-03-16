# Running on Railway

[Railway](https://railway.app/) is attempting to build software development infrastructure for humans. It's founded with the core ideology that building software should be "Take what you need, leave what you don't" and that the current iteration of tools for software development is far too complicated for current generations of developers, let alone the ones that come next. As a result, Railway handles your builds, deployments, scaling, and management of infrastructure; from development to production!

## Setup

## container and Database

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/rvH-WB&referralCode=l7uav7)

About the environment variables:

See [Configure L2D](../quickstart.md#configure-l2d)

## Setting Webhook

When Deploy is completed,
1. Access to https://railway.app/dashboard and select line2discord
2. Below line2discord, there will display a domain, please memorize it. (When you access, you will see `Hello, world!`)
3. Access to https://developers.line.biz/console/ and select the channel you created
4. In the `Messaging Api Setting`, set the `Webhook URL` to https://yourdomain.com/webhook
