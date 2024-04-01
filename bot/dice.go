package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

var flag = false

func DiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var discordName string

	if i.Member.Nick == "" {
		discordName = i.Member.User.Username
	} else {
		discordName = i.Member.Nick
	}

	color := 0x4cd8b9 // blue
	if flag {
		color = 0xfa4454 // red
	}
	flag = !flag

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s さんは %d を出しました！", discordName, rand.Intn(100)+1),
		Color: color,
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
