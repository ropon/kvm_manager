package core

import (
	"fmt"
	"github.com/ropon/kvm_manager/models"
	"libvirt.org/go/libvirt"
)

var libVirtMgrMap = map[string]*LibVirtMgr{}

type LibVirtMgr struct {
	Conn *libvirt.Connect
}

func GetLibVirtMgr(hostUUID string) (*LibVirtMgr, error) {
	libVirtMgr, ok := libVirtMgrMap[hostUUID]
	if !ok {
		var err error
		libVirtMgr, err = newLibVirtMgr(hostUUID)
		if err != nil {
			return nil, err
		}
		libVirtMgrMap[hostUUID] = libVirtMgr
	}
	return libVirtMgr, nil
}

func newLibVirtMgr(hostUUID string) (*LibVirtMgr, error) {
	host := models.Host{UUID: hostUUID}
	err := host.GetByUUID()
	if err != nil {
		return nil, err
	}
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+tcp://%s/system", host.Ipv4))
	if err != nil {
		return nil, err
	}
	return &LibVirtMgr{Conn: conn}, nil
}

func (s *LibVirtMgr) defineVm(xml string) (err error) {
	_, err = s.Conn.DomainDefineXML(xml)
	return err
}

func (s *LibVirtMgr) GetDomain(uuid string) (*libvirt.Domain, error) {
	vm, err := s.Conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, err
	}
	return vm, nil
}

func (s *LibVirtMgr) unDefineVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	err = vm.Undefine()
	return err
}

func (s *LibVirtMgr) StartVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	defer vm.Free()
	return vm.Create()
}

func (s *LibVirtMgr) ShutdownVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	defer vm.Free()
	return vm.Shutdown()
}

func (s *LibVirtMgr) DestroyVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	defer vm.Free()
	return vm.Destroy()
}

func (s *LibVirtMgr) SuspendVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	defer vm.Free()
	return vm.Suspend()
}

func (s *LibVirtMgr) ResumeVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	defer vm.Free()
	return vm.Resume()
}

func (s *LibVirtMgr) RebootVm(uuid string) error {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return err
	}
	defer vm.Free()
	return vm.Reboot(0)
}

func (s *LibVirtMgr) VmStatus(uuid string) (int, error) {
	vm, err := s.GetDomain(uuid)
	if err != nil {
		return 0, err
	}
	defer vm.Free()
	a, b, err := vm.GetState()
	fmt.Println("state:", a, b, err)
	return 0, nil
}

func (s *LibVirtMgr) GetDomains() ([]libvirt.Domain, error) {
	vms, err := s.Conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_PERSISTENT)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

func DefineVm(hostUUID, vmXml string) error {
	libVirtMgr, err := GetLibVirtMgr(hostUUID)
	if err != nil {
		return err
	}
	return libVirtMgr.defineVm(vmXml)
}

func UnDefineVm(hostUUID string, uuid string) error {
	libVirtMgr, err := GetLibVirtMgr(hostUUID)
	if err != nil {
		return err
	}
	return libVirtMgr.unDefineVm(uuid)
}
