package apis

import (
	"errors"
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type CompanyCoupon struct {
	api.Api
}

// GetPage 获取CompanyCoupon列表
// @Summary 获取CompanyCoupon列表
// @Description 获取CompanyCoupon列表
// @Tags CompanyCoupon
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "优惠卷名称"
// @Param type query string false "类型"
// @Param range query string false "使用范围"
// @Param startTime query time.Time false "开始使用时间"
// @Param endTime query time.Time false "截止使用时间"
// @Param inventory query string false "库存"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyCoupon}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-coupon [get]
// @Security Bearer
func (e CompanyCoupon) GetPage(c *gin.Context) {
	req := dto.CompanyCouponGetPageReq{}
	s := service.CompanyCoupon{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.CompanyCoupon, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyCoupon失败，\r\n失败信息 %s", err.Error()))
		return
	}
	result := make([]map[string]interface{}, 0)

	for _, row := range list {
		r := map[string]interface{}{
			"id":          row.Id,
			"name":        row.Name,
			"type":        row.Type,
			"expire_type": row.ExpireType,
			"expire_day":  row.ExpireDay,
			"reduce":      row.Reduce,
			"enable":      row.Enable,
			"discount":    row.Discount,
			"start_time": func() string {
				if row.StartTime.Valid {
					return row.StartTime.Time.Format("2006-01-02 15:04:05")
				}
				return ""
			}(),
			"end_time": func() string {
				if row.EndTime.Valid {
					return row.EndTime.Time.Format("2006-01-02 15:04:05")
				}
				return ""
			}(),
			"inventory":row.Inventory,
			"threshold":   row.Threshold,
			"receive_num": row.ReceiveNum,
			"layer":       row.Layer,
			"created_at":  row.CreatedAt,
			"updated_at":row.UpdatedAt,
			"first":row.First,
		}
		result = append(result, r)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e CompanyCoupon) Receive(c *gin.Context) {
	req := dto.CompanyCouponReceiveGetPageReq{}
	s := service.CompanyCoupon{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models2.ReceiveCouponLog, 0)
	var count int64

	err = s.GetReceivePage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领劵记录失败,%s", err.Error()))
		return
	}
	result := make([]map[string]interface{}, 0)

	for _, row := range list {
		var shopUser models2.SysShopUser
		e.Orm.Model(&shopUser).Select("user_id,username").Where("c_id = ? and user_id = ?",row.CId,row.UserId).Limit(1).Find(&shopUser)
		var counpon models.CompanyCoupon
		e.Orm.Model(&counpon).Select("id,name").Where("c_id = ? and id = ?",row.CId,row.CouponId).Limit(1).Find(&counpon)
		var typeName string
		switch row.Type {
		case global.ReceiveCoupon1:
			typeName = "订单领取"
		case global.ReceiveCoupon2:
			typeName = "直接领取"
		case global.ReceiveCoupon3:
			typeName = "活动领取"
		}
		var statusName string
		switch row.Status {
		case global.CouponState1:
			statusName = "未使用"
		case global.CouponState2:
			statusName = "已使用"
		case global.CouponState3:
			statusName = "过期"
		case global.CouponState4:
			statusName = "作废"
		}
		r := map[string]interface{}{
			"id":row.Id,
			"created_at":  row.CreatedAt.Format("2006-01-02 15:04:05"),
			"coupon_type":global.GetCouponType(row.CouponType),
			"type":typeName, //1订单2.直接领取3.活动领取
			"status_name":statusName,
			"status":row.Status,
		}
		if shopUser.UserId > 0 {
			r["username"] = shopUser.Username
		}
		if counpon.Id > 0 {
			r["name"] = counpon.Name
		}
		result = append(result, r)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}


// Get 获取CompanyCoupon
// @Summary 获取CompanyCoupon
// @Description 获取CompanyCoupon
// @Tags CompanyCoupon
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyCoupon} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-coupon/{id} [get]
// @Security Bearer
func (e CompanyCoupon) Get(c *gin.Context) {
	req := dto.CompanyCouponGetReq{}
	s := service.CompanyCoupon{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.CompanyCoupon

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyCoupon失败，\r\n失败信息 %s", err.Error()))
		return
	}
	if object.StartTime.Valid {

		object.Start = object.StartTime.Time.Format("2006-01-02 15:04:05")
	}
	if object.EndTime.Valid {
		object.End = object.EndTime.Time.Format("2006-01-02 15:04:05	")
	}

	e.OK(object, "查询成功")
}

// Insert 创建CompanyCoupon
// @Summary 创建CompanyCoupon
// @Description 创建CompanyCoupon
// @Tags CompanyCoupon
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyCouponInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-coupon [post]
// @Security Bearer
func (e CompanyCoupon) Insert(c *gin.Context) {
	req := dto.CompanyCouponInsertReq{}
	s := service.CompanyCoupon{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	var count int64
	var object models.CompanyCoupon
	e.Orm.Model(&models.CompanyCoupon{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where(" name = ?",req.Name).Count(&count)
	if count > 0 {

		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建优惠卷失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyCoupon
// @Summary 修改CompanyCoupon
// @Description 修改CompanyCoupon
// @Tags CompanyCoupon
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyCouponUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-coupon/{id} [put]
// @Security Bearer
func (e CompanyCoupon) Update(c *gin.Context) {
	req := dto.CompanyCouponUpdateReq{}
	s := service.CompanyCoupon{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyCoupon失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除CompanyCoupon
// @Summary 删除CompanyCoupon
// @Description 删除CompanyCoupon
// @Tags CompanyCoupon
// @Param data body dto.CompanyCouponDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-coupon [delete]
// @Security Bearer
func (e CompanyCoupon) Delete(c *gin.Context) {
	s := service.CompanyCoupon{}
	req := dto.CompanyCouponDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CompanyCoupon失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
