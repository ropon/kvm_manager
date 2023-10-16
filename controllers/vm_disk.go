package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/logics"
	"github.com/ropon/kvm_manager/utils"
)

// CreateVmDisk 创建虚拟磁盘接口
// @Summary 创建虚拟磁盘接口
// @Description 创建虚拟磁盘接口
// @Tags 虚拟磁盘相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUVmDiskReq true "请求参数"
// @Success 200 {object} models.VmDisk "创建成功返回结果"
// @Router /kvm_manager/api/v1/vm_disk [post]
func CreateVmDisk(c *gin.Context) {
	req := new(logics.CUVmDiskReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.CreateVmDisk(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "创建成功")
}
