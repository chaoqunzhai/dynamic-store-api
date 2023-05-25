package apis

import (
	"fmt"
	"go-admin/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/jwt/user"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Company struct {
	api.Api
}

// GetPage 获取Company列表
// @Summary 获取Company列表
// @Description 获取Company列表
// @Tags Company
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param name query string false "公司(大B)名称"
// @Param phone query string false "负责人联系手机号"
// @Param userName query string false "大B负责人名称"
// @Param shop query string false "自定义大B系统名称"
// @Param renewalTime query time.Time false "续费时间"
// @Param expirationTime query time.Time false "到期时间"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Company}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company [get]
// @Security Bearer
func (e Company) GetPage(c *gin.Context) {
	req := dto.CompanyGetPageReq{}
	s := service.Company{}
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
	list := make([]models.Company, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Company失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e Company) MonitorData(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//实时订单数据
	overview := make(map[string]interface{}, 0)
	overview = map[string]interface{}{
		"orderTotalPrice": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
		"orderTotal": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
		"newUserTotal": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
		"consumeUserTotal": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
	}
	//统计
	statistics := make(map[string]interface{}, 0)
	statistics = map[string]interface{}{
		"goodsTotal":       "12",
		"userTotal":        "1",
		"orderTotal":       "0",
		"consumeUserTotal": "0",
	}
	//待办
	pending := make(map[string]interface{}, 0)
	pending = map[string]interface{}{
		"goodsTotal":       "12",
		"userTotal":        "1",
		"orderTotal":       "0",
		"consumeUserTotal": "0",
	}
	//近七日交易走势
	tradeTrend := make(map[string]interface{}, 0)
	tradeTrend = map[string]interface{}{
		"date": []string{
			"2023-05-19",
			"2023-05-20",
			"2023-05-21",
			"2023-05-22",
			"2023-05-23",
			"2023-05-24",
			"2023-05-25",
		},
		"orderTotal": []string{
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
		},
		"orderTotalPrice": []string{
			"0.00",
			"0.00",
			"0.00",
			"0.00",
			"0.00",
			"0.00",
			"0.00",
		},
	}
	result := map[string]interface{}{
		"overview":   overview,
		"statistics": statistics,
		"pending":    pending,
		"tradeTrend": tradeTrend,
	}
	e.OK(result, "successful")
	return
}
func (e Company) Info(c *gin.Context) {
	req := dto.CompanyGetReq{}
	s := service.Company{}
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
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	storeInfo := map[string]interface{}{
		"store_id":      0,
		"store_name":    "暂无信息",
		"describe":      global.Describe,
		"logo_image_id": 0,
		"sort":          100,
		"is_recycle":    0,
		"is_delete":     0,
		"create_time":   time.Now().Format("2006-01-02 15:04:05"),
		"update_time":   time.Now().Format("2006-01-02 15:04:05"),
		"logoImage":     "",
	}
	if userDto.RoleId == global.RoleSuper {
		storeInfo = map[string]interface{}{
			"store_id":      0,
			"store_name":    global.SysName,
			"describe":      global.Describe,
			"logo_image_id": 0,
			"sort":          100,
			"is_recycle":    0,
			"is_delete":     0,
			"create_time":   time.Now().Format("2006-01-02 15:04:05"),
			"update_time":   time.Now().Format("2006-01-02 15:04:05"),
			"logoImage":     "",
		}
	} else {
		if userDto.CId == 0 {

			e.OK(storeInfo, "successful")
			return
		}
		var object models.Company
		e.Orm.Model(&models.Company{}).Where("enable = 1 and id = ?", userDto.CId).First(&object)

		if object.Id == 0 {
			storeInfo["store_name"] = "已经下线"
			e.OK(storeInfo, "successful")
			return
		}
		storeInfo = map[string]interface{}{
			"store_id":      object.Id,
			"store_name":    object.Name,
			"describe":      object.Desc,
			"logo_image_id": 0,
			"sort":          object.Layer,
			"is_recycle":    0,
			"is_delete":     0,
			"create_time":   object.CreatedAt.Format("2006-01-02 15:04:05"),
			"update_time":   object.UpdatedAt.Format("2006-01-02 15:04:05"),
			"logoImage":     "",
		}
	}
	//如果超管,那就返回超管的一些自定义信息

	//如果是大B,那就查询company

	e.OK(storeInfo, "successful")
	return
}

// Get 获取Company
// @Summary 获取Company
// @Description 获取Company
// @Tags Company
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Company} "{"code": 200, "data": [...]}"
// @Router /api/v1/company/{id} [get]
// @Security Bearer
func (e Company) Get(c *gin.Context) {
	req := dto.CompanyGetReq{}
	s := service.Company{}
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
	var object models.Company

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Company失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建Company
// @Summary 创建Company
// @Description 创建Company
// @Tags Company
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company [post]
// @Security Bearer
func (e Company) Insert(c *gin.Context) {
	req := dto.CompanyInsertReq{}
	s := service.Company{}
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建Company失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Company
// @Summary 修改Company
// @Description 修改Company
// @Tags Company
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company/{id} [put]
// @Security Bearer
func (e Company) Update(c *gin.Context) {
	req := dto.CompanyUpdateReq{}
	s := service.Company{}
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
		e.Error(500, err, fmt.Sprintf("修改Company失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除Company
// @Summary 删除Company
// @Description 删除Company
// @Tags Company
// @Param data body dto.CompanyDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company [delete]
// @Security Bearer
func (e Company) Delete(c *gin.Context) {
	s := service.Company{}
	req := dto.CompanyDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除Company失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
