package models

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type Host struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Cpu           uint      `json:"cpu" form:"cpu" gorm:"column:cpu"`                         //总核心数
	Mem           uint      `json:"mem" form:"mem" gorm:"column:mem"`                         //总内存容量，单位GB
	UsedCpu       uint      `json:"used_cpu" form:"used_cpu"`                                 //已使用cpu核心数
	UsedMem       uint      `json:"used_mem" form:"used_mem"`                                 //已使用内存
	MaxVms        uint      `json:"max_vms" form:"max_vms" gorm:"column:max_vms"`             //最大虚拟机数量
	CreatedVms    uint      `json:"created_vms" form:"created_vms" gorm:"column:created_vms"` //已创建虚拟机数量
	Status        int       `json:"status" form:"status" gorm:"column:status"`                //状态
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                      //宿主机UUID
	Ipv4          string    `json:"ipv4" form:"ipv4" gorm:"column:ipv4"`                      //ipv4地址
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`    //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time;type:datetime"`
}

type HostList []*Host

func (s *Host) TableName() string {
	return "host"
}

// FormatTime 特殊处理时间
func (s *Host) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

// Create 增(post /host)
func (s *Host) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

// Delete 删(delete /host/:h_id)
func (s *Host) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

// Update 改(put /host/:h_id)/全部
func (s *Host) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *Host) UpdateTx(db *gorm.DB) (err error) {
	if db == nil {
		db = conf.MysqlDb
	}
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

// Patch 改(patch /host/:h_id)/部分
func (s *Host) Patch(v interface{}) (err error) {
	tmp := v.(*Host)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

// Get 查(get /host/:h_id)一个
func (s *Host) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

func (s *Host) GetByUUID() (err error) {
	err = conf.MysqlDb.Where("uuid = ?", s.UUID).Find(s).Error
	return
}

// GetByIpv4 根据ipv4查询一个
func (s *Host) GetByIpv4() (err error) {
	err = conf.MysqlDb.Where("ipv4 = ?", s.Ipv4).Find(s).Error
	return
}

// List 查(get /host)多个
func (s *Host) List(ctx context.Context, PageSize, PageNum int64) (list HostList, count int64, err error) {
	sp, _ := utils.ExtractChildSpan("db:get hosts", ctx)
	defer sp.Finish()
	list = make(HostList, 0)
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
