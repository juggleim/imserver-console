## ADDED Requirements

### Requirement: Channel tabs present package configuration cards
系统 SHALL 在“应用管理 → 推送设置”的每个现有渠道标签页中，以卡片列表展示该应用在当前渠道下的全部包名配置。

#### Scenario: Channel has no configuration
- **WHEN** 用户打开一个尚未配置任何包名的渠道标签页
- **THEN** 页面仅展示一张可点击的“+”新增卡片

#### Scenario: Channel has existing configurations
- **WHEN** 当前渠道存在一组或多组包名配置
- **THEN** 页面为每组配置展示一张配置卡片，并同时展示可继续新增的“+”卡片

#### Scenario: User switches channels
- **WHEN** 用户从一个推送渠道切换到另一个渠道
- **THEN** 页面加载并仅展示新渠道的配置卡片，且不混用前一渠道的表单或文件状态

### Requirement: User can add a package configuration
系统 SHALL 允许用户从当前渠道的新增卡片打开配置弹窗，并 SHALL 按该渠道的实际字段保存一组新的包名配置。

#### Scenario: Add a Huawei configuration
- **WHEN** 用户在华为标签页点击“+”，填写唯一的包名、App ID 和 App Secret，并点击保存
- **THEN** 系统保存该华为配置、关闭弹窗，并在刷新后的华为卡片列表中展示该包名

#### Scenario: Add configuration for another channel
- **WHEN** 用户在小米、OPPO、VIVO、iOS、FCM、极光、荣耀或个推标签页点击“+”
- **THEN** 弹窗展示该渠道现有定义的包名、凭证、选项或文件字段，而不展示其他渠道专属字段

#### Scenario: Required value is missing
- **WHEN** 用户新增配置时未填写包名或该渠道要求的凭证/文件字段
- **THEN** 系统阻止提交并在对应字段附近显示可理解的校验提示

#### Scenario: Duplicate package in the same channel
- **WHEN** 用户尝试在同一应用和同一渠道下新增已存在的包名
- **THEN** 系统拒绝保存、保留弹窗输入，并提示该包名已配置

#### Scenario: Same package in a different channel
- **WHEN** 同一包名已存在于该应用的另一个推送渠道
- **THEN** 系统允许在当前渠道保存该包名配置

### Requirement: Configuration cards show channel-specific details safely
每张配置卡片 SHALL 以包名作为标题，并 SHALL 在正文中展示当前渠道适合展示的实际配置字段；密钥和密码类字段 MUST 以掩码呈现，不得向页面返回明文。

#### Scenario: Display Huawei card
- **WHEN** 华为配置列表加载成功
- **THEN** 每张华为卡片以包名为标题，并在正文显示 App ID 与掩码后的 App Secret

#### Scenario: Display channel-specific card fields
- **WHEN** 非华为渠道的配置列表加载成功
- **THEN** 卡片正文显示该渠道定义的凭证、选项或文件名，并使用对应字段标签

#### Scenario: Display secret values
- **WHEN** 配置包含 App Secret、Master Secret、证书密码或其他密钥类字段
- **THEN** 卡片仅显示统一掩码且浏览器端不能从列表响应中取得原始密钥

### Requirement: User can edit one package configuration
每张配置卡片 SHALL 在底部提供“设置”按钮；系统 SHALL 使用与新增相同的渠道弹窗编辑被选中的配置，并且只更新该配置。

#### Scenario: Open settings for an existing card
- **WHEN** 用户点击某张卡片底部的“设置”按钮
- **THEN** 弹窗展示该卡片的包名和非敏感配置值，敏感输入为空并提示留空将保留原值

#### Scenario: Save changed credentials
- **WHEN** 用户修改一个或多个字段并保存且服务端保存成功
- **THEN** 系统关闭弹窗、重新获取当前渠道列表，并在对应卡片中展示更新后的非敏感信息

#### Scenario: Preserve an unchanged secret
- **WHEN** 用户编辑配置但将敏感输入留空
- **THEN** 服务端保留该配置原有密钥且不会用空值或掩码字符串覆盖它

#### Scenario: Rename a package
- **WHEN** 用户把现有配置的包名改为当前渠道中尚未使用的新包名并保存
- **THEN** 系统以新包名更新同一配置，不额外保留旧包名卡片

#### Scenario: Rename to a duplicate package
- **WHEN** 用户把现有配置的包名改为当前渠道中另一配置已使用的包名
- **THEN** 系统拒绝更新并保持两组原配置不变

### Requirement: File-backed configurations support multiple packages
系统 SHALL 对 iOS 证书配置和 FCM 配置文件按包名独立存储，并 SHALL 在编辑未选择新文件时保留原文件内容。

#### Scenario: Add iOS configurations for multiple packages
- **WHEN** 用户分别为两个不同 iOS 包名上传符合要求的证书并保存
- **THEN** iOS 标签页展示两张独立卡片，且每张卡片关联各自的证书、密码、VoIP 信息和环境

#### Scenario: Edit iOS metadata without replacing certificates
- **WHEN** 用户编辑已有 iOS 配置且未选择新的普通证书或 VoIP 证书
- **THEN** 系统保留对应的原证书字节和文件名

#### Scenario: Replace one iOS certificate
- **WHEN** 用户只选择一个新的普通证书或 VoIP 证书并保存
- **THEN** 系统仅替换所选证书并保留另一证书

#### Scenario: Add FCM configurations for multiple packages
- **WHEN** 用户为不同包名分别上传 FCM 配置文件
- **THEN** FCM 标签页展示独立卡片，并且每张卡片仅关联对应包名的配置文件

### Requirement: Configuration persistence is package-scoped and backward compatible
服务端 SHALL 以应用、规范化渠道和包名作为 Android/FCM 配置的唯一范围，以应用和包名作为 iOS 配置的唯一范围，并 SHALL 让升级前的既有配置可被列表接口读取。

#### Scenario: List multiple Android configurations
- **WHEN** 客户端按应用和渠道请求 Android 配置列表
- **THEN** 服务端返回该范围下所有配置的轻量列表，不包含证书字节或明文密钥

#### Scenario: Existing single configuration after upgrade
- **WHEN** 数据库中存在升级前保存的一条渠道配置
- **THEN** 新列表接口将其作为一张普通配置返回且无需用户重新录入

#### Scenario: Concurrent duplicate creation
- **WHEN** 两个请求同时尝试为同一应用、渠道和包名新增配置
- **THEN** 数据库唯一约束只允许一条记录成功，另一请求返回可识别的冲突错误

#### Scenario: Query fails
- **WHEN** 当前渠道配置列表加载失败
- **THEN** 页面显示查询失败反馈且不把失败误呈现为空配置状态

#### Scenario: Save fails
- **WHEN** 新增或编辑请求保存失败
- **THEN** 页面保留弹窗和用户输入、显示失败反馈，并且不提前修改卡片列表
