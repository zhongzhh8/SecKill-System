package model

// 数据库实体
type Coupon struct {
	Id          int64     `gorm:"primary_key;auto_increment"`
	Username    string    `gorm:"type:varchar(20); not null"` // 用户名
	CouponName  string    `gorm:"type:varchar(60); not null"` // 优惠券名称
	Amount      int64     						              // 最大优惠券数
	Left        int64								          // 剩余优惠券数
	Stock       int64                                       // 面额
	Description string    `gorm:"type:varchar(60)"`           // 优惠券描述信息
}

type ReqCoupon struct {
	Name			string
	Amount 			int64
	Description     string
	Stock           int64
}

type ResCoupon struct {
	Name            string  `json:"name"`
	Stock           int64   `json:"stock"`
	Description     string  `json:"description"`
}

// 商家查询优惠券时，返回的数据结构
type SellerResCoupon struct {
	ResCoupon
	Amount int64  `json:"amount"`
	Left   int64  `json:"left"`
}

// 顾客查询优惠券时，返回的数据结构
type CustomerResCoupon struct {
	ResCoupon
}

func ParseSellerResCoupons(coupons []Coupon) []SellerResCoupon {
	var sellerCoupons []SellerResCoupon
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			SellerResCoupon{ResCoupon{coupon.CouponName, coupon.Stock, coupon.Description},
				coupon.Amount, coupon.Left})
	}
	return sellerCoupons
}

func ParseCustomerResCoupons(coupons []Coupon) []CustomerResCoupon {
	var sellerCoupons []CustomerResCoupon
	for _, coupon := range coupons {
		sellerCoupons = append(sellerCoupons,
			CustomerResCoupon{ResCoupon{coupon.CouponName, coupon.Stock, coupon.Description}})
	}
	return sellerCoupons
}