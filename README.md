# Campsite

This is an open source version of the [Campsite](https://www.campsite.com/) app. This codebase is available **as-is** and has been refactored to remove AWS dependencies and provide a Go-based backend alternative.

## 🔄 Major Changes

This fork includes significant architectural changes:

- **✅ AWS Dependencies Removed**: S3, ECS, and Transcribe services replaced with local alternatives
- **✅ Go Backend Available**: New Go API with local file storage and WebSocket support  
- **✅ Local Storage**: Files stored locally instead of AWS S3
- **✅ Background Jobs**: Local processing instead of AWS ECS
- **✅ WebSocket Support**: Native implementation instead of Pusher

See [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) for detailed migration instructions.

## 🚀 Quick Start Options

### Option 1: Ruby API with Local Storage (Modified Original)

The Ruby API has been updated to work without AWS:

```bash
cd api
bundle install
script/server
```

### Option 2: Go API (New Implementation)

A new Go-based API with core functionality:

```bash
cd go-api
make setup
make run
```

The Go API provides:
- Organizations, Projects, Posts management
- Local file upload/download
- WebSocket real-time messaging  
- MySQL database support

## Local dev setup

To get started with the **Ruby API** (modified to work without AWS), run this command from the workspace root:

```bash
script/setup
```

For the **Go API**, see the [Go API README](go-api/README.md).

## Service Dependencies

The application has been modified to work with **local alternatives** to cloud services:

### File Storage (Required)

Instead of S3, files are now stored locally:

**Ruby API**: Files stored in `api/storage/uploads/`
**Go API**: Files stored in `go-api/uploads/`

No additional configuration needed - directories are created automatically.

### Real-time Updates

**Ruby API**: Still requires Pusher configuration (see original instructions below)
**Go API**: Uses built-in WebSocket server (no external dependencies)

### Database

Both APIs use MySQL/MariaDB. Follow the original setup instructions for database configuration.

### Optional: Pusher (Ruby API only)

**Ruby API**: [Pusher](https://pusher.com/) is used to send realtime updates and events. After creating your account, configure this under the `pusher` key in your credentials.

**Go API**: Uses built-in WebSocket server (no external dependencies).

### Optional: Imgix (Ruby API only)

**Ruby API**: [Imgix](https://www.imgix.com/) was the CDN powering Campsite. This is now optional - files can be served directly from local storage.

**Go API**: Serves files directly from local storage.

If you want to use Imgix with the Ruby API, you can still set it up following the original instructions.

### 100ms

[100ms](https://www.100ms.live/) powers Campsite's video conferencing features. After creating an account, set these values under the `hms` credentials key.

> [!NOTE]
> Video calls with webhooks will only work when using a secure tunnel such as ngrok.

### OpenAI

We use the OpenAI API to generate summaries for calls and posts. You need both an API key and your organization ID (not the name, found in OpenAI platform settings). Add these under the `openai` credentials key.

### Others

There are many other services we use to power Campsite features. Create accounts and set up credentials as needed:

- `aws_ecs` - run data exports on AWS Elastic Container Service
- `cal_dot_com` - the Campsite [Cal.com](https://cal.com/) app
- `figma` - render frames via the [Figma REST API](https://www.figma.com/developers/api)
- `linear` - the Campsite Linear app (see [API docs](https://developers.linear.app/docs))
- `omniauth_google` - Google OAuth
- `plain` - customer feedback (likely not needed)
- `postmark` - sending emails via the [Postmark API](https://postmarkapp.com/developer)
- `sentry` - bug reports, just need the DSN
- `slack` - Campsite Slack app
- `vercel.revalidate_static_cache` - generate your own key to safely revalidate cached docs (ISR)
- `webpush_vapid` - VAPID keys necessary to send web push notifications ([docs](https://github.com/pushpad/web-push#generating-vapid-keys))
- `zapier` - The Campsite Zapier app
- `tenor` - GIF search ([docs](https://tenor.com/gifapi))

## Running Campsite in development

We use [Overmind](https://github.com/DarthSim/overmind) to run all of the services needed to develop Campsite locally. To get started, run `script/setup` to install tmux and Overmind. Then run:

```shell
script/dev
```

Our [Procfile](https://github.com/campsite/campsite/blob/main/Procfile) defines the services Overmind runs. Use `overmind connect` to view the logs of a specific service (e.g. `overmind connect sync-server`).

## Running the web app

1. Install npm packages

```shell
pnpm install
```

2. Connect Vercel

```shell
pnpm i -g vercel
```

3. Sign in to your Vercel account

```shell
npx vercel login
```

4. Link the apps to Vercel

```shell
npx vercel link --repo
? Set up and develop “~/<your_path_to_campsite>”? [Y/n] y
? Which scope should contain your project? Campsite
? Found project “campsite/campsite”. Link to it? [Y/n] y
🔗  Linked to campsite/campsite (created .vercel)
```

5. Pull environment variables from Vercel

```shell
cd apps/web && npx vercel env pull
```

6. Run the app

```shell
npx vercel dev
```

5. Open the app at `http://app.campsite.test:3000` — you will be redirected to the auth page with user credentials pre-filled.

## Running the marketing site

1. Follow the steps above to connect the `apps/site` repo to the `campsite-site` project on Vercel.

2. Pull environment variables from Vercel

```shell
cd apps/site && npx vercel env pull
```

3. Run the site

```shell
pnpm -F @campsite/site dev
```

4. Open the app at `http://localhost:3003`

## Running Storybook

```shell
pnpm -F @campsite/web storybook
```

## Installing new packages

To add new packages, you always install at the root of the repository and pass a filter flag to pnpm. For example, to install Tailwind to the marketing site, you would run:

```shell
pnpm -F @campsite/site i tailwindcss
```

Or to install lodash to the web app, you would run:

```shell
pnpm -F @campsite/web i lodash
```

To find an app or package name, navigate to the `package.json` for each `app/` or `package/` directory.

## Running the Desktop App locally

1. [Download the ToDesktop build macOS App](https://dl.todesktop.com/2108257l5oobyoe).
2. Look for the app switcher in the ToDesktop title bar. If you do not see our org apps here, you may need to create a dummy application.
3. Switch to "Campsite Dev" and click "Run" in the bottom left corner of the app.

<details>
  <summary>View screenshot</summary>

  <img width="927" alt="Screen Shot 2023-01-11 at 15 20 46@2x" src="https://user-images.githubusercontent.com/1923260/211939117-f9a2a6cf-99c2-4220-b732-b109c5c859e9.png">
</details>

## Running the API locally

Are you looking to setup and run the API locally? Checkout the [API README](api/README.md)

✨

## API secrets

Running `script/setup` will generate [Rails custom credential files](https://edgeguides.rubyonrails.org/security.html#custom-credentials) for development and production environments. From the API directory, run `VISUAL="code --wait" bin/rails credentials:edit --environment development` to edit development credentials, or run `VISUAL="code --wait" bin/rails credentials:edit --environment development` to edit production credentials.

When you deploy Campsite, `config/credentials/production.yml.enc` must be present on your production server, and you must set the `RAILS_MASTER_KEY` environment variable equal to the contents of `config/credentials/production.key`. If you deploy from a git repository, you can remove `config/credentials/production.yml.enc` from `.gitignore` and commit the file. Do not remove `config/credentials/production.key` from `.gitignore`, instead share this with your team in a safe place outside of the repository, such as in a password vault.

## Calls in development

To use calls in development, you'll need to treat the `http://app.campsite.test:3000` origin as secure. In a Chromium-based browser, you can do that at chrome://flags/#unsafely-treat-insecure-origin-as-secure.

For any call features that require webhooks (like getting chat bubbles in message threads when calls start), you should [use ngrok](https://github.com/campsite/campsite/tree/main/api#using-ngrok-for-publicly-accessible-https-development-urls). You can enable webhooks in the [developer tab of the 100ms dashboard](https://dashboard.100ms.live/developer). The webhook URL should be set to `https://api-dev.campsite.com/v1/integrations/hms/events`, and there should be a webhook header with the name `X-Passcode` and a value that comes from running `Rails.application.credentials.hms.webhook_passcode` in a development Rails console. Be sure that you configure webhooks in the "Development" workspace and disable them when you're done working in ngrok (webhooks in the "Production" workspace should always be enabled and the configuration shouldn't change).
