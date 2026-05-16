package routes

import (
	adhandler "gotribe/internal/admin/ad/handler"
	adservice "gotribe/internal/admin/ad/service"
	adscenehandler "gotribe/internal/admin/ad_scene/handler"
	adsceneservice "gotribe/internal/admin/ad_scene/service"
	adminhandler "gotribe/internal/admin/admin_user/handler"
	adminservice "gotribe/internal/admin/admin_user/service"
	aihandler "gotribe/internal/admin/ai/handler"
	aiservice "gotribe/internal/admin/ai/service"
	apihandler "gotribe/internal/admin/api/handler"
	apiservice "gotribe/internal/admin/api/service"
	categoryhandler "gotribe/internal/admin/category/handler"
	categoryservice "gotribe/internal/admin/category/service"
	columnhandler "gotribe/internal/admin/column/handler"
	columnservice "gotribe/internal/admin/column/service"
	commenthandler "gotribe/internal/admin/comment/handler"
	commentservice "gotribe/internal/admin/comment/service"
	confighandler "gotribe/internal/admin/config/handler"
	configservice "gotribe/internal/admin/config/service"
	feedbackhandler "gotribe/internal/admin/feedback/handler"
	feedbackservice "gotribe/internal/admin/feedback/service"
	indexhandler "gotribe/internal/admin/index/handler"
	indexservice "gotribe/internal/admin/index/service"
	jobhandler "gotribe/internal/admin/job/handler"
	menuhandler "gotribe/internal/admin/menu/handler"
	menuservice "gotribe/internal/admin/menu/service"
	operationloghandler "gotribe/internal/admin/operation_log/handler"
	operationlogservice "gotribe/internal/admin/operation_log/service"
	pointhandler "gotribe/internal/admin/point/handler"
	pointservice "gotribe/internal/admin/point/service"
	posthandler "gotribe/internal/admin/post/handler"
	postservice "gotribe/internal/admin/post/service"
	projecthandler "gotribe/internal/admin/project/handler"
	projectservice "gotribe/internal/admin/project/service"
	resourcehandler "gotribe/internal/admin/resource/handler"
	resourceservice "gotribe/internal/admin/resource/service"
	rolehandler "gotribe/internal/admin/role/handler"
	roleservice "gotribe/internal/admin/role/service"
	systemconfighandler "gotribe/internal/admin/system_config/handler"
	systemconfigservice "gotribe/internal/admin/system_config/service"
	taghandler "gotribe/internal/admin/tag/handler"
	tagservice "gotribe/internal/admin/tag/service"
	userhandler "gotribe/internal/admin/user/handler"
	userservice "gotribe/internal/admin/user/service"
	authhandler "gotribe/internal/auth/admin/handler"
	authservice "gotribe/internal/auth/admin/service"

	"gotribe/internal/auth/core"
	coreconfig "gotribe/internal/core/config"
	"gotribe/internal/core/database"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// AdminModules 汇总 Admin 端全部业务模块的 Handler。
type AdminModules struct {
	Auth         *authhandler.Handler
	AI           *aihandler.Handler
	Admin        *adminhandler.Handler
	Role         *rolehandler.Handler
	Menu         *menuhandler.Handler
	API          *apihandler.Handler
	OperationLog *operationloghandler.Handler
	Project      *projecthandler.Handler
	Config       *confighandler.Handler
	Tag          *taghandler.Handler
	Category     *categoryhandler.Handler
	Post         *posthandler.Handler
	User         *userhandler.Handler
	Resource     *resourcehandler.Handler
	Column       *columnhandler.Handler
	AdScene      *adscenehandler.Handler
	Ad           *adhandler.Handler
	Comment      *commenthandler.Handler
	Point        *pointhandler.Handler
	SystemConfig *systemconfighandler.Handler
	Feedback     *feedbackhandler.Handler
	Index        *indexhandler.Handler
	Job          *jobhandler.Handler
}

// BuildAdminModules 装配全部 Admin 业务模块。
func BuildAdminModules(tx *database.TransactionManager, enforcer *casbin.Enforcer, log *zap.SugaredLogger, authManager *core.Manager, cdnDomain string, uploadCfg resourceservice.UploadConfig, redisClient redis.UniversalClient, aiCfg coreconfig.AIConfig) *AdminModules {
	return &AdminModules{
		Auth:         authhandler.NewHandler(core.AudienceAdmin, authservice.NewService(core.AudienceAdmin, tx, authManager), authManager),
		AI:           aihandler.NewHandler(aiservice.NewService(aiCfg)),
		Admin:        adminhandler.NewHandler(adminservice.NewService(tx, enforcer)),
		Role:         rolehandler.NewHandler(roleservice.NewService(tx, enforcer)),
		Menu:         menuhandler.NewHandler(menuservice.NewService(tx)),
		API:          apihandler.NewHandler(apiservice.NewService(tx, enforcer)),
		OperationLog: operationloghandler.NewHandler(operationlogservice.NewService(tx, log)),
		Project:      projecthandler.NewHandler(projectservice.NewService(tx)),
		Config:       confighandler.NewHandler(configservice.NewService(tx)),
		Tag:          taghandler.NewHandler(tagservice.NewService(tx)),
		Category:     categoryhandler.NewHandler(categoryservice.NewService(tx, log)),
		Post:         posthandler.NewHandler(postservice.NewService(tx, log)),
		User:         userhandler.NewHandler(userservice.NewService(tx), cdnDomain),
		Resource:     resourcehandler.NewHandler(resourceservice.NewService(tx, uploadCfg, log), cdnDomain),
		Column:       columnhandler.NewHandler(columnservice.NewService(tx)),
		AdScene:      adscenehandler.NewHandler(adsceneservice.NewService(tx)),
		Ad:           adhandler.NewHandler(adservice.NewService(tx)),
		Comment:      commenthandler.NewHandler(commentservice.NewService(tx)),
		Point:        pointhandler.NewHandler(pointservice.NewService(tx)),
		SystemConfig: systemconfighandler.NewHandler(systemconfigservice.NewService(tx)),
		Feedback:     feedbackhandler.NewHandler(feedbackservice.NewService(tx), cdnDomain),
		Index:        indexhandler.NewHandler(indexservice.NewService(tx, redisClient)),
		Job:          jobhandler.NewHandler(),
	}
}
