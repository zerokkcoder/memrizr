package handler

import (
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Me 获取用户个人详情
func (h *Handler) Me(c *gin.Context) {

	// 检查上下文 是否 存在 user
	user, exists := c.Get("user")

	if !exists {
		log.Printf("Unable to exteact user from request context for unkonw reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	// 获取用户 id
	uid := user.(*model.User).UID

	// 获取用户
	u, err := h.UserService.Get(c, uid)
	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
