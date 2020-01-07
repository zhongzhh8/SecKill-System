package engine

import (
	"SecKill/api"
	"SecKill/conf"
	"SecKill/data"
	"SecKill/middleware/jwt"
	"SecKill/model"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Visible for test
const SessionHeaderKey =  "Authorization"

func SeckillEngine() *gin.Engine {
	router := gin.New()

	// 设置session为Redis存储（但是后来没有用到session，而是用jwt来做用户授权）
	config, err := conf.GetAppConfig()
	if err != nil {
		panic("failed to load redisService config" + err.Error())
	}
	store, _ := redis.NewStore(config.App.Redis.MaxIdle, config.App.Redis.Network,
		config.App.Redis.Address, config.App.Redis.Password, []byte("seckill"))
	router.Use(sessions.Sessions(SessionHeaderKey, store))
	gob.Register(&model.User{})

	// 设置路由（路由只需要严格按照接口文档来写就ok了）
	userRouter := router.Group("/api/users")
	userRouter.POST("", api.RegisterUser) //注册
	userRouter.Use(jwt.JWTAuth())//这些请求都需要通过jwt做用户授权
	{
		userRouter.PATCH("/:username/coupons/:name", api.FetchCoupon)
		userRouter.GET("/:username/coupons", api.GetCoupons)
		userRouter.POST("/:username/coupons", api.AddCoupon)
	}

	authRouter := router.Group("/api/auth") //登录和注销
	{
		authRouter.POST("", api.LoginAuth)
		authRouter.POST("/logout", api.Logout)
	}

	testRouter := router.Group("/test")
	{
		testRouter.GET("/", api.Welcome)
		testRouter.GET("/flush", func(context *gin.Context) {
			if _, err := data.FlushAll(); err != nil {
				println("Error when flushAll. " + err.Error())
			} else {
				println("Flushall succeed.")
			}
		})
	}

	// 启动秒杀功能的消费者（用来异步更新数据库）
	api.RunSecKillConsumer()

	return router
}
