package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func newDiscordSession() *discordgo.Session {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + Conf.DiscordToken)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}

	return session
}

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

	imageURL := fmt.Sprintf("http://staticns.ankama.com/dofus/www/game/items/200/%s.w75h75.png", todayAlmanax.ItemImage)
	msgImage := &discordgo.MessageEmbedImage{
		URL:    imageURL,
		Width:  75,
		Height: 75,
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

	session := newDiscordSession()
	defer session.Close()

	log.Println("Sending daily Almanax message to Discord...")
	message := formatAlmanaxDailyMessage(todayAlmanax)
	_, err := session.ChannelMessageSendEmbed(Conf.DiscordChannel, message)
	if err != nil {
		log.Fatal(err)
	}
}
