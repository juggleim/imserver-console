# admin-language-preference Specification

## Purpose

定义管理后台首次访问时的默认语言、语言切换行为以及用户语言偏好的浏览器本地持久化规则，确保用户再次打开管理后台时能够恢复上一次主动选择的界面语言。

## Requirements

### Requirement: Management console defaults to English

管理后台 SHALL 在没有有效本地语言偏好时优先使用浏览器首个受支持的语言，并 SHALL 在浏览器语言均不受支持时使用 English，同时使用英文作为缺失文案的回退语言。

#### Scenario: First visit has a supported browser language

- **WHEN** 用户打开管理后台、本地没有已保存的语言偏好，且浏览器语言包含管理后台支持的语言
- **THEN** 管理后台使用浏览器语言列表中首个受支持的语言显示，但不自动写入语言偏好

#### Scenario: First visit has no supported browser language

- **WHEN** 用户打开管理后台、本地没有已保存的语言偏好，且浏览器语言均不在管理后台支持范围内
- **THEN** 管理后台以 English 显示，但不自动写入语言偏好

#### Scenario: Saved preference is invalid

- **WHEN** 本地保存的语言值不是管理后台支持的语言
- **THEN** 管理后台忽略该值，继续按浏览器受支持语言优先、English 兜底的规则选择语言

### Requirement: User language selection persists locally

管理后台 SHALL 在用户主动选择支持的语言后将该语言保存到浏览器本地缓存，并 SHALL 在后续访问时优先恢复该选择。

#### Scenario: User selects Simplified Chinese

- **WHEN** 用户在语言选择器中选择简体中文
- **THEN** 页面立即切换为简体中文，并将 `zh-CN` 保存到本地缓存

#### Scenario: User returns after selecting a language

- **WHEN** 用户之前选择的有效语言仍存在于本地缓存并再次打开管理后台
- **THEN** 管理后台恢复该语言，而不是使用默认 English
