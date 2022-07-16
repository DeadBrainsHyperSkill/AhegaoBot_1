package discordmarket

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	dayAsSec        int64 = 86400
	hoursAsSec      int64 = 3600
	DefaultPanel          = 0
	ProfilePanel          = 1
	MarketPanel           = 2
	ReactionPanel         = 3
	CreateRolePanel       = 4
)

type Bot struct {
	Cfg            *Config
	MarketData     Database
	LotterySession *Lottery
}

type Lottery struct {
	Participants []string
	Expires      time.Time
}

type CreateRoleUser struct {
	UserId    string
	RoleColor int
	Duration  int64
}

var userIdMsgIdDict map[string]string = make(map[string]string)

var panelsCalledUsers = []string{}

var createRoleUsers = []CreateRoleUser{}

var setStatusUsers = []string{}

func openMainPanel(s *discordgo.Session, i *discordgo.InteractionCreate, t discordgo.InteractionResponseType) {

	openCustomPanel(s, i, t, &discordgo.MessageEmbed{
		Title:       "Панель",
		Description: "> **Тип:**⠀`Пользовательская`\n\n> **Версия:**⠀`1.0`",
		Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/838477959131430912/AhegaoWelcome.gif"},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858719063316496394/gear.gif"},
	},
		[]*discordgo.Button{{
			Label:    "⠀⠀⠀⠀⠀Профиль⠀⠀⠀⠀⠀⠀",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			CustomID: "Профиль",
			Emoji:    discordgo.ComponentEmoji{Name: "📁"},
		}, {
			Label:    "⠀⠀⠀⠀⠀⠀Маркет⠀⠀⠀⠀⠀⠀",
			Style:    discordgo.SuccessButton,
			Disabled: false,
			CustomID: "Маркет",
			Emoji:    discordgo.ComponentEmoji{Name: "🛒"},
		}}, DefaultPanel)
}

func openCustomPanel(s *discordgo.Session, i *discordgo.InteractionCreate, t discordgo.InteractionResponseType, embed *discordgo.MessageEmbed, btns []*discordgo.Button, panelType int) {

	var row discordgo.ActionsRow
	var components []discordgo.MessageComponent

	if !strings.Contains(embed.Title, "Панель") {
		row.Components = append(row.Components, discordgo.Button{
			Label:    "⠀⠀⠀<<⠀⠀⠀⠀",
			Style:    discordgo.SecondaryButton,
			Disabled: false,
			CustomID: "<<",
			Emoji:    discordgo.ComponentEmoji{Name: "⚙️"},
		})
	}

	j := 0
	for i := 2; i < len(btns)+2; i++ {
		if panelType == CreateRolePanel && i == 2 {
			components = append(components, row)
			row.Components = nil
		}
		row.Components = append(row.Components, btns[j])
		j++
		if panelType == CreateRolePanel && (i == 6 || i == 11 || i == 16 || i == 17) {
			components = append(components, row)
			row.Components = nil
			continue
		}
		if panelType == ProfilePanel && (i == 2 || i == 5 || i == 7) {
			components = append(components, row)
			row.Components = nil
			continue
		}
		if panelType == ReactionPanel && (i == 4 || i == 9 || i == 14 || i == 19 || i == 24) {
			components = append(components, row)
			row.Components = nil
			continue
		}
		if (i%5 == 0 && i >= 5 || i == len(btns)+1 && i >= 5) && panelType != CreateRolePanel || panelType == MarketPanel && i == 3 {
			if panelType == MarketPanel && i == 5 {
				continue
			}
			if panelType == ReactionPanel {
				continue
			}
			components = append(components, row)
			row.Components = nil
		}
	}
	if row.Components != nil {
		components = append(components, row)
	}

	if embed.Color == 0 {
		embed.Color = 3092790
	}
	embed.Footer = &discordgo.MessageEmbedFooter{
		IconURL: i.Member.User.AvatarURL("128"),
		Text:    fmt.Sprintf("%v · %v", i.Member.User.Username, i.Member.User.ID),
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: t,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

func (b *Bot) openPackPanel(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {

	var userID = i.Member.User.ID

	lbxsAmount, err := b.MarketData.LootboxAmount(userID)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get customer's lootbox amount info")
		return
	}

	isBtnOpenPackDisabled := true
	if lbxsAmount != 0 {
		isBtnOpenPackDisabled = false
	}

	isBtnBuyPackDisabled := true
	if isEnoughCurrency, _, err := b.isEnoughCurrency(userID, b.Cfg.LootboxPrice); isEnoughCurrency {
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
			return
		}
		isBtnBuyPackDisabled = false
	}

	embed.Title = "Паки"
	embed.Description += fmt.Sprintf("\n\n> **Паков в наличии:**\n\n%v `%v` ", b.Cfg.EmojiLootbox, lbxsAmount)
	openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, embed,
		[]*discordgo.Button{{
			Label:    "⠀Открыть⠀ ⠀",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnOpenPackDisabled,
			CustomID: "Открыть Пак",
			Emoji:    discordgo.ComponentEmoji{ID: "854346864211656705"},
		}, {
			Label:    "⠀ ⠀Купить⠀ ⠀",
			Style:    discordgo.DangerButton,
			Disabled: isBtnBuyPackDisabled,
			CustomID: "Купить Пак",
			Emoji:    discordgo.ComponentEmoji{ID: "854346916593795092"},
		}}, DefaultPanel)
}

