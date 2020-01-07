package httptest

import (
	"SecKill/api"
	"github.com/gavv/httpexpect"
	"net/http"
)

var fetchCouponPath = "/api/users/{username}/coupons/{name}"

func fetchDemoCouponSuccess(e *httpexpect.Expect) {
	e.PATCH(fetchCouponPath, demoSellerName, demoCouponName).
		Expect().
		Status(http.StatusCreated).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")
}

func fetchDemoCouponFail(e *httpexpect.Expect) {
	e.PATCH(fetchCouponPath, demoSellerName, demoCouponName).
		Expect().
		Status(http.StatusNoContent).
		Body().Empty()
}