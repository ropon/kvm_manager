package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/logics"
	"github.com/ropon/kvm_manager/utils"
)

// CreateOsInfo 创建镜像接口
// @Summary 创建镜像接口
// @Description 创建镜像接口
// @Tags 镜像相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUOsInfoReq true "请求参数"
// @Success 200 {object} models.OsInfo "创建成功返回结果"
// @Router /kvm_manager/api/v1/os_info [post]
func CreateOsInfo(c *gin.Context) {
	req := new(logics.CUOsInfoReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.CreateOsInfo(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "创建成功")
}
