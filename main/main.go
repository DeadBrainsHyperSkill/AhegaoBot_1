package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"discordmarket"

	"github.com/bwmarrin/discordgo"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	file, err := ioutil.ReadFile("../config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading config file %v", err)
		os.Exit(1)
	}

	cfg := &discordmarket.Config{}
	err = json.Unmarshal(file, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading config file %v", err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("Unable to connect to the database: %s", err.Error())
	}
	defer dbpool.Close()

	b := &discordmarket.Bot{
		Cfg:        cfg,
		MarketData: discordmarket.Database{Conn: dbpool},
	}

	dg, err := discordgo.New(cfg.Token)
	if err != nil {
		fmt.Println("Unable to connect to the bot")
	}

	dg.AddHandler(b.OnGuildMemberJoined)
	dg.AddHandler(b.OnGuildMemberVoiceUpd)
	dg.AddHandler(b.OnInteractionAct)
	dg.AddHandler(b.OnMessegeReceive)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)

	commands := []*discordgo.ApplicationCommand{{
		Name:        "панель",
		Description: "Индивидуальный интерфейс для взаимодействия с окружением.",
	},
		{
			Name:        "реакция",
			Description: "Применить реакцию.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "название",
					Description: "Выберите название реакции. Если реакция не приобретена, возникнет ошибка взаимодействия.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Дразнить",
							Value: 25,
						},
						{
							Name:  "Кусь",
							Value: 26,
						},
						{
							Name:  "На колени",
							Value: 27,
						},
						{
							Name:  "Оружие",
							Value: 28,
						},
						{
							Name:  "Отброс",
							Value: 29,
						},
						{
							Name:  "Отлизать",
							Value: 30,
						},
						{
							Name:  "Приветик",
							Value: 31,
						},
						{
							Name:  "Спать",
							Value: 32,
						},
						{
							Name:  "Суплекс",
							Value: 33,
						},
						{
							Name:  "Трусики",
							Value: 34,
						},
						{
							Name:  "Хз",
							Value: 35,
						},
					},
					Type: discordgo.ApplicationCommandOptionInteger,
				}, {
					Name:        "пользователь",
					Description: "Укажите, кому адресовать реакцию. Если не указано, то реакция применяется на всех.",
					Type:        discordgo.ApplicationCommandOptionUser,
				}}},
	}
	for _, cmnd := range commands {
		_, err := dg.ApplicationCommandCreate(b.Cfg.AppID, b.Cfg.GuildID, cmnd)
		if err != nil {
			dg.ChannelMessageSend(b.Cfg.AlertChannel, fmt.Sprintf("Cannot create '%v' command: %v", cmnd.Name, err))
		}
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Printf("Unable to open connection. Error: %s", err.Error())
	}
	defer dg.Close()

	quit := make(chan struct{})
	defer close(quit)
	b.LotteryInit(dg, quit)
	b.SweepExpiredOrders(dg, quit)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
