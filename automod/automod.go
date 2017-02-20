package automod

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func IsWordCensored(m *discordgo.Message) bool {
	//this will check through a preloaded map eventually
	var words [3]string
	words[0] = "dudu"
	words[1] = "brained"
	words[2] = "dorf"

	for i := 0; i < len(words); i++ {
		if strings.Contains(m.Content, words[i]) {
			return true
		}
	}
	return false
}
