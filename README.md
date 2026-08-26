# Recommendation Engine

推荐引擎后端服务，纯 Go 标准库实现，零第三方依赖。

## 运行

```bash
cd origin
go run ./cmd/server
```

服务默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 端点

| 实体 | 方法 | 路径 | 说明 |
|------|------|------|------|
| UserProfile | POST | /api/user-profiles | 创建用户画像 |
| UserProfile | GET | /api/user-profiles | 列表（支持 region、tag 筛选） |
| UserProfile | GET | /api/user-profiles/{id} | 详情 |
| UserProfile | PUT | /api/user-profiles/{id} | 更新 |
| UserProfile | DELETE | /api/user-profiles/{id} | 删除 |
| Item | POST | /api/items | 创建物品 |
| Item | GET | /api/items | 列表（支持 category、status、tag、keyword 筛选） |
| Item | GET | /api/items/{id} | 详情 |
| Item | PUT | /api/items/{id} | 更新 |
| Item | PUT | /api/items/{id}/status | 状态流转（online↔offline） |
| Item | DELETE | /api/items/{id} | 删除 |
| BehaviorEvent | POST | /api/behavior-events | 创建行为事件 |
| BehaviorEvent | POST | /api/behavior-events/batch | 批量创建行为事件 |
| BehaviorEvent | GET | /api/behavior-events | 列表（支持 user_id、item_id、event_type 筛选） |
| BehaviorEvent | GET | /api/behavior-events/{id} | 详情 |
| BehaviorEvent | DELETE | /api/behavior-events/{id} | 删除 |
| Strategy | POST | /api/strategies | 创建推荐策略 |
| Strategy | GET | /api/strategies | 列表（支持 name、type、status 筛选） |
| Strategy | GET | /api/strategies/{id} | 详情 |
| Strategy | PUT | /api/strategies/{id} | 更新 |
| Strategy | PUT | /api/strategies/{id}/status | 状态流转（enabled↔disabled） |
| Strategy | DELETE | /api/strategies/{id} | 删除 |
| Recommendation | POST | /api/recommendations | 创建推荐结果 |
| Recommendation | GET | /api/recommendations | 列表（支持 user_id、strategy_id、status 筛选） |
| Recommendation | GET | /api/recommendations/{id} | 详情 |
| Recommendation | PUT | /api/recommendations/{id} | 更新 |
| Recommendation | PUT | /api/recommendations/{id}/status | 状态流转（draft→published→expired） |
| Recommendation | POST | /api/recommendations/generate | 根据策略+特征权重+行为热度生成推荐 |
| Recommendation | DELETE | /api/recommendations/{id} | 删除 |
| FeatureWeight | POST | /api/feature-weights | 创建特征权重 |
| FeatureWeight | GET | /api/feature-weights | 列表（支持 feature、strategy_id、enabled 筛选） |
| FeatureWeight | GET | /api/feature-weights/{id} | 详情 |
| FeatureWeight | PUT | /api/feature-weights/{id} | 更新 |
| FeatureWeight | DELETE | /api/feature-weights/{id} | 删除 |
| Stats | GET | /api/stats/overview | 统计概览（用户数、物品数、事件分布、热门物品、推荐结果分布） |
