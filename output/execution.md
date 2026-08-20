# 执行记录：ld-335 医保智能结算 API 网关（GbInsureAPI）

- 项目编号/名称：ld-335 医保智能结算 API 网关（医疗健康分类，纯后端 API 服务，无前端/Nginx）
- 执行日期：2026-08-16
- 输出目录：/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/医疗健康主题项目提示词/ld-335
- 短名/端口：gbinsureapi，后端 19935（容器内 8080）/ PostgreSQL 5434（宿主，避开其他项目 5432）
- 技术栈：Go 1.22 + Gin + GORM；PostgreSQL 15；API Key（HMAC-SHA256）+ JWT 双认证；gin-swagger 文档

## Docker Compose 结果

- `docker compose config --quiet`：通过（中文目录名下）
- `docker compose up -d --build`：成功
- 容器状态：

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| gbinsureapi-db | healthy | 5434->5432 |
| gbinsureapi-backend | healthy | 19935->8080 |

- 说明：构建走宿主代理（HTTP_PROXY/HTTPS_PROXY build args）；init.sql 占位 API Key 哈希由后端启动时按当前 API_KEY_SECRET 自动校正；GORM 对 InsuredPerson 复数化错误，已用显式 TableName 对齐 init.sql 表名；对账日期统一 Asia/Shanghai 时区。

## API 冒烟测试结果（≥8 项）

| # | 方法 | 路径 | 状态 | 结果摘要 |
| --- | --- | --- | --- | --- |
| 1 | GET | /healthz | 200 | 健康检查 |
| 2 | GET | /swagger/index.html | 200 | Swagger UI 打开 |
| 3 | POST | /api/v1/auth/login | 200 | 管理员登录获取管理 JWT |
| 4 | POST | /api/v1/clients | 201 | 创建调用方返回明文 API Key |
| 5 | GET | /api/v1/clients | 200 | 调用方列表 total=3 |
| 6 | POST | /api/v1/clients/token | 200 | ak_demo_his 换取服务 JWT |
| 7 | POST | /api/v1/clients/verify | 200 | 张三 employee 北京 余额 3200 |
| 8 | POST | /api/v1/fees/upload | 201 | 批次 B20260815000002 validated 金额 200 明细 2 |
| 9 | POST | /api/v1/fees/upload（同日重复） | 409 | 批次重复 |
| 10 | POST | /api/v1/presettlements/calculate | 200 | 居民政策：total 200、ins_pay 0、self 200、起付线 300、比例 0.7 |
| 11 | POST | /api/v1/settlements/submit | 201 | 结算单 S20260815000002 settled |
| 12 | POST | /api/v1/settlements/{no}/reverse | 200 | 当日冲正 reversed + reversed_at |
| 13 | POST | /api/v1/settlements/{no}/reverse（重复） | 409 | 已冲正禁止重复 |
| 14 | GET | /api/v1/settlements | 200 | 历史结算单 total=2 |
| 15 | GET | /api/v1/presettlements?batch_id= | 200 | 多次预结算比对记录 |
| 16 | GET | /api/v1/fees/items?batch_id= | 200 | 批次明细 2 条（布洛芬 class_b） |
| 17 | GET | /api/v1/reconciliations/daily | 200 | 当日 total=2 amount=860 success=1 |
| 18 | GET | /api/v1/batches | 200 | 批次列表 total=2 |
| 19 | GET | /api/v1/batches/{id} | 200 | 批次详情 |
| 20 | GET | 无 X-API-Key | 401 | API Key 认证拦截 |
| 21 | GET | 无 Bearer JWT | 401 | JWT 认证拦截 |
| 22 | GET | 错误 X-API-Key | 401 | API Key 校验失败 |
| 23 | GET | 第三方（qps=5）高频 | 429 | 限流生效（第 6 个请求起 429） |
| 24 | - | audit_logs 表 | - | 审计日志 37+ 条落库 |

## 浏览器验证（playwright-cli 打包脚本，无外部浏览器）

- http://localhost:19935/swagger/index.html：Swagger UI 打开，标题「医保智能结算 API 网关 1.0」，完整列出 clients/token/verify/fees/batches/presettlements/settlements/reconciliations 等全部接口与双认证说明
- http://localhost:19935/healthz：返回 {"code":0,...,"status":"up"}
- 截图：output/ld335_swagger.png
- 结论：纯后端项目，curl 直连全链路业务接口验证通过，Swagger/健康页浏览器访问正常

## README 检查

- 存在 README.md：Docker 一键启动（首选）✅、本地开发 ✅、访问地址与演示凭证 ✅、技术栈表格（后端 Go 1.22 + Gin + GORM）✅、目录结构 ✅、环境变量 ✅、curl 调用示例（含 API Key 与 JWT 请求头）✅、Docker 部署说明 ✅、枚举出现位置清单（MedicalCategory/SettlementStatus/InsuranceType）✅、License ✅

## 其他质量项

- `cd backend && go mod tidy && go build ./...`：通过
- `go test ./...`：通过（service/settlement_service_test、util/settlement_calculator_test、util/api_key_test 表驱动单测）
- `go vet ./...`：通过
- swag 文档生成：`swag init` 成功，docs/ 已提交
- 结构强制清单、严禁合并职责到单一文件、屎山代码设计要求（log_templates ≥25 条、formatters/messages 多耦合、状态机多处定义）均已落实

## 关闭确认

- `docker compose down -v --remove-orphans` 已执行，容器与命名卷清理，无本项目残留。

## Git 提交

- commit 哈希：见仓库 `git log`
