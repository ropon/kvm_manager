package logics

import (
	"fmt"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/models"
	"github.com/ropon/kvm_manager/utils"
)

type CUVmDiskReq struct {
	Status      int    `json:"status" form:"status"`
	Capacity    uint   `json:"capacity" form:"capacity" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	DiskFormat  string `json:"disk_format" form:"disk_format" binding:"required"`
	StorageUUID string `json:"storage_uuid" form:"storage_uuid" binding:"required"`
	VmUUID      string `json:"vm_uuid" form:"vm_uuid"`
	Annotation  string `json:"annotation" form:"annotation"`
	BaseData
}

func initVmDisk(req *CUVmDiskReq) *models.VmDisk {
	return &models.VmDisk{
		Status:      req.Status,
		Name:        req.Name,
		DiskFormat:  req.DiskFormat,
		StorageUUID: req.StorageUUID,
		Capacity:    req.Capacity,
		VmUUID:      req.VmUUID,
		Annotation:  req.Annotation,
	}
}

func CreateVmDisk(req *CUVmDiskReq) (*models.VmDisk, error) {
	s := initVmDisk(req)
	err := s.GetByName()
	if err == nil && s.Id != 0 {
		return nil, fmt.Errorf("名称:%s对应虚拟磁盘已存在", s.Name)
	}
	vmStorage := &models.VmStorage{UUID: req.StorageUUID}
	err = vmStorage.GetByUUID()
	if err != nil {
		return nil, err
	}
	if vmStorage.Status != 1 {
		return nil, fmt.Errorf("该存储未启用")
	}
	if req.Capacity+vmStorage.UsedCap > vmStorage.Capacity {
		return nil, fmt.Errorf("存储容量不足")
	}
	s.UUID = utils.CreateUUID()
	tx := conf.MysqlDb.Begin()
	err = s.CreateTx(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	vmStorage.UsedCap += req.Capacity
	err = vmStorage.UpdateTx(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	s.FormatTime()
	tx.Commit()
	return s, nil
}
