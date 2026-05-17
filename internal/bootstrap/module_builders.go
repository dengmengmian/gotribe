// module_builders.go provides builder functions that assemble each business module
// from shared infrastructure dependencies.
package bootstrap

import (
	categoryhandler "gotribe/internal/api/category/handler"
	categoryrepo "gotribe/internal/api/category/repository"
	confighandler "gotribe/internal/api/config/handler"
	configrepo "gotribe/internal/api/config/repository"
	configservice "gotribe/internal/api/config/service"
	examplehandler "gotribe/internal/api/example/handler"
	examplerepo "gotribe/internal/api/example/repository"
	exampleservice "gotribe/internal/api/example/service"
	healthhandler "gotribe/internal/api/health/handler"
	healthservice "gotribe/internal/api/health/service"
	posthandler "gotribe/internal/api/post/handler"
	postrepo "gotribe/internal/api/post/repository"
	postservice "gotribe/internal/api/post/service"
	profilehandler "gotribe/internal/api/profile/handler"
	profilerepo "gotribe/internal/api/profile/repository"
	profileservice "gotribe/internal/api/profile/service"
	taghandler "gotribe/internal/api/tag/handler"
	tagrepo "gotribe/internal/api/tag/repository"
	usereventhandler "gotribe/internal/api/user_event/handler"
	usereventrepo "gotribe/internal/api/user_event/repository"
	usereventservice "gotribe/internal/api/user_event/service"
	"gotribe/internal/auth/core"
	authhandler "gotribe/internal/auth/user/handler"
	authrepo "gotribe/internal/auth/user/repository"
	authservice "gotribe/internal/auth/user/service"
)

// HealthModule 汇总 health 模块对外暴露的能力。
type HealthModule struct {
	Handler *healthhandler.Handler
}

// AuthModule 汇总 auth 模块对外暴露的能力。
type AuthModule struct {
	Handler *authhandler.Handler
}

// ProfileModule 汇总 profile 模块对外暴露的能力。
type ProfileModule struct {
	Service *profileservice.Service
	Handler *profilehandler.Handler
}

// PostModule 汇总 post 模块对外暴露的能力。
type PostModule struct {
	Service *postservice.Service
	Handler *posthandler.Handler
}

// ExampleModule 汇总 example 模块对外暴露的能力。
type ExampleModule struct {
	Handler *examplehandler.Handler
}

// TagModule 汇总 tag 模块对外暴露的能力。
type TagModule struct {
	Handler *taghandler.Handler
}

// CategoryModule 汇总 category 模块对外暴露的能力。
type CategoryModule struct {
	Handler *categoryhandler.Handler
}

// ConfigModule 汇总 config 模块对外暴露的能力。
type ConfigModule struct {
	Handler *confighandler.Handler
}

// UserEventModule 汇总 user_event 模块对外暴露的能力。
type UserEventModule struct {
	Handler *usereventhandler.Handler
}

// Modules 汇总全部业务模块依赖。
type Modules struct {
	Health    HealthModule
	Auth      AuthModule
	Profile   ProfileModule
	Post      PostModule
	Tag       TagModule
	Category  CategoryModule
	Config    ConfigModule
	Example   ExampleModule
	UserEvent UserEventModule
}

func buildModules(infra *Infra) *Modules {
	health := buildHealthModule(infra)
	auth := buildAuthModule(infra)
	profile := buildProfileModule(infra)
	post := buildPostModule(infra)
	tag := buildTagModule(infra)
	category := buildCategoryModule(infra)
	config := buildConfigModule(infra)
	example := buildExampleModule(infra, post.Service)
	userEvent := buildUserEventModule(infra)

	return &Modules{
		Health:    health,
		Auth:      auth,
		Profile:   profile,
		Tag:       tag,
		Category:  category,
		Config:    config,
		Post:      post,
		Example:   example,
		UserEvent: userEvent,
	}
}

func buildHealthModule(infra *Infra) HealthModule {
	service := healthservice.NewService(infra.DB, infra.Redis, infra.AppName)
	return HealthModule{
		Handler: healthhandler.NewHandler(service),
	}
}

func buildAuthModule(infra *Infra) AuthModule {
	userRepo := authrepo.NewUserAuthRepository(infra.Tx)
	service := authservice.NewService(core.AudienceUser, userRepo, infra.AuthTokens, infra.JWT)
	return AuthModule{
		Handler: authhandler.NewHandler(service),
	}
}

func buildProfileModule(infra *Infra) ProfileModule {
	repo := profilerepo.NewRepository(infra.Tx)
	service := profileservice.NewService(core.AudienceUser, infra.UserAuth.AccessTokenTTL(), repo, infra.Store, infra.AuthTokens, infra.Tx, infra.CacheTTL)
	return ProfileModule{
		Service: service,
		Handler: profilehandler.NewHandler(service),
	}
}

func buildPostModule(infra *Infra) PostModule {
	repo := postrepo.NewRepository(infra.Tx)
	tagRepository := tagrepo.NewRepository(infra.Tx)
	categoryRepository := categoryrepo.NewRepository(infra.Tx)
	service := postservice.NewService(repo, tagRepository, categoryRepository, infra.Store, infra.CacheTTL)
	return PostModule{
		Service: service,
		Handler: posthandler.NewHandler(service),
	}
}

func buildExampleModule(infra *Infra, posts exampleservice.PostSummaryReader) ExampleModule {
	repo := examplerepo.NewRepository(infra.Tx)
	service := exampleservice.NewService(repo, infra.Tx, posts)
	return ExampleModule{
		Handler: examplehandler.NewHandler(service),
	}
}

func buildTagModule(infra *Infra) TagModule {
	repo := tagrepo.NewRepository(infra.Tx)
	return TagModule{
		Handler: taghandler.NewHandler(repo),
	}
}

func buildCategoryModule(infra *Infra) CategoryModule {
	repo := categoryrepo.NewRepository(infra.Tx)
	return CategoryModule{
		Handler: categoryhandler.NewHandler(repo),
	}
}

func buildConfigModule(infra *Infra) ConfigModule {
	repo := configrepo.NewRepository(infra.Tx)
	service := configservice.NewService(repo)
	return ConfigModule{
		Handler: confighandler.NewHandler(service),
	}
}

func buildUserEventModule(infra *Infra) UserEventModule {
	repo := usereventrepo.NewRepository(infra.Tx)
	service := usereventservice.NewService(repo)
	return UserEventModule{
		Handler: usereventhandler.NewHandler(service),
	}
}
