package redisService

import (
	"SecKill/api/dbService"
	"SecKill/data"
	"SecKill/model"
	"fmt"
	"log"
	"strconv"
)


// 获取"用户持有优惠券"的key
func getHasCouponsKeyByName(userName string) string {
	return fmt.Sprintf("%s-has", userName)
}

// 获取"优惠券"的key
func getCouponKeyByCoupon(coupon model.Coupon) string {
	return getCouponKeyByName(coupon.CouponName)
}
func getCouponKeyByName(couponName string) string {
	return fmt.Sprintf("%s-info", couponName)
}

// 缓存用户拥有优惠券/商家创建优惠券的信息
func CacheHasCoupon(coupon model.Coupon) (int64, error) {
	key := getHasCouponsKeyByName(coupon.Username) //得到的key其实就是 coupon.Username-has
	val, err := data.SetAdd(key, coupon.CouponName)
	return val, err
}

// 缓存优惠券的完整信息
func CacheCoupon(coupon model.Coupon) (string, error) {
	key := getCouponKeyByCoupon(coupon)
	fields := map[string]interface{}{
		"id": coupon.Id,
		"username": coupon.Username,
		"couponName": coupon.CouponName,
		"amount": coupon.Amount,
		"left": coupon.Left,
		"stock": coupon.Stock,
		"description": coupon.Description,
	}
	val, err := data.SetMapForever(key, fields)
	return val, err
}

// 缓存优惠券
func CacheCouponAndHasCoupon(coupon model.Coupon) error {
	if _, err := CacheHasCoupon(coupon); err != nil {
		return err
	}

	// user = 根据优惠券的username查user
	if user, err := dbService.GetUser(coupon.Username); err != nil {
		log.Println("Database service error: ", err)
		return err
	} else {
		if user.IsSeller() {
			_, err = CacheCoupon(coupon)
		}
		return err
	}
}

// 从缓存获取优惠券
func GetCoupon(couponName string) model.Coupon {
	key := getCouponKeyByName(couponName)
	values, err := data.GetMap(key, "id", "username", "couponName", "amount", "left", "stock", "description")
	if err != nil {
		println("Error on getting coupon. " + err.Error())
	}
	// log.Println(values) TODO
	// values[0]类型是nil，说明key是不存在的？
	id, err := strconv.ParseInt(values[0].(string), 10, 64)
	if err != nil {
		println("Wrong type of id. " + err.Error())
	}
	amount, err := strconv.ParseInt(values[3].(string), 10, 64)
	if err != nil {
		println("Wrong type of amount. " + err.Error())
	}
	left, err := strconv.ParseInt(values[4].(string), 10, 64)
	if err != nil {
		println("Wrong type of left. " + err.Error())
	}
	stock, err := strconv.ParseInt(values[5].(string), 10, 64)
	if err != nil {
		println("Wrong type of stock. " + err.Error())
	}
	return model.Coupon{
		Id:          id,
		Username:    values[1].(string),
		CouponName:  values[2].(string),
		Amount:      amount,
		Left:        left,
		Stock:       stock,
		Description: values[6].(string),
	}

}

// 从缓存获取某个用户的所有优惠券
func GetCoupons(userName string) ([]model.Coupon, error) {
	var coupons []model.Coupon
	hasCouponsKey := getHasCouponsKeyByName(userName)
	couponNames, err := data.GetSetMembers(hasCouponsKey)
	if err != nil {
		println("Error when getting coupon members. " + err.Error())
		return nil, err
	}
	// TODO: 使用数组, 不使用slice append
	for _, couponName := range couponNames {
		coupon := GetCoupon(couponName)
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}