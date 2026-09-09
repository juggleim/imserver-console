## 1. 共享凭据和存储

- [x] 1.1 实现标准库共享验证与 AES-GCM helper，并添加临时密钥离线测试。
- [x] 1.2 添加六列 DAO、20260908 单列幂等迁移及同连接执行保障。
- [x] 1.3 实现事务锁内合并校验、版本递增、非激活凭据保留与 SQL mock 回归。

## 2. API 与界面

- [x] 2.1 扩展现有 API、安全 DTO、上传上限及已有应用访问校验，更新 handler 测试。
- [x] 2.2 扩展现有 iOS 表单、模式条件校验、上传参数及中英文文案，更新前端测试并构建 dist。

## 3. 验证与发布

- [x] 3.1 记录 Secret、共享 DB、发布回滚、轮换及五分钟刷新说明。
- [x] 3.2 运行离线 Go 测试/race/vet/build、npm test/build、OpenSpec strict validate 并记录结果。
- [ ] 3.3 可丢弃 MySQL 验证两仓库 20260908 迁移顺序与重复执行、控制台到运行时读取（需要基础设施）。
- [ ] 3.4 真实 APNs sandbox/production 普通/VoIP、撤销与轮换验收（需要授权凭据）。

## 4. 审查修复

- [x] 4.1 增加 JSON/multipart 乐观版本校验、409 重载 UI 与双管理员回归。
- [x] 4.2 恢复旧 P12 元数据、VoIP-only 和空密码兼容，并同步前端与测试。
- [x] 4.3 新增共享 ValidTopic、统一 100 字符规则与边界/改名回滚测试。
- [x] 4.4 重跑离线后端/race、前端测试/build 和 strict validate，记录结果。

## 5. 用户要求原始 P8 存储修订（历史迭代）

- [x] 5.1 更新本仓库同名 OpenSpec 并 strict validate 后再实现。
- [x] 5.2 删除共享加解密和主密钥依赖，新增 raw DAO 列及仅 ADD 的幂等 20260909.sql。
- [x] 5.3 Save 原样写 raw，元数据验证 raw，旧密文仅提示重传；保留静默 SQL、版本锁、P12、missing-file 行为及安全 DTO。
- [x] 5.4 UI 区分 raw 可用与 legacy 重传，移除主密钥文案，更新中英文、回归测试和 dist。
- [x] 5.5 更新架构/升级/回滚文档，保留历史部署事实；运行离线 full test/race/vet/build、npm test/build 与 strict validate。

第 1 至 4 节属于加密迭代，第 5 节记录首次 raw 修订的已完成工作，均不代表本轮简化已验证或部署。本仓库不做真实 DB、密钥读取、服务器操作或 Git commit。

## 6. P2 审查修复（历史迭代）

- [x] 6.1 Upgrade 返回错误并在版本读取、SQL 执行、版本写入失败时停止；仅初始 MySQL 1146/记录不存在允许引导，独立启动检查错误，补 SQL mock。
- [x] 6.2 P8 切回 P12 允许保留/新上传 VoIP-only 凭据，测试保持真实 _originalAuthType=p8。
- [x] 6.3 离线全量 tests/race、vet/build、前端 tests/build 和 strict validate 复验。

P2 复验：继续使用缓存 Go 1.25.13，`GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off RUN_DB_TEST=0`。全模块 `go test -mod=readonly ./... -count=1`、`go test -mod=readonly -race ./... -count=1`、`go vet -mod=readonly ./...` 和最终 dist 后 `go build -mod=readonly ./...` 全部通过。`npm test` 25/25，通过 `npm run build` 并重建 dist；实现前后 strict validate、`git diff --check` 通过，既有构建警告未变化。

SQL mock 覆盖初始无记录/ErrRecordNotFound/类型化与包装后的 MySQL 1146、其他权限/连接/缺数据库错误、伪装为 1146 的文本错误、非法版本、已到最新版本，以及 20260908 SQL 失败或版本写失败后禁止执行 20260909 和推进版本（含后续 1146 不能当 bootstrap），成功时顺序写两个版本。最早历史 SQL 文件为空，bootstrap 测试在首次版本写处注入错误验证已进入引导且停止，不声称本仓库能创建真实空库。

前端测试不再改写 _originalAuthType 绕过校验，覆盖原 P8 保留 VoIP-only、新上传 VoIP-only、无 P12 拒绝、移除新文件恢复必填及新增 P12 要求不变。主 admingateway 的 Upgrade 返回值检查由主 agent 修改，本轮未改主仓库；未接真实 DB/服务器、读实际密钥、执行厂商推送或提交 Git。

## 7. 未发布功能简化

- [x] 7.1 更新当前 proposal/design/specs 为仅 raw、单一 20260908 六列契约，保留历史测试事实；重新读取 apply instructions 并 strict validate。
- [x] 7.2 删除 P8 密文 DAO 字段、重传错误与 DTO/UI 分支及对应测试；20260908 逐列幂等直接添加 raw，删除 20260909，保留 metadata/版本锁/P12 回归。
- [x] 7.3 清理当前共享包文档，不保留加密兼容或转换接口；运行全量离线 Go 普通/race/build/vet、前端测试/build 重建 dist、strict validate，记录本轮结果。

### 本轮简化验证

2026-09-08，仅修改 console 仓库；重新读取 apply 流程及 CLI instructions，先更新当前规范并 strict validate，再实施简化。当前六列仅为 auth_type、p8_key_id、p8_team_id、p8_private_key BLOB NULL、p8_key_name、config_version，唯一 P8 迁移版本为 20260908。第 5、6 节及下方首轮记录是历史执行事实，不定义当前兼容行为。

