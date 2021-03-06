package handler

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"gorm.io/gorm"
)

// Set a SoT RAT cookie
func UserSetRatCookie(d *gorm.DB, c *viper.Viper, u *user.User, r string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.UserSetRatCookie",
	})

	if err := database.UserSetPrefEnc(d, c, u.UserInfo.ID, "rat_cookie", r); err != nil {
		l.Errorf("Failed to store RAT cookie in DB: %v", err)
		return "", err
	}

	if err := database.UserDelPref(d, u.UserInfo.ID, "failed_rat_notify"); err != nil {
		l.Errorf("Failed to delete 'failed_rat_notify' preference: %v", err)
	}

	if err := database.UserDelPref(d, u.UserInfo.ID, "failed_rat_tries"); err != nil {
		l.Errorf("Failed to delete 'failed_rat_tries' userpref in DB: %v", err)
	}

	u.RatCookie = r
	responseMsg := "Thanks for setting/updating your RAT cookie."
	return responseMsg, nil
}
