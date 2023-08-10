package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/kvm_manager/logics"
	"github.com/ropon/kvm_manager/utils"
)

// CreateHost 创建宿主机接口
// @Summary 创建宿主机接口
// @Description 创建宿主机接口
// @Tags 宿主机相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUHostReq true "请求参数"
// @Success 200 {object} models.Host "创建成功返回结果"
// @Router /kvm_manager/api/v1/host [post]
func CreateHost(c *gin.Context) {
	req := new(logics.CUHostReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.CreateHost(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "创建成功")
}

// DeleteHost 删除宿主机接口
// @Summary 删除宿主机接口
// @Description 删除宿主机接口
// @Tags 宿主机相关接口
// @Produce application/json
// @Param id path uint true "id"
// @Success 200
// @Router /kvm_manager/api/v1/host/{id} [delete]
func DeleteHost(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	err := logics.DeleteHost(id)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, "", "删除成功")
}

// UpdateHost 更新宿主机全部参数接口
// @Summary 更新宿主机全部参数接口
// @Description 更新宿主机全部参数接口
// @Tags 宿主机相关接口
// @Accept application/json
// @Produce application/json
// @Param id path uint true "id"
// @Param data body logics.CUHostReq true "请求参数"
// @Success 200 {object} models.Host "更新成功返回结果"
// @Router /kvm_manager/api/v1/host/{id} [put]
func UpdateHost(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	req := new(logics.CUHostReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.UpdateHost(id, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "更新成功")
}

// PatchUpdateHost 更新宿主机部分参数接口
// @Summary 更新宿主机部分参数接口
// @Description 更新宿主机部分参数接口
// @Tags 宿主机相关接口
// @Accept application/json
// @Produce application/json
// @Param id path uint true "id"
// @Param data body logics.HostReq true "请求参数"
// @Success 200 {object} models.Host "更新成功返回结果"
// @Router /kvm_manager/api/v1/host/{id} [patch]
func PatchUpdateHost(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	req := new(logics.HostReq)
	if !checkData(c, req) {
		return
	}

	res, err := logics.PatchUpdateHost(id, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "更新成功")
}

// GetHosts 获取宿主机列表接口
// @Summary 获取宿主机列表接口
// @Description 获取宿主机列表接口
// @Tags 宿主机相关接口
// @Produce application/json
// @Param data query logics.HostReq true "请求参数"
// @Success 200 {object} logics.HostRes "服务列表返回结果"
// @Router /kvm_manager/api/v1/host [get]
func GetHosts(c *gin.Context) {
	req := new(logics.HostReq)
	if !checkData(c, req) {
		return
	}

	ctx := utils.ExtractStdContext(nil, c.Request.Header)
	resList, err := logics.GetHosts(ctx, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, resList, "获取列表成功")
}

// GetHost 获取单个宿主机接口
// @Summary 获取单个宿主机接口
// @Description 获取单个宿主机接口
// @Tags 宿主机相关接口
// @Produce application/json
// @Param id path uint true "id"
// @Success 200 {object} models.Host "服务返回结果"
// @Router /kvm_manager/api/v1/host/{id} [get]
func GetHost(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	res, err := logics.GetHost(id)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "获取成功")
}