func (b *Bot) addChannels(s *discordgo.Session, i *discordgo.InteractionCreate, textChannelID string, voiceChannelID string, ch chan map[string]string) {

	userId := i.Member.User.ID

	channels := make(map[string]string)

	if textChannelID == "" {

		category, err := s.GuildChannelCreate(i.GuildID, i.Member.User.Username, discordgo.ChannelTypeGuildCategory)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error creating category channel")
			return
		}
		channels["categoryChannel"] = category.Name

		textChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{Name: "『💬』сообщения", Type: discordgo.ChannelTypeGuildText, ParentID: category.ID})
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error creating text channel")
			return
		}
		channels["textChannel"] = textChannel.Mention()

		voiceChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{Name: "『🎤』голос", Type: discordgo.ChannelTypeGuildVoice, ParentID: category.ID})
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error creating voice channel")
			return
		}
		channels["voiceChannel"] = voiceChannel.Mention()

		ch <- channels
		s.ChannelMessageSend("859539629048856578", fmt.Sprintf("!Синхронизировать %v %v", userId, voiceChannel.ID))
		err = s.ChannelPermissionSet(textChannel.ID, i.Member.User.ID, discordgo.PermissionOverwriteTypeMember, 805829713, 0)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error setting permissions to text channel")
			return
		}
		err = s.ChannelPermissionSet(voiceChannel.ID, i.Member.User.ID, discordgo.PermissionOverwriteTypeMember, 334497553, 0)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error setting permissions to voice channel")
			return
		}
		role, err := s.GuildRoleCreate(i.GuildID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error creating role")
			return
		}
		_, err = s.GuildRoleEdit(i.GuildID, role.ID, i.Member.User.Username, 8357810, false, 0, false)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error editing role")
			return
		}
		guild, err := s.Guild(i.GuildID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting guild")
			return
		}
		spacesRolePosition := 0
		for _, v := range guild.Roles {
			if v.ID == "836937639831142401" {
				spacesRolePosition = v.Position
			}
		}
		for _, v := range guild.Roles {
			if v.ID == role.ID {
				v.Position = spacesRolePosition + 1
			}
		}
		_, err = s.GuildRoleReorder(i.GuildID, guild.Roles)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error editing role")
			return
		}
		err = s.GuildMemberRoleAdd(i.GuildID, userId, role.ID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error adding role")
			return
		}
		err = s.GuildMemberRoleAdd(i.GuildID, userId, "836937639831142401")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error adding role")
			return
		}
		err = b.MarketData.CustomerChannelsAdd(userId, textChannel.ID, voiceChannel.ID, role.ID, false)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to add user's channels")
			return
		}
	} else {
		textChannel, err := s.Channel(textChannelID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting text channel in guild")
			return
		}
		channels["textChannel"] = textChannel.Mention()
		voiceChannel, err := s.Channel(voiceChannelID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting voice channel in guild")
			return
		}
		channels["voiceChannel"] = voiceChannel.Mention()
		category, err := s.Channel(textChannel.ParentID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting category in guild")
			return
		}
		channels["categoryChannel"] = category.Name
		err = b.MarketData.CustomerChannelsAdd(userId, "", "", "", true)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to add user's channels")
			return
		}
		ch <- channels
	}
}
func (b *Bot) OnMessegeReceive(s *discordgo.Session, m *discordgo.MessageCreate) {

	userId := m.Author.ID
	msgHasEmbeds := len(m.Embeds) != 0
	msgHasFooter := false
	if msgHasEmbeds && m.Embeds[0].Footer != nil && strings.Contains(m.Embeds[0].Footer.Text, " · ") {
		msgHasFooter = true
	}
	for i, v := range panelsCalledUsers {
		if msgHasEmbeds && msgHasFooter {
			if v == strings.Split(m.Embeds[0].Footer.Text, " · ")[1] {
				userIdMsgIdDict[strings.Split(m.Embeds[0].Footer.Text, " · ")[1]] = m.Message.ID
				panelsCalledUsers[i] = panelsCalledUsers[len(panelsCalledUsers)-1]
				panelsCalledUsers = panelsCalledUsers[:len(panelsCalledUsers)-1]
				return
			}
		}
	}

	for i, v := range setStatusUsers {
		if v == userId {
			if len(m.Message.Content) > 100 {
				s.ChannelMessageSendEmbed(b.Cfg.LotteryChannelID,
					&discordgo.MessageEmbed{
						Title:       "Статус",
						Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/866316842879811584/status.png"},
						Description: "*Длина статуса не может быть больше 100 символов!*\n\n*Ожидаю повторный ввод статуса, в виде текстового сообщения.*",
						Footer: &discordgo.MessageEmbedFooter{
							IconURL: m.Author.AvatarURL("128"),
							Text:    fmt.Sprintf("%v · %v", m.Author.Username, userId),
						},
					})
			} else {
				setStatusUsers[i] = setStatusUsers[len(setStatusUsers)-1]
				setStatusUsers = setStatusUsers[:len(setStatusUsers)-1]
				err := b.MarketData.CustomerStatusUpdate(userId, m.Message.Content, b.Cfg.StatusPrice)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update customer's status")
					return
				}
				s.ChannelMessageSendEmbed(b.Cfg.LotteryChannelID,
					&discordgo.MessageEmbed{
						Title:       "Статус",
						Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/866316842879811584/status.png"},
						Description: fmt.Sprintf("> **Статус:**\n```fix\n%v```\n*Загляни скорее в профиль, чтобы убедиться в этом.*\n\n*Сестичка явно недооценила братика в фантазии...*", m.Message.Content),
						Footer: &discordgo.MessageEmbedFooter{
							IconURL: m.Author.AvatarURL("128"),
							Text:    fmt.Sprintf("%v · %v", m.Author.Username, userId),
						},
					})
			}
			return
		}
	}
	for i, v := range createRoleUsers {
		if v.UserId == userId {
			if len(m.Message.Content) > 100 {
				s.ChannelMessageSendEmbed(b.Cfg.LotteryChannelID,
					&discordgo.MessageEmbed{
						Color:       3092790,
						Title:       "Роли",
						Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858704257402404904/icon_role_create.png"},
						Description: "*Название роли не может иметь больше 100 символов!*\n\n*Ожидаю повторного ввода названия роли, в виде текстового сообщения.*",
						Footer: &discordgo.MessageEmbedFooter{
							IconURL: m.Author.AvatarURL("128"),
							Text:    fmt.Sprintf("%v · %v", m.Author.Username, userId),
						},
					})
			} else {
				randImage := randSlice([]string{
					"https://cdn.discordapp.com/attachments/838477787877998624/858703916555698196/role.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703919718989844/role1.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703921458970645/role2.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703924038598686/role3.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703925662056508/role4.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858704180276625438/role6.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703907211706388/role7.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703908607361024/role8.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703909321441290/role9.gif",
					"https://cdn.discordapp.com/attachments/838477787877998624/858703914077519872/role10.gif",
				})

				createRoleUsers[i] = createRoleUsers[len(createRoleUsers)-1]
				createRoleUsers = createRoleUsers[:len(createRoleUsers)-1]
				role, err := s.GuildRoleCreate(m.GuildID)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error creating role")
					return
				}
				role, err = s.GuildRoleEdit(m.GuildID, role.ID, m.Message.Content, v.RoleColor, false, 0, false)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error editing role")
					return
				}
				s.ChannelMessageSendEmbed(b.Cfg.LotteryChannelID,
					&discordgo.MessageEmbed{
						Color:       3092790,
						Title:       "Роли",
						Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858704257402404904/icon_role_create.png"},
						Image:       &discordgo.MessageEmbedImage{URL: randImage},
						Description: fmt.Sprintf("*Ооо, да!*\n\n*Теперь у тебя есть уникальная роль:*\n\n%v\n\n*Они-чан, самое время искать тяночку!*", role.Mention()),
						Footer: &discordgo.MessageEmbedFooter{
							IconURL: m.Author.AvatarURL("128"),
							Text:    fmt.Sprintf("%v · %v", m.Author.Username, userId),
						},
					})
				var value int
				switch int(v.Duration / dayAsSec) {
				case 1:
					value = -1500
				case 3:
					value = -3900
				case 7:
					value = -8100
				case 31:
					value = -30000
				}
				err = b.MarketData.BalanceUpdate(userId, value)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
				}
				err = b.MarketData.CustomRoleAdd(role.ID, m.Message.Content, userId, time.Now().Unix()+v.Duration)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error adding custom role")
					return
				}
				err = b.roleGive(s, m.GuildID, userId, role.ID)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error giving role")
					return
				}
				guild, err := s.Guild(m.GuildID)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting guild")
					return
				}
				pos := 0
				for _, v := range guild.Roles {
					if v.ID == "834432874945839160" {
						pos = v.Position
						break
					}
				}
				for _, v := range guild.Roles {
					if v.ID == role.ID {
						v.Position = pos
						break
					}
				}
				_, err = s.GuildRoleReorder(m.GuildID, guild.Roles)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error editing role")
					return
				}
				err = s.GuildMemberRoleAdd(m.GuildID, userId, role.ID)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error adding role")
					return
				}
				err = s.GuildMemberRoleAdd(m.GuildID, userId, "836937639831142401")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error adding role")
					return
				}
			}
			return
		}
	}
}
func (b *Bot) OnInteractionAct(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var userId = i.Member.User.ID

	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == "панель" {
			if i.ChannelID == "941236439146979378" || i.ChannelID == "836964413298704395" || i.ChannelID == "834432875129733138" || i.ChannelID == "844343580314566696" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   1 << 6,
						Content: fmt.Sprintf("Вы увидели это сообщение, потому что пытались вызвать панель в канале <#%v>\n\nПанель доступна лишь в канале <#941259025939595305> либо в личном пространстве.", i.ChannelID),
					},
				})
				return
			}
			panelsCalledUsers = append(panelsCalledUsers, userId)
			openMainPanel(s, i, discordgo.InteractionResponseChannelMessageWithSource)
		} else if i.ApplicationCommandData().Name == "реакция" {
			b.reactionBase(s, i)
		}
	}

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	if userIdMsgIdDict[userId] != i.Message.ID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "```Вы увидели это сообщение, потому что пытались взаимодействовать с чужой либо неактуальной панелью.```\n__Вы всегда можете вызвать новую панель командой:__ `/панель`",
			},
		})
		return
	}

	var customId = i.MessageComponentData().CustomID

	switch customId {

	case "<<":
		for i, v := range createRoleUsers {
			if v.UserId == userId {
				createRoleUsers[i] = createRoleUsers[len(createRoleUsers)-1]
				createRoleUsers = createRoleUsers[:len(createRoleUsers)-1]
				break
			}
		}
		for i, v := range setStatusUsers {
			if v == userId {
				setStatusUsers[i] = setStatusUsers[len(setStatusUsers)-1]
				setStatusUsers = setStatusUsers[:len(setStatusUsers)-1]
				break
			}
		}
		openMainPanel(s, i, discordgo.InteractionResponseUpdateMessage)
	case "Фарм":

		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/941081208773828678/farm.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/941081208278909059/farm2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/941082016848424960/farm3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/941082017410474084/farm4.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/941082015581733044/farm5.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/941082016231870526/farm6.gif",
		})

		nextFarmTime := time.Now().Add(12 * time.Hour).Unix()
		minUntilFarm := time.Until(time.Unix(nextFarmTime, 0)).Minutes()
		hours := int(minUntilFarm) / 60
		mins := int(minUntilFarm) - hours*60
		err := b.MarketData.Farm(userId, b.Cfg.FarmRate, nextFarmTime)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Can not update user's farm data")
			return
		}

		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Фарм",
			Image:       &discordgo.MessageEmbedImage{URL: randImage},
			Description: fmt.Sprintf("*Ты успешно украл у сестрички `%v` %s\n\nОни-чан, приходи через `%v` час. `%v` мин.*", b.Cfg.FarmRate, b.Cfg.Currency, hours, mins),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/853796492476809226/mining_trusiki.png"},
		}, nil, DefaultPanel)

	case "Профиль":

		balance, err := b.MarketData.Balance(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
			return
		}

		nextFarmTime, err := b.MarketData.NextFarm(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's farm data")
			return
		}

		btnLabel := "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ Фарм ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀"
		isBtnFarmDisabled := false
		if nextFarmTime > time.Now().Unix() {
			minUntilFarm := time.Until(time.Unix(nextFarmTime, 0)).Minutes()
			hours := int(minUntilFarm) / 60
			mins := int(minUntilFarm) - hours*60
			isBtnFarmDisabled = true
			btnLabel = fmt.Sprintf("⠀⠀%v:%v⠀⠀", hours, mins)
		}

		status, voiceT, spent, lbxs, err := b.MarketData.CustomerProfileData(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get customer's profile info")
			return
		}

		isBtnSetStatusDisabled := true
		if balance >= b.Cfg.StatusPrice {
			isBtnSetStatusDisabled = false
		}
		if status == "" {
			status = "Статус не установлен."
		}

		var r, e int
		is, err := b.MarketData.UserOrders(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's items")
			return
		}

		for _, i := range is {
			if i.Product.Type == "role" {
				r += 1
			} else if i.Product.Type == "reaction" {
				e += 1
			}
		}

		hasSpace := "Нет"
		textChannelID, _, _, _, err := b.MarketData.HasCustomerChannels(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's channels")
			return
		}
		if textChannelID != "" {
			hasSpace = "Есть"
		}

		hasTicket := "Нет"
		for _, v := range b.LotterySession.Participants {
			if v == userId {
				hasTicket = "Есть"
				break
			}
		}

		isBtnManageRolesDisabled, err := b.MarketData.HasRoles(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's roles")
			return
		}
		isBtnManageRolesDisabled = !isBtnManageRolesDisabled

		gamesEmoji := ""
		genderEmoji := ""
		isFemale := false
		isMale := false

		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		for _, roleID := range u.Roles {
			if roleID == "836888796443574313" {
				gamesEmoji += "<:cs_go:941493420541763654>⠀"
			} else if roleID == "836883811782361109" {
				gamesEmoji += "<:dota_2:941493420705329152>⠀"
			} else if roleID == "836884309985067070" {
				gamesEmoji += "<:lol:941493420663382087>⠀"
			}
			if roleID == "834432874928799767" {
				isFemale = true
				genderEmoji = "♀️"
			} else if roleID == "836723670813638677" {
				isMale = true
				genderEmoji = "♂️"
			}
		}
		if !isFemale && !isMale {
			genderEmoji = "❓"
		}
		if gamesEmoji == "" {
			gamesEmoji = "`Нет`"
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:     "Профиль",
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: i.Member.User.AvatarURL("128")},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "> Баланс:",
					Value:  fmt.Sprintf("`%v` %v", balance, b.Cfg.Currency),
					Inline: true,
				},
				{
					Name:   "> Гендер:",
					Value:  genderEmoji,
					Inline: true,
				},
				{
					Name:   "> Гейминг:",
					Value:  gamesEmoji,
					Inline: true,
				},
				{
					Name:  "> Статус:",
					Value: fmt.Sprintf("```fix\n%s ```", status),
				},
				{
					Name:  "> Проведено часов в голосовых каналах:",
					Value: fmt.Sprintf("```ini\n[ %d ]```", voiceT/hoursAsSec),
				},
				{
					Name:   "Всего потрачено:",
					Value:  fmt.Sprintf(" **%d** %s", spent, b.Cfg.Currency),
					Inline: true,
				},
				{
					Name: "Инвентарь:",
					Value: fmt.Sprintf("%v `%v`  %v `%v`  %v `%v`  %v `%v`  %v `%v`", b.Cfg.EmojiPrivateChannel, hasSpace, b.Cfg.EmojiRole, r, b.Cfg.EmojiReaction,
						e, b.Cfg.EmojiLootbox, lbxs, b.Cfg.EmojiLottery, hasTicket),
					Inline: true,
				},
			},
		},
			[]*discordgo.Button{{
				Label:    btnLabel,
				Style:    discordgo.SuccessButton,
				Disabled: isBtnFarmDisabled,
				Emoji:    discordgo.ComponentEmoji{ID: "843056414602952734"},
				CustomID: "Фарм",
			}, {
				Label:    "⠀ ⠀Статус⠀⠀",
				Style:    discordgo.PrimaryButton,
				Disabled: isBtnSetStatusDisabled,
				Emoji:    discordgo.ComponentEmoji{Name: "🏷️"},
				CustomID: "Статус",
			}, {
				Label:    "⠀⠀ Игровые Роли⠀⠀⠀",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				Emoji:    discordgo.ComponentEmoji{Name: "🎮"},
				CustomID: "Игровые Роли",
			},
				{
					Label:    "⠀⠀⠀Гендер⠀⠀⠀",
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					Emoji:    discordgo.ComponentEmoji{ID: "941485792096833536"},
					CustomID: "Гендер",
				}, {
					Label:    "⠀⠀⠀⠀⠀⠀Маркет Топ-10⠀⠀⠀⠀⠀",
					Style:    discordgo.DangerButton,
					Disabled: false,
					Emoji:    discordgo.ComponentEmoji{Name: "💹"},
					CustomID: "Топ Маркет",
				}, {
					Label:    "⠀⠀⠀⠀⠀⠀ Войс Топ-10⠀⠀⠀⠀⠀⠀",
					Style:    discordgo.DangerButton,
					Disabled: false,
					Emoji:    discordgo.ComponentEmoji{Name: "🎤"},
					CustomID: "Топ Войс",
				}, {
					Label:    "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ Управление Ролями ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
					Style:    discordgo.SecondaryButton,
					Disabled: isBtnManageRolesDisabled,
					Emoji:    discordgo.ComponentEmoji{Name: "🔧"},
					CustomID: "Управление Ролями",
				}}, ProfilePanel)
	case "Статус":
		setStatusUsers = append(setStatusUsers, userId)
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Статус",
			Description: "*Установите статус, который будет отображаться в профиле.*\n\n**Ограничения:**\n\n> *Максимальное количество символов:* `100`\n\n> *Кастомные эмодзи не приемлемы, используйте эмодзи и символы [Юникода](https://unicode-table.com/ru/).*\n\n*Ожидаю ввода, в виде текстового сообщения.*",
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/866316842879811584/status.png"},
		}, nil, DefaultPanel)
	case "Гендер":
		genderEmoji := ""
		isBtnFemaleDisabled := false
		isBtnMaleDisabled := false

		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}

		for _, roleID := range u.Roles {
			if roleID == "834432874928799767" {
				isBtnFemaleDisabled = true
				genderEmoji = "♀️"
			} else if roleID == "836723670813638677" {
				isBtnMaleDisabled = true
				genderEmoji = "♂️"
			}
		}
		if !isBtnFemaleDisabled && !isBtnMaleDisabled {
			genderEmoji = "❓"
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Гендер",
			Description: fmt.Sprintf("*Используйте кнопки для получения роли. Чтобы убрать роль, кликните⠀🚫*\n\n<@&834432874928799767> ***— девушка.***\n<@&836723670813638677> ***— парень.\n\n> Текущий гендер:⠀%v***", genderEmoji),
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnFemaleDisabled,
			Emoji:    discordgo.ComponentEmoji{Name: "♀️"},
			CustomID: "Тян",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnMaleDisabled,
			Emoji:    discordgo.ComponentEmoji{Name: "♂️"},
			CustomID: "Кун",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: !isBtnFemaleDisabled && !isBtnMaleDisabled,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Неизвестный",
		}}, DefaultPanel)
	case "Тян":
		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		for _, roleID := range u.Roles {
			if roleID == "836723670813638677" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "836723670813638677")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			}
		}
		err = s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "834432874928799767")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
			return
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Гендер",
			Description: "*Используйте кнопки для получения роли. Чтобы убрать роль, кликните⠀🚫*\n\n<@&834432874928799767> ***— девушка.***\n<@&836723670813638677> ***— парень.\n\n> Текущий гендер:⠀♀️***",
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{Name: "♀️"},
			CustomID: "Тян",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "♂️"},
			CustomID: "Кун",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Неизвестный",
		}}, DefaultPanel)
	case "Кун":
		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		for _, roleID := range u.Roles {
			if roleID == "834432874928799767" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "834432874928799767")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			}
		}
		err = s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836723670813638677")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
			return
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Гендер",
			Description: "*Используйте кнопки для получения роли. Чтобы убрать роль, кликните⠀🚫*\n\n<@&834432874928799767> ***— девушка.***\n<@&836723670813638677> ***— парень.\n\n> Текущий гендер:⠀♂️***",
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "♀️"},
			CustomID: "Тян",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{Name: "♂️"},
			CustomID: "Кун",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Неизвестный",
		}}, DefaultPanel)
	case "Неизвестный":
		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		for _, roleID := range u.Roles {
			if roleID == "834432874928799767" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "834432874928799767")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			} else if roleID == "836723670813638677" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "836723670813638677")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			}
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Гендер",
			Description: "*Используйте кнопки для получения роли. Чтобы убрать роль, кликните⠀🚫*\n\n<@&834432874928799767> ***— девушка.***\n<@&836723670813638677> ***— парень.\n\n> Текущий гендер:⠀❓***",
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "♀️"},
			CustomID: "Тян",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "♂️"},
			CustomID: "Кун",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Неизвестный",
		}}, DefaultPanel)
	case "Игровые Роли":
		gamesEmoji := ""
		isBtnDota2Disabled := false
		isBtnLolDisabled := false
		isBtnCsgoDisabled := false

		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		for _, roleID := range u.Roles {
			if roleID == "836883811782361109" {
				isBtnDota2Disabled = true
				gamesEmoji += "<:dota_2:941493420705329152>⠀"
			} else if roleID == "836884309985067070" {
				isBtnLolDisabled = true
				gamesEmoji += "<:lol:941493420663382087>⠀"
			} else if roleID == "836888796443574313" {
				isBtnCsgoDisabled = true
				gamesEmoji += "<:cs_go:941493420541763654>⠀"
			}
		}
		if gamesEmoji == "" {
			gamesEmoji = "`Нет`"
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Игровые Роли",
			Description: fmt.Sprintf("*Используйте кнопки для получения роли. Чтобы убрать роли, кликните⠀🚫*\n\n<@&836888796443574313>***,*** <@&836884309985067070>***,*** <@&836883811782361109>*** — пользователь с одной из этих ролей получает уведомления в ***<#834432875129733138>\n\n> ***Текущие игровые роли:***⠀%v", gamesEmoji),
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnDota2Disabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420705329152"},
			CustomID: "Dota 2",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnLolDisabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420663382087"},
			CustomID: "LoL",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnCsgoDisabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420541763654"},
			CustomID: "CS:GO",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: !isBtnCsgoDisabled && !isBtnLolDisabled && !isBtnDota2Disabled,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Удалить Роли",
		}}, DefaultPanel)
	case "Dota 2":
		err := s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836883811782361109")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
			return
		}
		gamesEmoji := ""
		isBtnLolDisabled := false
		isBtnCsgoDisabled := false
		isDelimiterRoleAssigned := false

		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		gamesEmoji += "<:dota_2:941493420705329152>⠀"
		for _, roleID := range u.Roles {
			if roleID == "836884309985067070" {
				isBtnLolDisabled = true
				gamesEmoji += "<:lol:941493420663382087>⠀"
			} else if roleID == "836888796443574313" {
				isBtnCsgoDisabled = true
				gamesEmoji += "<:cs_go:941493420541763654>⠀"
			} else if roleID == "836862989872922675" {
				isDelimiterRoleAssigned = true
			}
		}
		if !isDelimiterRoleAssigned {
			err := s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836862989872922675")
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
				return
			}
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Игровые Роли",
			Description: fmt.Sprintf("*Используйте кнопки для получения роли. Чтобы убрать роли, кликните⠀🚫*\n\n<@&836888796443574313>***,*** <@&836884309985067070>***,*** <@&836883811782361109>*** — пользователь с одной из этих ролей получает уведомления в ***<#834432875129733138>\n\n> ***Текущие игровые роли:***⠀%v", gamesEmoji),
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420705329152"},
			CustomID: "Dota 2",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnLolDisabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420663382087"},
			CustomID: "LoL",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnCsgoDisabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420541763654"},
			CustomID: "CS:GO",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Удалить Роли",
		}}, DefaultPanel)
	case "LoL":
		err := s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836884309985067070")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
			return
		}
		gamesEmoji := ""
		isBtnDota2Disabled := false
		isBtnCsgoDisabled := false
		isDelimiterRoleAssigned := false

		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		gamesEmoji += "<:lol:941493420663382087>⠀"
		for _, roleID := range u.Roles {
			if roleID == "836883811782361109" {
				isBtnDota2Disabled = true
				gamesEmoji += "<:dota_2:941493420705329152>⠀"
			} else if roleID == "836888796443574313" {
				isBtnCsgoDisabled = true
				gamesEmoji += "<:cs_go:941493420541763654>⠀"
			} else if roleID == "836862989872922675" {
				isDelimiterRoleAssigned = true
			}
		}
		if !isDelimiterRoleAssigned {
			err := s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836862989872922675")
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
				return
			}
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Игровые Роли",
			Description: fmt.Sprintf("*Используйте кнопки для получения роли. Чтобы убрать роли, кликните⠀🚫*\n\n<@&836888796443574313>***,*** <@&836884309985067070>***,*** <@&836883811782361109>*** — пользователь с одной из этих ролей получает уведомления в ***<#834432875129733138>\n\n> ***Текущие игровые роли:***⠀%v", gamesEmoji),
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnDota2Disabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420705329152"},
			CustomID: "Dota 2",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420663382087"},
			CustomID: "LoL",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnCsgoDisabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420541763654"},
			CustomID: "CS:GO",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Удалить Роли",
		}}, DefaultPanel)
	case "CS:GO":
		err := s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836888796443574313")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
			return
		}
		gamesEmoji := ""
		isBtnDota2Disabled := false
		isBtnLolDisabled := false
		isDelimiterRoleAssigned := false
		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		gamesEmoji += "<:cs_go:941493420541763654>⠀"
		for _, roleID := range u.Roles {
			if roleID == "836883811782361109" {
				isBtnDota2Disabled = true
				gamesEmoji += "<:dota_2:941493420705329152>⠀"
			} else if roleID == "836884309985067070" {
				isBtnLolDisabled = true
				gamesEmoji += "<:lol:941493420663382087>⠀"
			} else if roleID == "836862989872922675" {
				isDelimiterRoleAssigned = true
			}
		}
		if !isDelimiterRoleAssigned {
			err := s.GuildMemberRoleAdd(b.Cfg.GuildID, userId, "836862989872922675")
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
				return
			}
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Игровые Роли",
			Description: fmt.Sprintf("*Используйте кнопки для получения роли. Чтобы убрать роли, кликните⠀🚫*\n\n<@&836888796443574313>***,*** <@&836884309985067070>***,*** <@&836883811782361109>*** — пользователь с одной из этих ролей получает уведомления в ***<#834432875129733138>\n\n> ***Текущие игровые роли:***⠀%v", gamesEmoji),
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnDota2Disabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420705329152"},
			CustomID: "Dota 2",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: isBtnLolDisabled,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420663382087"},
			CustomID: "LoL",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420541763654"},
			CustomID: "CS:GO",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Удалить Роли",
		}}, DefaultPanel)
	case "Удалить Роли":
		gamesEmoji := "`Нет`"

		u, err := s.GuildMember(b.Cfg.GuildID, userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user")
			return
		}
		for _, roleID := range u.Roles {
			if roleID == "836883811782361109" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "836883811782361109")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			} else if roleID == "836884309985067070" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "836884309985067070")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			} else if roleID == "836888796443574313" {
				err := s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "836888796443574313")
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
					return
				}
			}
		}
		err = s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, "836862989872922675")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the role to the user")
			return
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Игровые Роли",
			Description: fmt.Sprintf("*Используйте кнопки для получения роли. Чтобы убрать роли, кликните⠀🚫*\n\n<@&836888796443574313>***,*** <@&836884309985067070>***,*** <@&836883811782361109>*** — пользователь с одной из этих ролей получает уведомления в ***<#834432875129733138>\n\n> ***Текущие игровые роли:***⠀%v", gamesEmoji),
		}, []*discordgo.Button{{
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420705329152"},
			CustomID: "Dota 2",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420663382087"},
			CustomID: "LoL",
		}, {
			Label:    "",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji:    discordgo.ComponentEmoji{ID: "941493420541763654"},
			CustomID: "CS:GO",
		}, {
			Label:    "",
			Style:    discordgo.DangerButton,
			Disabled: true,
			Emoji:    discordgo.ComponentEmoji{Name: "🚫"},
			CustomID: "Удалить Роли",
		}}, DefaultPanel)
	case "Топ Маркет":

		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/845729225204957294/top_spender_1.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729224312487966/top_spender_2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729224135278612/top_spender_3.jpg",
			"https://cdn.discordapp.com/attachments/838477787877998624/940202319507759124/top_spender3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729230645886986/top_spender_4.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729222386647051/top_spender_5.jpg",
			"https://cdn.discordapp.com/attachments/838477787877998624/941079746773672026/top_spender6.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/941079747453124618/top_spender7.gif",
		})
		fields := []*discordgo.MessageEmbedField{}
		sp, err := b.MarketData.TopTenSpenders()
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get top 10 spenders")
			return
		}
		for n, s := range sp {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%v место:", n+1),
				Value: fmt.Sprintf("<@!%v> — **%v** %v", s.UserID, s.Spent, b.Cfg.Currency),
			})
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:  "Топ-10 (Маркет)",
			Fields: fields,
			Image:  &discordgo.MessageEmbedImage{URL: randImage},
		}, nil, DefaultPanel)
	case "Топ Войс":
		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/845729222386647051/top_spender_5.jpg",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729224312487966/top_spender_2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729225204957294/top_spender_1.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/845729230645886986/top_spender_4.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/940202319507759124/top_spender3.gif",
		})
		fields := []*discordgo.MessageEmbedField{}
		vs, err := b.MarketData.TopTenVoice()
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get top 10 voice")
			return
		}

		for n, v := range vs {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%v место:", n+1),
				Value: fmt.Sprintf("<@!%v> — **%v** ч. ", v.UserID, v.Spent/hoursAsSec),
			})
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:  "Топ-10 (Войс)",
			Fields: fields,
			Image:  &discordgo.MessageEmbedImage{URL: randImage},
		}, nil, DefaultPanel)
	case "Управление Ролями":
		var rolesID []string
		var rolesDials []int
		var rolesIsHidden []bool
		rolesList := ""
		ps, err := b.MarketData.UserOrders(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get products")
		}
		for _, p := range ps {
			if p.Product.Type == "role" {
				rolesID = append(rolesID, p.Product.RoleID)
				rolesDials = append(rolesDials, p.Product.Dial)
				rolesIsHidden = append(rolesIsHidden, p.IsHidden)
				showOrHide := "Отображается"
				if p.IsHidden {
					showOrHide = "Скрыта"
				}
				rolesList += fmt.Sprintf("\n\n<@&%v> — `%v`", p.Product.RoleID, showOrHide)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to determine whether user has a product")
					return
				}

			}
		}
		btnsRoleHide := []*discordgo.Button{}
		for i := 1; i < len(rolesID)+1; i++ {
			label := ""
			emoji := ""
			customID := ""
			if rolesIsHidden[i-1] {
				label = "⠀ Показать ⠀"
				customID = "Показать"
			} else {
				label = "⠀⠀Скрыть⠀⠀"
				customID = "Скрыть"
			}
			switch rolesDials[i-1] {
			case 1:
				emoji = "♟️"
			case 2:
				emoji = "🌸"
			case 3:
				emoji = "❄️"
			case 4:
				emoji = "🔫"
			case 5:
				emoji = "🍾"
			case 6:
				emoji = "🐾"
			case 7:
				emoji = "🗡️"
			case 8:
				emoji = "🤪"
			case 9:
				emoji = "⭐"
			case 10:
				emoji = "🧬"
			case 11:
				emoji = "🏵️"
			case 12:
				emoji = "🧻"
			case 13:
				emoji = "🌈"
			case 14:
				emoji = "💥"
			case 15:
				emoji = "💊"
			}

			btnsRoleHide = append(btnsRoleHide, &discordgo.Button{
				Label:    label,
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: fmt.Sprintf(customID+" Роль %v", rolesDials[i-1]),
				Emoji:    discordgo.ComponentEmoji{Name: emoji},
			})
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Управление Ролями",
			Description: fmt.Sprintf("```Приобритенная роль может стоять выше других ролей, тем самым, лишая вас возможности изменить цвет никнейма.```\n```Кастомные роли всегда выше в иерархии ролей, поэтому их нельзя показать / скрыть.```\n*Выберите роль для показа / сокрытия.*\n\n> **Приобретенные роли:**%v", rolesList),
		}, btnsRoleHide, DefaultPanel)
	case "Купить Пак":

		_, balance, err := b.isEnoughCurrency(userId, b.Cfg.LootboxPrice)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
			return
		}
		isbtnLtbxAmountDisabled1 := true
		isbtnLtbxAmountDisabled3 := true
		isbtnLtbxAmountDisabled5 := true
		isbtnLtbxAmountDisabled10 := true
		if balance >= b.Cfg.LootboxPrice*1 {
			isbtnLtbxAmountDisabled1 = false
		}
		if balance >= b.Cfg.LootboxPrice*3 {
			isbtnLtbxAmountDisabled3 = false
		}
		if balance >= b.Cfg.LootboxPrice*5 {
			isbtnLtbxAmountDisabled5 = false
		}
		if balance >= b.Cfg.LootboxPrice*10 {
			isbtnLtbxAmountDisabled10 = false
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Паки",
			Description: "*Выберите, сколько паков вы хотите приобрести.*",
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/868859306051768340/icon_lootbox_buy.png"},
		},
			[]*discordgo.Button{{
				Style:    discordgo.SuccessButton,
				Disabled: isbtnLtbxAmountDisabled1,
				CustomID: "Количество Паков 1",
				Emoji:    discordgo.ComponentEmoji{ID: "859212034515533864"},
			}, {
				Style:    discordgo.SuccessButton,
				Disabled: isbtnLtbxAmountDisabled3,
				CustomID: "Количество Паков 3",
				Emoji:    discordgo.ComponentEmoji{ID: "859212034662727680"},
			}, {
				Style:    discordgo.SuccessButton,
				Disabled: isbtnLtbxAmountDisabled5,
				CustomID: "Количество Паков 5",
				Emoji:    discordgo.ComponentEmoji{ID: "859212034095710239"},
			}, {
				Style:    discordgo.SuccessButton,
				Disabled: isbtnLtbxAmountDisabled10,
				CustomID: "Количество Паков 10",
				Emoji:    discordgo.ComponentEmoji{ID: "859210540398936094"},
			},
			}, DefaultPanel)
	case "Открыть Пак":
		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/854813596747563028/pack_open.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/854813575084638248/pack_open2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/854813573334564864/pack_open3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/854813567097241631/open_pack4.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/854813707820597248/open_pack5.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/854813569510408252/open_pack6.gif",
		})
		err := b.MarketData.LootboxRemove(userId, 1)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove a lootbox from the user")
		}

		embed := &discordgo.MessageEmbed{
			Description: "```yaml\nТы отрываешь верхушку от пака и видишь...```",
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/868859317934252072/icon_lootbox_open.png"},
			Image: &discordgo.MessageEmbedImage{
				URL: randImage,
			},
		}

		rand.Seed(time.Now().UnixNano())
		r := rand.Float64()
		// User gets currency
		if r < 0.7 {
			amount := int(float64(rand.Intn(b.Cfg.LootboxPrice)) * 1.3)
			err = b.MarketData.BalanceUpdate(userId, amount)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
			}
			embed.Description += fmt.Sprintf("\n*Ммм... Свежие трусики, в количестве:*\n\n`%v` %v", amount, b.Cfg.Currency)
			b.openPackPanel(s, i, embed)
			return
		}

		ps, err := b.MarketData.Products()
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get products")
		}

		r = rand.Float64()
		dur := 0
		if r < 0.33 {
			dur = 1
		} else if r < 0.66 {
			dur = 3
		} else {
			dur = 7
		}
		duration := int64(dur) * dayAsSec
		textDuration := fmt.Sprintf("*на* `%d` *д.*", dur)

		if r < 0.65 {
			// User gets a reaction
			var reactions []product
			for _, p := range ps {
				if p.Type == "reaction" {
					reactions = append(reactions, p)
				}
			}
			p := &reactions[rand.Intn(len(reactions))]
			p.Price = 0
			embed.Description += fmt.Sprintf("\n*Повезло, повезло!*\n\n*Набор аниме-реакций —*\n\n**%v** %v\n\n*Самое время на ком-нибудь опробовать!* 🤪", p.Name, b.Cfg.EmojiReaction)
			err = b.MarketData.Order(userId, p, -1)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error while processing order")
				return
			}

		} else if r < 0.95 {
			// User gets a role
			var roles []product
			for _, p := range ps {
				if p.Type == "role" {
					roles = append(roles, p)
				}
			}

			p := &roles[rand.Intn(len(roles))]
			p.Price = 0
			embed.Description += fmt.Sprintf("\n*Какая красивенькая, миленькая роль* %v\n\n %v <@&%s> \n\n*Что может быть лучше?*", textDuration, b.Cfg.EmojiRole, p.RoleID)

			err = b.MarketData.Order(userId, p, duration)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error while processing order")
				return
			}

			err = b.roleGive(s, i.GuildID, userId, p.RoleID)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error giving role")
				return
			}

		} else {
			ch := make(chan map[string]string)
			defer close(ch)
			textChannelID, voiceChannelID, _, _, err := b.MarketData.HasCustomerChannels(userId)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's channels")
				return
			}
			go b.addChannels(s, i, textChannelID, voiceChannelID, ch)
			c := <-ch
			if textChannelID == "" {
				embed.Description += fmt.Sprintf("\n***ЧТО!? Поздравляю, они-чан!***\n\n*Тебе выпало приватное пространство* %v\n\n> **Категория:**\n\n`%v`\n\n> **Текстовый канал:**\n\n%v\n\n> **Голосовой канал:**\n\n%v\n\n*В течение `6` месяцев теперь принадлежит тебе!*", b.Cfg.EmojiPrivateChannel, c["categoryChannel"], c["textChannel"], c["voiceChannel"])
			} else {
				embed.Description += fmt.Sprintf("\n***ЧТО!? Поздравляю, они-чан!***\n\n*Тебе выпало приватное пространство* %v\n\n*Однако, оно у тебя уже есть.*\n\n*Увеличиваем срок текущего на* `6` *месяцев!*\n\n> **Категория:**\n\n`%v`\n\n> **Текстовый канал:**\n\n%v\n\n> **Голосовой канал:**\n\n%v", b.Cfg.EmojiPrivateChannel, c["categoryChannel"], c["textChannel"], c["voiceChannel"])
			}
		}
		b.openPackPanel(s, i, embed)
	case "Маркет":
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Маркет",
			Description: "```Выберите интересующий вас товар.```",
			Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/846497994702192680/AhegaoStore.gif"},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/857292340876017694/cart_buy.png"},
		},
			[]*discordgo.Button{{
				Label:    "⠀Пространства⠀",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: "Маркет Пространства",
				Emoji:    discordgo.ComponentEmoji{ID: "941315452523675699"},
			}, {
				Label:    "⠀⠀Роли⠀⠀",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: "Маркет Роли",
				Emoji:    discordgo.ComponentEmoji{ID: "941315452536234084"},
			}, {
				Label:    "⠀Реакции⠀⠀",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: "Маркет Реакции",
				Emoji:    discordgo.ComponentEmoji{ID: "941315452196495380"},
			}, {
				Label:    "⠀⠀⠀⠀Паки⠀⠀⠀ ⠀",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: "Маркет Паки",
				Emoji:    discordgo.ComponentEmoji{ID: "941315452574007296"},
			}, {
				Label:    "⠀ Лотерея",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: "Маркет Лотерея",
				Emoji:    discordgo.ComponentEmoji{ID: "941315452427178025"},
			}}, MarketPanel)
	case "Маркет Пространства":
		description := "```У пользователя в наличии может быть лишь одно пространство.```\n```Название пространства, при покупке, устанавливается на основе имени пользователя.```\n```Вы всегда можете продлить срок существования текущего пространства, что дешевле его создания.```"
		buySpaceButtonDisabled := true
		isEnoughCurrency, balance, err := b.isEnoughCurrency(userId, 15000)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
			return
		}
		if isEnoughCurrency {
			buySpaceButtonDisabled = false
		}
		textChannelID, voiceChannelID, _, channelsExpires, err := b.MarketData.HasCustomerChannels(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's channels")
			return
		}
		if textChannelID == "" {
			description += "\n> **Имеется ли в наличии пространство:**\n\n`Нет`"
			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Пространства",
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/846465333593768017/AhegaoSpaces.gif"},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718325709635614/icon_space.png"},
			},
				[]*discordgo.Button{{
					Label:    "⠀Приобрести⠀|⠀15000 т.⠀|⠀6 мес.⠀⠀",
					Style:    discordgo.SuccessButton,
					Disabled: buySpaceButtonDisabled,
					CustomID: "Приобрести Пространство",
					Emoji:    discordgo.ComponentEmoji{ID: "941476541144137819"},
				}}, DefaultPanel)
		} else {
			extendSpaceButtonDisabled := true
			if balance > 12000 {
				extendSpaceButtonDisabled = false
			}
			textChannel, err := s.Channel(textChannelID)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting text channel in guild")
				return
			}
			voiceChannel, err := s.Channel(voiceChannelID)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting voice channel in guild")
				return
			}
			category, err := s.Channel(voiceChannel.ParentID)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error getting category channel in guild")
				return
			}
			format := "02/01/06"
			description += fmt.Sprintf("\n*Найдено пространство, принадлежащее вам.*\n\n> **Категория:**\n\n`%v`\n\n> **Текстовый канал:**\n\n%v\n\n> **Голосовой канал:**\n\n%v\n\n> **Срок истечения:**\n\n`%v`", category.Name, textChannel.Mention(), voiceChannel.Mention(), time.Unix(channelsExpires, 0).Format(format))
			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Пространства",
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/846465333593768017/AhegaoSpaces.gif"},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718325709635614/icon_space.png"},
			},
				[]*discordgo.Button{{
					Label:    "⠀⠀Продлить⠀|⠀12000 т.⠀|⠀6 мес.⠀ ⠀",
					Style:    discordgo.SuccessButton,
					Disabled: extendSpaceButtonDisabled,
					CustomID: "Продлить Пространство",
					Emoji:    discordgo.ComponentEmoji{ID: "849720295861780480"},
				}}, DefaultPanel)
		}
	case "Маркет Роли":

		items, err := b.MarketData.UserOrders(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's items")
			return
		}
		description := "```В разделе \"Приобрести\" для покупки доступны лишь стандартные роли.```\n```Вы можете создать свою роль, со своим цветом и названием, нажав на кпопку \"Создать\".```"
		format := "02/01/06"
		d := ""
		for _, i := range items {
			if i.Product.Type == "role" || i.Product.Type == "customrole" {
				d += fmt.Sprintf("\n<@&%s>  —  `%s`\n", i.Product.RoleID, time.Unix(i.Expires, 0).Format(format))
			}
		}
		if d != "" {
			description += fmt.Sprintf("\n> **Роль — Срок истечения:**\n%v", d)
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Роли",
			Description: description,
			Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/846512025433538610/AhegaoRoles.gif"},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718327135567872/icon_role.png"},
		},
			[]*discordgo.Button{{
				Label:    "⠀Приобрести⠀",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: "Приобрести Роль",
				Emoji:    discordgo.ComponentEmoji{ID: "858763110285574154"},
			}, {
				Label:    "⠀Создать⠀⠀",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: "Создать Роль R|+0",
				Emoji:    discordgo.ComponentEmoji{ID: "858694330075185182"},
			}}, DefaultPanel)
	case "Маркет Реакции":

		items, err := b.MarketData.UserOrders(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's items")
			return
		}
		description := "```Вы можете применять реакции, не приобретая их навсегда, однако одноразовое использование обходится дорого.```\n```Наборы реакций являются односторонними и двухсторонними, например:```\n`/обнять` и `/обнять @пользователь`"
		d := ""
		for _, i := range items {
			if i.Product.Type == "reaction" {
				expires := time.Unix(i.Expires, 0).Format("02/01/06")
				if i.Expires == -1 {
					expires = "Навсегда"
				}
				d += fmt.Sprintf("\n**%s**  —  `%s`\n", i.Product.Name, expires)
			}
		}
		if d != "" {
			description += fmt.Sprintf("\n\n> **Реакция — Срок истечения:**\n%v", d)
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Реакции",
			Description: description,
			Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/848963272619327498/AhegaoReactions.gif"},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718328456511488/icon_reaction.png"},
		},
			[]*discordgo.Button{{
				Label:    "⠀⠀⠀⠀⠀⠀⠀⠀⠀Приобрести⠀⠀⠀⠀⠀⠀⠀ ⠀",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: "Приобрести Реакцию",
				Emoji:    discordgo.ComponentEmoji{ID: "859160144571400193"},
			}}, DefaultPanel)
	case "Маркет Паки":
		b.openPackPanel(s, i, &discordgo.MessageEmbed{
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718331191197746/icon_lootbox.png"},
			Image:     &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/861645969778540544/AhegaoPacks.gif"}})
	case "Маркет Лотерея":
		isTicketHas := "Нет"
		isTicketHasIsBtnDisabled := false
		for _, v := range b.LotterySession.Participants {
			if v == userId {
				isTicketHas = "Да"
				isTicketHasIsBtnDisabled = true
				break
			}
		}
		if len(b.LotterySession.Participants) == 10 {
			isTicketHasIsBtnDisabled = true
		}
		lotteryTimeMinutesRemaining := time.Until(b.LotterySession.Expires).Minutes()
		hours := int(lotteryTimeMinutesRemaining) / 60
		mins := int(lotteryTimeMinutesRemaining) - hours*60
		description := fmt.Sprintf("```Каждые 12 часов запускается новая лотерея, по окончанию которой между участниками разыграваются трусики в количестве: %v * число участников.``````Максимальное количество участников в одной лотерее — 10.```\n> **Приобретён ли лотерейный билет?**\n\n `%v`\n\n> **Билетов приобретено:**\n\n `%v из 10`\n\n> **До окончания текущей лотереи:**\n\n `%v:%v`", b.Cfg.LotteryPrice, isTicketHas, len(b.LotterySession.Participants), hours, mins)
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Лотерея",
			Description: description,
			Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/861645967702491136/AhegaoLottery.gif"},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718749275619428/icon_ticket.png"},
		}, []*discordgo.Button{{
			Label:    "⠀⠀⠀⠀Приобрести Билет⠀|⠀100 т.⠀⠀⠀",
			Style:    discordgo.SuccessButton,
			Disabled: isTicketHasIsBtnDisabled,
			CustomID: "Приобрести Билет",
			Emoji:    discordgo.ComponentEmoji{ID: "861217535286181908"},
		}}, DefaultPanel)
	case "Приобрести Пространство":
		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/858440892582068264/space.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440871891566612/space1.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440873083928616/space2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440877294878750/space3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440881156390922/space4.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440886320365608/space5.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440885364326410/space6.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440888601411620/space7.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440890616119336/space8.gif",
		})
		err := b.MarketData.BalanceUpdate(userId, -15000)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
			return
		}
		textChannelID, voiceChannelID, _, _, err := b.MarketData.HasCustomerChannels(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's channels")
			return
		}
		ch := make(chan map[string]string)
		defer close(ch)
		go b.addChannels(s, i, textChannelID, voiceChannelID, ch)
		c := <-ch
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Пространства",
			Description: fmt.Sprintf("*Они-чан может многое себе позволить!*\n\n*Ведь приватное пространство не из дешевых.*\n\n*Рассекаем сквозь пространство и время!*\n\n> **Категория:**\n\n`%v`\n\n> **Текстовый канал:**\n\n%v\n\n> **Голосовой канал:**\n\n%v\n\n*В течение `6` месяцев теперь принадлежит тебе!*", c["categoryChannel"], c["textChannel"], c["voiceChannel"]),
			Image:       &discordgo.MessageEmbedImage{URL: randImage},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/909130471122825276/icon_space_buy.png"},
		}, nil, DefaultPanel)
	case "Продлить Пространство":
		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/858440892582068264/space.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440871891566612/space1.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440873083928616/space2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440877294878750/space3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440881156390922/space4.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440886320365608/space5.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440885364326410/space6.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440888601411620/space7.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/858440890616119336/space8.gif",
		})
		err := b.MarketData.BalanceUpdate(userId, -12000)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
			return
		}
		textChannelID, voiceChannelID, _, _, err := b.MarketData.HasCustomerChannels(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's channels")
			return
		}
		ch := make(chan map[string]string)
		defer close(ch)
		go b.addChannels(s, i, textChannelID, voiceChannelID, ch)
		c := <-ch
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Пространства",
			Description: fmt.Sprintf("*Созерцая миры, расширяем границы вместе с братиком.*\n\n> **Категория:**\n\n`%v`\n\n> **Текстовый канал:**\n\n%v\n\n> **Голосовой канал:**\n\n%v\n\n*Увеличиваем срок существования текущего пространства на `6` месяцев* %v", c["categoryChannel"], c["textChannel"], c["voiceChannel"], b.Cfg.EmojiPrivateChannel),
			Image:       &discordgo.MessageEmbedImage{URL: randImage},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/909130471122825276/icon_space_buy.png"},
		}, nil, DefaultPanel)
	case "Приобрести Роль":
		var rolesID []string
		var isHasProducts map[string]bool = make(map[string]bool)
		rolesList := ""
		allRoles, err := b.MarketData.Products()
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get products")
		}
		for _, p := range allRoles {
			if p.Type == "role" {
				rolesID = append(rolesID, p.RoleID)
				rolesList += fmt.Sprintf("\n\n<@&%v>", p.RoleID)
			}
		}
		ps, err := b.MarketData.UserOrders(userId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's orders")
		}
		for _, p := range ps {
			if p.Product.Type == "role" {
				isHasProducts[p.Product.RoleID] = true
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to determine whether user has a product")
					return
				}
			}
		}
		buyRoleButtonDisabled := true
		isEnoughCurrency, _, err := b.isEnoughCurrency(userId, 2700)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
			return
		}
		if isEnoughCurrency {
			buyRoleButtonDisabled = false
		}
		btnsRoleBuy := []*discordgo.Button{}
		for i := 1; i < len(rolesID)+1; i++ {
			label := ""
			emoji := ""
			customID := ""
			if isHasProducts[rolesID[i-1]] {
				label = "⠀Продлить⠀"
				customID = "Продлить"
			} else {
				label = "⠀⠀Купить⠀⠀"
				customID = "Купить"
			}
			switch i {
			case 1:
				emoji = "♟️"
			case 2:
				emoji = "🌸"
			case 3:
				emoji = "❄️"
			case 4:
				emoji = "🔫"
			case 5:
				emoji = "🍾"
			case 6:
				emoji = "🐾"
			case 7:
				emoji = "🗡️"
			case 8:
				emoji = "🤪"
			case 9:
				emoji = "⭐"
			case 10:
				emoji = "🧬"
			case 11:
				emoji = "🏵️"
			case 12:
				emoji = "🧻"
			case 13:
				emoji = "🌈"
			case 14:
				emoji = "💥"
			case 15:
				emoji = "💊"
			}

			btnsRoleBuy = append(btnsRoleBuy, &discordgo.Button{
				Label:    label,
				Style:    discordgo.SuccessButton,
				Disabled: buyRoleButtonDisabled,
				CustomID: fmt.Sprintf(customID+" Роль %v", strconv.FormatInt(int64(i), 10)),
				Emoji:    discordgo.ComponentEmoji{Name: emoji},
			})
		}
		description := fmt.Sprintf("**Стоимость всех стандартных ролей условно едина —\n\n** `2700` %v на срок `1` недели.\n\n\n> **Ассортимент ролей:**%v", b.Cfg.EmojiRole, rolesList)
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Роли",
			Description: description,
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858764424160608306/icon_role_buy.png"},
		}, btnsRoleBuy, DefaultPanel)
	case "Приобрести Реакцию":

		var reactionsNames []string
		var isHasProducts map[string]bool = make(map[string]bool)
		reactionsList := ""
		ps, err := b.MarketData.Products()
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get products")
		}
		j := 0
		for _, p := range ps {
			if p.Type == "reaction" {
				j++
				reactionsNames = append(reactionsNames, p.Name)
				reactionsList += fmt.Sprintf("\n\n**%v.** `%v` %v", j, p.Name, b.Cfg.EmojiReaction)
				hasProduct, err := b.MarketData.HasProduct(userId, p.Dial)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to determine whether user has a product")
					return
				}
				if hasProduct {
					isHasProducts[p.Name] = true
				}
			}
		}
		buyReactionButtonDisabled := true
		isEnoughCurrency, _, err := b.isEnoughCurrency(userId, 1350)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
			return
		}
		btnsRoleBuy := []*discordgo.Button{}
		for i := 1; i < len(reactionsNames)+1; i++ {
			buyReactionButtonDisabled = false
			emojiID := ""
			if isHasProducts[reactionsNames[i-1]] || !isEnoughCurrency {
				buyReactionButtonDisabled = true
			}
			switch i {
			case 1:
				emojiID = "859212034515533864"
			case 2:
				emojiID = "859212034528903210"
			case 3:
				emojiID = "859212034662727680"
			case 4:
				emojiID = "859212034338717747"
			case 5:
				emojiID = "859212034095710239"
			case 6:
				emojiID = "859212034712141834"
			case 7:
				emojiID = "859212034578055168"
			case 8:
				emojiID = "859212034604400690"
			case 9:
				emojiID = "859212217274859531"
			case 10:
				emojiID = "859210540398936094"
			case 11:
				emojiID = "859210540023152711"
			case 12:
				emojiID = "859210540395921459"
			case 13:
				emojiID = "859210540689260564"
			case 14:
				emojiID = "859210540517556254"
			case 15:
				emojiID = "859210540487147531"
			case 16:
				emojiID = "859210540458442772"
			case 17:
				emojiID = "859210540441927750"
			case 18:
				emojiID = "859210540475088916"
			case 19:
				emojiID = "859210540499992606"
			case 20:
				emojiID = "859210540491210752"
			case 21:
				emojiID = "859210540136136715"
			case 22:
				emojiID = "859210540491210772"
			case 23:
				emojiID = "859210540520570890"
			case 24:
				emojiID = "859210540231950337"

			}

			btnsRoleBuy = append(btnsRoleBuy, &discordgo.Button{
				Label:    "",
				Style:    discordgo.SuccessButton,
				Disabled: buyReactionButtonDisabled,
				CustomID: fmt.Sprintf("Реакция %v", strconv.FormatInt(int64(i)+24, 10)),
				Emoji:    discordgo.ComponentEmoji{ID: emojiID},
			})
		}
		description := fmt.Sprintf("**Стоимость всех реакций условно едина —\n\n** `1350` %v `Навсегда`.\n\n\n> **Ассортимент реакций:**%v", b.Cfg.EmojiReaction, reactionsList)
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Реакции",
			Description: description,
			Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/848963272619327498/AhegaoReactions.gif"},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/859160231376453652/icon_reaction_buy.png"},
		}, btnsRoleBuy, ReactionPanel)
	case "Приобрести Билет":
		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/861315613221453834/lottery.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/861316510928994334/lottery1.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/861315631113830400/lottery2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/861315632423632916/lottery3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/861315643160920104/lottery5.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/861316162352971776/lottery6.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/861321268782825472/lottery7.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/940201422945914920/lottery8.gif",
		})
		b.LotterySession.Participants = append(b.LotterySession.Participants, userId)
		descriptionAdditional := "*Но пока, кроме тебя, они-чан, никто не купил билет.\n\nНадо подождать других!*"
		if len(b.LotterySession.Participants) > 1 {
			descriptionAdditional = fmt.Sprintf("*В данный момент у вас `1` шанс из `%v`.\n\nВозможно, удача улыбнётся тебе, они-чан!*", len(b.LotterySession.Participants))
		}
		err := b.MarketData.BalanceUpdate(userId, -b.Cfg.LotteryPrice)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
		}
		openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
			Title:       "Лотерея",
			Description: fmt.Sprintf("*Поздравляю, с приобретением %v, теперь вы участник лотереи!*\n\n%v", b.Cfg.EmojiLottery, descriptionAdditional),
			Image:       &discordgo.MessageEmbedImage{URL: randImage},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/861310418735726622/icon_ticket_buy.png"},
		}, nil, DefaultPanel)
	default:
		if strings.Contains(customId, "Скрыть") || strings.Contains(customId, "Показать") {
			action := customId[:strings.Index(customId, " ")]
			productDial, _ := strconv.Atoi(customId[strings.LastIndex(customId, " ")+1:])
			product, err := b.MarketData.Product(productDial)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get info about product")
				return
			}
			isHide := true
			description := fmt.Sprintf("*Братик спрятал от всех роль* <@&%v>*, обязательно проверь отображаемый цвет никнейма!*", product.RoleID)
			if action == "Показать" {
				err = b.roleGive(s, i.GuildID, userId, product.RoleID)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Error giving role")
					return
				}
				isHide = false
				description = fmt.Sprintf("*Братик успешно раскрыл всем роль* <@&%v>*, обязательно проверь отображаемый цвет никнейма!*", product.RoleID)
			} else {
				err = s.GuildMemberRoleRemove(b.Cfg.GuildID, userId, product.RoleID)
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove role from the user")
					return
				}
				b.roleSweepDelimiters(s, b.Cfg.GuildID, userId)
			}
			err = b.MarketData.IsHiddenUpdate(userId, isHide, productDial)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to hide user's role")
				return
			}

			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Управление Ролями",
				Description: description,
			}, nil, DefaultPanel)

		} else if strings.Contains(customId, "Создать Роль") {
			colorAndOperation := strings.Split(strings.ReplaceAll(customId, "Создать Роль ", ""), "|")
			red := 0
			green := 0
			blue := 0
			if i.Message.Embeds[0].Fields != nil {
				red, _ = strconv.Atoi(strings.ReplaceAll(i.Message.Embeds[0].Fields[0].Value, "`", ""))
				green, _ = strconv.Atoi(strings.ReplaceAll(i.Message.Embeds[0].Fields[1].Value, "`", ""))
				blue, _ = strconv.Atoi(strings.ReplaceAll(i.Message.Embeds[0].Fields[2].Value, "`", ""))
			}
			if strings.Contains(colorAndOperation[1], "+") {
				value, _ := strconv.Atoi(strings.ReplaceAll(colorAndOperation[1], "+", ""))
				switch colorAndOperation[0] {
				case "R":
					red += value
				case "G":
					green += value
				case "B":
					blue += value
				}
			} else {
				value, _ := strconv.Atoi(strings.ReplaceAll(colorAndOperation[1], "-", ""))
				switch colorAndOperation[0] {
				case "R":
					red -= value
				case "G":
					green -= value
				case "B":
					blue -= value
				}
			}
			operation := ""
			btnsRoleCreate := []*discordgo.Button{}
			for j := 5; j < 20; j++ {
				btnRgbDisabled := false
				var style discordgo.ButtonStyle
				emojiId := ""
				if j < 10 {
					style = discordgo.DangerButton
					operation = "R"
				} else if j < 15 {
					style = discordgo.SuccessButton
					operation = "G"
				} else {
					style = discordgo.PrimaryButton
					operation = "B"
				}
				switch j % 5 {
				case 0:
					emojiId = "862080160962052096"
					if red+50 > 255 && operation == "R" {
						btnRgbDisabled = true
					} else if green+50 > 255 && operation == "G" {
						btnRgbDisabled = true
					} else if blue+50 > 255 && operation == "B" {
						btnRgbDisabled = true
					}
					operation += "|+50"
				case 1:
					emojiId = "862080160970571786"
					if red+10 > 255 && operation == "R" {
						btnRgbDisabled = true
					} else if green+10 > 255 && operation == "G" {
						btnRgbDisabled = true
					} else if blue+10 > 255 && operation == "B" {
						btnRgbDisabled = true
					}
					operation += "|+10"
				case 2:
					emojiId = "862080160991936543"
					if red+1 > 255 && operation == "R" {
						btnRgbDisabled = true
					} else if green+1 > 255 && operation == "G" {
						btnRgbDisabled = true
					} else if blue+1 > 255 && operation == "B" {
						btnRgbDisabled = true
					}
					operation += "|+1"
				case 3:
					emojiId = "862080161079099422"
					if red-10 < 0 && operation == "R" {
						btnRgbDisabled = true
					} else if green-10 < 0 && operation == "G" {
						btnRgbDisabled = true
					} else if blue-10 < 0 && operation == "B" {
						btnRgbDisabled = true
					}
					operation += "|-10"
				case 4:
					emojiId = "862080160776978473"
					if red-50 < 0 && operation == "R" {
						btnRgbDisabled = true
					} else if green-50 < 0 && operation == "G" {
						btnRgbDisabled = true
					} else if blue-50 < 0 && operation == "B" {
						btnRgbDisabled = true
					}
					operation += "|-50"
				}
				btnsRoleCreate = append(btnsRoleCreate, &discordgo.Button{
					Label:    "",
					Style:    style,
					Disabled: btnRgbDisabled,
					CustomID: fmt.Sprintf("Создать Роль %v", operation),
					Emoji:    discordgo.ComponentEmoji{ID: emojiId},
				})
			}
			btnsRoleCreate = append(btnsRoleCreate, &discordgo.Button{
				Label:    "⠀⠀⠀>>⠀⠀⠀⠀",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: "Создать СрокРоль",
				Emoji:    discordgo.ComponentEmoji{Name: "✔️"},
			})
			var rgb int = (red << 16) | (green << 8) | blue
			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Роли",
				Color:       rgb,
				Description: "*Установите цвет для роли, руководствуясь цветовой моделью [RGB](https://www.google.com/search?q=Выбор+Цвета).*\n\n*Предварительный цвет роли отображается слева, в виде вертикальной линии.*",
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858704257402404904/icon_role_create.png"},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "🔴",
						Value:  fmt.Sprintf("`%v`", red),
						Inline: true,
					},
					{
						Name:   "🟢",
						Value:  fmt.Sprintf("`%v`", green),
						Inline: true,
					},
					{
						Name:   "🔵",
						Value:  fmt.Sprintf("`%v`", blue),
						Inline: true,
					},
				}}, btnsRoleCreate, CreateRolePanel)
		} else if strings.Contains(customId, "Создать СрокРоль") {
			randImage := randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/859026312710848512/time.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026322214617088/time2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026326123315200/time3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026276803411968/time4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026278647595018/time5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026291138887700/time7.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026297380536360/time8.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026300265431040/time9.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026304190775316/time10.gif",
			})

			createRoleUsers = append(createRoleUsers, CreateRoleUser{UserId: userId, RoleColor: i.Message.Embeds[0].Color})

			createRoleButtonOneDayDisabled, createRoleButtonThreeDaysDisabled, createRoleButtonSevenDaysDisabled, createRoleButtonOneMouthDisabled := true, true, true, true
			isEnoughCurrency, balance, err := b.isEnoughCurrency(userId, 1500)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
				return
			}
			if isEnoughCurrency {
				createRoleButtonOneDayDisabled = false
			}
			if balance > 3900 {
				createRoleButtonThreeDaysDisabled = false
			}
			if balance > 8100 {
				createRoleButtonSevenDaysDisabled = false
			}
			if balance > 30000 {
				createRoleButtonOneMouthDisabled = false
			}

			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Роли",
				Description: "*Выберите, на какой срок вы хотите приобрести кастомную роль.*",
				Image:       &discordgo.MessageEmbedImage{URL: randImage},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858704257402404904/icon_role_create.png"},
			}, []*discordgo.Button{{
				Label:    "⠀1 д.⠀|⠀1500 т.",
				Style:    discordgo.SuccessButton,
				Disabled: createRoleButtonOneDayDisabled,
				CustomID: "Создать ИмяРоль 1",
				Emoji:    discordgo.ComponentEmoji{ID: "858694330075185182"},
			}, {
				Label:    "⠀3 д.⠀| 3900 т.",
				Style:    discordgo.SuccessButton,
				Disabled: createRoleButtonThreeDaysDisabled,
				CustomID: "Создать ИмяРоль 3",
				Emoji:    discordgo.ComponentEmoji{ID: "858694330075185182"},
			}, {
				Label:    "⠀7 д.⠀|⠀8100 т.",
				Style:    discordgo.SuccessButton,
				Disabled: createRoleButtonSevenDaysDisabled,
				CustomID: "Создать ИмяРоль 7",
				Emoji:    discordgo.ComponentEmoji{ID: "858694330075185182"},
			}, {
				Label:    "⠀1 м.⠀|⠀30000 т.",
				Style:    discordgo.SuccessButton,
				Disabled: createRoleButtonOneMouthDisabled,
				CustomID: "Создать ИмяРоль 31",
				Emoji:    discordgo.ComponentEmoji{ID: "858694330075185182"},
			},
			}, DefaultPanel)
		} else if strings.Contains(customId, "Создать ИмяРоль") {
			dur, _ := strconv.Atoi(customId[strings.LastIndex(customId, " ")+1:])
			for i := 0; i < len(createRoleUsers); i++ {
				if createRoleUsers[i].UserId == userId {
					createRoleUsers[i].Duration = int64(dur) * dayAsSec
					break
				}
			}
			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Роли",
				Color:       i.Message.Embeds[0].Color,
				Description: "*Установите имя для роли.*\n\n**Ограничения:**\n\n> *Максимальное количество символов:* `100`\n\n> *Кастомные эмодзи не приемлемы, используйте эмодзи и символы [Юникода](https://unicode-table.com/ru/).*\n\n*Ожидаю ввода названия роли, в виде текстового сообщения.*",
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858704257402404904/icon_role_create.png"},
			}, nil, DefaultPanel)

		} else if strings.Contains(customId, "Количество Паков") {
			randImage := randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/854813534302109756/buy_pack.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/854813570457927691/pack_open1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/854813535879430154/buy_pack2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/854813535892275210/buy_pack3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/854813531127676938/buy_pack4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/854813532902260806/buy_pack6.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/854813531399520256/buy_pack7.gif",
			})

			packsAmount, _ := strconv.Atoi(customId[strings.LastIndex(customId, " ")+1:])
			err := b.MarketData.LootboxAdd(userId, packsAmount, b.Cfg.LootboxPrice)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's lootbox amount")
				return
			}
			b.openPackPanel(s, i, &discordgo.MessageEmbed{
				Description: "*Когда братик будет открывать, ему обязательно повезёт!*",
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/854114891311874048/icon_lootbox_hq.png"},
				Image:       &discordgo.MessageEmbedImage{URL: randImage},
			})
		} else if strings.Contains(customId, "Купить Роль") || strings.Contains(customId, "Продлить Роль") {
			randImage := randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/859026312710848512/time.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026322214617088/time2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026326123315200/time3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026276803411968/time4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026278647595018/time5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026291138887700/time7.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026297380536360/time8.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026300265431040/time9.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859026304190775316/time10.gif",
			})

			productDial, _ := strconv.Atoi(customId[strings.LastIndex(customId, " ")+1:])
			product, err := b.MarketData.Product(productDial)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get info about product")
				return
			}

			buyRoleButtonOneDayDisabled, buyRoleButtonThreeDaysDisabled, buyRoleButtonSevenDaysDisabled, buyRoleButtonOneMouthDisabled := true, true, true, true
			isEnoughCurrency, balance, err := b.isEnoughCurrency(userId, 500)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user balance")
				return
			}
			if isEnoughCurrency {
				buyRoleButtonOneDayDisabled = false
			}
			if balance > 1300 {
				buyRoleButtonThreeDaysDisabled = false
			}
			if balance > 2700 {
				buyRoleButtonSevenDaysDisabled = false
			}
			if balance > 10000 {
				buyRoleButtonOneMouthDisabled = false
			}

			buyOrExtend := customId[:strings.Index(customId, " ")+1]
			description := "*Выберите, на какой срок вы хотите приобрести роль.*"
			if strings.Contains(buyOrExtend, "Продлить") {
				description = "*Выберите, на какой срок вы хотите продлить роль.*"
			}

			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Роли",
				Description: fmt.Sprintf("%v\n\n> **Роль:**\n\n<@&%v>\n", description, product.RoleID),
				Image:       &discordgo.MessageEmbedImage{URL: randImage},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858764424160608306/icon_role_buy.png"},
			}, []*discordgo.Button{{
				Label:    "⠀1 д.⠀|⠀500 т.",
				Style:    discordgo.SuccessButton,
				Disabled: buyRoleButtonOneDayDisabled,
				CustomID: fmt.Sprintf("%vСрок Роль %v|1", buyOrExtend, productDial),
				Emoji:    discordgo.ComponentEmoji{ID: "858763110285574154"},
			}, {
				Label:    "⠀3 д.⠀|⠀1300 т.",
				Style:    discordgo.SuccessButton,
				Disabled: buyRoleButtonThreeDaysDisabled,
				CustomID: fmt.Sprintf("%vСрок Роль %v|3", buyOrExtend, productDial),
				Emoji:    discordgo.ComponentEmoji{ID: "858763110285574154"},
			}, {
				Label:    "⠀7 д.⠀|⠀2700 т.",
				Style:    discordgo.SuccessButton,
				Disabled: buyRoleButtonSevenDaysDisabled,
				CustomID: fmt.Sprintf("%vСрок Роль %v|7", buyOrExtend, productDial),
				Emoji:    discordgo.ComponentEmoji{ID: "858763110285574154"},
			}, {
				Label:    "⠀1 м.⠀|⠀10000 т.",
				Style:    discordgo.SuccessButton,
				Disabled: buyRoleButtonOneMouthDisabled,
				CustomID: fmt.Sprintf("%vСрок Роль %v|31", buyOrExtend, productDial),
				Emoji:    discordgo.ComponentEmoji{ID: "858763110285574154"},
			},
			}, DefaultPanel)
		} else if strings.Contains(customId, "Срок Роль") {
			randImage := randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/858703916555698196/role.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703919718989844/role1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703921458970645/role2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703924038598686/role3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703925662056508/role4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858704180276625438/role6.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703907211706388/role7.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703908607361024/role8.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703909321441290/role9.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/858703914077519872/role10.gif",
			})

			dialAndDur := customId[strings.LastIndex(customId, " ")+1:]
			productDial, _ := strconv.Atoi(dialAndDur[:strings.Index(dialAndDur, "|")])
			productDurDays, _ := strconv.Atoi(dialAndDur[strings.Index(dialAndDur, "|")+1:])
			product, err := b.MarketData.Product(productDial)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get info about product")
				return
			}

			switch productDurDays {
			case 1:
				product.Price = 500
			case 3:
				product.Price = 1300
			case 7:
				product.Price = 2700
			case 31:
				product.Price = 10000
			}

			err = b.MarketData.Order(userId, product, int64(productDurDays)*dayAsSec)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to order product")
				return
			}

			buyOrExtend := customId[:strings.Index(customId, " ")+1]
			description := fmt.Sprintf("*Поздравляю, теперь у вас в наличии новая роль —\n\n<@&%v> на `%v` д.*\n\n*Сверкайте ею, как никогда.*", product.RoleID, productDurDays)
			if strings.Contains(buyOrExtend, "Продлить") {
				description = fmt.Sprintf("*Поздравляю, вы продлили свою роль —\n\n<@&%v> на `%v` д.*\n\n*Продолжайте демонстрировать её всем.*", product.RoleID, productDurDays)
			}

			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Роли",
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: randImage},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/858764424160608306/icon_role_buy.png"},
			}, nil, DefaultPanel)

			err = b.roleGive(s, i.GuildID, userId, product.RoleID)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Error giving role")
				return
			}
		} else if strings.Contains(customId, "Реакция") {

			randImage := randSlice([]string{
				"https://cdn.discordapp.com/attachments/838477787877998624/859433907014860830/reaction.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433917207937034/reaction1.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433850107461692/reaction2.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433857082327051/reaction3.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433861107941416/reaction4.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433873267097620/reaction5.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433877528903711/reaction6.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433881471811614/reaction7.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433889473495043/reaction8.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433895531773972/reaction9.gif",
				"https://cdn.discordapp.com/attachments/838477787877998624/859433903969402920/reaction10.gif",
			})

			productDial, _ := strconv.Atoi(customId[strings.LastIndex(customId, " ")+1:])

			product, err := b.MarketData.Product(productDial)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get info about product")
				return
			}

			err = b.MarketData.Order(userId, product, -1)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to order product")
				return
			}

			openCustomPanel(s, i, discordgo.InteractionResponseUpdateMessage, &discordgo.MessageEmbed{
				Title:       "Реакции",
				Description: fmt.Sprintf("*Поздравляю, теперь у вас в наличии новая реакция —\n\n%v `%v` навсегда!*\n\n*Попробуйте испытать на ком-нибудь.*", b.Cfg.EmojiReaction, product.Name),
				Image:       &discordgo.MessageEmbedImage{URL: randImage},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/859160231376453652/icon_reaction_buy.png"},
			}, nil, DefaultPanel)
		}
	}
}