- 使用缓存 Go 1.25.13，全部后端命令设置 `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off RUN_DB_TEST=0`，未下载依赖或工具链。
- PASS: 全模块 `go test -mod=readonly ./... -count=1`、`go test -mod=readonly -race ./... -count=1`、`go vet -mod=readonly ./...`，以及最终 dist 重建后的 `go build -mod=readonly ./...`。
- PASS: `npm test` 25/25；`npm run build` 已重建嵌入 dist。既有 Sass/CoreUI 弃用、clipboard.js eval、两处图片路径和大 chunk 警告仍存在。
- PASS: `openspec validate add-apns-p8-auth --strict`、`git diff --check`；源码/schema/tests/当前规范及生成 dist 不含 P8 legacy 字段、重传状态或主密钥兼容实现。
- 回归覆盖 raw 字节/空白原样保存、metadata/缺文件保留、缺失或非法 raw 拒绝、版本锁/HTTP 409、P12 元数据/空密码/VoIP-only/非激活凭据保留、安全 DTO 与 SQL 日志，以及六列独立迁移检查、重复语句执行和失败停止。
- 未运行真实 MySQL、APNs 或浏览器 E2E；3.3/3.4 保持未完成。历史加密部署与推送事实不作为当前 raw 实现验收依据。
- 未修改主仓库，未读取或修改 opencode.jsonc，未读取实际凭据、执行线上操作或 Git commit。

## Raw 首轮验证

2026-09-08，仅在本 console 仓库写入代码和执行测试，不修改主仓库。先更新规格并 strict validate 通过，再实现和复验。

- 使用已缓存 `/Users/yuwnloy/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-amd64/bin/go`，设置 `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off RUN_DB_TEST=0`，没有下载依赖或工具链。系统启动器的 local 模式实际为 Go 1.23.2，自动选择模式在关闭 checksum 服务时失败，故直接调用缓存二进制。
- PASS: `go test -mod=readonly ./... -count=1`。
- PASS: `go test -mod=readonly -race ./... -count=1`，本轮为 console 全模块 race，不是仅五个包。
- PASS: `go vet -mod=readonly ./...`、最终 dist 构建后的 `go build -mod=readonly ./...`。
- PASS: `npm test` 24/24；`npm run build`，已更新嵌入 dist。原有 Sass/CoreUI 弃用、clipboard.js eval、两处图片路径、大 chunk 警告仍存在。
- PASS: 实现前后 `openspec validate add-apns-p8-auth --strict` 与 `git diff --check`。
- 回归覆盖原始 PEM 字节/空白保持、新上传与轮换、raw/legacy 四态、invalid raw 不回退、legacy 重传、安全 JSON/DTO、SQL debug 日志抑制、JSON/multipart 缺文件保留、版本冲突、P12 元数据/VoIP-only/空密码/非激活凭据保留，以及两个迁移的重复语句执行和失败停止。SQL 精确断言保证编辑不更新旧密文列。
- 未运行真实 DB、读取实际 P8/主密钥、连接服务器或执行推送；未递归测试主仓库 pushmanager 外部厂商。未提交 Git。
- 当时 runtime raw loader/digest 与独立转换工具仍待主 agent 完成；当时的 20260909 真实 MySQL、raw 版本 APNs 与桌面/移动浏览器 E2E 均未验收，不能由该轮 mock/构建结果替代。此处仅保留历史验证范围，不定义当前实现或待办。

## 历史验证记录

以下保留加密迭代原记录。后续部署事实以主仓库同名 change 为准：其 4.3 已通过真实控制台/运行时往返，hk-ali sandbox 普通通知由用户确认收到，ali 已部署但未发送；其余 live 验收仍未完成。这些历史事实不表示 raw 修订已部署。

审查修复复验：同一缓存 Go 1.25.13、离线与 RUN_DB_TEST=0 设置下，`go test ./... -count=1`、原五个包的 `go test -race ... -count=1` 均通过；`npm test` 更新为 22/22 通过，`npm run build` 已重新生成 dist，strict validate 通过。新增 SQL mock 覆盖两管理员先轮换后旧 ID 提交、P8 缺版本/P12 显式版本、JSON/multipart HTTP 409、非法版本、旧 VoIP-only/空密码元数据与替换、包名 100/101 和非法改名回滚。UI 重载流程覆盖 draft helper 和源代码契约测试，未宣称真实浏览器 E2E 验收。

2026-09-08，所有命令在控制台仓库执行，未修改主仓库，未读取或修改 opencode.jsonc，未提交 Git。

- 使用缓存 Go 1.25.13，设置 `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off RUN_DB_TEST=0`。
- `go test ./... -count=1` 通过；需真实 DB 的测试按已有开关跳过。
- `go test -race ./commons/apnscredentials ./commons/dbcommons ./dbs ./apis ./services -count=1` 通过。
- `go vet ./...` 与 `go build ./...` 通过。
- `npm test` 18/18 通过；包含 P8 条件字段、Key ID/Team ID、16 KiB、凭据保留和旧 P12 密码只写回归。
- `npm run build` 通过，已生成嵌入 dist。现有 Sass/CoreUI 弃用、clipboard.js eval、两处图片路径和大 chunk 警告未在此变更处理。
- `openspec validate add-apns-p8-auth --strict` 与 `git diff --check` 通过。
- SQL mock 验证六列各自检查与迁移语句顺序、重复调用及错误停止，不替代真实 MySQL DDL 验收。
- 尚未进行桌面/移动浏览器交互和视觉验收；本次沿用现有组件与布局，前端结果限于逻辑测试和生产构建。
- 部署、Secret、共享 DB、刷新、回滚和轮换说明见 `commons/apnscredentials/README.md`。
