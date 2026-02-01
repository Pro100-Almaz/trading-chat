package route

import (
	"github.com/gorilla/mux"
	"github.com/Pro100-Almaz/trading-chat/api/controller"
)

func NewEmojiRouter(r *mux.Router) {
	ec := &controller.EmojiController{}

	r.HandleFunc("/emojis", ec.GetEmojis).Methods("GET")
}
