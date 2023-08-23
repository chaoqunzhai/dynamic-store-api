package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/global"
	"gorm.io/gorm"
	"strings"
	"time"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CompanyCoupon struct {
	service.Service
}

// GetPage 获取CompanyCoupon列表
func (e *CompanyCoupon) GetPage(c *dto.CompanyCouponGetPageReq, p *actions.DataPermission, list *[]models.CompanyCoupon, count *int64) error {
	var err error
	var data models.CompanyCoupon

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyCouponService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyCoupon对象
func (e *CompanyCoupon) Get(d *dto.CompanyCouponGetReq, p *actions.DataPermission, model *models.CompanyCoupon) error {
	var data models.CompanyCoupon

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyCoupon error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyCoupon对象
func (e *CompanyCoupon) Insert(cid int, c *dto.CompanyCouponInsertReq) error {
	var err error
	var data models.CompanyCoupon
	c.Generate(&data)
	data.CId = cid
	//todo:时间处理
	fmt.Println("expr", c.ExpireType)
	if c.ExpireType == 1 {
		if len(c.BetweenTime) != 2 {
			return errors.New("请填入开始结束时间")
		}
		startTime := c.BetweenTime[0]
		endTime := c.BetweenTime[1]
		if startTime != "" {
			startTime = strings.Replace(startTime, "T", " ", -1)
			startTime = strings.Replace(startTime, "Z", "", -1)
			t, _ := time.Parse("2006-01-02 15:04:05", startTime)
			//fmt.Println("startTime", startTime, "t", t)
			data.StartTime = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		if endTime != "" {
			endTime = strings.Replace(endTime, "T", " ", -1)
			endTime = strings.Replace(endTime, "Z", "", -1)
			t, _ := time.Parse("2006-01-02 15:04:05", endTime)
			//fmt.Println("endTime", endTime, "t", t)
			data.EndTime = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
	}

	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyCouponService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyCoupon对象
func (e *CompanyCoupon) Update(c *dto.CompanyCouponUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CompanyCoupon{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)
	if c.ExpireType == 1 {
		startTime := c.BetweenTime[0]
		endTime := c.BetweenTime[1]
		//fmt.Println("BetweenTime", c.BetweenTime)
		if startTime != "" {
			startTime = strings.Replace(startTime, "T", " ", -1)
			startTime = strings.Replace(startTime, "Z", "", -1)
			t, _ := time.Parse("2006-01-02 15:04:05", startTime)
			data.StartTime = sql.NullTime{
				Time:  t,
				Valid: true,
			}
			fmt.Println("startTime", startTime, "t", t)
		}
		if endTime != "" {
			endTime = strings.Replace(endTime, "T", " ", -1)
			endTime = strings.Replace(endTime, "Z", "", -1)
			t, _ := time.Parse("2006-01-02 15:04:05", endTime)
			fmt.Println("endTime", endTime, "t", t)
			data.EndTime = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		data.ExpireDay = 7
	}

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CompanyCouponService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CompanyCoupon
func (e *CompanyCoupon) Remove(d *dto.CompanyCouponDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyCoupon

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCompanyCoupon error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
