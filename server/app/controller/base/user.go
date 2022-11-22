package base

import (
	"net/http"
	"time"

	"gitee.com/jiang-xia/gin-zone/server/app/model"
	"gitee.com/jiang-xia/gin-zone/server/app/service"
	"gitee.com/jiang-xia/gin-zone/server/middleware"
	"gitee.com/jiang-xia/gin-zone/server/pkg/response"
	"gitee.com/jiang-xia/gin-zone/server/pkg/tip"
	"github.com/gin-gonic/gin"
	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// User类，go类的写法
type User struct {
	model.User
}

// Register godoc
//
// @Summary     注册接口
// @Description 用户注册接口
// @Tags        用户模块
// @Accept      json
// @Produce     json
// @Param       user body     model.MainUser false "需要上传的json"
// @Success     200  {object} model.MainUser
// @Router      /base/register [post]
func (u *User) Register(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}

	err := service.User.Create(user)
	if err != nil {
		logrus.Error("新增失败", err)
		response.Fail(c, tip.Msg(tip.ErrorInsert), err)
		return
	}

	response.Success(c, user)
}

// Login godoc
//
// @Summary     登录接口
// @Description 用户登录接口
// @Tags        用户模块
// @Accept      json
// @Produce     json
// @Param       user body     model.LoginForm false "需要上传的json"
// @Success     200  {string} token
// @Router      /base/login [post]
func (u *User) Login(c *gin.Context) {
	var login = &model.LoginForm{}
	if err := c.ShouldBind(&login); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	user, errCode := service.User.SignIn(login.UserName, login.Password)

	if errCode == 0 {
		generateToken(c, user)
	} else if errCode == tip.AuthUserNotFound {
		response.Fail(c, tip.Msg(tip.AuthUserNotFound), "")
	} else {
		response.Fail(c, tip.Msg(tip.AuthUserPasswordError), "")
	}
}

// Login godoc
//
// @Summary     用户信息
// @Description 用户信息接口
// @Tags        用户模块
// @Security	Authorization
// @Accept      json
// @Produce     json
// @Success     200 {object} User
// @Router      /base/users/info [get]
func (u *User) UserInfo(c *gin.Context) {
	token := c.GetHeader("authorization")
	user, err := middleware.NewJWT().ParseToken(token)
	if err != nil {
		response.Fail(c, err.Error(), user)
		return
	}
	// user.UserId
	response.Success(c, "用户信息")
}

// UserList godoc
//
// @Summary     用户列表
// @Description 用户列表接口
// @Tags        用户模块
// @Security	Authorization
// @Accept      json
// @Produce     json
// @Param       q   query    string false "username search by q" Format(email)
// @Success     200 {object} User
// @Router      /base/users [get]
func (u *User) UserList(c *gin.Context) {
	c.JSON(http.StatusOK, "用户列表")
}

// DeleteUser godoc
//
// @Summary     删除用户
// @Description 删除用户接口
// @Tags        用户模块
// @Security	Authorization
// @Accept      json
// @Produce     json
// @Param       id  path     int true "用户id" (param name,param type,data type,is mandatory（是否鉴权）?,comment attribute(optional))
// @Success     200 {object} User
// @Router      /base/delete/{id} [delete]
func (u *User) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, "删除用户")
}

// generateToken 生成token
func generateToken(c *gin.Context, user *model.User) {
	data := map[string]interface{}{} // 任意接口
	now := time.Now()

	j := middleware.NewJWT()
	claims := middleware.JWTCustomClaims{
		UserId:   cast.ToString(user.UserId),
		UserName: user.UserName,
		RegisteredClaims: jwtgo.RegisteredClaims{
			IssuedAt:  jwtgo.NewNumericDate(now),                                               // 签发时间
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(12 * time.Hour * time.Duration(1))), // 过期时间12小时
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		response.Fail(c, tip.Msg(tip.AuthFailedGenerateToken), err)
		return
	}

	data["token"] = token

	response.Success(c, data)
}
