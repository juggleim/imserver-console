## Why

控制台需要配合 IM 服务支持 Apple P8 认证，并修复 iOS 凭据回显及事务外合并导致的安全和并发问题。

## What Changes

- 扩展现有 iOS 表单及 API，显式选择 P12/P8，保留非激活凭据。
- 新增标准库共享验证包、逐列幂等的 20260908 六列迁移和事务锁内合并校验，将上传的 P8 原始字节直接保存到 p8_private_key。
- P8 功能未发布，不保留加密兼容、密文列、重传状态或转换接口；仅支持原始 P8 校验保存。
- **BREAKING**: iOS 查询不再返回证书密码；空密码编辑保留原值。Android 契约不变。
- 增加双语 UI、离线测试和嵌入 dist。

## Capabilities

### New Capabilities
- `apns-p8-credentials`: 原始 P8 共享验证契约、只写保存与迁移部署。

### Modified Capabilities
- `multi-package-push-settings`: iOS 认证选择、只写凭据、编辑保留及安全列表。

## Impact

仅修改本仓库 API、DAO、commons、迁移、Vue、测试及文档。关联只读契约为 `/Users/yuwnloy/Documents/gitee/im-server-cluster/openspec/changes/add-apns-p8-auth/`，不修改主仓库或未跟踪 opencode.jsonc。不增加下载接口、不重构全局认证。
