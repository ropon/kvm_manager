package routers

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/controllers"
	_ "github.com/ropon/kvm_manager/docs"
	"github.com/ropon/kvm_manager/utils"
	"github.com/ropon/logger"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// setupRouter 初始化路由
func setupRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(ginLogger())
	engine.Use(cors())
	engine.Use(utils.TraceHttpRoot(conf.SERVERNAME, conf.Cfg.External["JaegerAgentAddr"]))
	engine.Use(gin.Recovery())
	pprof.Register(engine)

	engine.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := engine.Group("/kvm_manager/api/v1")
	{
		v1.GET("/hi", controllers.Hi)

		v1.POST("/host", controllers.CreateHost)
		v1.DELETE("/host/:id", controllers.DeleteHost)
		v1.PUT("/host/:id", controllers.UpdateHost)
		v1.PATCH("/host/:id", controllers.PatchUpdateHost)
		v1.GET("/host", controllers.GetHosts)
		v1.GET("/host/:id", controllers.GetHost)

		v1.POST("/vm", controllers.CreateVm)
		v1.DELETE("/vm/:id", controllers.DeleteVm)
		v1.PUT("/vm/:id", controllers.UpdateVm)
		v1.PATCH("/vm/:id", controllers.PatchUpdateVm)
		v1.GET("/vm", controllers.GetVms)

		v1.POST("/vm_storage", controllers.CreateVmStorage)

		v1.POST("/vm_disk", controllers.CreateVmDisk)

		v1.POST("/os_info", controllers.CreateOsInfo)
	}

	engine.NoRoute(func(c *gin.Context) {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, "Page not found")
	})
	return engine
}

func Run(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: setupRouter(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen error: %s\n", err.Error())
		}
	}()
	fmt.Printf("[GIN-debug] Listening and serving HTTP on %s\n", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	<-quit
	logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown error: %s", err.Error())
	}
	fmt.Printf("\nServer exiting\n")
}
