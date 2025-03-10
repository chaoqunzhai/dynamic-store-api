package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type CompanyArticle struct {
	api.Api
}
type Message struct {

	Context string `json:"context"`
}

func (e CompanyArticle) Message(c *gin.Context) {

	err := e.MakeContext(c).
		MakeOrm().
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
	var object models2.Message

	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Limit(1).Find(&object)


	e.OK(object.Context,"操作成功")
	return
}


func (e CompanyArticle) UpdateMessage(c *gin.Context) {
	req := Message{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
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
	var object models2.Message

	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Limit(1).Find(&object)
	if object.Id > 0 {

		object.Context = req.Context
		e.Orm.Save(&object)
	}else {
		object = models2.Message{
			Context: req.Context,
		}
		object.CId = userDto.CId
		object.CreateBy = userDto.UserId
		e.Orm.Create(&object)
	}
	e.OK("","操作成功")
	return

}
func (e CompanyArticle) Enable(c *gin.Context) {
	req := dto.UpdateEnableReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).

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
	e.Orm.Model(&models.CompanyArticle{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Updates(map[string]interface{}{
		"enable":req.Enable,
	})
	e.OK("", "更新成功")
}


func (e CompanyArticle) GetPage(c *gin.Context) {
    req := dto.CompanyArticleGetPageReq{}
    s := service.CompanyArticle{}
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
	list := make([]models.CompanyArticle, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyArticle失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CompanyArticle
// @Summary 获取CompanyArticle
// @Description 获取CompanyArticle
// @Tags CompanyArticle
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyArticle} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-article/{id} [get]
// @Security Bearer
func (e CompanyArticle) Get(c *gin.Context) {
	req := dto.CompanyArticleGetReq{}
	s := service.CompanyArticle{}
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
	var object models.CompanyArticle

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyArticle失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CompanyArticle
// @Summary 创建CompanyArticle
// @Description 创建CompanyArticle
// @Tags CompanyArticle
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyArticleInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-article [post]
// @Security Bearer
func (e CompanyArticle) Insert(c *gin.Context) {
    req := dto.CompanyArticleInsertReq{}
    s := service.CompanyArticle{}
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
	req.SetCreateBy(userDto.UserId)
	req.CId = userDto.CId
	req.Enable = true

	var count int64
	e.Orm.Model(&models.CompanyArticle{}).Where("c_id = ? and title = ?",userDto.CId,req.Title).Count(&count)
	if count > 0 {

		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建CompanyArticle失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyArticle
// @Summary 修改CompanyArticle
// @Description 修改CompanyArticle
// @Tags CompanyArticle
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyArticleUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-article/{id} [put]
// @Security Bearer
func (e CompanyArticle) Update(c *gin.Context) {
    req := dto.CompanyArticleUpdateReq{}
    s := service.CompanyArticle{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var count int64
	e.Orm.Model(&models.CompanyArticle{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.CompanyArticle
	e.Orm.Model(&models.CompanyArticle{}).Where("title = ? and c_id = ?", req.Title, userDto.CId).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	req.CId = userDto.CId
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyArticle失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CompanyArticle
// @Summary 删除CompanyArticle
// @Description 删除CompanyArticle
// @Tags CompanyArticle
// @Param data body dto.CompanyArticleDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-article [delete]
// @Security Bearer
func (e CompanyArticle) Delete(c *gin.Context) {
    s := service.CompanyArticle{}
    req := dto.CompanyArticleDeleteReq{}
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

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CompanyArticle失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
