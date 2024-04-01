package bot

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*ApplicationCommand{
	{
		Name: "dice",
		description: &CommandDescription{
			descriptionJa: "ダイスを振ります",
			options:       nil,
		},
		handler: DiceHandler,
	},
}

var MessageActions = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{}

func replyEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, m string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
