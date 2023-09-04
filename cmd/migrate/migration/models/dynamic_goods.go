package models

// todo:商品信息
type Goods struct {
	BigBRichGlobal
	Name      string       `gorm:"size:50;comment:商品名称"`
	Subtitle  string       `gorm:"size:8;comment:商品广告"`
	Image     string       `gorm:"size:100;comment:商品图片路径"`
	Quota     bool         `gorm:"comment:是否限购"`
	VipSale   bool         `gorm:"comment:会员价"`
	Inventory int          `gorm:"comment:库存"`
	Sale      int          `json:"sale" gorm:"comment:销量"`
	Recommend bool         `json:"recommend" gorm:"comment:是否推荐"`
	SpecName  string       `gorm:"size:8;comment:规格命名,例如是:颜色,重量,系列"`
	Money     string       `gorm:"size:12;comment:价格区间显示"`
	Tag       []GoodsTag   `gorm:"many2many:goods_mark_tag;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:tag_id;"`
	Class     []GoodsClass `gorm:"many2many:goods_mark_class;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:class_id;"`
}

func (Goods) TableName() string {
	return "goods"
}

//todo:商品详情
type GoodsDesc struct {
	Model
	GoodsId int
	Desc      string //描述内容
}
func (GoodsDesc) TableName() string {
	return "goods_desc"
}
// todo: 规格名称
type GoodsSpecs struct {
	BigBMiniGlobal
	GoodsId   int     `gorm:"index;comment:商品ID"`
	Name      string  `gorm:"size:30;comment:规格名称"`
	Price     float32 `gorm:"comment:售价"`
	Original  float32 `gorm:"comment:原价"`
	Inventory int     `gorm:"comment:库存"`
	Unit      string  `gorm:"size:8;comment:单位"`
	Limit     int     `gorm:"comment:起售量"`
	Code      string  `gorm:"size:30;comment:条形码"`
	Image     string  `gorm:"size:100;comment:商品图片路径"`
}

func (GoodsSpecs) TableName() string {
	return "goods_specs"
}

// todo: 商品VIP价格
type GoodsVip struct {
	BigBMiniGlobal
	GoodsId     int     `gorm:"index;comment:商品ID"`
	SpecsId     int     `gorm:"index;comment:规格ID"`
	GradeId     int     `gorm:"index;comment:VipId"`
	CustomPrice float32 `gorm:"index;comment:自定义价格"`
}

func (GoodsVip) TableName() string {
	return "goods_vip"
}

// todo: 记录产品的销量和库存
type GoodsSales struct {
	BigBMiniGlobal
	ProductId   int    `gorm:"index;comment:产品ID"`
	ProductName string `gorm:"size:30;comment:产品名称"`
	Sales       int    `gorm:"comment:当时销量"`
	Inventory   int    `gorm:"comment:当时剩余库存"`
}

func (GoodsSales) TableName() string {
	return "goods_sales"
}

// todo:商品分类
type GoodsClass struct {
	BigBRichGlobal
	Name  string `gorm:"index;size:8;comment:商品分类名称"`
	Image string `gorm:"size:60;comment:商品分类图片路径"`
}

func (GoodsClass) TableName() string {
	return "goods_class"
}

// todo:商品标签
type GoodsTag struct {
	BigBRichGlobal
	Name  string `gorm:"index;size:8;comment:商品标签名称"`
	Color string `gorm:"size:10;comment:标签颜色"`
}

func (GoodsTag) TableName() string {
	return "goods_tag"
}
