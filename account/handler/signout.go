package handler

import (
	"memrizr/model"
	"memrizr/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Signout(c *gin.Context) {
	// 检查上下文 是否 存在 user
	user := c.MustGet("user")

	ctx := c.Request.Context()

	if err := h.TokenService.Signout(ctx, user.(*model.User).UID); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully!",
	})
}
