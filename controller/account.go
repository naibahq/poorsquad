package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"

	"github.com/naiba/poorsquad/model"
)

// AccountController ..
type AccountController struct {
}

// ServeAccount ..
func ServeAccount(r gin.IRoutes) {
	ac := AccountController{}
	r.POST("/account", ac.addOrEdit)
}

type accountForm struct {
	CompanyID uint64 `binding:"required" json:"company_id,omitempty"`
	Token     string `binding:"required" json:"token,omitempty"`
}

func (ac *AccountController) addOrEdit(c *gin.Context) {
	var af accountForm
	if err := c.ShouldBindJSON(&af); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("格式错误：%s", err),
		})
		return
	}
	u := c.MustGet(model.CtxKeyAuthorizedUser).(*model.User)

	var uc model.UserCompany
	if err := db.Where("user_id = ? AND company_id = ?", u.ID, af.CompanyID).First(&uc).Error; err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("您不是该企业的雇员：%s", err),
		})
		return
	}

	if uc.Permission < model.UCPManager {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "您不是该企业的管理人员",
		})
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: af.Token,
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	gu, _, err := client.Users.Get(ctx, "")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("验证Token失败：%s", err),
		})
		return
	}
	a := model.NewAccountFromGitHub(gu)
	a.Token = af.Token
	a.CompanyID = af.CompanyID
	if err := db.Save(&a).Error; err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("数据库错误：%s", err),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: a,
	})
}
