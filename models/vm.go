package models

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type Vm struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Cpu           uint      `json:"cpu" form:"cpu" gorm:"column:cpu"`                      //CPU核心数
	Mem           uint      `json:"mem" form:"mem" gorm:"column:mem"`                      //内存容量
	Status        int       `json:"status" form:"status" gorm:"column:status"`             //状态 0 关机 1开机 2暂停
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                   //虚拟机UUID
	HostUUID      string    `json:"host_uuid" form:"host_uuid" gorm:"column:host_uuid"`    //宿主机uuid
	Name          string    `json:"name" form:"name" gorm:"column:name"`                   //虚拟机名称
	Ipv4          string    `json:"ipv4" form:"ipv4" gorm:"column:ipv4"`                   //ipv4地址
	VmXml         string    `json:"vm_xml" form:"vm_xml" gorm:"column:vm_xml;type:text"`   //虚拟机xml配置文件
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"` //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time;type:datetime"`
}

type VmList []*Vm

func (s *Vm) TableName() string {
	return "vm"
}

// FormatTime 特殊处理时间
func (s *Vm) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

// Create 增(post /vm)
func (s *Vm) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

func (s *Vm) CreateTx(db *gorm.DB) (err error) {
	if db == nil {
		db = conf.MysqlDb
	}
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

// Delete 删(delete /vm/:h_id)
func (s *Vm) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

// Update 改(put /vm/:h_id)/全部
func (s *Vm) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

// Patch 改(patch /vm/:h_id)/部分
func (s *Vm) Patch(v interface{}) (err error) {
	tmp := v.(*Vm)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

// Get 查(get /vm/:h_id)一个
func (s *Vm) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

// GetByIpv4 根据ip查询一个
func (s *Vm) GetByIpv4() (err error) {
	err = conf.MysqlDb.Where("ipv4 = ?", s.Ipv4).Find(s).Error
	return
}

// List 查(get /vm)多个
func (s *Vm) List(ctx context.Context, PageSize, PageNum int64) (list VmList, count int64, err error) {
	sp, _ := utils.ExtractChildSpan("db:get vms", ctx)
	defer sp.Finish()
	list = make(VmList, 0)
	//默认精确匹配
	db := conf.MysqlDb.Where(s)
	//可以自定义查询
	if s.Ipv4 != "" {
		db = conf.MysqlDb.Where("ipv4 like ?", fmt.Sprintf(`%%%s%%`, s.Ipv4))
	}
	offset, limit := utils.GetOffsetAndLimit(PageSize, PageNum)
	err = db.Model(s).Count(&count).Offset(offset).Limit(limit).Find(&list).Error
	return list, count, err
}
