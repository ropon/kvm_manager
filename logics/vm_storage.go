package logics

import (
	"fmt"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/core"
	"github.com/ropon/kvm_manager/models"
	"github.com/ropon/kvm_manager/utils"
)

type CUVmStorageReq struct {
	Status     int    `json:"status" form:"status"`
	Name       string `json:"name" form:"name" binding:"required"`
	Type       string `json:"type" form:"type" binding:"required"`
	Pool       string `json:"pool" form:"pool" binding:"required"`
	ExtraKeys  string `json:"extra_keys" form:"extra_keys" binding:"required"`
	Hosts      string `json:"hosts" form:"hosts"`
	Port       string `json:"port" form:"port"`
	HostUUID   string `json:"host_uuid" form:"host_uuid"`
	Capacity   uint   `json:"capacity" form:"capacity"`
	Annotation string `json:"annotation" form:"annotation"`
	BaseData
}

func initVmStorage(req *CUVmStorageReq) *models.VmStorage {
	return &models.VmStorage{
		Name:       req.Name,
		Type:       req.Type,
		Pool:       req.Pool,
		Status:     req.Status,
		Hosts:      req.Hosts,
		Port:       req.Port,
		Capacity:   req.Capacity,
		ExtraKeys:  req.ExtraKeys,
		Annotation: req.Annotation,
	}
}

func CreateVmStorage(req *CUVmStorageReq) (*models.VmStorage, error) {
	if req.Type != models.StorageTypeLocal {
		if req.Hosts == "" || req.Port == "" {
			return nil, fmt.Errorf("hosts/port参数不完整")
		}
	} else {
		if req.Capacity == 0 {
			return nil, fmt.Errorf("本地盘容量不能为空")
		}
	}
	s := initVmStorage(req)
	err := s.GetByName()
	if err == nil && s.Id != 0 {
		return nil, fmt.Errorf("名称:%s对应存储已存在", s.Name)
	}
	s.UUID = utils.CreateUUID()
	tx := conf.MysqlDb.Begin()
	err = s.CreateTx(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	storageXmlMgr := core.StorageXmlMgr{
		Name:      s.Name,
		Type:      s.Type,
		ExtraKeys: s.ExtraKeys,
	}
	storageXml, err := storageXmlMgr.PoolCreateXml()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = core.CreatePool(req.HostUUID, storageXml)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	s.FormatTime()
	tx.Commit()
	return s, nil
}
