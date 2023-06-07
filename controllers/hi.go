package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

func Hi(c *gin.Context) {
	time.Sleep(time.Second * 10)
	utils.GinOKRsp(c, "hi ropon", "ok")
}
