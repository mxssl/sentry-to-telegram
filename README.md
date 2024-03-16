# sentry-to-telegram

This is a simple app to send Sentry alerts to a Telegram group.

## Usage

1. Create a Telegram bot and get the bot token.
2. Create a Telegram group and add the bot with permissions to post messages.
3. Get the chat ID of the group.
4. Set the environment variables `TELEGRAM_BOT_TOKEN` and `TELEGRAM_GROUP_ID` with the bot token and group ID.
5. Deploy the app and set the Sentry webhook to the app URL.
6. Done! Now you will receive alerts from Sentry in your Telegram channel.

Example of a `docker-compose.yaml` file:

```yaml
version: '3'
services:
  sentry-to-telegram:
    image: mxssl/sentry-to-telegram:v0.0.4
    ports:
      - "9999:9999"
    restart: always
    environment:
      TELEGRAM_BOT_TOKEN: ""
      TELEGRAM_CHAT_ID: ""
```
