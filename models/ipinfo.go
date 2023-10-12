package models

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type IpInfo struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                       //虚拟机UUID
	Status        int       `json:"status" form:"status" gorm:"column:status"`                 //状态 0未使用 1已使用
	Ipv4          string    `json:"ipv4" form:"ipv4" gorm:"column:ipv4" sql:"unique;not null"` //ipv4地址
	MacAddr       string    `json:"mac_addr" form:"mac_addr" gorm:"column:mac_addr"`           //MAC地址
	BridgeName    string    `json:"bridge_name" form:"bridge_name"`                            //网桥名称
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`     //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time" sql:"type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time" sql:"type:datetime"`
}

type IpInfoList []*IpInfo

func (s *IpInfo) TableName() string {
	return "ip_info"
}

// FormatTime 特殊处理时间
func (s *IpInfo) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

func (s *IpInfo) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

func (s *IpInfo) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

func (s *IpInfo) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *IpInfo) UpdateTx(db *gorm.DB) (err error) {
	if db == nil {
		db = conf.MysqlDb
	}
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *IpInfo) Patch(v interface{}) (err error) {
	tmp := v.(*Host)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

func (s *IpInfo) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

func (s *IpInfo) GetByIpv4() (err error) {
	err = conf.MysqlDb.Where("ipv4 = ?", s.Ipv4).Find(s).Error
	return
}

func (s *IpInfo) List(ctx context.Context, PageSize, PageNum int64) (list VmList, count int64, err error) {
	sp, _ := utils.ExtractChildSpan("db:get hosts", ctx)
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
