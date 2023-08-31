package handler

import (
	"httpserver/database"
	"httpserver/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *database.Database
}

func New(db *database.Database) (*Handler, error) {
	return &Handler{
		db: db,
	}, nil
}

// Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

func (h *Handler) GetUserSegments(c *gin.Context) {

	ctx := c.Request.Context()

	userID := c.Param("id")

	segments, err := h.db.GetUserSegments(ctx, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, segments)
}

// Метод создания сегмента. Принимает slug (название) сегмента.

func (h *Handler) CreateSegment(c *gin.Context) {

	ctx := c.Request.Context()

	var seg model.SegName

	if err := c.ShouldBindJSON(&seg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !isValidSegmentName(seg.SegName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment name"})
		return
	}

	err := h.db.CreateSegment(ctx, seg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment created successfully"})
}

// Метод удаления сегмента. Принимает slug (название) сегмента.

func (h *Handler) DeleteSegment(c *gin.Context) {
	ctx := c.Request.Context()

	var seg model.SegName
	if err := c.ShouldBindJSON(&seg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !isValidSegmentName(seg.SegName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment name"})
		return
	}

	err := h.db.DeleteSegment(ctx, seg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully"})
}

// Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю,
// список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.

func (h *Handler) UpdateUserSegments(c *gin.Context) {
	ctx := c.Request.Context()

	var us model.UserSegments
	if err := c.ShouldBindJSON(&us); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if len(us.SegmentsToAdd) == 0 && len(us.SegmentsToRemove) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No values to update"})
		return
	}

	err := h.db.UpdateUserSegments(ctx, us)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User segments updated successfully"})
}
