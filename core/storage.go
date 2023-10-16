package core

import (
	"libvirt.org/go/libvirtxml"
)

type StorageXmlMgr struct {
	Name      string
	Type      string //存储类型
	Hosts     string //服务列表
	Port      string //服务端口
	Pool      string //存储池
	ExtraKeys string //扩展字段 存储ceph/nfs等密钥
}

func (s *StorageXmlMgr) PoolCreateXml() (string, error) {
	poolXml := new(libvirtxml.StoragePool)
	poolXml.Name = s.Name
	poolXml.Type = s.Type
	poolXml.Target = s.PoolCreateTarget()
	xmlStr, err := poolXml.Marshal()
	if err != nil {
		return "", err
	}
	return xmlStr, nil
}

func (s *StorageXmlMgr) PoolCreateTarget() *libvirtxml.StoragePoolTarget {
	return &libvirtxml.StoragePoolTarget{
		Path: s.ExtraKeys,
	}
}

func (s *StorageXmlMgr) Clone() {

}

func (s *StorageXmlMgr) Delete() {

}

func (s *StorageXmlMgr) Update() {

}

func CreatePool(hostUUID, poolXml string) error {
	libVirtMgr, err := GetLibVirtMgr(hostUUID)
	if err != nil {
		return err
	}
	_, err = libVirtMgr.Conn.StoragePoolDefineXML(poolXml, 0)
	if err != nil {
		return err
	}
	return nil
}
