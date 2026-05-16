package errs

// MsgKey 统一消息键类型。
type MsgKey string

// 通用操作消息键
const (
	MsgLoginSuccess        MsgKey = "login_success"
	MsgLogoutSuccess       MsgKey = "logout_success"
	MsgRefreshTokenSuccess MsgKey = "refresh_token_success"
	MsgListSuccess         MsgKey = "list_success"
	MsgListFail            MsgKey = "list_fail"
	MsgGetSuccess          MsgKey = "get_success"
	MsgGetFail             MsgKey = "get_fail"
	MsgCreateSuccess       MsgKey = "create_success"
	MsgCreateFail          MsgKey = "create_fail"
	MsgUpdateSuccess       MsgKey = "update_success"
	MsgUpdateFail          MsgKey = "update_fail"
	MsgDeleteSuccess       MsgKey = "delete_success"
	MsgDeleteFail          MsgKey = "delete_fail"
	MsgJWTAuthFail         MsgKey = "jwt_auth_fail"
	MsgUserSerializeFail   MsgKey = "user_serialize_fail"

	MsgUserNotFound      MsgKey = "user_not_found"
	MsgUserDisabled      MsgKey = "user_disabled"
	MsgUserRoleDisabled  MsgKey = "user_role_disabled"
	MsgPasswordIncorrect MsgKey = "password_incorrect"
	MsgUserNotLoggedIn   MsgKey = "user_not_logged_in"
	MsgUserNotFoundByID  MsgKey = "user_not_found_by_id"
	MsgNoUsersFound      MsgKey = "no_users_found"
	MsgRoleInfoFailed    MsgKey = "role_info_failed"
	MsgNoUsersWithRole   MsgKey = "no_users_with_role"

	MsgGetRoleApisFailed    MsgKey = "get_role_apis_failed"
	MsgLoadRolePolicyFailed MsgKey = "load_role_policy_failed"
	MsgUpdateRoleApisFailed MsgKey = "update_role_apis_failed"
	MsgDeleteRoleApisFailed MsgKey = "delete_role_apis_failed"

	MsgGetApiInfoFailed    MsgKey = "get_api_info_failed"
	MsgUpdateApiFailed     MsgKey = "update_api_failed"
	MsgLoadApiPolicyFailed MsgKey = "load_api_policy_failed"
	MsgGetApiListFailed    MsgKey = "get_api_list_failed"
	MsgNoApiListFound      MsgKey = "no_api_list_found"
	MsgDeleteApiFailed     MsgKey = "delete_api_failed"
)

var i18nMessages = map[string]map[MsgKey]string{
	"zh": {
		MsgLoginSuccess:         "登录成功",
		MsgLogoutSuccess:        "退出成功",
		MsgRefreshTokenSuccess:  "刷新token成功",
		MsgListSuccess:          "获取列表成功",
		MsgListFail:             "获取列表失败",
		MsgGetSuccess:           "获取成功",
		MsgGetFail:              "获取失败",
		MsgCreateSuccess:        "创建成功",
		MsgCreateFail:           "创建失败",
		MsgUpdateSuccess:        "更新成功",
		MsgUpdateFail:           "更新失败",
		MsgDeleteSuccess:        "删除成功",
		MsgDeleteFail:           "删除失败",
		MsgJWTAuthFail:          "JWT认证失败",
		MsgUserSerializeFail:    "用户信息序列化失败",
		MsgUserNotFound:         "用户不存在",
		MsgUserDisabled:         "用户被禁用",
		MsgUserRoleDisabled:     "用户角色被禁用",
		MsgPasswordIncorrect:    "密码错误",
		MsgUserNotLoggedIn:      "用户未登录",
		MsgUserNotFoundByID:     "未获取到ID为%d的用户",
		MsgNoUsersFound:         "未获取到任何用户信息",
		MsgRoleInfoFailed:       "根据角色ID获取角色信息失败",
		MsgNoUsersWithRole:      "根据角色ID未获取到拥有该角色的用户",
		MsgGetRoleApisFailed:    "获取角色的权限接口失败",
		MsgLoadRolePolicyFailed: "角色的权限接口策略加载失败",
		MsgUpdateRoleApisFailed: "更新角色的权限接口失败",
		MsgDeleteRoleApisFailed: "删除角色关联权限接口失败",
		MsgGetApiInfoFailed:     "根据接口ID获取接口信息失败",
		MsgUpdateApiFailed:      "更新权限接口失败",
		MsgLoadApiPolicyFailed:  "权限接口策略加载失败",
		MsgGetApiListFailed:     "根据接口ID获取接口列表失败",
		MsgNoApiListFound:       "根据接口ID未获取到接口列表",
		MsgDeleteApiFailed:      "删除权限接口失败",
	},
	"en": {
		MsgLoginSuccess:         "Login successful",
		MsgLogoutSuccess:        "Logout successful",
		MsgRefreshTokenSuccess:  "Refresh token successful",
		MsgListSuccess:          "Get list successfully",
		MsgListFail:             "Failed to get list",
		MsgGetSuccess:           "Get successfully",
		MsgGetFail:              "Failed to get",
		MsgCreateSuccess:        "Create successful",
		MsgCreateFail:           "Failed to create",
		MsgUpdateSuccess:        "Update successful",
		MsgUpdateFail:           "Failed to update",
		MsgDeleteSuccess:        "Delete successful",
		MsgDeleteFail:           "Failed to delete",
		MsgJWTAuthFail:          "JWT authentication failed",
		MsgUserSerializeFail:    "User information serialization failed",
		MsgUserNotFound:         "User not found",
		MsgUserDisabled:         "User is disabled",
		MsgUserRoleDisabled:     "User role is disabled",
		MsgPasswordIncorrect:    "Incorrect password",
		MsgUserNotLoggedIn:      "User not logged in",
		MsgUserNotFoundByID:     "User not found by ID %d",
		MsgNoUsersFound:         "No users found",
		MsgRoleInfoFailed:       "Failed to get role information by role ID",
		MsgNoUsersWithRole:      "No users found with this role",
		MsgGetRoleApisFailed:    "Failed to get role APIs",
		MsgLoadRolePolicyFailed: "Failed to load role policy",
		MsgUpdateRoleApisFailed: "Failed to update role APIs",
		MsgDeleteRoleApisFailed: "Failed to delete role APIs",
		MsgGetApiInfoFailed:     "Failed to get API information by ID",
		MsgUpdateApiFailed:      "Failed to update API",
		MsgLoadApiPolicyFailed:  "Failed to load API policy",
		MsgGetApiListFailed:     "Failed to get API list by ID",
		MsgNoApiListFound:       "No API list found by ID",
		MsgDeleteApiFailed:      "Failed to delete API",
	},
}

// T 根据语言返回消息键对应的翻译文本。
func T(lang string, key MsgKey) string {
	if m, ok := i18nMessages[lang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	return i18nMessages["zh"][key]
}
