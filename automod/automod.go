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

	tokens := strings.Split(m.Content, " ")
	for i := 0; i < len(words); i++ {
		for j := 0; j < len(tokens); j++ {
			if strings.EqualFold(words[i], tokens[j]) {
				return true
			}
		}
	}
	return false
}

func IsWordOnTimer(m *discordgo.Message) bool {
	var word = "ian"
	tokens := strings.Split(m.Content, " ")
	for i := 0; i < len(tokens); i++ {
		if strings.EqualFold(tokens[i], word) {
			return true
		}
	}
	return false
}
