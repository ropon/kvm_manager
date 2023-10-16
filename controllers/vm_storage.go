package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/logics"
	"github.com/ropon/kvm_manager/utils"
)

// CreateVmStorage 创建存储接口
// @Summary 创建存储接口
// @Description 创建存储接口
// @Tags 存储相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUVmStorageReq true "请求参数"
// @Success 200 {object} models.VmStorage "创建成功返回结果"
// @Router /kvm_manager/api/v1/vm_storage [post]
func CreateVmStorage(c *gin.Context) {
	req := new(logics.CUVmStorageReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.CreateVmStorage(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "创建成功")
}
