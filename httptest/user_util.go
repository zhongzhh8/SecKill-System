package httptest

import (
	"SecKill/api"
	"SecKill/model"
	"github.com/gavv/httpexpect"
	"net/http"
)

// 本文件存放了一些demo用户信息，demo用户的注册/登录函数
// 定义了用户登出函数
// 还定义了注册/登录用户的表格

const demoSellerName = "kiana"
const demoCustomerName = "jinzili"
const demoArCustomerName = "karsa"  // name of another customer
const demoPassword = "shen6508"

type RegisterForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Kind     string  `form:"kind"`
}

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func registerDemoUsers(e *httpexpect.Expect)  {
	e.POST("/api/users/").
		WithForm(RegisterForm{demoSellerName, demoPassword, model.NormalSeller}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

	e.POST("/api/users/").
		WithForm(RegisterForm{demoCustomerName, demoPassword, model.NormalCustomer}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")


	e.POST("/api/users/").
		WithForm(RegisterForm{demoArCustomerName, demoPassword, model.NormalCustomer}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "")

}


func demoCustomerLogin(e *httpexpect.Expect)  {
	e.POST("/api/auth/").
		WithForm(LoginForm{demoCustomerName, demoPassword}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "").
		ValueEqual("kind", model.NormalCustomer)
}

func demoArCustomerLogin(e *httpexpect.Expect)  {
	e.POST("/api/auth/").
		WithForm(LoginForm{demoArCustomerName, demoPassword}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "").
		ValueEqual("kind", model.NormalCustomer)
}

func demoSellerLogin(e *httpexpect.Expect)  {
	e.POST("/api/auth/").
		WithForm(LoginForm{demoSellerName, demoPassword}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "").
		ValueEqual("kind", model.NormalSeller)
}

func logout(e *httpexpect.Expect)  {
	e.POST("/api/auth/logout").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual(api.ErrMsgKey, "log out.")
}
