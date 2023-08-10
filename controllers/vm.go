package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/logics"
	"github.com/ropon/kvm_manager/utils"
)

// CreateVm 创建虚拟机接口
// @Summary 创建虚拟机接口
// @Description 创建虚拟机接口
// @Tags 虚拟机相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUVmReq true "请求参数"
// @Success 200 {object} models.Vm "创建成功返回结果"
// @Router /kvm_manager/api/v1/vm [post]
func CreateVm(c *gin.Context) {
	req := new(logics.CUVmReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.CreateVm(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "创建成功")
}

// DeleteVm 删除虚拟机接口
// @Summary 删除虚拟机接口
// @Description 删除虚拟机接口
// @Tags 虚拟机相关接口
// @Produce application/json
// @Param id path uint true "id"
// @Success 200
// @Router /kvm_manager/api/v1/vm/{id} [delete]
func DeleteVm(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	err := logics.DeleteVm(id)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, "", "删除成功")
}

// UpdateVm 更新虚拟机全部参数接口
// @Summary 更新虚拟机全部参数接口
// @Description 更新虚拟机全部参数接口
// @Tags 虚拟机相关接口
// @Accept application/json
// @Produce application/json
// @Param id path uint true "id"
// @Param data body logics.CUVmReq true "请求参数"
// @Success 200 {object} models.Vm "更新成功返回结果"
// @Router /kvm_manager/api/v1/vm/{id} [put]
func UpdateVm(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	req := new(logics.CUVmReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.UpdateVm(id, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "更新成功")
}

// PatchUpdateVm 更新虚拟机部分参数接口
// @Summary 更新虚拟机部分参数接口
// @Description 更新虚拟机部分参数接口
// @Tags 虚拟机相关接口
// @Accept application/json
// @Produce application/json
// @Param id path uint true "id"
// @Param data body logics.VmReq true "请求参数"
// @Success 200 {object} models.Vm "更新成功返回结果"
// @Router /kvm_manager/api/v1/vm/{id} [patch]
func PatchUpdateVm(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	req := new(logics.VmReq)
	if !checkData(c, req) {
		return
	}

	res, err := logics.PatchUpdateVm(id, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "更新成功")
}
