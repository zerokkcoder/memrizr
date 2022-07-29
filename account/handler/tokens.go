package handler

import (
	"log"
	"memrizr/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokensReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Tokens handler
func (h *Handler) Tokens(c *gin.Context) {
	var req tokensReq

	// 检验请求数据
	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	// 验证 refresh JWT
	refreshToken, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// 获取 用户
	u, err := h.UserService.Get(ctx, refreshToken.UID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// 创建 fresh 的token对
	tokens, err := h.TokenService.NewTokenPairFromUser(ctx, u, refreshToken.ID.String())
	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", u, err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
