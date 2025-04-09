package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
)

type I18nHandler struct {
	localesDir string
}

// NewI18nHandler creates a new i18n handler
func NewI18nHandler() *I18nHandler {
	return &I18nHandler{
		localesDir: "locales", // Directory where translation files are stored
	}
}

// GetTranslations handles the get translations request
// @Summary Get translations for a language
// @Description Get all translation keys for a specific language
// @Tags i18n
// @Accept json
// @Produce json
// @Param language path string true "Language code (e.g., en, vi)"
// @Success 200 {object} utils.Response{data=object} "Success"
// @Failure 404 {object} utils.Response "Language not found"
// @Failure 500 {object} utils.Response "Server error"
// @Router /i18n/{language} [get]
func (h *I18nHandler) GetTranslations(c *gin.Context) {
	language := c.Param("language")
	if language == "" {
		language = "en" // Default language
	}

	// Check if the language file exists
	filePath := filepath.Join(h.localesDir, language+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.NotFoundResponse(c, "language not found")
		return
	}

	// Read the language file
	data, err := os.ReadFile(filePath)
	if err != nil {
		utils.ServerErrorResponse(c)
		return
	}

	// Parse the JSON
	var translations map[string]any
	if err := json.Unmarshal(data, &translations); err != nil {
		utils.ServerErrorResponse(c)
		return
	}

	utils.SuccessResponse(c, translations)
}
