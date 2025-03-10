package version

import (
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"runtime"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if config.DatabaseConfig.Driver == "mysql" {
			tx = tx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
		err := tx.Migrator().AutoMigrate(
			new(models.SysDept),
			new(models.SysConfig),
			new(models.SysTables),
			new(models.SysColumns),
			new(models.SysMenu),
			new(models.SysLoginLog),
			new(models.SysOperaLog),
			new(models.SysRoleDept),
			new(models.SysUser),
			new(models.SysShopUser),
			new(models.SysRole),
			new(models.SysPost),
			new(models.DictData),
			new(models.DictType),
			new(models.SysJob),
			new(models.SysConfig),
			new(models.SysApi),

			new(models.DyNamicMenu),
			new(models.Line),
			new(models.Driver),
			new(models.GoodsSales),
			new(models.GoodsClass),
			new(models.GoodsTag),
			new(models.GoodsCollect),
			new(models.ShopTag),
			new(models.ExtendUser),
			new(models.Company),
			new(models.CompanyCategory),

			//物流方式
			new(models.CompanyExpress),
			//运费配置
			new(models.CompanyFreight),
			new(models.CompanyExpressStore),
			new(models.CompanyRole),
			//大B配置
			new(models.CompanyRegisterRule),
			new(models.CompanyRegisterUserVerify),
			//new(models.CompanyRegisterCnf),
			new(models.CompanyWeAppCnf),
			new(models.CompanyQuotaCnf),
			new(models.CompanyLineCnf),
			new(models.CompanyLineCnfLog),
			new(models.Line),
			new(models.CompanySmsQuotaCnf),
			new(models.CompanySmsQuotaCnfLog),
			new(models.CompanySmsRecordLog),
			new(models.CompanyRenewalTimeLog),
			new(models.WeChatAppIdCnf),
			new(models.WeChatOfficialPay),
			new(models.AliPay),
			new(models.OrderTrade),
			new(models.SplitTableMap),
			new(models.GradeVip),
			new(models.PayCnf),
			new(models.DebitCard),
			new(models.OfflinePay),
			//商品
			new(models.Shop),
			new(models.ShopRechargeLog),
			new(models.ShopBalanceLog),
			new(models.ShopIntegralLog),
			new(models.ShopCreditLog),
			new(models.ShopOrderRecord),
			new(models.ShopOrderBindRecord),
			new(models.Coupon),
			new(models.ReceiveCouponLog),
			new(models.UserAmountStore),
			new(models.CycleTimeConf),
			//订单
			new(models.Orders),
			new(models.OrderSpecs),
			new(models.OrderExtend),

			new(models.OrderToRedisMap),
			new(models.OrderCycleCnf),
			new(models.Goods),
			new(models.GoodsDesc),
			new(models.GoodsSpecs),
			new(models.GoodsVip),

			new(models.Article),
			new(models.Message),
			new(models.Ads),
			//小程序配置
			new(models.WeAppGlobalNavCnf),
			new(models.CompanyNavCnf),
			new(models.WeAppQuickTools),
			new(models.CompanyQuickTools),
			new(models.WeAppExtendCnf),
			new(models.VipShowEnable),

			new(models.DynamicUserAddress),
		)
		if err != nil {
			return err
		}
		//初始化DB的时候执行
		//if err := models.InitDb(tx); err != nil {
		//	return err
		//}
		return nil
		//return tx.Create(&common.Migration{
		//	Version: version,
		//}).Error
	})
}
