package http

import (
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/errors"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

// Auth handlers
type handlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	log    *logger.Logger
}

// Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log *logger.Logger) auth.Handlers {
	return &handlers{cfg, authUC, log}
}

// Crate new user
func (h *handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create user", zap.String("ReqID", utils.GetRequestID(c)))

		var user models.User
		if err := c.Bind(&user); err != nil {
			h.log.Error("Create c.Bind", zap.String("ReqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		createdUser, err := h.authUC.Create(ctx, &user)
		if err != nil {
			h.log.Error("auth repo create", zap.String("reqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info("Created user", zap.String("reqID", utils.GetRequestID(c)), zap.String("ID", createdUser.ID.String()))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Update existing user
func (h *handlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update user", zap.String("ReqID", utils.GetRequestID(c)))

		var user models.UserUpdate
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			h.log.Error("Update uuid.Parse", zap.String("ReqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}
		user.ID = uID

		if err := c.Bind(&user); err != nil {
			h.log.Error("Update c.Bind", zap.String("ReqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		updatedUser, err := h.authUC.Update(ctx, &user)
		if err != nil {
			h.log.Error("auth repo update", zap.String("reqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info("Update user", zap.String("reqID", utils.GetRequestID(c)), zap.String("ID", updatedUser.ID.String()))

		return c.JSON(http.StatusCreated, updatedUser)
	}
}

// Get user by id
func (h *handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update user", zap.String("ReqID", utils.GetRequestID(c)))

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			h.log.Error("Update uuid.Parse", zap.String("ReqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		user, err := h.authUC.GetByID(ctx, uID)
		if err != nil {
			h.log.Error("auth repo get by id", zap.String("reqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Delete user handler
func (h *handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update user", zap.String("ReqID", utils.GetRequestID(c)))

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			h.log.Error("Update uuid.Parse", zap.String("ReqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		if err := h.authUC.Delete(ctx, uID); err != nil {
			h.log.Error("auth repo delete", zap.String("reqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Find users by name
func (h *handlers) FindByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info(
			"FindByName",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.String("name", c.QueryParam("name")),
		)

		if c.QueryParam("name") == "" {
			return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("name query param is required"))
		}

		users, err := h.authUC.FindByName(ctx, c.QueryParam("name"))
		if err != nil {
			h.log.Error(
				"auth repo find by name",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info("FindByName", zap.String("ReqID", utils.GetRequestID(c)), zap.Int("Found", len(users)))

		return c.JSON(http.StatusOK, users)
	}
}

// Gat all users with pagination page and size query params
func (h *handlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create user", zap.String("ReqID", utils.GetRequestID(c)))

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			h.log.Error(
				"GetPaginationFromCtx",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		usersList, err := h.authUC.GetUsers(ctx, paginationQuery)
		if err != nil {
			h.log.Error("GetUsers", zap.String("reqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"GetUsers",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.Int("Found", len(usersList.Users)),
			zap.String("Query", fmt.Sprintf("%#v", paginationQuery)),
		)

		return c.JSON(http.StatusOK, usersList)
	}
}
