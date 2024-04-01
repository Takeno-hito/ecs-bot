package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"os"
)

var registeredCommands []*discordgo.ApplicationCommand
var session *discordgo.Session

func Init() {
	discordToken := os.Getenv("DISCORD_TOKEN")

	_session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		logrus.Panicf("Invalid bot parameters: %v", err)
	}

	_session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	_session.AddHandler(OnInteractionCreate)

	err = _session.Open()

	if err != nil {
		logrus.Panicf("Cannot open the session: %v", err)
	}

	registeredCommands = make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		registeredCommands[i] = v.Command()
	}
	registeredCommands, err = _session.ApplicationCommandBulkOverwrite(_session.State.User.ID, "", registeredCommands)
	session = _session
}

func Close() {
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", []*discordgo.ApplicationCommand{})
	if err != nil {
		logrus.Panicf("Cannot delete all commands: %v", err)
	}

	if err := session.Close(); err != nil {
		logrus.Errorf("Cannot close Discord connection: %v", err)
	}
}
