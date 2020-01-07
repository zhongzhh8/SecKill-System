package httptest

import (
	"SecKill/api"
	"SecKill/data"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

const loginPath = "/api/users/"

// 测试登录不存在的用户或错误的密码
func testWrongLogin(e *httpexpect.Expect) {
	wrongUserName := "wrongUserName"
	e.POST(loginPath).
		WithForm(LoginForm{wrongUserName, "whatever_pw"}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueNotEqual(api.ErrMsgKey, "No such queryUser.")

	wrongPassword := "sysucs515"
	e.POST(loginPath).
		WithForm(LoginForm{demoSellerName, wrongPassword}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueNotEqual(api.ErrMsgKey, "Password mismatched.")
}

// 测试登录demo商家和demo顾客
func testUsersLogin(e *httpexpect.Expect) {
	demoSellerLogin(e)
	demoCustomerLogin(e)
}

func isGetCouponUserNotExist(e *httpexpect.Expect, notExistUsername string)  {
	e.GET(getCouponPath, notExistUsername).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual(api.ErrMsgKey, "Record not found.")
}

func testAddAndGetCoupon(e *httpexpect.Expect) {
	// 本函数按顺序做以下测试：
	// 1. 商家未添加优惠券时，查不到优惠券
	// 2. 商家添加优惠券成功
	// 3. 添加优惠券后，用户能查到优惠券，且格式合法

	veryLargePage := 10000
	invalidCustomerName := "someImpCustomer__"  // Some impossible customer

	/* 作为商家和用户分别获取一次优惠券信息 */
	// --顾客查询顾客/商家的优惠券--
	demoCustomerLogin(e)
	// 自己没抢过优惠券，查询不到
	isEmptyBody(e, demoCustomerName, -1)
	isEmptyBody(e, demoCustomerName, 0)
	isEmptyBody(e, demoCustomerName, veryLargePage)
	// 商家没创建过优惠券，查询不到
	isEmptyBody(e, demoSellerName, -1)
	isEmptyBody(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	// 不可查询其它用户
	isGetCouponUnauthorized(e, demoArCustomerName, 0)
	// 查不到不存在的用户
	isGetCouponUserNotExist(e, invalidCustomerName)

	// --商家查询商家的优惠券--
	demoSellerLogin(e)
	isEmptyBody(e, demoSellerName, -1)
	isEmptyBody(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	// 不可查询其它用户
	isGetCouponUnauthorized(e, demoArCustomerName, 0)
	// 查不到不存在的用户
	isGetCouponUserNotExist(e, invalidCustomerName)

	// --创建demo优惠券--
	demoAddCoupon(e)

	// --顾客查询该商家创建的优惠券信息--
	demoCustomerLogin(e)
	isNonEmptyCoupons(e, demoSellerName, -1)
	isNonEmptyCoupons(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	isCustomerSchema(e, demoSellerName, 0)
	// 自己没抢过优惠券，查询不到
	isEmptyBody(e, demoCustomerName, -1)
	isEmptyBody(e, demoCustomerName, 0)
	isEmptyBody(e, demoCustomerName, veryLargePage)
	// 不可查询其它用户
	isGetCouponUnauthorized(e, demoArCustomerName, 0)
	// 查不到不存在的用户
	isGetCouponUserNotExist(e, invalidCustomerName)

	// --商家查询到自己创建的优惠券信息--
	demoSellerLogin(e)
	isNonEmptyCoupons(e, demoSellerName, -1)
	isNonEmptyCoupons(e, demoSellerName, 0)
	isEmptyBody(e, demoSellerName, veryLargePage)
	isSellerSchema(e, demoSellerName, 0)
	// 不可查询其它用户
	isGetCouponUnauthorized(e, demoArCustomerName, 0)
	// 查不到不存在的用户
	isGetCouponUserNotExist(e, invalidCustomerName)
}

func testFetchCoupon(e *httpexpect.Expect, couponAmount int) {
	// demo顾客登录
	demoCustomerLogin(e)

	// 抢一张优惠券
	fetchDemoCouponSuccess(e)

	// 商家优惠券数量-1 顾客可看到该优惠券。
	isCouponExpectedLeft(e,  demoSellerName, 0, 0, couponAmount - 1)
	isNonEmptyCoupons(e, demoSellerName, 0)
	isSellerSchema(e, demoSellerName, 0)
	isNonEmptyCoupons(e, demoCustomerName, 0)
	isCustomerSchema(e, demoCustomerName, 0)

	// 不可重复抢优惠券
	fetchDemoCouponFail(e)
}

// 进行普通的测试，用户注册、登录后进行常规操作
func TestNormal(t *testing.T) {
	_, e := startServer(t)
	defer data.Close()

	// 注册用户,商家
	registerDemoUsers(e)

	// 用户登录错误
	testWrongLogin(e)

	// 用户登录
	testUsersLogin(e)

	// 测试查看、添加优惠券功能
	testAddAndGetCoupon(e)

	// 优惠券已添加，测试抢购优惠券、查看优惠券功能
	testFetchCoupon(e, demoAmount)
}
