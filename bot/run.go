package bot

import (
	"embed"
	"encoding/csv"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	LeaderRoleId      = "1196348597340491849"
	ParticipantRoleId = "1196348556190171156"
	CoachRoleId       = "1196348690936385606"
	StaffRoleId       = "1196348936106016869"
	GuildId           = "1196348379161174076"
	PlayoffRoleId     = "1214792826592956476"
	NoRoleRoleId      = "1207204509055721522"
)

//go:embed data.csv
var dataCsv embed.FS

func Run() {
	data, err := dataCsv.Open("data.csv")
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(data)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	setRoleToMembers(records)

	setNoRoleToMembers()

	//generateTeams(records)
}

func setNoRoleToMembers() {
	members, err := session.GuildMembers(GuildId, "", 100)
	if err != nil {
		logrus.Errorf("cannot get members: %v", err)
		return
	}

	for _, v := range members {
		if v.Roles == nil || len(v.Roles) == 0 {
			err := session.GuildMemberRoleAdd(GuildId, v.User.ID, NoRoleRoleId)
			if err != nil {
				logrus.Errorf("cannot add role: %v, %v", err, v.User.Username)
			} else {
				logrus.Infof("added role: %v", v.User.Username)
			}
		}
	}
}

func setRoleToMembers(data [][]string) {
	roles, err := session.GuildRoles(GuildId)
	if err != nil {
		logrus.Errorf("cannot get roles: %v", err)
		return
	}
	members, err := session.GuildMembers(GuildId, "", 1000)
	if err != nil {
		logrus.Errorf("cannot get members: %v", err)
		return
	}
	membersMap := make(map[string]*discordgo.Member)

	for _, v := range members {
		membersMap[v.User.String()] = v
	}

	for _, v := range data[1:] {
		if len(v) != 8 {
			logrus.Errorf("invalid data: %v", v)
			continue
		}
		teamRole := v[2]
		teamFullName := v[0] + " - " + v[1]
		displayName := v[4]
		discordUserName := v[6]

		member, ok := membersMap[discordUserName]
		if !ok {
			logrus.Warnf("cannot find user: %v", discordUserName)
			continue
		}

		//var hasRole bool
		//if member.Roles != nil {
		//	for _, r := range member.Roles {
		//		if r == LeaderRoleId || r == ParticipantRoleId || r == CoachRoleId {
		//			logrus.Debugf("already has role: %v, skipping...", member.User.Username)
		//			hasRole = true
		//			break
		//		}
		//	}
		//}
		//if hasRole {
		//	continue
		//}

		var role *discordgo.Role

		for _, r := range roles {
			if r.Name == teamFullName {
				role = r
				break
			}
		}

		if role == nil {
			logrus.Errorf("cannot find role: %v", teamFullName)
			return
		}

		err = updateMember(role.ID, member, teamRole, displayName)
		if err != nil {
			logrus.Errorf("cannot update member: %v, %v", err, member.User.Username)
			return
		} else {
			logrus.Infof("updated member: %v", member.User.Username)
		}
	}
}

func updateMember(teamRoleId string, member *discordgo.Member, role string, displayName string) error {
	roles := make([]string, 0)

	roles = append(roles, teamRoleId)

	roles = append(roles, PlayoffRoleId)

	switch role {
	case "Leader":
		roles = append(roles, LeaderRoleId)
		roles = append(roles, ParticipantRoleId)
	case "Player&Coach":
		roles = append(roles, ParticipantRoleId)
		roles = append(roles, CoachRoleId)
	case "Coach":
		roles = append(roles, CoachRoleId)
	case "Player":
		roles = append(roles, ParticipantRoleId)
	default:
		return errors.New("unknown role: " + role)
	}

	params := &discordgo.GuildMemberParams{
		Nick:                       displayName,
		Roles:                      &roles,
		ChannelID:                  nil,
		Mute:                       nil,
		Deaf:                       nil,
		CommunicationDisabledUntil: nil,
	}
	_, err := session.GuildMemberEdit(GuildId, member.User.ID, params)
	if err != nil {
		logrus.Errorf("cannot update member: %v, %v", err, params)
	}
	return err
}

func generateTeams(data [][]string) {
	roles, err := session.GuildRoles(GuildId)
	if err != nil {
		logrus.Errorf("cannot get roles: %v", err)
		return
	}

	for _, v := range data {
		if len(v) != 2 {
			logrus.Errorf("invalid data: %v", v)
			continue
		}
		teamFullName := v[0] + " - " + v[1]

		var role *discordgo.Role
		for _, r := range roles {
			if r.Name == teamFullName {
				role = r
				break
			}
		}

		if role == nil {
			logrus.Warnf("cannot find role, creating: %v", teamFullName)

			roleParam := &discordgo.RoleParams{
				Name:        teamFullName,
				Color:       nil,
				Hoist:       nil,
				Permissions: nil,
				Mentionable: nil,
			}
			role, err = session.GuildRoleCreate(GuildId, roleParam)
			if err != nil {
				logrus.Errorf("cannot create role: %v", err)

				continue
			}
		}

		_, err = session.GuildChannelCreateComplex(GuildId, discordgo.GuildChannelCreateData{
			Name:             teamFullName,
			Type:             discordgo.ChannelTypeGuildVoice,
			Topic:            "",
			Bitrate:          0,
			UserLimit:        0,
			RateLimitPerUser: 0,
			Position:         0,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    role.ID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel,
				},
				{
					ID:    StaffRoleId,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel,
				},
				{
					ID:   GuildId,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
			},
			ParentID: "1196350726985744394",
			NSFW:     false,
		})

		if err != nil {
			logrus.Errorf("cannot create channel: %v", err)
		}

		_, err = session.GuildChannelCreateComplex(GuildId, discordgo.GuildChannelCreateData{
			Name:             "連絡-" + v[0],
			Type:             discordgo.ChannelTypeGuildText,
			Topic:            "",
			Bitrate:          0,
			UserLimit:        0,
			RateLimitPerUser: 0,
			Position:         0,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    role.ID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel,
				},
				{
					ID:    StaffRoleId,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel,
				},
				{
					ID:   GuildId,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
			},
			ParentID: "1208040837456334878",
			NSFW:     false,
		})
	}
}
