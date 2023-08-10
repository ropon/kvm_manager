package logics

import (
	"context"
	"fmt"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/models"
)

type CUHostReq struct {
	Ipv4   string `json:"ipv4" form:"ipv4" binding:"required"`
	Cpu    uint   `json:"cpu" form:"cpu" binding:"required"`
	Mem    uint   `json:"mem" form:"mem" binding:"required"`
	MaxVms uint   `json:"max_vms" form:"max_vms" binding:"required"`
	BaseData
}

type HostReq struct {
	models.Host
	PageSize int64 `json:"page_size" form:"page_size"`
	PageNum  int64 `json:"page_num" form:"page_num"`
}

type HostRes struct {
	TotalCount int64           `json:"total_count"`
	HostList   models.HostList `json:"host_list"`
}

func initHost(req *CUHostReq) *models.Host {
	return &models.Host{
		Ipv4:   req.Ipv4,
		Cpu:    req.Cpu,
		Mem:    req.Mem,
		MaxVms: req.MaxVms,
	}
}

//创建宿主机返回详情
func CreateHost(req *CUHostReq) (*models.Host, error) {
	s := initHost(req)
	err := s.GetByIpv4()
	if err == nil && s.Id != 0 {
		return nil, fmt.Errorf("ipv4地址:%s对应宿主机已存在", s.Ipv4)
	}

	err = s.Create()
	if err != nil {
		return nil, err
	}

	s.FormatTime()
	return s, nil
}

//通过服务ID删除指定宿主机
func DeleteHost(id uint) error {
	do := &DbObj{
		Id:  id,
		Obj: &models.Host{Id: id},
	}
	return do.delete()
}

//通过服务ID更新指定宿主机全部信息
func UpdateHost(id uint, req *CUHostReq) (*models.Host, error) {
	do := &DbObj{
		Id:  id,
		Obj: &models.Host{Id: id},
	}
	if err := do.get(); err != nil {
		return nil, err
	}

	s := initHost(req)
	s.Id = id
	s.CreateTime = do.Obj.(*models.Host).CreateTime

	return s, do.update(s)
}

//通过服务ID更新指定宿主机部分信息
func PatchUpdateHost(id uint, req *HostReq) (interface{}, error) {
	s := req.Host
	do := &DbObj{
		Id:  id,
		Obj: &models.Host{Id: id},
	}
	if err := do.patch(&s); err != nil {
		return nil, err
	}
	do.Obj.(*models.Host).UpdateTimeStr = s.UpdateTimeStr
	return do.Obj, nil
}

//获取宿主机列表
func GetHosts(ctx context.Context, req *HostReq) (*HostRes, error) {
	s := req.Host
	hl, count, err := s.List(ctx, req.PageSize, req.PageNum)
	if err != nil {
		return nil, err
	}
	for _, tempS := range hl {
		tempS.FormatTime()
	}
	res := &HostRes{
		TotalCount: count,
		HostList:   hl,
	}
	return res, nil
}

//获取单个宿主机详情
func GetHost(id uint) (interface{}, error) {
	do := &DbObj{
		Id:  id,
		Obj: &models.Host{Id: id},
	}
	return do.Obj, do.get()
}

func Migrate() {
	conf.MysqlDb.AutoMigrate(&models.Host{}, &models.IpInfo{}, &models.OsInfo{}, &models.Vm{})
}