// OnGuildMemberJoined registers a member when they join the guild. If member rejoined, would reassign roles owned by them.
func (b *Bot) OnGuildMemberJoined(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	userId := m.User.ID
	err := b.MarketData.CustomerAdd(userId)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to add user to the database")
		return
	}

	is, err := b.MarketData.UserOrders(userId)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user's items")
	}

	for _, i := range is {
		if i.Product.Type == "role" {
			err = b.roleGive(s, m.GuildID, userId, i.Product.RoleID)
			if err != nil {
				s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign role to a user")
				return
			}
		}
	}

	textChannelId, _, roleId, _, _ := b.MarketData.HasCustomerChannels(userId)
	if textChannelId != "" {
		err = s.GuildMemberRoleAdd(m.GuildID, userId, roleId)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign role to a user")
			return
		}
		err = s.GuildMemberRoleAdd(m.GuildID, userId, "836937639831142401")
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Error adding role")
			return
		}
	}
}

// OnGuildMemberVoiceUpd tracks user's voice activity and adds currency to the balance.
// Amount depends on the voice farm rate.
func (b *Bot) OnGuildMemberVoiceUpd(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	ID := m.UserID
	user, err := s.GuildMember(b.Cfg.GuildID, ID)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get user info")
		return
	}

	if user.User.Bot {
		return
	}

	if m.BeforeUpdate == nil {
		t := time.Now().Unix()
		err := b.MarketData.VoiceLogAdd(ID, t)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to insert voice log")
		}
		return
	}

	if m.ChannelID == "" {
		t, err := b.MarketData.VoiceLogGetAndClean(ID)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to sweep voice log")
			return
		}

		deltaM := time.Since(time.Unix(t, 0))
		farmed := int(deltaM.Minutes()) / b.Cfg.VoiceDelta * b.Cfg.VoiceRate
		err = b.MarketData.BalanceUpdate(ID, farmed)
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's balance")
			return
		}

		err = b.MarketData.CustomerTotalVoiceUpdate(ID, int64(deltaM.Seconds()))
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to update user's total voice time")
			return
		}
	}
}

