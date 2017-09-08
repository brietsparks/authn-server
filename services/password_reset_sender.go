package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/keratin/authn-server/config"
	"github.com/keratin/authn-server/models"
	"github.com/keratin/authn-server/tokens/resets"
)

func PasswordResetSender(cfg *config.Config, account *models.Account) error {
	if cfg.AppPasswordResetURL == nil {
		return fmt.Errorf("AppPasswordResetURL unconfigured")
	}
	if account == nil {
		return FieldErrors{{"account", ErrMissing}}
	}
	if account.Locked {
		return FieldErrors{{"account", ErrLocked}}
	}

	reset, err := resets.New(cfg, account.ID, account.PasswordChangedAt)
	if err != nil {
		return err
	}
	resetStr, err := reset.Sign(cfg.ResetSigningKey)
	if err != nil {
		return err
	}

	res, err := http.PostForm(cfg.AppPasswordResetURL.String(), url.Values{
		"account_id": []string{strconv.Itoa(account.ID)},
		"token":      []string{resetStr},
	})
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Status Code: %v", res.StatusCode)
	}
	return nil
}
