package apis

import (
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type GoodsSpecs struct {
	api.Api
}

// GetPage 获取GoodsSpecs列表
// @Summary 获取GoodsSpecs列表
// @Description 获取GoodsSpecs列表
// @Tags GoodsSpecs
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param goodsId query string false "商品ID"
// @Param name query string false "规格名称"
// @Param unit query string false "单位"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsSpecs}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-specs [get]
// @Security Bearer
func (e GoodsSpecs) GetPage(c *gin.Context) {
    req := dto.GoodsSpecsGetPageReq{}
    s := service.GoodsSpecs{}
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
	list := make([]models.GoodsSpecs, 0)
	var count int64


	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsSpecs失败，\r\n失败信息 %s", err.Error()))
        return
	}

	result:=make([]map[string]interface{},0)
	for _,row:=range list {
		dd:=map[string]interface{}{
			"goods_id":row.GoodsId,
			"id":row.Id,
			"unit":row.Unit,
			"price":row.Price, //销售的价格
			"original":row.Original,
			"name":row.Name,
			"inventory":row.Inventory,
			"market":row.Market,
		}
		var vipSpecs []models.GoodsVip
		e.Orm.Model(&vipSpecs).Select("grade_id,custom_price").Where("goods_id = ? and specs_id = ?",row.GoodsId,row.Id).Find(&vipSpecs)

		vipSpecList:=make([]string,0)
		for _,vip_spec:=range vipSpecs{
			if vip_spec.CustomPrice == 0 {
				continue
			}
			var vipRow models.GradeVip
			e.Orm.Model(&vipRow).Select("name,id").Where("c_id = ? and id = ?",row.CId,vip_spec.GradeId).Limit(1).Find(&vipRow)

			if vipRow.Id > 0 {
				vipSpecList = append(vipSpecList,fmt.Sprintf("%v: ¥%v",vipRow.Name,vip_spec.CustomPrice))
			}
		}
		dd["vip_spec_list"] = vipSpecList
		result = append(result,dd)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsSpecs
// @Summary 获取GoodsSpecs
// @Description 获取GoodsSpecs
// @Tags GoodsSpecs
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsSpecs} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-specs/{id} [get]
// @Security Bearer
func (e GoodsSpecs) Get(c *gin.Context) {
	req := dto.GoodsSpecsGetReq{}
	s := service.GoodsSpecs{}
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
	var object models.GoodsSpecs

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsSpecs失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建GoodsSpecs
// @Summary 创建GoodsSpecs
// @Description 创建GoodsSpecs
// @Tags GoodsSpecs
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsSpecsInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-specs [post]
// @Security Bearer
func (e GoodsSpecs) Insert(c *gin.Context) {
    req := dto.GoodsSpecsInsertReq{}
    s := service.GoodsSpecs{}
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
		e.Error(500, err, fmt.Sprintf("创建GoodsSpecs失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsSpecs
// @Summary 修改GoodsSpecs
// @Description 修改GoodsSpecs
// @Tags GoodsSpecs
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsSpecsUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-specs/{id} [put]
// @Security Bearer
func (e GoodsSpecs) Update(c *gin.Context) {
    req := dto.GoodsSpecsUpdateReq{}
    s := service.GoodsSpecs{}
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
		e.Error(500, err, fmt.Sprintf("修改GoodsSpecs失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除GoodsSpecs
// @Summary 删除GoodsSpecs
// @Description 删除GoodsSpecs
// @Tags GoodsSpecs
// @Param data body dto.GoodsSpecsDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-specs [delete]
// @Security Bearer
func (e GoodsSpecs) Delete(c *gin.Context) {
    s := service.GoodsSpecs{}
    req := dto.GoodsSpecsDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除GoodsSpecs失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
