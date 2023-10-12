package models

import (
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type VmStorage struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                               //存储UUID
	Name          string    `json:"name" form:"name" gorm:"column:name" sql:"unique;not null"`         //存储名称
	Type          string    `json:"type" form:"type" gorm:"column:storage_type" sql:"unique;not null"` //存储类型 local/nfs/ceph
	Status        int       `json:"status" form:"status" gorm:"column:status"`                         //状态 0 不可用 1可用
	Capacity      uint      `json:"capacity" form:"capacity" gorm:"column:capacity" sql:"not null"`    //容量
	UsedCap       uint      `json:"used_cap" form:"used_cap" gorm:"column:used_cap" sql:"not null"`    //已使用量
	Hosts         string    `json:"hosts" form:"hosts" gorm:"column:hosts" sql:"not null"`             //服务列表
	Port          string    `json:"port" form:"port" gorm:"column:port" sql:"not null"`                //服务端口
	ExtraKeys     string    `json:"extra_keys" form:"extra_keys" gorm:"column:extra_keys"`             //扩展字段 存储ceph/nfs等密钥
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`             //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time" sql:"type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time" sql:"type:datetime"`
}

type VmStorageList []*VmStorage

func (s *VmStorage) TableName() string {
	return "vm_storage"
}

// FormatTime 特殊处理时间
func (s *VmStorage) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}
