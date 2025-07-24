package handlers

import (
	"net/http"
	"time"

	"subscriptions_service/internal/models/dto"
	"subscriptions_service/internal/models/entities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func toResponse(s entities.Subscription) dto.SubscriptionResponse {
	start := s.StartDate.Format("01-2006")
	var end *string
	if s.EndDate != nil {
		e := s.EndDate.Format("01-2006")
		end = &e
	}
	return dto.SubscriptionResponse{
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   start,
		EndDate:     end,
	}
}

func toResponses(list []entities.Subscription) []dto.SubscriptionResponse {
	out := make([]dto.SubscriptionResponse, 0, len(list))
	for _, s := range list {
		out = append(out, toResponse(s))
	}
	return out
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

// CreateSubscription godoc
// @Summary      Создать подписку
// @Description  Добавляет новую запись о подписке
// @Accept       json
// @Produce      json
// @Param        subscription  body      dto.SubscriptionPayload  true  "data"
// @Success      201  {object}  nil
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Router       /subscriptions [post]
func (h *Handler) CreateSubscription(ctx *gin.Context) {
	var in dto.SubscriptionPayload
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	uid, err := uuid.Parse(in.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}
	start, err := parseMonthYear(in.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid start_date format"})
		return
	}
	var endPtr *time.Time
	if in.EndDate != nil {
		end, err := parseMonthYear(*in.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid end_date format"})
			return
		}
		endPtr = &end
	}

	var exists entities.Subscription
	if err := h.db.Where("user_id = ? AND service_name = ?", uid, in.ServiceName).First(&exists).Error; err == nil {
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{Error: "subscription already exists for this user and service"})
		return
	} else if err != gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sub := entities.Subscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     endPtr,
	}

	if err := h.db.Create(&sub).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

// GetSubscription godoc
// @Summary      Получить подписку
// @Param        user_id       path      string  true  "ID пользователя"
// @Param        service_name  path      string  true  "Название сервиса"
// @Success      200  {object}  dto.SubscriptionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /subscriptions/{user_id}/{service_name} [get]
func (h *Handler) GetSubscription(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	serviceName := ctx.Param("service_name")

	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	var sub entities.Subscription
	if err := h.db.Where("user_id = ? AND service_name = ?", uid, serviceName).First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, toResponse(sub))
}

// ListSubscriptions godoc
// @Summary      Список подписок пользователя
// @Param        user_id       query     string  true  "ID пользователя"
// @Success      200  {array}   dto.SubscriptionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /subscriptions [get]
func (h *Handler) ListSubscriptions(ctx *gin.Context) {
	uidStr := ctx.Query("user_id")
	if uidStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "user_id is required"})
		return
	}

	var subs []entities.Subscription
	if err := h.db.Where("user_id = ?", uidStr).Find(&subs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, toResponses(subs))
}

// UpdateSubscription godoc
// @Summary      Обновить подписку
// @Param        user_id       path      string  true  "ID пользователя"
// @Param        service_name  path      string  true  "Название сервиса"
// @Param        subscription  body      dto.SubscriptionUpdatePayload  true  "Обновлённая подписка"
// @Success      200  {object}  dto.SubscriptionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /subscriptions/{user_id}/{service_name} [put]
func (h *Handler) UpdateSubscription(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	serviceName := ctx.Param("service_name")

	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	var sub entities.Subscription
	if err := h.db.Where("user_id = ? AND service_name = ?", uid, serviceName).First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var in dto.SubscriptionUpdatePayload
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	start, err := parseMonthYear(in.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid start_date format"})
		return
	}
	var endPtr *time.Time
	if in.EndDate != nil {
		end, err := parseMonthYear(*in.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid end_date format"})
			return
		}
		endPtr = &end
	}

	sub.Price = in.Price
	sub.StartDate = start
	sub.EndDate = endPtr

	if err := h.db.Save(&sub).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, toResponse(sub))
}

// DeleteSubscription godoc
// @Summary      Удалить подписку
// @Param        user_id       path      string  true  "ID пользователя"
// @Param        service_name  path      string  true  "Название сервиса"
// @Success      204  {object}  nil
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /subscriptions/{user_id}/{service_name} [delete]
func (h *Handler) DeleteSubscription(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	serviceName := ctx.Param("service_name")

	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	if err := h.db.Where("user_id = ? AND service_name = ?", uid, serviceName).
		Delete(&entities.Subscription{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SumSubscriptions godoc
// @Summary      Сумма подписок
// @Description  Считает суммарную стоимость подписок за период
// @Param        user_id       query     string  true   "ID пользователя"
// @Param        from          query     string  true   "Начало периода MM-YYYY"
// @Param        to            query     string  true   "Конец периода MM-YYYY"
// @Param        service_name  query     string  false  "Название сервиса"
// @Success      200  {object}  dto.TotalResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /subscriptions/summary [get]
func (h *Handler) SumSubscriptions(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	serviceName := ctx.Query("service_name")
	if userIDStr == "" || fromStr == "" || toStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "user_id, from and to are required"})
		return
	}

	_, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	fromTime, err := parseMonthYear(fromStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid from format"})
		return
	}
	toTime, err := parseMonthYear(toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid to format"})
		return
	}

	fromMonth := time.Date(fromTime.Year(), fromTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	toMonth := time.Date(toTime.Year(), toTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	toExclusive := toMonth.AddDate(0, 1, 0)

	var subs []entities.Subscription
	query := h.db.Where("user_id = ? AND start_date < ? AND (end_date IS NULL OR end_date >= ?)",
		userIDStr, toExclusive, fromMonth)
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	if err := query.Find(&subs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var total int64
	for _, s := range subs {
		sStart := time.Date(s.StartDate.Year(), s.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		start := maxMonth(sStart, fromMonth)

		var lastInPeriod time.Time
		if s.EndDate != nil {
			e := time.Date(s.EndDate.Year(), s.EndDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			lastInPeriod = minMonth(e, toExclusive.AddDate(0, 0, -1))
		} else {
			lastInPeriod = toExclusive.AddDate(0, 0, -1)
		}
		end := time.Date(lastInPeriod.Year(), lastInPeriod.Month(), 1, 0, 0, 0, 0, time.UTC)

		if end.Before(start) {
			continue
		}

		months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month()) + 1
		total += int64(s.Price * months)
	}

	ctx.JSON(http.StatusOK, dto.TotalResponse{Total: total})
}

func maxMonth(a, b time.Time) time.Time {
	if a.Before(b) {
		return b
	}
	return a
}

func minMonth(a, b time.Time) time.Time {
	if a.After(b) {
		return b
	}
	return a
}