// SweepExpiredOrders removes expired orders from the database and unassignes expired roles from the users.
func (b *Bot) SweepExpiredOrders(s *discordgo.Session, quit chan struct{}) {
	t := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-t.C:
				textChannels, voiceChannels, rolesDelete, usersID, err := b.MarketData.CustomerChannelsExpired()
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to delete expired channels")
					return
				}
				for _, v := range textChannels {
					ch, err := s.Channel(v)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to get channel from guild")
						return
					}
					_, err = s.ChannelDelete(ch.ParentID)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to delete channel from guild")
						return
					}
					_, err = s.ChannelDelete(v)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to delete channel from guild")
						return
					}
				}
				for _, v := range voiceChannels {
					_, err = s.ChannelDelete(v)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to delete channel from guild")
						return
					}
				}
				for _, v := range rolesDelete {
					err = s.GuildRoleDelete(b.Cfg.GuildID, v)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove role from the user")
						return
					}
				}
				for _, v := range usersID {
					b.roleSweepDelimiters(s, b.Cfg.GuildID, v)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove delimeter from the user")
						return
					}
				}
				es, err := b.MarketData.OrdersDeleteExpired()
				if err != nil {
					s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to perform sweep and get expired roles")
					return
				}
				var u []string
				for _, e := range es {

					for _, v := range u {
						if v != e.UserID {
							u = append(u, e.UserID)
						}
					}

					err = s.GuildMemberRoleRemove(b.Cfg.GuildID, e.UserID, e.RoleID)
					if err != nil {
						s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove role from the user")
					}
				}

				for _, userID := range u {
					b.roleSweepDelimiters(s, b.Cfg.GuildID, userID)
				}

			case <-quit:
				t.Stop()
				return
			}
		}
	}()
}

