package bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Bilingual(Ja string, En string, l discordgo.Locale) string {
	if Ja == "" {
		logrus.Warn("Bilingual: Ja is undefined")
		return En
	}
	if En == "" {
		logrus.Warn("Bilingual: En is undefined")
		return En
	}
	if l == discordgo.Japanese {
		return Ja
	}
	return En
}

type CommandDescription struct {
	descriptionEn string
	descriptionJa string
	options       []*discordgo.ApplicationCommandOption
}

type ApplicationCommand struct {
	Name        string
	description *CommandDescription
	handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type ApplicationCommandList map[string]ApplicationCommand

func (c *ApplicationCommand) Command() *discordgo.ApplicationCommand {
	if c.description.descriptionEn == "" {
		return &discordgo.ApplicationCommand{
			Name:        c.Name,
			Description: c.description.descriptionJa,
			Options:     c.description.options,
		}
	}
	if c.description.descriptionJa == "" {
		return &discordgo.ApplicationCommand{
			Name:        c.Name,
			Description: c.description.descriptionEn,
			Options:     c.description.options,
		}
	}

	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.description.descriptionEn,
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.Japanese: c.description.descriptionJa,
		},
		Options: c.description.options,
	}
}

func searchCommand(l []*ApplicationCommand, name string) (command *ApplicationCommand, ok bool) {
	for _, c := range l {
		if c.Name == name {
			return c, true
		}
	}
	return nil, false
}

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		interactionId := i.ApplicationCommandData().Name
		if c, ok := searchCommand(Commands, interactionId); ok {
			err := c.handler(s, i)
			if err != nil {
				_ = replyError(s, i, err)
			}
		} else {
			logrus.Errorf("undefined command: %v", interactionId)
			_ = replyError(s, i, ErrUnknownCommand)
		}
		return
	case discordgo.InteractionModalSubmit:
		interactionId := i.ModalSubmitData().CustomID
		if h, ok := MessageActions[interactionId]; ok {
			err := h(s, i)
			if err != nil {
				_ = replyError(s, i, err)
			}
		} else {
			logrus.Errorf("undefined command: %v", interactionId)
			_ = replyError(s, i, ErrUnknownCommand)
		}
		return

	case discordgo.InteractionMessageComponent:
		interactionId := i.MessageComponentData().CustomID
		if h, ok := MessageActions[interactionId]; ok {
			err := h(s, i)
			if err != nil {
				_ = replyError(s, i, err)
			}
		} else {
			logrus.Errorf("undefined command: %v", interactionId)
			_ = replyError(s, i, ErrUnknownCommand)
		}
		return
	default:
		logrus.Errorf("Unknown Type: %v", i.Type)
		if err := replyEphemeral(s, i, "Sorry, but the command is undefined"); err != nil {
			logrus.Error(err)
		}
		return
	}
}

func SendMessage(s *discordgo.Session, channelID string, msg string) error {
	_, err := s.ChannelMessageSend(channelID, msg)
	logrus.Debug("Message sent: ", msg)
	return err
}

func SendMessageComplex(s *discordgo.Session, channelID string, msg *discordgo.MessageSend) error {
	_, err := s.ChannelMessageSendComplex(channelID, msg)
	logrus.Debug("Message sent: ", msg.Content)
	return err
}

func replyError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) error {
	switch {
	case errors.Is(err, ErrUnknownCommand):
		return replyEphemeral(s, i, Bilingual(
			"コマンドが定義されていません。",
			"Command is undefined.",
			i.Locale,
		))
	case errors.Is(err, ErrUnderConstruction):
		return replyEphemeral(s, i, Bilingual(
			"現在このコマンドは使用できません… :(",
			"Sorry, but this command is now under construction... :(",
			i.Locale,
		))
	case errors.Is(err, ErrNeedRegister):
		return replyEphemeral(s, i, Bilingual(
			// TODO: /register のリンクを作る
			"登録が必要です！ /register から登録してください。",
			"You need to register your summoner profile! Please type /register.",
			i.Locale,
		))
	default:
		logrus.Errorf("Internal Error: %v", err)
		return replyEphemeral(s, i, "サーバーエラーが発生しました。管理者に連絡してください。")
	}
}

func Session() *discordgo.Session {
	return session
}

func makePointer[T any](t T) *T {
	return &t
}
