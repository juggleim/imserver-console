## Why

实时消息统计后端接口和应用级采集开关已经存在，但当前管理后台没有对应入口，管理员只能通过手工请求查看数据。需要参考 JuggleChat Console 的交互，在当前管理后台补齐可配置、可查询的实时消息统计页面。

## What Changes

- 在数据统计菜单中新增实时消息统计入口和独立路由。
- 新增实时消息统计页面，支持单聊、群聊、聊天室切换。
- 支持最近 15 分钟至 3 天的快捷范围以及自定义日期时间范围。
- 接入现有 `apps/statistic/msgrealtime` 接口，以折线图展示上行、下行、分发的每秒平均消息数，并补齐无数据时间桶。
- 在页面中读取和更新 `open_real_time_msg_statistic` 应用配置开关。
- 补齐中英文文案、前端测试和嵌入式 `dist` 构建产物。

## Capabilities

### New Capabilities
- `realtime-message-statistics-ui`: 管理员启停实时消息采集并按频道和时间范围查看实时消息统计。

### Modified Capabilities

无。

## Impact

- 前端路由、菜单、统计 API service、应用配置调用、ECharts 页面及双语资源。
- 新增纯逻辑测试并更新 `npm test` 显式测试列表。
- 重新生成 `webconsole/web/dist`，供 Go `embed` 打包。
- 不修改现有实时统计后端接口、数据库结构和普通按天统计页面。
