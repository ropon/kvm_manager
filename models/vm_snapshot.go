package models

import (
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type VmSnapShot struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Capacity      uint      `json:"capacity" form:"capacity" gorm:"column:capacity" sql:"not null"` //容量
	Status        int       `json:"status" form:"status" gorm:"column:status"`                      //状态 0 不可用 1可用
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                            //快照UUID
	VmUUID        string    `json:"vm_uuid" form:"vm_uuid" gorm:"column:vm_uuid"`                   //虚拟机uuid
	Name          string    `json:"name" form:"name" gorm:"column:name" sql:"unique;not null"`      //快照名称
	StorageUUID   string    `json:"storage_uuid" form:"storage_uuid" gorm:"column:storage_uuid"`    //存储信息
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`          //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time;type:datetime"`
}

type VmSnapShotList []*VmSnapShot

func (s *VmSnapShot) TableName() string {
	return "vm_snapshot"
}

//FormatTime 特殊处理时间
func (s *VmSnapShot) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}
