<div align="center">
  <img src="webconsole/web/src/assets/images/login/logo.png" alt="JuggleIM" width="240" />

  <h1>IMServer Console</h1>

  <p><strong>An open-source, self-hosted admin console for operating JuggleIM.</strong></p>
  <p>Manage apps, users, groups, messages, bots, integrations, monitoring, and analytics from one web UI.</p>

  <p>
    <a href="https://github.com/juggleim/imserver-console/stargazers"><img src="https://img.shields.io/github/stars/juggleim/imserver-console?style=flat-square&logo=github&label=Stars" alt="GitHub Stars" /></a>
    <a href="https://github.com/juggleim/imserver-console/releases"><img src="https://img.shields.io/github/v/release/juggleim/imserver-console?style=flat-square&label=Release" alt="Latest release" /></a>
    <a href="https://github.com/juggleim/imserver-console/blob/master/LICENSE"><img src="https://img.shields.io/github/license/juggleim/imserver-console?style=flat-square" alt="License" /></a>
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.25+" />
    <img src="https://img.shields.io/badge/Vue-3-42b883?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue 3" />
  </p>

  <p>
    <a href="README.zh-CN.md">简体中文</a> ·
    <a href="https://www.juggle.im/">Website</a> ·
    <a href="https://www.juggle.im/docs/guide/intro/">Documentation</a> ·
    <a href="https://github.com/juggleim/im-server">IM Server</a> ·
    <a href="https://gitee.com/juggleim/imserver-console">Gitee</a>
  </p>
</div>

---

IMServer Console is the operations layer for [JuggleIM](https://github.com/juggleim/im-server), a self-hosted instant messaging platform. It packages a Vue 3 admin UI, Go management APIs, and an API gateway into one deployable service—so teams can operate an IM system without building an admin backend from scratch.

> If this project saves you time, please give it a ⭐ **Star**. It helps more developers discover the project and supports continued maintenance.

## Why IMServer Console?

- **One place for daily operations** — apps, accounts, users, groups, bots, conversations, and historical messages.
- **Production-oriented controls** — sensitive words, message interceptors, webhooks, client logs, and role-based accounts.
- **Integrations without scattered config** — iOS/Android push, object storage, RTC, SMS, email, and translation providers.
- **Built-in visibility** — user activity, message statistics, connection counts, and node performance metrics.
- **Simple deployment** — the compiled frontend is embedded in the Go service; no separate web server is required.
- **Open and extensible** — Apache-2.0 licensed, with clear API, service, and data-access layers.

## Features

| Area | Capabilities |
| --- | --- |
| Application management | Create/import apps, configure service switches, callbacks, and app credentials |
| Users and accounts | Admin accounts, app permissions, user search, ban/unban, groups, and bots |
| Message operations | Conversation inspection, history search, recall/delete, sensitive words, and custom interceptors |
| Push and storage | APNs, FCM/Android push, file storage providers, and client log collection |
| Communication services | RTC providers (Agora, ZEGO, LiveKit), SMS, email, and translation |
| Analytics and monitoring | User activity, private/group/chatroom messages, connections, and node performance |
| Developer tools | IM API debugging and connection inspection from the console |
| Internationalization | Built-in English and Simplified Chinese UI |

## Quick Start

### Prerequisites

- Go 1.25+
- MySQL
- A running [JuggleIM server](https://github.com/juggleim/im-server) (see the [deployment guide](https://www.juggle.im/docs/guide/deploy/quickdeploy/))

### 1. Clone the repository

```bash
git clone https://github.com/juggleim/imserver-console.git
cd imserver-console
```

For users in mainland China:

```bash
git clone https://gitee.com/juggleim/imserver-console.git
cd imserver-console
```

### 2. Create the database

```sql
CREATE DATABASE jim_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;
```

Tables are migrated automatically when the service starts.

### 3. Configure the console

Edit [`conf/config.yml`](conf/config.yml):

```yaml
port: 8091
adminSecret: "replace-with-a-strong-random-secret"

log:
  logPath: ./logs
  logName: imserver-console

mysql:
  user: root
  password: your_mysql_password
  address: 127.0.0.1:3306
  name: jim_db

imApiDomain: http://127.0.0.1:9001
imAdminDomain: http://127.0.0.1:8090
```

`imApiDomain` is the IM service API address; `imAdminDomain` is the JuggleIM admin API address.

### 4. Run

```bash
go run .
```

Open **http://127.0.0.1:8091** and sign in with:

```text
Username: admin
Password: 123456
```

> For any non-local deployment, set `adminSecret` and change the default password immediately after the first login. Do not commit production credentials to `conf/config.yml`.

## Frontend Development

The production frontend is already embedded in the Go binary. Run the Vue app separately only when working on the web UI:

```bash
cd webconsole/web
npm ci
npm run dev
```

The Vite dev server proxies `/admingateway` to `http://127.0.0.1:8090` by default. To rebuild the embedded frontend:

```bash
npm run build
```

## Architecture

```text
Browser
   │
   ▼
Vue 3 Admin UI (embedded)
   │  /admingateway
   ▼
Gin API + Auth + API Gateway
   ├── MySQL (console configuration and statistics)
   └── JuggleIM APIs (IM operations and runtime data)
```

## Project Structure

```text
.
├── apis/           # HTTP handlers and request models
├── services/       # Business logic
├── dbs/            # GORM data-access layer
├── commons/        # Config, auth, logging, migrations, and utilities
├── routers/        # Gin routes and compiled admin assets
├── webconsole/     # Vue 3 + Vite admin console and Go embed loader
├── conf/           # Runtime configuration
└── main.go         # Application entry point
```

## JuggleIM Ecosystem

- [im-server](https://github.com/juggleim/im-server) — high-performance, self-hosted IM server
- [Official documentation](https://www.juggle.im/docs/guide/intro/) — deployment, integration, SDK, and server API guides
- [Server API reference](https://www.juggle.im/docs/server/api/) — integrate users, groups, messages, chatrooms, and more

## Contributing

Issues and pull requests are welcome. Good first contributions include bug fixes, documentation improvements, new integration providers, tests, and UI refinements.

1. Fork the repository.
2. Create a feature branch.
3. Add tests where appropriate and verify the web/backend build.
4. Open a pull request with a clear description and screenshots for UI changes.

## License

Released under the [Apache License 2.0](LICENSE).

---

<div align="center">
  <strong>Build and operate your own IM platform with confidence.</strong><br />
  <sub>If IMServer Console is useful to you, a ⭐ Star is the simplest way to support it.</sub>
</div>
