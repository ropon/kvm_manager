package models

import (
	"context"
	"fmt"
	"github.com/ropon/kvm_manager/conf"
	"github.com/ropon/kvm_manager/utils"
	"time"
)

type OsInfo struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UUID          string    `json:"uuid" form:"uuid" gorm:"column:uuid"`                         //镜像UUID
	Name          string    `json:"name" form:"name" gorm:"column:name" sql:"not null"`          //镜像名称
	Status        int       `json:"status" form:"status" gorm:"column:status"`                   //状态 0 未启用 1启用
	OsType        string    `json:"os_type" form:"os_type" gorm:"column:os_type" sql:"not null"` //镜像类型
	Storage       string    `json:"storage" form:"storage" gorm:"column:storage" sql:"not null"` //存储信息
	OsXml         string    `json:"os_xml" form:"os_xml" gorm:"column:os_xml"`                   //镜像xml配置文件
	Annotation    string    `json:"annotation" form:"annotation" gorm:"column:annotation"`       //备注
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time" sql:"type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time" sql:"type:datetime"`
}

type OsInfoList []*OsInfo

func (s *OsInfo) TableName() string {
	return "os_info"
}

// FormatTime 特殊处理时间
func (s *OsInfo) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

func (s *OsInfo) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

func (s *OsInfo) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

func (s *OsInfo) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

func (s *OsInfo) Patch(v interface{}) (err error) {
	tmp := v.(*Host)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

func (s *OsInfo) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

func (s *OsInfo) GetByName() (err error) {
	err = conf.MysqlDb.Where("name = ?", s.Name).Find(s).Error
	return
}

func (s *OsInfo) List(ctx context.Context, PageSize, PageNum int64) (list VmList, count int64, err error) {
	sp, _ := utils.ExtractChildSpan("db:get os_info", ctx)
	defer sp.Finish()
	list = make(VmList, 0)
	//默认精确匹配
	db := conf.MysqlDb.Where(s)
	//可以自定义查询
	if s.Name != "" {
		db = conf.MysqlDb.Where("name like ?", fmt.Sprintf(`%%%s%%`, s.Name))
	}
	offset, limit := utils.GetOffsetAndLimit(PageSize, PageNum)
	err = db.Model(s).Count(&count).Offset(offset).Limit(limit).Find(&list).Error
	return list, count, err
}
