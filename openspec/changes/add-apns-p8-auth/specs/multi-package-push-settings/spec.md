## ADDED Requirements

### Requirement: iOS forms support explicit P8 authentication
系统 SHALL 扩展既有 iOS 表单以选择 P12 或 P8，P8 要求 Key ID、Team ID、环境和首次上传 P8 文件，编辑未选文件时保留原文件。P8 模式 SHALL 不要求普通/VoIP P12 文件或密码，双语界面显示保存的 P8 文件名与状态。

#### Scenario: Add P8 configuration
- **WHEN** 用户选择 P8 并提供合法元数据与不超过 16 KiB 的 P8 文件
- **THEN** 前端通过现有上传路径提交 auth_type、p8_key_id、p8_team_id 和 p8_file；服务端验证后原样保存到 p8_private_key

#### Scenario: Edit retained key
- **WHEN** 用户编辑已有 raw P8 配置且不选新文件
- **THEN** 表单不要求重新上传，服务端使用保留密钥校验新元数据和环境

#### Scenario: Missing key requires upload
- **WHEN** 配置无 raw 文件
- **THEN** 中英文卡片和编辑器显示未保存，不将文件名显示为可用密钥；P8 保存要求新文件

### Requirement: iOS edits detect stale drafts and preserve legacy metadata
新 UI SHALL 在 JSON 和 multipart 编辑中提交读取到的 config_version，发生 HTTP 409 时保留输入、显示冲突并提供明确的重新加载操作，不自动用旧草稿覆盖新记录。旧 P12 纯元数据编辑 SHALL 支持 VoIP-only、无文件名及合法空密码记录，编辑留空密码保留原值，不通过非空密码判断已有凭据有效性。新建仍保持普通证书和密码要求，已有记录替换文件验证文件而不强制非空密码。

#### Scenario: Two administrators
- **WHEN** B 更新密钥及 Key ID 后 A 提交旧版本环境编辑
- **THEN** 显示冲突并保留 A 的输入，A 明确重新加载时获得 B 的最新安全元数据和版本，密码及文件选择为空

#### Scenario: Legacy VoIP-only metadata
- **WHEN** 用户只改已有 P12 VoIP-only 或空密码配置的环境或合法包名
- **THEN** 前端允许提交，服务端保留所有凭据并增加版本，不要求新增普通证书或密码

#### Scenario: P8 rollback with VoIP-only credentials
- **WHEN** 用户编辑原认证模式为 P8 的配置并切回 P12，存在保留的 VoIP 证书或新上传的 VoIP 文件
- **THEN** 前端 SHALL 与 DAO 一致允许提交，不强制普通证书或非空密码；无任何 P12 凭据时仍拒绝，新增配置要求不变

## MODIFIED Requirements

### Requirement: Channel forms use the current provider fields

系统 SHALL 按以下渠道字段渲染新增与编辑弹窗，并 SHALL 按必填与条件必填规则进行校验：

| 渠道 | 必填字段 | 可选字段 |
| --- | --- | --- |
| 华为 | 包名、App ID、App Secret | Badge Class |
| 小米 | 包名、App Secret | Channel ID |
| OPPO | 包名、App Key、Master Secret | Channel ID |
| VIVO | 包名、App ID、App Key、App Secret | 无 |
| iOS | 包名、认证类型、环境；P12 首次要求普通证书及密码；P8 首次要求 P8 文件、Key ID、Team ID | P12 VoIP 证书及密码；编辑不替换的密码留空 |
| FCM | 包名、配置文件 | 无 |
| 极光 | 包名、App Key、Master Secret | Classification、Badge Class、华为/小米/荣耀/OPPO/VIVO/魅族渠道参数 |
| 荣耀 | 包名、App ID、App Key、App Secret | Badge Class |
| 个推 | 包名、App ID、App Key、Master Secret | 无 |

#### Scenario: Optional Badge Class is omitted
- **WHEN** 用户新增华为或荣耀配置时未填写 Badge Class，或仅输入空格
- **THEN** 前端提交参数和服务端保存的渠道配置中均不包含 badge_class

#### Scenario: Existing Badge Class is edited
- **WHEN** 用户打开已配置 Badge Class 的华为或荣耀卡片
- **THEN** 弹窗回填 badge_class，保存非空值时去除首尾空格；编辑时将该字段清空则服务端保留原值

#### Scenario: Configure JPush common options
- **WHEN** 用户在极光弹窗填写可选的 Classification 或 Badge Class
- **THEN** 系统分别以整数 options.classification 和与 app_key、master_secret 同级的字符串 badge_class 保存；Classification 未填写时不包含对应键，非整数时阻止提交并提示必须为整数

#### Scenario: Configure JPush third-party channel options
- **WHEN** 用户在极光弹窗切换华为、小米、荣耀、OPPO、VIVO 或魅族标签页并填写参数
- **THEN** 系统按 options.third_party_channel 的小写渠道键保存：华为 importance/category；小米 channel_id/mi_template_id/mi_template_param；荣耀 importance；OPPO distribution/channel_id/category/整数 notify_level/可选 badge_operation_type（0 覆盖、1 增加）/private_msg_template_id/JSON 字符串映射 private_content_parameters/private_title_parameters；VIVO distribution/category/可选整数 Push Mode（0 正式、1 测试）/布尔 add_badge；魅族 distribution

