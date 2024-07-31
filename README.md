#

## deploy
 
need `tinygo`, `npm`

```shell
cd  cmd/ai
npm install
vim wrangler.toml
```

create wrangler.toml with

```yaml
name = "tg-workers-ai" # your worker name
main = "./build/worker.mjs"
compatibility_date = "2024-07-31"

[build]
command = "make build"

[vars]
telegram_token = "telegram_bot_token" # telegram bot token
worker_url = "https://tg-workers-ai.xxxxxx.workers.dev" # your cloudflare worker url
telegram_ids = "10********,293********" # who can access bot, can get by /user_id command, separated by comma

[ai]
binding = "AI"
```

deploy

```shell
make deploy
```

register telegram bot webhook

```shell
curl https://tg-workers-ai.xxxxxx.workers.dev/tgbot/register
```

![image](https://raw.githubusercontent.com/Asutorufa/telegram-workers-ai/main/assets/screenshot.png)
