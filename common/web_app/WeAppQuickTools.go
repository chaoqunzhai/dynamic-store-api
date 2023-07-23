/**
@Author: chaoqun
* @Date: 2023/7/22 01:09
*/
package web_app

import (
	"encoding/json"
	"go-admin/common/redis_db"
)

type ValueRow struct {

}
type MakeWeAppQuickTools struct {
	CId       int
	ToolsData []interface{}
}

func (m *MakeWeAppQuickTools)LoadRedis()  {
	navLibMap := DiyDefineRow{
		Global:make(map[string]interface{},0),
		Value:make([]map[string]interface{},0),
	}

	_ = json.Unmarshal([]byte(DIY_VIEW_MEMBER_INDEX), &navLibMap)


	for _,row:=range navLibMap.Value {
		//先把自定义的工具列表加入进去
		if row["mode"] == "graphic"{
			row["list"] = m.ToolsData
		}
	}
	redis_db.SetMemberInfo(m.CId,navLibMap)

}
func NewMakeWeAppQuickTools(cid int) *MakeWeAppQuickTools {
	m := &MakeWeAppQuickTools{
		CId: cid,
		ToolsData:make([]interface{},0),
	}

	return m
}
