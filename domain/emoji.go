package domain

// Available avatar emojis - index corresponds to the avatar_emoji field in User
var AvatarEmojis = []string{
	"😀", // 0
	"😎", // 1
	"🤖", // 2
	"👨‍💻", // 3
	"👩‍💻", // 4
	"🦊", // 5
	"🐱", // 6
	"🐶", // 7
	"🦁", // 8
	"🐼", // 9
	"🦄", // 10
	"🐲", // 11
	"🦅", // 12
	"🐬", // 13
	"🦋", // 14
	"🌟", // 15
	"🔥", // 16
	"💎", // 17
	"🚀", // 18
	"⚡", // 19
}

// EmojiResponse represents a single emoji option
type EmojiResponse struct {
	Index int    `json:"index" example:"0"`
	Emoji string `json:"emoji" example:"😀"`
}

// EmojiListResponse represents the list of available emojis
type EmojiListResponse struct {
	Emojis []EmojiResponse `json:"emojis"`
}

// GetEmojiList returns all available avatar emojis
func GetEmojiList() EmojiListResponse {
	var emojis []EmojiResponse
	for i, emoji := range AvatarEmojis {
		emojis = append(emojis, EmojiResponse{
			Index: i,
			Emoji: emoji,
		})
	}
	return EmojiListResponse{Emojis: emojis}
}

// IsValidEmojiIndex checks if the given index is valid
func IsValidEmojiIndex(index int) bool {
	return index >= 0 && index < len(AvatarEmojis)
}
