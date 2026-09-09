## Purpose

定义控制台与 IM 推送服务共同依赖的 Apple P8 凭据保护协议，保证应用隔离、格式校验、原子更新及升级兼容，并明确失败行为与需要部署环境验证的安全边界。

## ADDED Requirements

### Requirement: Shared raw P8 credentials
系统 SHALL 校验至多 16 KiB 的 PKCS#8 ECDSA P-256 私钥与合法 Key ID/Team ID，并将原始上传字节直接保存到 p8_private_key。共享包 SHALL 提供 Validate/ValidTopic/MaxPrivateKeySize，不依赖主密钥且不提供加解密或转换接口。

#### Scenario: Invalid input
- **WHEN** 私钥格式、大小、曲线或元数据不合法
- **THEN** 系统拒绝保存且不返回或记录任何凭据内容

#### Scenario: Raw bytes are authoritative
- **WHEN** 新上传通过校验或已有 raw 密钥进行元数据编辑
- **THEN** 系统原样写入或保留 p8_private_key，并验证合并后的 raw 配置；保留应用访问校验与静默 SQL

#### Scenario: Missing key
- **WHEN** P8 配置没有已保存的 raw 且未上传新文件
- **THEN** 系统拒绝保存，安全 DTO 的 has_p8_key 仅按 raw 是否非空返回，不由文件名判断

### Requirement: Atomic explicit configuration
系统 SHALL 在同一加锁事务中读取、合并、验证和写入配置，每次成功保存增加版本，未知认证模式被拒绝；未指定认证模式的新增默认 P12，编辑保留原模式。

#### Scenario: Concurrent metadata edit
- **WHEN** 一个编辑未上传文件而另一个编辑刚刚轮换密钥
- **THEN** 旧草稿的 config_version 与锁内版本不匹配时返回 HTTP 409 和安全错误，不合并旧 ID 或修改新密钥；重新加载后才允许提交

#### Scenario: Missing edit version
- **WHEN** 已存 P8 或涉及 P8 的编辑没有提供 config_version
- **THEN** 系统 SHALL 返回 HTTP 409；旧 P12 不涉及 P8 的无版本编辑保持兼容，任何显式版本均须匹配

### Requirement: Shared bundle validation
系统 SHALL 只接受非空、至多 100 ASCII 字符且字符集为 [A-Za-z0-9.-] 的配置包名。VoIP 派生 topic 的 .voip 后缀不计入配置包名长度。

#### Scenario: Invalid rename
- **WHEN** 改名包含内部空格、斜杠、非 ASCII 字符或超过 100 字符
- **THEN** 保存失败且原记录与版本不变；100 字符合法包名可保存

#### Scenario: Explicit rollback
- **WHEN** 用户从 P8 显式切换到 P12
- **THEN** 系统验证保留的 P12 凭据并保存，仍保留 P8 凭据用于后续显式切换

### Requirement: Additive shared schema
系统 SHALL 仅通过 20260908.sql 逐列幂等添加 auth_type、p8_key_id、p8_team_id、p8_private_key BLOB NULL、p8_key_name、config_version 六列，保留所有旧证书列，默认旧记录为 P12 版本 1，并允许控制台和 IM 顺序重复执行各自迁移。未发布的 P8 功能 SHALL 不包含密文列或加密兼容迁移。

#### Scenario: Repeated shared-table migration
- **WHEN** 另一仓库已添加部分或全部新列
- **THEN** 本迁移仅添加缺失列，不清空已有数据

#### Scenario: Migration failure stops startup
- **WHEN** 读取版本、执行 SQL 文件或写入版本失败
- **THEN** Upgrade SHALL 返回错误并停止所有后续迁移，独立控制台不启动 HTTP 服务；仅初始版本查询的 MySQL 1146 或记录不存在允许从版本零引导，不将其他读取错误或非法版本当成空库
