# Database Migrations

## 目录结构

```
migrations/
├── postgres/          # PostgreSQL migration files
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
└── mysql/             # MySQL migration files (TBD)
```

## 设计原则

- **开发环境**：允许 AutoMigrate（由 `configs/admin.yaml` 中 `enable-migrate: true` 控制）。
- **生产环境**：**禁止 AutoMigrate**。release 模式下无论配置如何，AutoMigrate 被强制跳过。
- **生产环境**：使用本目录下的版本化 SQL migration，由 CI/CD 或运维手动执行。

## 生产环境部署流程

1. 在部署前，先执行 SQL migration：
   ```bash
   # PostgreSQL
   psql -U <user> -d <database> -f migrations/postgres/000001_init_schema.up.sql
   ```

2. 启动服务时，`mode` 必须为 `release`，`enable-migrate` 建议显式设为 `false`：
   ```yaml
   system:
     mode: release
     enable-migrate: false
   ```

## 新增 migration

命名规范：`{version:06d}_{description}.up.sql` / `{version:06d}_{description}.down.sql`

例如：
```
000002_add_user_email_index.up.sql
000002_add_user_email_index.down.sql
```

## 初始 schema 生成方式

`000001_init_schema` 通过在空 PostgreSQL 数据库上运行 GORM `AutoMigrate` + `migrate.go` 中的索引修复逻辑后，使用 `pg_dump --schema-only` 导出。

## 注意事项

- `migrate.go` 中的数据修复、重复检查、legacy 迁移等逻辑属于**数据治理**，不属于 schema migration。这些逻辑在 release 模式下不再随启动自动执行，如有需要应转为独立的 admin 工具或一次性脚本。
- `init-data`（种子数据）与 migration 无关，仍由 `configs/admin.yaml` 中 `init-data` 控制。

## 开发环境反复清库

开发阶段需要频繁清空数据库并重跑 migration / seed 时，优先使用显式的 reset 命令，不要手动只删业务表，避免 `schema_migrations` 留下 dirty 或版本错位状态。

```bash
make dev-db-reset
```

默认连接 docker-compose 中的 `gotribe-postgres`。如果使用共享本地 PostgreSQL 容器，可覆盖参数：

```bash
make dev-db-reset DB_CONTAINER=common-postgres DB_USER=develop DB_NAME=develop
```

执行完成后重新启动 Admin，启动流程会从空 schema 重新执行 SQL migration 和初始化数据。

## 初始化数据同步策略

`admin.init_data` 控制是否执行 seed。默认 seed 只补齐缺失数据，不覆盖已有业务数据。

如需把代码中的内置菜单、API 权限、默认角色同步到已有数据库，可临时打开：

```yaml
admin:
  sync_seed_data: true
```

同步完成后应改回 `false`，避免后续重启覆盖后台手工调整过的菜单名称、排序或权限描述。
