package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
)

// Get current SoT balance
func (b *Bot) GetBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetBalance",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!balance") {
		l.Debugf("Received '!balance' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}
		if userObj.ID <= 0 {
			replyMsg := fmt.Sprintf("%v, sorry but your are not a registered user.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}

		userBalance, err := database.GetBalance(b.Db, userObj.ID)
		if err != nil {
			replyMsg := fmt.Sprintf("Sorry, %v but there was an error fetching your balance from the SoT API: %v",
				m.Author.Mention(), err)
			AnswerUser(s, m, replyMsg)
			return
		}

		p := message.NewPrinter(language.German)
		replyMsg := fmt.Sprintf("%v, your current SoT balance is: %v gold, %v doubloons and %v ancient coins",
			m.Author.Mention(),
			p.Sprintf("%d", userBalance.Gold),
			p.Sprintf("%d", userBalance.Doubloons),
			p.Sprintf("%d", userBalance.AncientCoins),
		)
		AnswerUser(s, m, replyMsg)
	}
}