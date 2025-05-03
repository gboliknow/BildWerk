package processing

import (
	"errors"
	"net/http"

	"github.com/gboliknow/bildwerk/internal/cache"
	"github.com/gboliknow/bildwerk/internal/middleware"
	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ImageHandler struct {
	imageService *ImageService
	logger       zerolog.Logger
	cache        *cache.RedisCache
}

func NewImageHandler(imageService *ImageService, logger zerolog.Logger, c *cache.RedisCache) *ImageHandler {
	return &ImageHandler{imageService: imageService, logger: logger, cache: c}
}

func (h *ImageHandler) RegisterImageRoutes(r *gin.RouterGroup) {
	// rateLimiter := middleware.RateLimitMiddleware(time.Second, 5, h.cache)

	imageGroup := r.Group("/image")
	imageGroup.Use(middleware.AuthMiddleware())
	{
		imageGroup.POST("/upload", h.HandleUploadImage)
	}

}

// @Summary      Upload an image
// @Description  Upload an image file (JPEG, PNG, WebP) with metadata
// @Tags         Image
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData  file    true  "Image file"
// @Param        filename formData  string  true  "File name"
// @Param        Authorization header string true "Bearer Token"
// @Security     BearerAuth
// @Success      200 {object} models.ImageResponse
// @Failure      400 {object} map[string]string "Invalid form data"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /api/v1/image/upload [post]
func (h *ImageHandler) HandleUploadImage(c *gin.Context) {
	var input UploadImageDTO
	if err := c.ShouldBind(&input); err != nil {
		h.logger.Error().Err(err).Msg("Invalid form data")
		utility.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := c.FormFile("file")
	userID := c.GetString("userID")
	if err != nil {
		h.logger.Error().Err(err).Msg("Missing file in request")
		utility.RespondWithError(c, http.StatusBadRequest, "File is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to open uploaded file")
		utility.RespondWithError(c, http.StatusInternalServerError, "Unable to read file")
		return
	}
	defer file.Close()

	filename := c.PostForm("filename")

	imageResp, appErr := h.imageService.UploadImage(file, filename, userID, fileHeader.Size)
	if appErr != nil {
		h.logger.Error().Err(errors.New(appErr.Message)).Msg("Upload failed")
		utility.RespondWithError(c, appErr.StatusCode, appErr.Message)
		return
	}
	utility.WriteJSON(c.Writer, http.StatusOK, "Image saved", imageResp)

}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE3NDYzNzU3MjgsInVzZXJJRCI6ImNlNDFjOTZjLThjNjAtNGEyYi05YTNkLWRhZDllYzYyYTQzMCJ9.OY1NRbnpvfG0UUfgmzzFENmeG7_s51leKASTdPTeK2w