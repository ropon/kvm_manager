package models

import (
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type VmDisk struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                            //磁盘UUID
	Name          string    `json:"name" form:"name" gorm:"column:name" sql:"unique;not null"`      //磁盘名称
	Status        int       `json:"status" form:"status" gorm:"column:status"`                      //状态 0 不可用 1可用
	Capacity      uint      `json:"capacity" form:"capacity" gorm:"column:capacity" sql:"not null"` //容量
	VmId          uint      `json:"vm_id" form:"vm_id" gorm:"column:vm_id" sql:"not null"`          //虚拟机id
	Storage       string    `json:"storage" form:"storage" gorm:"column:storage" sql:"not null"`    //存储信息
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`          //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time" sql:"type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time" sql:"type:datetime"`
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
