package controller

import (
	"net/http"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/usecase"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/system"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/message"
	coreClaims "github.com/cool9850311/StreamPlatformLite-Core/pkg/claims"
	"github.com/gin-gonic/gin"
)

type SystemSettingController struct {
	Log                  logger.Logger
	systemSettingUseCase *usecase.SystemSettingUseCase
}

func NewSystemSettingController(log logger.Logger, systemSettingUseCase *usecase.SystemSettingUseCase) *SystemSettingController {
	return &SystemSettingController{
		Log:                  log,
		systemSettingUseCase: systemSettingUseCase,
	}
}

func (c *SystemSettingController) GetSetting(ctx *gin.Context) {
	claims := ctx.Request.Context().Value("claims").(*coreClaims.Claims)
	setting, err := c.systemSettingUseCase.GetSetting(ctx, claims.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": message.MsgInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, setting)
}

func (c *SystemSettingController) SetSetting(ctx *gin.Context) {
	claims := ctx.Request.Context().Value("claims").(*coreClaims.Claims)
	setting := &system.Setting{}
	if err := ctx.BindJSON(setting); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": message.MsgBadRequest})
		return
	}
	err := c.systemSettingUseCase.SetSetting(ctx, setting, claims.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": message.MsgInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": message.MsgOK})
}
