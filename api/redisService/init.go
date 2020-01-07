package redisService

import (
	"SecKill/api/dbService"
	"SecKill/data"
)

const secKillScript = `
    -- Check if User has coupon
    -- KEYS[1]: hasCouponKey "{username}-has"
    -- KEYS[2]: couponName   "{couponName}"
    -- KEYS[3]: couponKey    "{couponName}-info"
    -- 返回值有-1, -2, -3, 都代表抢购失败
    -- 返回值为1代表抢购成功

    -- Check if coupon exists and is cached
	local couponLeft = redis.call("hget", KEYS[3], "left");
	if (couponLeft == false)
	then
		return -2;  -- No such coupon
	end
	if (tonumber(couponLeft) == 0)  --- couponLeft是字符串类型
    then
		return -3;  --  No Coupon Left.
	end

    -- Check if the user has got the coupon --
	local userHasCoupon = redis.call("SISMEMBER", KEYS[1], KEYS[2]);
	if (userHasCoupon == 1)
	then
		return -1;
	end

    -- User gets the coupon --
	redis.call("hset", KEYS[3], "left", couponLeft - 1);
	redis.call("SADD", KEYS[1], KEYS[2]);
	return 1;
`
var secKillSHA string  // SHA expression of secKillScript

// 将数据加载到缓存预热，防止缓存穿透
// 预热加载了商品库存key
func preHeatKeys()  {
	coupons, err := dbService.GetAllCoupons()
	if err != nil {
		panic("Error when getting all coupons." + err.Error())
	}

	for _, coupon := range coupons {
		err := CacheCouponAndHasCoupon(coupon)
		if err != nil {
			panic("Error while setting redis keys of coupons. " + err.Error())
		}
	}
	println("---Set redis keys of coupons success.---")
}

func init() {
	// 让redis加载秒杀的lua脚本
	secKillSHA = data.PrepareScript(secKillScript)

	// 预热
	preHeatKeys()
}