// uRIP gets all user roles beetwen specified positions.
func uRIP(s *discordgo.Session, guildID, userID, lowID, highID string) (discordgo.Roles, error) {
	var r discordgo.Roles
	g, err := s.Guild(guildID)
	if err != nil {
		return r, err
	}

	sort.Slice(g.Roles, func(i, j int) bool {
		return g.Roles[i].Position <= g.Roles[j].Position
	})

	var rolesInRegion discordgo.Roles
	inBound := false
	for _, role := range g.Roles {
		if role.ID == highID {
			break
		}

		if inBound {
			rolesInRegion = append(rolesInRegion, role)

		}
		if role.ID == lowID {
			inBound = true
		}
	}

	gm, err := s.GuildMember(guildID, userID)
	if err != nil {
		return r, err
	}

	for _, roleID := range gm.Roles {
		for _, role := range rolesInRegion {
			if role.ID == roleID {
				r = append(r, role)
			}
		}
	}
	return r, nil
}

func (b *Bot) isEnoughCurrency(userID string, comparedProductPrice int) (bool, int, error) {
	balance, err := b.MarketData.Balance(userID)
	if err != nil {
		return false, balance, err
	}
	if balance < comparedProductPrice {
		return false, balance, nil
	}
	return true, balance, nil
}

