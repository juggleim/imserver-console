# multi-package-push-settings Specification

## Purpose

定义“应用管理 → 推送设置”的多包名配置能力，包括九个推送渠道的实际字段、卡片列表、新增与编辑弹窗、敏感字段显示策略、iOS/FCM 文件保留规则，以及服务端按应用、渠道和包名进行持久化的兼容性要求。

## Requirements

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
- **WHEN** 用户在华为标签页点击“+”，填写唯一的包名、App ID 和 App Secret，可选填写 Badge Class，并点击保存
- **THEN** 系统保存该华为配置、关闭弹窗，并在刷新后的华为卡片列表中展示该包名；Badge Class 为空或仅包含空格时，保存的配置中不包含 `badge_class`

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

### Requirement: Channel forms use the current provider fields

系统 SHALL 按以下渠道字段渲染新增与编辑弹窗，并 SHALL 按必填与条件必填规则进行校验：

| 渠道 | 必填字段 | 可选字段 |
| --- | --- | --- |
| 华为 | 包名、App ID、App Secret | Badge Class |
| 小米 | 包名、App Secret | Channel ID |
| OPPO | 包名、App Key、Master Secret | Channel ID |
| VIVO | 包名、App ID、App Key、App Secret | 无 |
| iOS | 包名、普通证书文件、普通证书密码、证书环境 | VoIP 证书文件、VoIP 证书密码 |
| FCM | 包名、配置文件 | 无 |
| 极光 | 包名、App Key、Master Secret | Classification、Badge Class、华为/小米/荣耀/OPPO/VIVO/魅族渠道参数 |
| 荣耀 | 包名、App ID、App Key、App Secret | Badge Class |
| 个推 | 包名、App ID、App Key、Master Secret | 无 |

#### Scenario: Optional Badge Class is omitted

- **WHEN** 用户新增华为或荣耀配置时未填写 Badge Class，或仅输入空格
- **THEN** 前端提交参数和服务端保存的渠道配置中均不包含 `badge_class`

#### Scenario: Existing Badge Class is edited

- **WHEN** 用户打开已配置 Badge Class 的华为或荣耀卡片
- **THEN** 弹窗回填 `badge_class`，保存非空值时去除首尾空格；编辑时将该字段清空则服务端保留原值

#### Scenario: Configure JPush common options

- **WHEN** 用户在极光弹窗填写可选的 Classification 或 Badge Class
- **THEN** 系统分别以整数 `options.classification` 和与 `app_key`、`master_secret` 同级的字符串 `badge_class` 保存；Classification 输入框内提示请输入整数且为选填，字段未填写时不包含对应键，Classification 填写非整数时阻止提交并提示必须为整数

#### Scenario: Configure JPush third-party channel options

- **WHEN** 用户在极光弹窗的渠道可选参数区域切换华为、小米、荣耀、OPPO、VIVO 或魅族标签页并填写参数
- **THEN** 系统按照 `options.third_party_channel` 下对应的小写渠道键保存参数：华为支持 `importance`、`category`；小米支持 `channel_id`、`mi_template_id`、`mi_template_param`；荣耀支持 `importance`；OPPO 支持 `channel_id`、`category`、整数 `notify_level`、可选的 `badge_operation_type` 下拉框（`0` 为覆盖、`1` 为增加）、`private_msg_template_id`，以及 JSON 字符串映射 `private_content_parameters`、`private_title_parameters`；VIVO 支持 `distribution`、`category`、布尔值 `add_badge`；魅族支持 `distribution`

#### Scenario: Omit empty JPush channel options

- **WHEN** 极光的某个渠道参数标签页未填写任何值，或所有极光可选参数均未填写
- **THEN** 系统不保存该空渠道对象；所有可选参数均为空时，新增配置不提交 `options`，编辑配置则清除原有 `options`，最终保存结果不包含 `options`

#### Scenario: JPush-only option interface

- **WHEN** 用户打开极光推送新增或编辑弹窗
- **THEN** 弹窗使用比其他渠道更宽的布局，在 Master Secret 后先展示 Badge Class，再展示默认折叠的 Options；折叠时弹窗高度随基础内容自适应，展开后弹窗提升到视口可用高度并允许正文竖向滚动，同时依次展示 Classification 和六个渠道参数标签页；打开其他推送渠道弹窗时不展示这些极光专属字段、标签页或扩展尺寸

#### Scenario: VoIP password is conditionally required

- **WHEN** 用户为 iOS 配置选择新的 VoIP 证书文件
- **THEN** 系统要求存在 VoIP 证书密码；未选择 VoIP 证书文件时该密码不是新增配置的必填项

### Requirement: Configuration cards show channel-specific details safely

每张配置卡片 SHALL 以包名作为标题，并 SHALL 在正文中展示当前渠道适合展示的实际配置字段；卡片中的 App Secret、Master Secret、普通证书密码和 VoIP 证书密码 MUST 统一显示为 `********`。为了支持管理端编辑，列表接口 SHALL 返回这些敏感字段的真实值，编辑弹窗 SHALL 将其回填到密码输入框并默认隐藏。

#### Scenario: Display Huawei card
- **WHEN** 华为配置列表加载成功
- **THEN** 每张华为卡片以包名为标题，并在正文显示 App ID 与掩码后的 App Secret

#### Scenario: Display channel-specific card fields
- **WHEN** 非华为渠道的配置列表加载成功
- **THEN** 卡片正文显示该渠道定义的凭证、选项或文件名，并使用对应字段标签

#### Scenario: Display secret values

- **WHEN** 配置包含 App Secret、Master Secret、普通证书密码或 VoIP 证书密码
- **THEN** 卡片仅显示统一掩码，编辑弹窗回填真实值并默认以密码形式隐藏

#### Scenario: Toggle secret visibility

- **WHEN** 敏感输入框处于密码隐藏状态
- **THEN** 输入框显示掩码字符且小眼睛图标带斜线；用户点击后输入框显示明文且小眼睛图标去掉斜线，再次点击恢复隐藏状态

### Requirement: User can edit one package configuration
每张配置卡片 SHALL 在底部提供“设置”按钮；系统 SHALL 使用与新增相同的渠道弹窗编辑被选中的配置，并且只更新该配置。

#### Scenario: Open settings for an existing card
- **WHEN** 用户点击某张卡片底部的“设置”按钮
- **THEN** 弹窗展示该卡片的包名和配置值，所有敏感输入均回填并默认隐藏，每个敏感输入均提供独立的小眼睛切换按钮

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
- **THEN** 服务端返回该范围下所有配置的轻量列表，包含编辑回填所需的 App Secret 或 Master Secret 真实值，不包含文件字节

#### Scenario: List multiple iOS configurations

- **WHEN** 客户端按应用请求 iOS 配置列表
- **THEN** 服务端返回全部包名配置、环境、文件名、普通证书密码和 VoIP 证书密码，不返回普通证书或 VoIP 证书的文件字节

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
