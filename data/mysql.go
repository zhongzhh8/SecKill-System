package data

import (
	"SecKill/conf"
	"SecKill/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var Db *gorm.DB

// 初始化连接，确保创建表、索引等
func initMysql(config conf.AppConfig) {
	fmt.Println("Load dbService config...")

	// 设置连接相关的参数
	dbType := config.App.Database.Type
	usr := config.App.Database.User
	pwd := config.App.Database.Password
	address := config.App.Database.Address
	dbName := config.App.Database.DbName
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		usr, pwd, address, dbName)

	//创建一个数据库的连接，因为docker中的mysql服务启动时延，一开始需要尝试重试连接
	fmt.Println("Init dbService connections...")
	var err error
	for Db, err = gorm.Open(dbType, dbLink); err != nil; Db, err = gorm.Open(dbType, dbLink) {
		log.Println("Failed to connect database: ", err.Error())
		log.Println("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	// 设置连接池连接数
	Db.DB().SetMaxOpenConns(config.App.Database.MaxOpen)
	Db.DB().SetMaxIdleConns(config.App.Database.MaxIdle)

	// 初始化数据库
	user := model.User{}
	coupon := &model.Coupon{}

	// 创建表
	tables := []interface{}{user, coupon}

	for _, table := range tables {
		if !Db.HasTable(table) {
			Db.AutoMigrate(table)
		}
	}

	if config.App.FlushAllForTest {
		println("FlushAllForTest is true. Delete records of all tables.")
		for _, table := range tables {
			Db.Delete(table)
		}
	}

	// 创建唯一索引
	Db.Model(user).AddUniqueIndex("username_index", "username")  // 用户的用户名唯一
	Db.Model(coupon).AddUniqueIndex("coupon_index", "username", "coupon_name")  // 优惠券的(用户名, 优惠券名)唯一

	println("---Mysql connection is initialized.---")
	// 添加外键的demo代码
	// Db.Model(credit_card).
	//	 AddForeignKey("owner_id", "users(id)", "RESTRICT", "RESTRICT").
	//	 AddUniqueIndex("unique_owner", "owner_id")
}

