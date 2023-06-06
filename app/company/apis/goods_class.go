package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type GoodsClass struct {
	api.Api
}

// GetPage 获取GoodsClass列表
// @Summary 获取GoodsClass列表
// @Description 获取GoodsClass列表
// @Tags GoodsClass
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "商品分类名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsClass}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-class [get]
// @Security Bearer
func (e GoodsClass) GetPage(c *gin.Context) {
	req := dto.GoodsClassGetPageReq{}
	s := service.GoodsClass{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.GoodsClass, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsClass失败，\r\n失败信息 %s", err.Error()))
		return
	}

	result := make([]interface{}, 0)
	for _, row := range list {
		var bindCount int64
		whereSql := fmt.Sprintf("SELECT COUNT(*) as count from goods_mark_class where class_id = %v", row.Id)
		e.Orm.Raw(whereSql).Scan(&bindCount)
		row.GoodsCount = bindCount
		result = append(result, row)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsClass
// @Summary 获取GoodsClass
// @Description 获取GoodsClass
// @Tags GoodsClass
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsClass} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-class/{id} [get]
// @Security Bearer
func (e GoodsClass) Get(c *gin.Context) {
	req := dto.GoodsClassGetReq{}
	s := service.GoodsClass{}
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
	var object models.GoodsClass

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsClass失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建GoodsClass
// @Summary 创建GoodsClass
// @Description 创建GoodsClass
// @Tags GoodsClass
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsClassInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-class [post]
// @Security Bearer
func (e GoodsClass) Insert(c *gin.Context) {
	req := dto.GoodsClassInsertReq{}
	s := service.GoodsClass{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var countAll int64
	e.Orm.Model(&models.GoodsClass{}).Where("c_id = ?", userDto.CId).Count(&countAll)

	CompanyCnf := business.GetCompanyCnf(userDto.CId, "good_class", e.Orm)
	MaxNumber, _ := strconv.Atoi(CompanyCnf["good_class"])

	if countAll > int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("分类最多只可创建%v个", MaxNumber)), fmt.Sprintf("分类最多只可创建%v个", MaxNumber))
		return
	}
	var count int64
	e.Orm.Model(&models.GoodsClass{}).Where("c_id = ? and name = ?", userDto.CId, req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}

	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建分类失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsClass
// @Summary 修改GoodsClass
// @Description 修改GoodsClass
// @Tags GoodsClass
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsClassUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-class/{id} [put]
// @Security Bearer
func (e GoodsClass) Update(c *gin.Context) {
	req := dto.GoodsClassUpdateReq{}
	s := service.GoodsClass{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var count int64
	e.Orm.Model(&models.GoodsClass{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.GoodsClass
	e.Orm.Model(&models.GoodsClass{}).Where("name = ? and c_id = ?", req.Name, userDto.CId).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("分类修改失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除GoodsClass
// @Summary 删除GoodsClass
// @Description 删除GoodsClass
// @Tags GoodsClass
// @Param data body dto.GoodsClassDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-class [delete]
// @Security Bearer
func (e GoodsClass) Delete(c *gin.Context) {
	s := service.GoodsClass{}
	req := dto.GoodsClassDeleteReq{}
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
	newIds := make([]int, 0)
	for _, t := range req.Ids {
		var count int64
		whereSql := fmt.Sprintf("SELECT COUNT(*) as count from goods_mark_class where class_id = %v", t)
		e.Orm.Raw(whereSql).Scan(&count)
		if count == 0 {
			newIds = append(newIds, t)
		}
	}
	if len(newIds) == 0 {
		e.Error(500, errors.New("存在关联不可删除！"), "存在关联不可删除！")
		return
	}
	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除分类失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
