package main

import (
	"fmt"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/logics"
	"github.com/ropon/kvm_manager/routers"
)

// @title kvm_manager
// @version 1.0
// @description 后端快速Api脚手架

// @contact.name Ropon
// @contact.url https://www.ropon.top
// @contact.email ropon@xxx.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host work-api.xxx.com:8989
// @BasePath /
func main() {
	err := conf.Init()
	if err != nil {
		fmt.Printf("init failed, err: %v\n", err)
		return
	}

	logics.Migrate()
	routers.Run(conf.Cfg.Listen)
}
