<div align="center">
  <img src="webconsole/web/src/assets/images/login/logo.png" alt="JuggleIM" width="240" />

  <h1>IMServer Console</h1>

  <p><strong>面向 JuggleIM 的开源、自托管 IM 管理控制台</strong></p>
  <p>用一个 Web 后台统一管理应用、用户、群组、消息、机器人、第三方服务、监控与数据分析。</p>

  <p>
    <a href="https://github.com/juggleim/imserver-console/stargazers"><img src="https://img.shields.io/github/stars/juggleim/imserver-console?style=flat-square&logo=github&label=Stars" alt="GitHub Stars" /></a>
    <a href="https://github.com/juggleim/imserver-console/releases"><img src="https://img.shields.io/github/v/release/juggleim/imserver-console?style=flat-square&label=Release" alt="最新版本" /></a>
    <a href="https://github.com/juggleim/imserver-console/blob/master/LICENSE"><img src="https://img.shields.io/github/license/juggleim/imserver-console?style=flat-square" alt="开源协议" /></a>
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.25+" />
    <img src="https://img.shields.io/badge/Vue-3-42b883?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue 3" />
  </p>

  <p>
    <a href="README.md">English</a> ·
    <a href="https://www.juggle.im/">官网</a> ·
    <a href="https://www.juggle.im/docs/guide/intro/">文档</a> ·
    <a href="https://github.com/juggleim/im-server">IM 服务端</a> ·
    <a href="https://gitee.com/juggleim/imserver-console">Gitee</a>
  </p>
</div>

---

IMServer Console 是 [JuggleIM](https://github.com/juggleim/im-server) 自托管即时通讯系统的运维管理层。项目将 Vue 3 管理后台、Go 管理 API 和 API 网关打包为一个可部署服务，帮助团队免去从零开发 IM 管理后台的成本。

> 如果这个项目为你节省了时间，欢迎点一个 ⭐ **Star**。你的支持能让更多开发者发现它，也是项目持续维护的重要动力。

## 为什么选择 IMServer Console？

- **日常运维一站完成**：集中管理应用、账号、用户、群组、机器人、会话和历史消息。
- **面向生产环境设计**：提供敏感词、消息拦截器、Webhook、客户端日志和分角色账号管理。
- **第三方能力统一配置**：集中接入 iOS/Android 推送、对象存储、RTC、短信、邮件和翻译服务。
- **数据和状态直观可见**：查看用户活跃、消息量、连接数和节点性能指标。
- **部署链路简单**：编译后的前端内嵌在 Go 服务中，无需额外部署 Web 服务器。
- **开放且易于扩展**：采用 Apache-2.0 协议，API、服务和数据访问分层清晰。

## 核心功能

| 模块 | 能力 |
| --- | --- |
| 应用管理 | 创建/导入应用、服务开关、回调地址和应用凭证配置 |
| 用户与账号 | 管理员账号、应用权限、用户查询与封禁、群组和机器人管理 |
| 消息治理 | 会话查询、历史消息查询与撤回/删除、敏感词和自定义拦截规则 |
| 推送与存储 | APNs、FCM/Android 推送、文件存储服务和客户端日志收集 |
| 通信服务 | Agora、ZEGO、LiveKit RTC，以及短信、邮件和翻译配置 |
| 统计与监控 | 用户活跃、单聊/群聊/聊天室消息、连接数和节点性能 |
| 开发者工具 | 在控制台调试 IM API、检查连接状态 |
| 国际化 | 内置简体中文和英文界面 |

## 快速开始

### 环境要求

- Go 1.25+
- MySQL
- 已运行的 [JuggleIM 服务端](https://github.com/juggleim/im-server)（参考[一键部署文档](https://www.juggle.im/docs/guide/deploy/quickdeploy/)）

### 1. 克隆项目

国内推荐使用 Gitee：

```bash
git clone https://gitee.com/juggleim/imserver-console.git
cd imserver-console
```

或使用 GitHub：

```bash
git clone https://github.com/juggleim/imserver-console.git
cd imserver-console
```

### 2. 创建数据库

```sql
CREATE DATABASE jim_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;
```

服务启动时会自动执行数据表迁移。

### 3. 修改配置

编辑 [`conf/config.yml`](conf/config.yml)：

```yaml
port: 8091
adminSecret: "请替换为高强度随机密钥"

log:
  logPath: ./logs
  logName: imserver-console

mysql:
  user: root
  password: 你的数据库密码
  address: 127.0.0.1:3306
  name: jim_db

imApiDomain: http://127.0.0.1:9001
imAdminDomain: http://127.0.0.1:8090
```

`imApiDomain` 是 IM 服务 API 地址，`imAdminDomain` 是 JuggleIM 管理 API 地址。

### 4. 启动服务

```bash
go run .
```

浏览器访问 **http://127.0.0.1:8091**，首次登录账号为：

```text
账号：admin
密码：123456
```

> 非本地环境务必配置 `adminSecret`，并在首次登录后立即修改默认密码。请勿将生产环境密码提交到 `conf/config.yml`。

## 前端开发

生产所需前端资源已经内嵌在 Go 服务中。仅在开发管理界面时，才需要单独启动 Vue 项目：

```bash
cd webconsole/web
npm ci
npm run dev
```

Vite 开发服务器默认将 `/admingateway` 代理到 `http://127.0.0.1:8090`。完成前端修改后重新构建内嵌资源：

```bash
npm run build
```

## 架构

```text
浏览器
   │
   ▼
Vue 3 管理后台（内嵌）
   │  /admingateway
   ▼
Gin API + 鉴权 + API 网关
   ├── MySQL（控制台配置与统计数据）
   └── JuggleIM API（IM 管理操作与运行数据）
```

## 项目结构

```text
.
├── apis/           # HTTP 接口和请求模型
├── services/       # 业务逻辑
├── dbs/            # GORM 数据访问层
├── commons/        # 配置、鉴权、日志、数据库迁移和工具包
├── routers/        # Gin 路由和编译后的管理后台资源
├── webconsole/     # Vue 3 + Vite 控制台及 Go 内嵌加载器
├── conf/           # 运行配置
└── main.go         # 应用入口
```

## JuggleIM 生态

- [im-server](https://github.com/juggleim/im-server)：高性能、自托管的开源 IM 服务端
- [官方文档](https://www.juggle.im/docs/guide/intro/)：部署、集成、客户端 SDK 和服务端指南
- [服务端 API](https://www.juggle.im/docs/server/api/)：用户、群组、消息、聊天室等接口文档

## 参与贡献

欢迎提交 Issue 和 Pull Request。适合作为首次贡献的内容包括：Bug 修复、文档完善、新增第三方服务适配、测试补充和界面体验优化。

1. Fork 本仓库。
2. 创建功能分支。
3. 补充必要测试，并验证前后端构建。
4. 提交说明清晰的 Pull Request；界面改动请附截图。

## 开源协议

本项目基于 [Apache License 2.0](LICENSE) 开源。

---

<div align="center">
  <strong>更轻松地搭建并运营属于自己的 IM 平台。</strong><br />
  <sub>如果 IMServer Console 对你有帮助，欢迎点亮 ⭐ Star 支持项目。</sub>
</div>
