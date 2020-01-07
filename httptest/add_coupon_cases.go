package httptest

import (
	"SecKill/api"
	"github.com/gavv/httpexpect"
	"net/http"
	"strconv"
)

/*
该文件下依赖于注册过的demo用户，需要先调用registerDemoUsers
该文件定义了添加优惠券的各种函数
*/

/* 定义了添加优惠券的表格，函数等等 */
const addCouponPath = "/api/users/{username}/coupons"
type AddCouponForm struct {
	Name        string `form:"name"`
	Amount      string `form:"amount"`      // 应当int
	Description string `form:"description"`
	Stock       string `form:"stock"`       // 应当int
}

// 定义了demo优惠券
var demoCouponName = "my_coupon"
var demoAmount     = 10
var demoStock      = 50
var demoAddCouponForm AddCouponForm = AddCouponForm{
	Name:        demoCouponName,
	Amount:      strconv.Itoa(demoAmount) ,
	Stock:       strconv.Itoa(demoStock),
	Description: "kiana: this is my good coupon",
}

// 测试添加优惠券时的表格格式
func testAddCouponWrongFormat(e *httpexpect.Expect) {
	// 登录商家
	logout(e)
	demoSellerLogin(e)

	// amount值不是数字
	amountNotNumberForm := demoAddCouponForm
	amountNotNumberForm.Amount = "blah-blah"
	e.POST(addCouponPath, demoSellerName).
		WithForm(amountNotNumberForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Amount field wrong format.")

	// stock值不是数字
	stockNotNumberForm := demoAddCouponForm
	stockNotNumberForm.Stock = "blah-blah"
	e.POST(addCouponPath, demoSellerName).
		WithForm(stockNotNumberForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Stock field wrong format.")

	// 优惠券名为空
	emptyCouponNameForm := demoAddCouponForm
	emptyCouponNameForm.Name = ""
	e.POST(addCouponPath, demoSellerName).
		WithForm(emptyCouponNameForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Coupon name or description should not be empty.")

	// 优惠券描述为空
	emptyDescriptionForm := demoAddCouponForm
	emptyDescriptionForm.Description = ""
	e.POST(addCouponPath, demoSellerName).
		WithForm(emptyDescriptionForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Coupon name or description should not be empty.")
}

// 测试非商家添加优惠券或为其它用户添加优惠券
func testAddCouponWrongUser(e *httpexpect.Expect) {
	// 登录顾客
	demoCustomerLogin(e)
	// 顾客不可添加优惠券
	e.POST(addCouponPath, demoCustomerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Only sellers can create coupons.")

	// 登录商家
	demoSellerLogin(e)
	// 不可为其它用户添加优惠券
	e.POST(addCouponPath, demoCustomerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Cannot create coupons for other users.")
}

// 测试未登录添加优惠券
func testAddCouponNotLogIn(e * httpexpect.Expect)  {
	logout(e)

	e.POST(addCouponPath, demoSellerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Not Logged in.")
}

func testAddCouponDuplicate(e *httpexpect.Expect)  {
	demoSellerLogin(e)

	e.POST(addCouponPath, demoSellerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

	// 添加重复优惠券失败
	e.POST(addCouponPath, demoSellerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Create failed. Maybe (username,coupon name) duplicates")
}

// 添加demo优惠券(事先不得添加过)
func demoAddCoupon(e *httpexpect.Expect)  {
	demoSellerLogin(e)

	e.POST(addCouponPath, demoSellerName).
		WithForm(demoAddCouponForm).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}
