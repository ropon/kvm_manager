package logics

import (
	"fmt"
	"github.com/ropon/kvm_manager/models"
	"github.com/ropon/kvm_manager/utils"
)

type CUOsInfoReq struct {
	Status     int    `json:"status" form:"status"`
	Name       string `json:"name" form:"name" binding:"required"`
	OsType     string `json:"os_type" form:"os_type" binding:"required"`
	VmDiskUUID string `json:"vm_disk_uuid" form:"vm_disk_uuid" binding:"required"`
	OsXml      string `json:"os_xml" form:"os_xml" binding:"required"`
	Annotation string `json:"annotation" form:"annotation"`
	BaseData
}

func initOsInfo(req *CUOsInfoReq) *models.OsInfo {
	return &models.OsInfo{
		Status:     req.Status,
		Name:       req.Name,
		OsType:     req.OsType,
		VmDiskUUID: req.VmDiskUUID,
		OsXml:      req.OsXml,
		Annotation: req.Annotation,
	}
}

func CreateOsInfo(req *CUOsInfoReq) (*models.OsInfo, error) {
	s := initOsInfo(req)
	err := s.GetByName()
	if err == nil && s.Id != 0 {
		return nil, fmt.Errorf("名称:%s对应镜像已存在", s.Name)
	}
	vmDisk := &models.VmDisk{UUID: req.VmDiskUUID}
	err = vmDisk.GetByUUID()
	if err != nil {
		return nil, err
	}
	if vmDisk.Status != 1 {
		return nil, fmt.Errorf("该磁盘不可用")
	}
	s.UUID = utils.CreateUUID()
	err = s.Create()
	if err != nil {
		return nil, err
	}

	s.FormatTime()
	return s, nil
}
