package common

// 内置默认 RBAC 模型（作为字符串常量内置）
// 与仓库根的 rbac_model.conf 保持一致，用作兜底默认
var embeddedRBACModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`
