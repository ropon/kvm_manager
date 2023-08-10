package logics

import (
	"fmt"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/core"
	"github.com/ropon/kvm_manager/models"
	"github.com/ropon/kvm_manager/utils"
	"github.com/ropon/logger"
)

type BaseData struct {
	UserEmail string `json:"user_email" form:"user_email"`
	OpsAdmin  string `json:"ops_admin" form:"ops_admin"`
}

type CUVmReq struct {
	Name   string `json:"name" form:"name" binding:"required"`
	Ipv4   string `json:"ipv4" form:"ipv4" binding:"required"`
	OsName string `json:"os_name" form:"os_name" binding:"required"`
	Cpu    uint   `json:"cpu" form:"cpu" binding:"required"`
	Mem    uint   `json:"mem" form:"mem" binding:"required"`
	HostId uint   `json:"host_id" form:"host_id" binding:"required"`
	BaseData
}

type VmReq struct {
	models.Vm
	OpType   string `json:"op_type" form:"op_type"`
	PageSize int64  `json:"page_size" form:"page_size"`
	PageNum  int64  `json:"page_num" form:"page_num"`
}

type VmRes struct {
	TotalCount int64         `json:"total_count"`
	VmList     models.VmList `json:"vm_list"`
}

func (bp *BaseData) Init(userEmail, opsAdmin string) {
	if bp.UserEmail == "" {
		bp.UserEmail = userEmail
	}
	if bp.OpsAdmin == "" {
		bp.OpsAdmin = opsAdmin
	}
}

func initVm(req *CUVmReq) *models.Vm {
	return &models.Vm{
		Name:   req.Name,
		Ipv4:   req.Ipv4,
		Cpu:    req.Cpu,
		Mem:    req.Mem,
		HostId: req.HostId,
	}
}

//创建虚拟机返回详情
func CreateVm(req *CUVmReq) (*models.Vm, error) {
	s := initVm(req)
	err := s.GetByIpv4()
	if err == nil && s.Id != 0 {
		return nil, fmt.Errorf("ipv4地址:%s对应虚拟机已存在", s.Ipv4)
	}

	//1.检查宿主机资源是否够
	host := models.Host{Id: s.HostId}
	err = host.Get()
	if err != nil {
		return nil, err
	}
	if req.Cpu+host.UsedCpu > host.Cpu {
		return nil, fmt.Errorf("cpu资源不足")
	}
	if req.Mem+host.UsedMem > host.Mem {
		return nil, fmt.Errorf("内存资源不足")
	}
	if 1+host.CreatedVms > host.MaxVms {
		return nil, fmt.Errorf("虚拟机数量超限")
	}

	//2.获取网络配置
	ipInfo := models.IpInfo{Ipv4: req.Ipv4}
	err = ipInfo.GetByIpv4()
	if err != nil {
		logger.Error("ipInfo.GetByIpv4 failed: %v", err)
		return nil, err
	}
	if ipInfo.Status != 1 {
		return nil, fmt.Errorf("IP:%s不可用", req.Ipv4)
	}

	//3.获取系统详情
	osInfo := models.OsInfo{Name: req.OsName}
	err = osInfo.GetByName()
	if err != nil {
		logger.Error("osInfo.GetByName failed: %v", err)
		return nil, err
	}
	if ipInfo.Status != 1 {
		return nil, fmt.Errorf("系统:%s不可用", req.OsName)
	}

	//4.准备创建虚拟机信息
	/*
		uuid/name/cpu/mem/network/disk
	*/
	s.UUID = utils.CreateUUID()
	xmlMgr := core.LibVirtXmlMgr{
		UUIDStr:    s.UUID,
		VmName:     req.Name,
		OsXml:      osInfo.OsXml,
		MacAddr:    ipInfo.MacAddr,
		BridgeName: ipInfo.BridgeName,
		Cpu:        req.Cpu,
		Mem:        req.Mem * 1024 * 1024,
	}
	defineXml, err := xmlMgr.CreateXml()
	if err != nil {
		logger.Error("xmlMgr.CreateXml failed: %v", err)
		return nil, err
	}
	s.VmXml = defineXml

	err = core.DefineVm(req.HostId, defineXml)
	if err != nil {
		return nil, err
	}

	tx := conf.MysqlDb.Begin()
	err = s.CreateTx(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	s.FormatTime()
	//5.更换host 更新ipinfo
	host.UsedCpu = host.UsedCpu + req.Cpu
	host.UsedMem = host.UsedMem + req.Mem
	host.CreatedVms = host.CreatedVms + 1
	err = host.UpdateTx(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ipInfo.Status = 0
	err = ipInfo.UpdateTx(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return s, nil
}

//通过服务ID删除指定虚拟机
func DeleteVm(id uint) error {
	//通过id获取vm
	vm := &models.Vm{Id: id}
	err := vm.Get()
	if err != nil {
		return err
	}
	if vm.Status != 0 {
		return fmt.Errorf("请先关机")
	}
	//host := &models.Host{Id: vm.HostId}
	//err = host.Get()
	//if err != nil {
	//	return err
	//}
	err = core.UnDefineVm(vm.HostId, vm.UUID)
	if err != nil {
		return err
	}
	return vm.Delete()
}

//通过服务ID更新指定虚拟机全部信息
func UpdateVm(id uint, req *CUVmReq) (*models.Vm, error) {
	do := &DbObj{
		Id:  id,
		Obj: &models.Vm{Id: id},
	}
	if err := do.get(); err != nil {
		return nil, err
	}

	s := initVm(req)
	s.Id = id
	s.CreateTime = do.Obj.(*models.Vm).CreateTime

	return s, do.update(s)
}

//通过服务ID更新指定虚拟机部分信息
func PatchUpdateVm(id uint, req *VmReq) (interface{}, error) {
	vm := &models.Vm{Id: id}
	err := vm.Get()
	if err != nil {
		return nil, err
	}

	//host := &models.Host{Id: vm.HostId}
	//err = host.Get()
	//if err != nil {
	//	return nil, err
	//}
	libvirtMgr, err := core.GetLibVirtMgr(vm.HostId)
	if err != nil {
		return nil, err
	}
	if req.OpType != "" {
		switch req.OpType {
		case "start":
			_ = libvirtMgr.StartVm(vm.UUID)
		case "shutdown":
			_ = libvirtMgr.ShutdownVm(vm.UUID)
		case "destroy":
			_ = libvirtMgr.DestroyVm(vm.UUID)
		case "suspend":
			_ = libvirtMgr.SuspendVm(vm.UUID)
		case "resume":
			_ = libvirtMgr.ResumeVm(vm.UUID)
		case "reboot":
			_ = libvirtMgr.RebootVm(vm.UUID)
		default:
			return nil, fmt.Errorf("不支持此操作")
		}
	}
	return nil, nil
	//s := req.Vm
	//do := &DbObj{
	//	Id:  id,
	//	Obj: &models.Vm{Id: id},
	//}
	//if err := do.patch(&s); err != nil {
	//	return nil, err
	//}
	//do.Obj.(*models.Vm).UpdateTimeStr = s.UpdateTimeStr
	//return do.Obj, nil
}
