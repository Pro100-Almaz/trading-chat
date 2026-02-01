package controller

import (
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"
)

type EmojiController struct{}

// GetEmojis godoc
// @Summary Get available avatar emojis
// @Description Returns a list of all available emojis that can be used as user avatars
// @Tags Emojis
// @Produce json
// @Success 200 {object} domain.EmojiListResponse "List of available emojis"
// @Router /emojis [get]
func (ec *EmojiController) GetEmojis(w http.ResponseWriter, r *http.Request) {
	emojis := domain.GetEmojiList()
	utils.JSON(w, http.StatusOK, emojis)
}
