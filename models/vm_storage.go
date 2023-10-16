package models

import (
	"github.com/jinzhu/gorm"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

const (
	vmStorageTableName = "vm_storage"

	StorageTypeLocal = "dir"
	StorageTypeNfs   = "netfs"
	StorageTypeCeph  = "ceph"
)

type VmStorage struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Capacity      uint      `json:"capacity" form:"capacity" gorm:"column:capacity"`       //容量
	UsedCap       uint      `json:"used_cap" form:"used_cap" gorm:"column:used_cap"`       //已使用量
	Status        int       `json:"status" form:"status" gorm:"column:status"`             //状态 0 不可用 1可用
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                   //存储UUID
	Name          string    `json:"name" form:"name" gorm:"column:name"`                   //存储名称
	Type          string    `json:"type" form:"type" gorm:"column:storage_type"`           //存储类型 local/nfs/ceph
	Hosts         string    `json:"hosts" form:"hosts" gorm:"column:hosts"`                //服务列表
	Port          string    `json:"port" form:"port" gorm:"column:port"`                   //服务端口
	Pool          string    `json:"pool" form:"pool" gorm:"column:pool"`                   //存储池
	ExtraKeys     string    `json:"extra_keys" form:"extra_keys" gorm:"column:extra_keys"` //扩展字段 存储ceph/nfs等密钥
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"` //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time;type:datetime"`
}

type Disk2StorageInfo struct {
	Type       string `json:"type" form:"type" gorm:"column:storage_type"`              //存储类型 dir/nfs/ceph
	Hosts      string `json:"hosts" form:"hosts" gorm:"column:hosts"`                   //服务列表
	Port       string `json:"port" form:"port" gorm:"column:port"`                      //服务端口
	Pool       string `json:"pool" form:"pool" gorm:"column:pool"`                      //存储池
	ExtraKeys  string `json:"extra_keys" form:"extra_keys" gorm:"column:extra_keys"`    //扩展字段 存储ceph/nfs等密钥
	Name       string `json:"name" form:"name" gorm:"column:name"`                      //磁盘名称
	DiskFormat string `json:"disk_format" form:"disk_format" gorm:"column:disk_format"` //磁盘格式raw qcow2
	Capacity   uint   `json:"capacity" form:"capacity" gorm:"column:capacity"`          //容量
	Status     int    `json:"status" form:"status" gorm:"column:status;default 0"`      //状态 0 未启用 1启用
	UUID       string `json:"uuid" form:"uuid" gorm:"column:uuid"`                      //镜像UUID
	OsXml      string `json:"os_xml" form:"os_xml" gorm:"column:os_xml;type:text"`      //镜像xml配置文件

}

type VmStorageList []*VmStorage

func (s *VmStorage) TableName() string {
	return vmStorageTableName
}

// FormatTime 特殊处理时间
func (s *VmStorage) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

func (s *VmStorage) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

func (s *VmStorage) CreateTx(db *gorm.DB) (err error) {
	if db == nil {
		db = conf.MysqlDb
	}
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}
func (s *VmStorage) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

func (s *VmStorage) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *VmStorage) UpdateTx(db *gorm.DB) (err error) {
	if db == nil {
		db = conf.MysqlDb
	}
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *VmStorage) Patch(v interface{}) (err error) {
	tmp := v.(*VmStorage)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

func (s *VmStorage) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

func (s *VmStorage) GetByUUID() (err error) {
	err = conf.MysqlDb.Where("uuid = ?", s.UUID).Find(s).Error
	return
}

func (s *VmStorage) GetByName() (err error) {
	err = conf.MysqlDb.Where("name = ?", s.Name).Find(s).Error
	return
}

func (ds *Disk2StorageInfo) GetByUUID() (err error) {
	err = conf.MysqlDb.Debug().Table("vm_storage a").Select("a.storage_type,a.pool,a.extra_keys,b.name,b.disk_format,b.capacity,c.status,c.uuid,c.os_xml").
		Joins("left join vm_disk b on a.uuid = b.storage_uuid").Joins("left join os_info c on b.uuid = c.vm_disk_uuid").
		Where("c.uuid = ?", ds.UUID).Find(ds).Error
	return
}
