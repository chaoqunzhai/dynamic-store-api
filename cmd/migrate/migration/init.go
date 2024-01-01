package migration

import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"log"
	"path/filepath"
	"sort"
	"sync"

	"gorm.io/gorm"
)

var Migrate = &Migration{
	version: make(map[string]func(db *gorm.DB, version string) error),
}

type Migration struct {
	db      *gorm.DB
	version map[string]func(db *gorm.DB, version string) error
	mutex   sync.Mutex
}

func (e *Migration) GetDb() *gorm.DB {
	return e.db
}

func (e *Migration) SetDb(db *gorm.DB) {
	e.db = db
}

func (e *Migration) SetVersion(k string, f func(db *gorm.DB, version string) error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.version[k] = f
}

func (e *Migration) Migrate() {
	versions := make([]string, 0)
	for k := range e.version {
		versions = append(versions, k)
	}
	if !sort.StringsAreSorted(versions) {
		sort.Strings(versions)
	}
	var err error
	//var count int64
	for _, v := range versions {
		//err = e.db.Table("sys_migration").Where("version = ?", v).Count(&count).Error
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//if count > 0 {
		//	log.Println(count)
		//	count = 0
		//	continue
		//}
		err = (e.version[v])(e.db, v)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (e *Migration) SelectTableMigrateTable(tables []string) {
	// 输出获取到的表名列表
	for _, table := range tables {
		fmt.Printf("table %v,开始迁移\n",table)
		res:=e.db.Table(table).AutoMigrate(&models.Orders{})
		if res != nil {
			fmt.Printf("table %v,迁移失败 err:%v\n",table,res.Error())
			continue
		}
		fmt.Printf("table %v,迁移成功\n",table)

	}
}
func GetFilename(s string) string {
	s = filepath.Base(s)
	return s[:13]
}
