package controllers

import (
	"github.com/gin-gonic/gin"
)

type IAppUtils interface {
	GinContext() *gin.Context

	GetUserID() (int, bool)
	GetRoleID() (int, bool)

	BindJSON(obj interface{}) error
	AbortWithError(code int, err error)
	JSON(code int, obj interface{})
}

type AppUtils struct {
	c *gin.Context
}

func NewAppUtils(c *gin.Context) *AppUtils {
	return &AppUtils{c: c}
}

func (u *AppUtils) GinContext() *gin.Context {
	return u.c
}

func (u *AppUtils) GetUserID() (int, bool) {
	userID, exists := u.c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int)
	return id, ok
}

func (u *AppUtils) GetRoleID() (int64, bool) {
	roleID, exists := u.c.Get("role_id")
	if !exists {
		return 0, false
	}
	id, ok := roleID.(int64)
	return id, ok
}

func (u *AppUtils) BindJSON(obj interface{}) error {
	return u.c.ShouldBindJSON(obj)
}

func (u *AppUtils) AbortWithError(code int, err error) {
	u.c.AbortWithStatusJSON(code, err)
}

func (u *AppUtils) JSON(code int, obj interface{}) {
	u.c.JSON(code, obj)
}
