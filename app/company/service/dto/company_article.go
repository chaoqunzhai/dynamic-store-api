package dto

import (
     
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyArticleGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company_article" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_article" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:company_article" comment:"大BID"`
    Title string `form:"title"  search:"type:exact;column:title;table:company_article" comment:"文章标题"`
    CompanyArticleOrder
}

type CompanyArticleOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company_article"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company_article"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_article"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_article"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_article"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_article"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company_article"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company_article"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:company_article"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:company_article"`
    CoverImage string `form:"coverImageOrder"  search:"type:order;column:cover_image;table:company_article"`
    Title string `form:"titleOrder"  search:"type:order;column:title;table:company_article"`
    Class string `form:"classOrder"  search:"type:order;column:class;table:company_article"`
    Context string `form:"contextOrder"  search:"type:order;column:context;table:company_article"`
    
}

func (m *CompanyArticleGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyArticleInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    CId string `json:"cId" comment:"大BID"`
    CoverImage string `json:"coverImage" comment:"封面图片"`
    Title string `json:"title" comment:"文章标题"`
    Class string `json:"class" comment:"文章分类"`
    Context string `json:"context" comment:"文章内容"`
    common.ControlBy
}

func (s *CompanyArticleInsertReq) Generate(model *models.CompanyArticle)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.CId = s.CId
    model.CoverImage = s.CoverImage
    model.Title = s.Title
    model.Class = s.Class
    model.Context = s.Context
}

func (s *CompanyArticleInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyArticleUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    CId string `json:"cId" comment:"大BID"`
    CoverImage string `json:"coverImage" comment:"封面图片"`
    Title string `json:"title" comment:"文章标题"`
    Class string `json:"class" comment:"文章分类"`
    Context string `json:"context" comment:"文章内容"`
    common.ControlBy
}

func (s *CompanyArticleUpdateReq) Generate(model *models.CompanyArticle)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.CId = s.CId
    model.CoverImage = s.CoverImage
    model.Title = s.Title
    model.Class = s.Class
    model.Context = s.Context
}

func (s *CompanyArticleUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyArticleGetReq 功能获取请求参数
type CompanyArticleGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyArticleGetReq) GetId() interface{} {
	return s.Id
}

// CompanyArticleDeleteReq 功能删除请求参数
type CompanyArticleDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyArticleDeleteReq) GetId() interface{} {
	return s.Ids
}
