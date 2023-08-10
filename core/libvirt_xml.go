package core

import (
	"fmt"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

const (
	Device   string = "disk"
	DiskName string = "qemu"
	DiskType string = "qcow2"
	DiskDev  string = "vda"
	DiskBus  string = "virtio"
)

type LibVirtXmlMgr struct {
	UUIDStr    string
	VmName     string
	OsXml      string
	MacAddr    string
	BridgeName string
	Cpu        uint
	Mem        uint
}

func (s *LibVirtXmlMgr) CreateXml() (string, error) {
	vmXml := new(libvirtxml.Domain)
	err := vmXml.Unmarshal(s.OsXml)
	if err != nil {
		return "", err
	}
	vmXml.UUID = s.UUIDStr
	vmXml.Name = s.VmName
	vmXml.Devices.Disks = append(vmXml.Devices.Disks, s.createDisk())
	vmXml.VCPU = s.createVCpu()
	vmXml.Memory = s.createMem()
	vmXml.CurrentMemory = s.createCurrMem()
	vmXml.Devices.Interfaces = append(vmXml.Devices.Interfaces, s.createNet())
	xmlStr, err := vmXml.Marshal()
	if err != nil {
		return "", err
	}
	return xmlStr, nil
}

func (s *LibVirtXmlMgr) createDisk() libvirtxml.DomainDisk {
	diskXml := libvirtxml.DomainDisk{}
	diskXml.Device = Device
	diskXml.Driver = &libvirtxml.DomainDiskDriver{
		Name: DiskName,
		Type: DiskType,
	}
	diskXml.Source = &libvirtxml.DomainDiskSource{
		File: &libvirtxml.DomainDiskSourceFile{
			File: fmt.Sprintf("/codoon/%s/%s_os.qcow2", s.VmName, s.VmName),
		},
	}
	diskXml.Target = &libvirtxml.DomainDiskTarget{
		Dev: DiskDev,
		Bus: DiskBus,
	}
	diskXml.Boot = s.createDeviceBoot(1)
	diskXml.Address = s.createDomainAddress(10)
	return diskXml
}

func (s *LibVirtXmlMgr) createDomainAddress(address uint) *libvirtxml.DomainAddress {
	return &libvirtxml.DomainAddress{
		PCI: s.createPCI(address),
	}
}

func (s *LibVirtXmlMgr) createDeviceBoot(order uint) *libvirtxml.DomainDeviceBoot {
	return &libvirtxml.DomainDeviceBoot{
		Order: order,
	}
}

func (s *LibVirtXmlMgr) createPCI(slot uint) *libvirtxml.DomainAddressPCI {
	var (
		Domain   uint = 00
		Bus      uint = 00
		Function uint = 0
	)
	return &libvirtxml.DomainAddressPCI{
		Domain:   &Domain,
		Bus:      &Bus,
		Slot:     &slot,
		Function: &Function,
	}
}

func (s *LibVirtXmlMgr) createVCpu() *libvirtxml.DomainVCPU {
	return &libvirtxml.DomainVCPU{
		Placement: "static",
		Value:     s.Cpu,
	}
}

func (s *LibVirtXmlMgr) createMem() *libvirtxml.DomainMemory {
	return &libvirtxml.DomainMemory{
		Unit:  "KiB",
		Value: s.Mem,
	}
}

func (s *LibVirtXmlMgr) createCurrMem() *libvirtxml.DomainCurrentMemory {
	return &libvirtxml.DomainCurrentMemory{
		Unit:  "KiB",
		Value: s.Mem,
	}
}

func (s *LibVirtXmlMgr) createNet() libvirtxml.DomainInterface {
	netXml := libvirtxml.DomainInterface{}
	netXml.Model = &libvirtxml.DomainInterfaceModel{
		Type: DiskBus,
	}
	netXml.Boot = s.createDeviceBoot(2)
	netXml.MAC = s.createNetMac()
	netXml.Source = s.createNetSource()
	netXml.Address = s.createDomainAddress(16)
	return netXml
}

func (s *LibVirtXmlMgr) createNetMac() *libvirtxml.DomainInterfaceMAC {
	return &libvirtxml.DomainInterfaceMAC{
		Address: s.MacAddr,
	}
}

func (s *LibVirtXmlMgr) createNetSource() *libvirtxml.DomainInterfaceSource {
	return &libvirtxml.DomainInterfaceSource{
		Bridge: &libvirtxml.DomainInterfaceSourceBridge{
			Bridge: s.BridgeName,
		},
	}
}
