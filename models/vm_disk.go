package models

import (
	"github.com/jinzhu/gorm"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type VmDisk struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Capacity      uint      `json:"capacity" form:"capacity" gorm:"column:capacity"`             //容量
	Status        int       `json:"status" form:"status" gorm:"column:status"`                   //状态 0 不可用 1可用
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                         //磁盘UUID
	VmUUID        string    `json:"vm_uuid" form:"vm_uuid" gorm:"column:vm_uuid"`                //虚拟机uuid
	Name          string    `json:"name" form:"name" gorm:"column:name"`                         //磁盘名称
	DiskFormat    string    `json:"disk_format" form:"disk_format" gorm:"column:disk_format"`    //磁盘格式raw qcow2
	StorageUUID   string    `json:"storage_uuid" form:"storage_uuid" gorm:"column:storage_uuid"` //存储信息
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`       //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time;type:datetime"`
}

type VmDiskList []*VmDisk

func (s *VmDisk) TableName() string {
	return "vm_disk"
}

// FormatTime 特殊处理时间
func (s *VmDisk) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

func (s *VmDisk) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

func (s *VmDisk) CreateTx(db *gorm.DB) (err error) {
	if db == nil {
		db = conf.MysqlDb
	}
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

func (s *VmDisk) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

func (s *VmDisk) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *VmDisk) Patch(v interface{}) (err error) {
	tmp := v.(*VmDisk)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

func (s *VmDisk) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

func (s *VmDisk) GetByUUID() (err error) {
	err = conf.MysqlDb.Where("uuid = ?", s.UUID).Find(s).Error
	return
}

func (s *VmDisk) GetByName() (err error) {
	err = conf.MysqlDb.Where("name = ?", s.Name).Find(s).Error
	return
}
