package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
)

type PayApi struct {
	api.Api
}
type PayCnfInsertReq struct {
	BalanceDeduct bool `json:"balance_deduct" gorm:"size:1;comment:是否开启余额支付"`
	Alipay bool `json:"alipay" gorm:"size:1;comment:是否开启阿里支付"`
	WeChat bool `json:"we_chat" gorm:"size:1;comment:是否开启微信支付"`
	Credit bool `json:"credit" gorm:"size:1;comment:支持授信减扣"`
	CashOn bool `json:"cash_on"`

}
func (e *PayApi) Create(c *gin.Context) {
	req := PayCnfInsertReq{}
	err := e.MakeContext(c).
		Bind(&req, binding.JSON, nil).
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

	var PayCnf models.PayCnf
	e.Orm.Model(&models.PayCnf{}).Scopes(actions.PermissionSysUser(PayCnf.TableName(), userDto)).Limit(1).Find(&PayCnf)


	if PayCnf.Id > 0 {

		PayCnf.Ali = req.Alipay
		PayCnf.CashOn = req.CashOn
		PayCnf.BalanceDeduct = req.BalanceDeduct
		PayCnf.WeChat = req.WeChat
		PayCnf.Credit = req.Credit
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.PayCnf{
			Ali: req.Alipay,
			CashOn: req.CashOn,
			WeChat: req.WeChat,
			BalanceDeduct: req.BalanceDeduct,
			Credit: req.Credit,
		}
		trade.CId = userDto.CId
		trade.Enable = true
		e.Orm.Create(&trade)
	}
	e.OK("","successful")
	return
}
func (e PayApi) Detail(c *gin.Context) {
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
	var PayCnf models.PayCnf
	e.Orm.Model(&models.PayCnf{}).Scopes(actions.PermissionSysUser(PayCnf.TableName(), userDto)).Limit(1).Find(&PayCnf)


	if PayCnf.Id == 0 {
		object := models.PayCnf{
			Credit: true,
		}
		object.Enable = true
		object.CId = userDto.CId
		e.Orm.Create(&object)
		e.OK(object,"successful")
	}
	e.OK(PayCnf,"successful")
	return
}
