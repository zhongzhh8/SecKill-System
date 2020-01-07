package httptest

import (
	"SecKill/data"
	"testing"
)

// 接连测试多个函数
func TestAddCouponWrongCases(t *testing.T) {
	_, e := startServer(t)
	defer data.Close()

	registerDemoUsers(e)

	testAddCouponWrongFormat(e)
	testAddCouponWrongUser(e)
	testAddCouponNotLogIn(e)
	testAddCouponDuplicate(e)
}
