package handler

import (
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type detailsReq struct {
	Name    string `json:"name" binding:"omitempty,max=40"`
	Email   string `json:"email" binding:"omitempty,email"`
	Website string `json:"website" binding:"omitempty,url"`
}

// Details handler
func (h *Handler) Details(c *gin.Context) {
	// 检查上下文 是否 存在 user
	authUser := c.MustGet("user").(*model.User)

	var req detailsReq

	// 检验请求数据
	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		UID:     authUser.UID,
		Name:    req.Name,
		Email:   req.Email,
		Website: req.Website,
	}

	ctx := c.Request.Context()
	if err := h.UserService.UpdateDetails(ctx, u); err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
