package apis

import (
	"errors"
	"fmt"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/utils"
	"go-admin/global"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
)

type PaymentOrder struct {
	api.Api
}
func (e PaymentOrder) GetPage(c *gin.Context) {
	req := dto.PayMetOrderGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.UserApplyPaymentOrder, 0)
	var count int64
	var data models.UserApplyPaymentOrder
	orm :=e.Orm
	if req.Before > 0 {
		orm = orm.Where("status > 0 ")
	}else {
		orm = orm.Where("status = 0 ")
	}
	err = orm.Model(&data).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("id desc").
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error

	result:=make([]models.UserApplyPaymentOrder,0)
	for _,row:=range list{
		var userRow sys.SysShopUser
		e.Orm.Model(&userRow).Select("user_id,username").Where("user_id = ? and c_id = ?",row.CreateBy,row.CId).Limit(1).Find(&userRow)
		if userRow.UserId > 0 {
			row.UserName = userRow.Username
		}
		result = append(result,row)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e PaymentOrder) Update(c *gin.Context) {
	req := dto.PayMetOrderUpdateReq{}
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

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var paymentObject models.UserApplyPaymentOrder
	e.Orm.Model(&models.UserApplyPaymentOrder{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Limit(1).Find(&paymentObject)
	if paymentObject.Id == 0 {

		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	updateMap:=make(map[string]interface{},0)
	updateMap["status"] = req.Status
	updateMap["use_to"] = req.UseTo
	updateMap["approve_msg"] = req.Desc
	//作废,更新即可
	if req.Status == 2 {
		e.Orm.Model(&models.UserApplyPaymentOrder{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Updates(&updateMap)
		e.OK("", "修改成功")
		return
	}
	if req.Status != 1 {
		e.Error(500, errors.New("非法状态"), "非法状态")
		return
	}
	updateErr:=e.Orm.Model(&models.UserApplyPaymentOrder{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Updates(&updateMap).Error
	if updateErr != nil{
		e.Error(500, errors.New("更新错误"), "更新错误")
		return
	}
	var shopObject models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("c_id = ? and user_id = ?",userDto.CId,paymentObject.CreateBy).Limit(1).Find(&shopObject)
	if shopObject.Id == 0 {

		e.Error(500, errors.New("用户商户平台不存在"), "用户商户平台不存在")
		return
	}
	switch req.UseTo {
	case 0://只是记录

	case 1://加入余额中
		//价格加入到用户余额中
		balance := shopObject.Balance + paymentObject.Money
		Money,_:=utils.RoundDecimal(balance).Float64()
		e.Orm.Model(&models2.Shop{}).Where("c_id = ? and user_id = ?",userDto.CId,paymentObject.CreateBy).Updates(map[string]interface{}{
			"balance":Money,
		})
		//增加余额明细
		row:=models2.ShopBalanceLog{
			CId: userDto.CId,
			ShopId: shopObject.Id,
			Money: paymentObject.Money,
			Scene:fmt.Sprintf("用户[%v] 提交付款单,审批通过,增加余额:%v",userDto.Username,paymentObject.Money),
			Action: global.UserNumberAdd, //增加
			Type: global.ScanAdmin,
		}
		row.CreateBy = userDto.UserId
		e.Orm.Create(&row)
	case 2://加入到授信额中
		//价格加入到用户授信额中
		credit := shopObject.Credit + paymentObject.Money
		Money,_:=utils.RoundDecimal(credit).Float64()
		e.Orm.Model(&models2.Shop{}).Where("c_id = ? and user_id = ?",userDto.CId,paymentObject.CreateBy).Updates(map[string]interface{}{
			"credit":Money,
		})
		//增加授信明细
		row:=models2.ShopCreditLog{
			CId: userDto.CId,
			ShopId: shopObject.Id,
			Number: paymentObject.Money,
			Scene:fmt.Sprintf("用户[%v] 提交付款单,审批通过,增加授信额:%v",userDto.Username,paymentObject.Money),
			Action: global.UserNumberAdd, //增加
			Type: global.ScanAdmin,
		}
		row.CreateBy = userDto.UserId
		e.Orm.Create(&row)
	default:
		e.Error(500, errors.New("非法状态"), "非法状态")
		return
	}

	e.OK("", "修改成功")
	return
}