// LotteryInit initializes lottery and makes it run periodically in a goroutine.
func (b *Bot) LotteryInit(s *discordgo.Session, quit chan struct{}) {

	if b.LotterySession != nil && len(b.LotterySession.Participants) != 0 {

		randImage := randSlice([]string{
			"https://cdn.discordapp.com/attachments/838477787877998624/908317282227347536/winner.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/908317277865275452/winner2.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/908317273893249044/winner3.gif",
			"https://cdn.discordapp.com/attachments/838477787877998624/908317277634560060/winner4.gif",
		})

		randWinner := randSlice(b.LotterySession.Participants)
		b.MarketData.BalanceUpdate(randWinner, len(b.LotterySession.Participants)*b.Cfg.LotteryPrice)

		s.ChannelMessageSendEmbed(b.Cfg.LotteryChannelID,
			&discordgo.MessageEmbed{
				Color:       3092790,
				Title:       "Лотерея",
				Image:       &discordgo.MessageEmbedImage{URL: randImage},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718749275619428/icon_ticket.png"},
				Description: fmt.Sprintf("*Юххуу!* 🎉\n\n<@!%v> выиграл в лотерее, став обладателем `%v` чистых трусиков %v", randWinner, len(b.LotterySession.Participants)*b.Cfg.LotteryPrice, b.Cfg.Currency),
			})
	}

	t := time.NewTicker(12 * time.Hour)

	b.LotterySession = &Lottery{Participants: []string{}, Expires: time.Now().Add(time.Hour * 12)}
	s.ChannelMessageSendEmbed(b.Cfg.LotteryChannelID,
		&discordgo.MessageEmbed{
			Color:     3092790,
			Title:     "Лотерея",
			Image:     &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/838477787877998624/861645967702491136/AhegaoLottery.gif"},
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://cdn.discordapp.com/attachments/838477787877998624/849718749275619428/icon_ticket.png"},
			Description: fmt.Sprintf("```Лотерейные билет можно приобрести через панель в разделе \"Маркет\".```\n*У тебя есть 12 часов, чтобы купить билет.*\n\n*Поторопись, %v всего `10`.*",
				b.Cfg.EmojiLottery),
		})

	go func() {
		for {
			select {
			case <-t.C:
				t.Stop()
				b.LotteryInit(s, quit)
				return
			case <-quit:
				t.Stop()
				return
			}
		}
	}()
}

// roleGive assigns the role to the user and gives the delitmeter role.
func (b *Bot) roleGive(s *discordgo.Session, guildID, userID, roleID string) error {
	err := s.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the role to the user")
		return err
	}

	err = s.GuildMemberRoleAdd(guildID, userID, "839970605724860416") // Delitmiter role
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to assign the delimiter role to the user")
		return err
	}
	return nil
}

// roleSweepDelimiters removes the delimiter role if user doesn't have any roles belonging to the delimiter region.
func (b *Bot) roleSweepDelimiters(s *discordgo.Session, guildID, userID string) {

	rs, err := uRIP(s, guildID, userID, "839970605724860416", "834432874945839160")
	if err != nil {
		s.ChannelMessageSend(b.Cfg.AlertChannel, "uRIP failed")
		return
	}
	if len(rs) < 1 {
		err := s.GuildMemberRoleRemove(guildID, userID, "839970605724860416") //Delitmiter role
		if err != nil {
			s.ChannelMessageSend(b.Cfg.AlertChannel, "Unable to remove the delimiter role to the user")
			return
		}
	}
}

// randSlice returns random element form the slice.
func randSlice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(s))
	return s[n]
}
