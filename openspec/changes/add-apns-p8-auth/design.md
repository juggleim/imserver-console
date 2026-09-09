## Context

参见 proposal.md。主仓库同名 change 已固定跨仓库协议。本仓库已有 multi-package-push-settings 规格，当前 iOS 编辑先在 handler 读旧值再进入 DAO 锁，查询还回传密码。

## Goals / Non-Goals

**Goals:** P12 兼容、P8 显式认证、应用隔离、原始 P8 只写存储、原子编辑、安全查询、既有 UI 扩展。

**Non-Goals:** 主仓库运行时实现、全局认证重构、密钥下载、自动轮换或 APNs 在线权限校验。

## Decisions

1. `commons/apnscredentials` 仅使用标准库，提供 `Validate(privateKey []byte, keyID, teamID string) error`、`ValidTopic(topic string) bool` 和 `MaxPrivateKeySize`。P8 限制 16 KiB，必须是单一 PKCS#8 PEM ECDSA P-256，Key ID/Team ID 为 Apple 的 10 位大写字母数字标识。不提供加解密或转换接口。
2. `P8PrivateKey []byte` 映射 GORM 列 `p8_private_key`，JSON `-`。上传内容不修剪、不重新编码，原样保存。P8 功能未发布，不需要密文存储兼容。
3. DAO 接收增量字段及临时 P8 文件，事务内 `FOR UPDATE` 读取后合并、校验并写入，config_version 每次保存递增。无 auth_type 时新增默认 p12、编辑保留；空字符串不删除密码/元数据，缺失文件保留 raw；非激活模式材料原样保留。P8 元数据保存验证保留的 raw，无自动回落 P12。继续静默 SQL，防止原始密钥进入参数日志。
4. 现有 set/upload/get/list 路径继续使用。安全 DTO 输出 p8_key_name/has_p8_key/config_version 和非敏感配置，has_p8_key 仅按 raw 非空，不由文件名判断；密码、证书和密钥均不返回。保留现有应用权限检查。
5. 仅使用 20260908.sql 逐列幂等添加 auth_type、p8_key_id、p8_team_id、p8_private_key BLOB NULL、p8_key_name、config_version 六列。使用 information_schema、PREPARE 和固定连接，保留 P12 列和版本机制。
6. 扩展原有 Vue 配置定义和弹窗，P8 必须有新上传或 raw 可用状态。编辑密码始终留空，中英文显示原始文件已保存或未保存。展示数据库/备份保护和五分钟刷新提示，不引入新框架。

## Risks / Trade-offs

审查补充：JSON/multipart 编辑传入读取到的 config_version，DAO 锁内在合并前比较。所有传入版本必须匹配；已存 P8、切入 P8 或编辑 P8 材料时缺版本也返回 HTTP 409（body code 409），旧 P12 无版本编辑兼容。UI 不自动重试，保留输入并允许明确丢弃草稿、重载最新配置。旧 P12 纯元数据编辑不新增普通证书/文件名/非空密码门槛；新建保留原上传要求，已有记录替换文件验证非空文件和文件名，不将空密码视为无效证书。共享新增 ValidTopic(topic string) bool：原始配置 bundle 必须匹配 [A-Za-z0-9.-]，非空且至多 100 ASCII 字节；运行时应在添加 .voip 前校验 bundle，不对派生 topic 套用 100 字节限制。

- 数据库读权限和备份现在能够访问原始 P8 → 运维必须保护数据库和备份访问，API 只写不代表静态加密。
- 不支持 P8 的旧二进制不能使用 P8 密钥 → 降级前须显式切回保留的 P12。
- 两仓库重复执行共享表迁移 → 逐列幂等；迁移仍须运维串行执行。
- 本地校验不证明 Apple 授权 → sandbox/production/VoIP 保留真实环境验收。
- 不改全局认证 → 仅复用已有应用级权限逻辑，不扩大权限体系。

## Migration Plan

先执行 20260908 迁移并核实六列，协调升级所有发送节点与控制台后再配置 P8。保存最多五分钟后在后续发送生效。回滚先显式选保留的 P12，等待刷新或重启后再降级，不删列。离线测试用临时生成的密钥与 sqlmock，RUN_DB_TEST=0，不接触真实基础设施或实际密钥。

## Historical Deployment

2026-09-08 的前一轮实现使用 AES-GCM 和部署 Secret。主仓库同名 change 记录 hk-ali 与 ali 的部署，以及 hk-ali sandbox 普通通知已由用户确认收到。这些事实及 tasks 中标注的加密迭代测试不证明当前 raw 实现已部署或验收。用户明确功能未发布，本轮直接简化当前 schema，不保留试验版本兼容；不连接服务器或操作远程 Secret。
