package handler

// 任务管理路由由 internal/admin/routes/routes.go 的 registerJobRoutes 统一注册：
// 在 /api/jobs 组上统一施加 jwt + adminLoader + Casbin 中间件。
// 本模块不再单独定义 RegisterRoutes，避免与 registerJobRoutes「双重定义」造成路由漂移。
