package dbService

import (
	"SecKill/data"
	"SecKill/model"
	"fmt"
)

func GetAllCoupons() ([]model.Coupon, error) {
	var coupons []model.Coupon
	result := data.Db.Find(&coupons)
	return coupons, result.Error
}
//下面原本应该使用gorm封装好的函数来操作数据库的，但是由于当时出了bug，没时间处理，所以直接写数据库命令来操作。
// 插入用户拥有优惠券的数据
func UserHasCoupon(userName string, coupon model.Coupon) error {
	return data.Db.Exec(fmt.Sprintf("INSERT IGNORE INTO coupons " +
		"(`username`,`coupon_name`,`amount`,`left`,`stock`,`description`) " +
		"values('%s', '%s', %d, %d, %f, '%s')",
		userName, coupon.CouponName, 1, 1, coupon.Stock, coupon.Description)).Error
}

// 优惠券库存自减1
func DecreaseOneCouponLeft(sellerName string, couponName string) error {
	return data.Db.Exec(fmt.Sprintf("UPDATE coupons c SET c.left=c.left-1 WHERE " +
		"c.username='%s' AND c.coupon_name='%s' AND c.left>0", sellerName, couponName)).Error
}