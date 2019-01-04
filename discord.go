package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func formatAlmanaxDailyMessage(todayAlmanax AlmanaxEvent) *discordgo.MessageEmbed {
	msgFieldEffect := &discordgo.MessageEmbedField{
		Name:   "Effet",
		Value:  todayAlmanax.Effect,
		Inline: false,
	}

	msgFieldOffering := &discordgo.MessageEmbedField{
		Name:   "Offrande",
		Value:  todayAlmanax.Offering,
		Inline: false,
	}

	msgFields := []*discordgo.MessageEmbedField{msgFieldOffering, msgFieldEffect}
	msgImage := &discordgo.MessageEmbedImage{
		URL:    "http://staticns.ankama.com/comm/news/dofus/www/08_2012/carrou-almanax1.jpg",
		Width:  20,
		Height: 20,
	}

	message := &discordgo.MessageEmbed{
		Title:  todayAlmanax.Quest,
		Image:  msgImage,
		Fields: msgFields,
	}

	return message
}

func SendDailyAlmanaxMessage() {
	todayAlmanax := getDailyAlmanax()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Conf.DiscordToken)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}
	defer dg.Close()

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	message := formatAlmanaxDailyMessage(todayAlmanax)
	_, err = dg.ChannelMessageSendEmbed(Conf.DiscordChannel, message)
	if err != nil {
		log.Fatal(err)
	}

}