#### Scenario: Omit empty JPush channel options
- **WHEN** 极光渠道参数为空或所有极光可选参数为空
- **THEN** 系统不保存空渠道对象；所有可选参数为空时新增不提交 options，编辑清除原 options

#### Scenario: JPush-only option interface
- **WHEN** 用户打开极光新增或编辑弹窗
- **THEN** 弹窗比其他渠道更宽，Master Secret 后展示 Badge Class 和默认折叠的 Options；折叠高度自适应，展开提升到可用高度且正文可滚动，依次展示 Classification 和六个渠道标签页；其他渠道不展示这些字段及尺寸

#### Scenario: VoIP password is conditionally required
- **WHEN** 用户在 P12 模式选择新的 VoIP 证书文件
- **THEN** 新建配置要求 VoIP 密码；编辑已有配置允许保留合法空密码，未选文件时不要求新增 VoIP 密码

### Requirement: Configuration cards show channel-specific details safely
每张配置卡片 SHALL 以包名为标题并显示渠道字段；App Secret、Master Secret 和证书密码以统一掩码展示。Android 的列表与编辑契约保持不变，iOS 查询 SHALL 不返回任何密码、证书字节或私钥；iOS 编辑密码留空，界面提示留空保留。

#### Scenario: Display Huawei card
- **WHEN** 华为配置列表加载成功
- **THEN** 每张卡片以包名为标题，显示 App ID 和掩码 App Secret

#### Scenario: Display channel-specific card fields
- **WHEN** 非华为配置列表加载成功
- **THEN** 卡片显示该渠道凭证元数据、选项或文件名及对应字段标签；P8 显示认证类型与文件状态

#### Scenario: Display secret values
- **WHEN** 配置包含敏感字段
- **THEN** 卡片仅显示掩码；Android 编辑回填真实值并隐藏，iOS 密码不回填

#### Scenario: Toggle secret visibility
- **WHEN** 密码输入框处于隐藏状态
- **THEN** 输入内容为掩码且眼睛图标带斜线，点击显示当前输入明文并去掉斜线，再次点击恢复隐藏

### Requirement: User can edit one package configuration
每张卡片 SHALL 提供设置按钮，以同一渠道弹窗只编辑选中配置；iOS 的读取、合并、校验和保存 SHALL 在加锁事务中进行。

#### Scenario: Open settings for an existing card
- **WHEN** 用户点击设置
- **THEN** 回填配置，Android 敏感输入回填并隐藏，iOS 密码留空并提示保留，各密码输入独立切换可见性

#### Scenario: Save changed credentials
- **WHEN** 保存修改成功
- **THEN** 关闭弹窗并刷新对应卡片的非敏感信息，iOS config_version 增加

#### Scenario: Preserve an unchanged secret
- **WHEN** 用户编辑时敏感输入留空
- **THEN** 保留原有凭据，不用空值或掩码覆盖

#### Scenario: Rename a package
- **WHEN** 用户改为当前渠道未使用的新包名
- **THEN** 更新同一配置，不保留旧卡片，保留的 raw P8 字节保持不变

#### Scenario: Rename to a duplicate package
- **WHEN** 用户改为另一配置已使用的包名
- **THEN** 拒绝更新并保持原配置不变

### Requirement: Configuration persistence is package-scoped and backward compatible
服务端 SHALL 以应用、规范化渠道、包名作为 Android/FCM 唯一范围，以应用和包名作为 iOS 唯一范围，升级前记录可被列表读取；iOS 未指定认证类型时新增默认 P12，编辑保留原类型。

#### Scenario: List multiple Android configurations
- **WHEN** 按应用和渠道查询 Android 列表
- **THEN** 返回该范围全部轻量配置，包含编辑所需 App Secret/Master Secret，不含文件字节

#### Scenario: List multiple iOS configurations
- **WHEN** 按应用查询 iOS 列表或兼容 get 接口
- **THEN** 返回包名、环境、文件名、auth_type、p8_key_id、p8_team_id、p8_key_name、has_p8_key、config_version；has_p8_key 仅按 raw 列；不返回密码或任何证书/密钥字节

#### Scenario: Existing single configuration after upgrade
- **WHEN** 存在升级前保存的单条配置
- **THEN** 列表作为普通卡片返回，无需重新录入；iOS 默认 P12

#### Scenario: Concurrent duplicate creation
- **WHEN** 同时新增同一应用、渠道和包名
- **THEN** 唯一约束只允许一条成功，另一条返回冲突

#### Scenario: Query fails
- **WHEN** 列表加载失败
- **THEN** 显示失败反馈，不呈现为空列表

#### Scenario: Save fails
- **WHEN** 新增或编辑失败
- **THEN** 保留弹窗输入并显示失败，不提前修改卡片
