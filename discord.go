package main

import (
	"fmt"
	"log"
	"time"

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

func cronRSSNews() {
	var newsToSend []*FeedItem
	lastNews := getLastNews(10)
	now := time.Now()

	for _, new := range lastNews {
		if now.Sub(*new.PublishedParsed) < time.Second*60 {
			newsToSend = append(newsToSend, new)
		}
	}

	sendRSSMessage(newsToSend)
}

func sendLastNNews(n int) {
	lastNews := getLastNews(n)

	sendRSSMessage(lastNews)
}

func sendRSSMessage(news []*FeedItem) {
	// Reverse news (to have the more recent at the end of the slice)
	for left, right := 0, len(news)-1; left < right; left, right = left+1, right-1 {
		news[left], news[right] = news[right], news[left]
	}

	session := newDiscordSession()
	defer session.Close()

	for _, new := range news {
		message := formatRSSDailyMessage(new)
		_, err := session.ChannelMessageSendEmbed(Conf.DiscordChannel, message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func formatRSSDailyMessage(new *FeedItem) *discordgo.MessageEmbed {
	msgFieldTitle := &discordgo.MessageEmbedField{
		Name:   "Title",
		Value:  new.Title,
		Inline: true,
	}

	msgFieldLink := &discordgo.MessageEmbedField{
		Name:   "Link",
		Value:  new.Link,
		Inline: true,
	}

	msgFields := []*discordgo.MessageEmbedField{msgFieldLink, msgFieldTitle}

	message := &discordgo.MessageEmbed{
		Title:       new.Source,
		Description: new.Description,
		Fields:      msgFields,
	}

	return message
}
