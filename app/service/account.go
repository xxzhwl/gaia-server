package service

import (
	"context"
	"fmt"

	"github.com/xxzhwl/gaia/framework/account"
)

var acct *account.Manager

func InitAccount() error {
	m, err := account.NewFramework()
	if err != nil {
		return fmt.Errorf("create account manager: %w", err)
	}
	if err := m.Bootstrap(context.Background()); err != nil {
		return fmt.Errorf("account bootstrap: %w", err)
	}
	acct = m
	return nil
}

func GetAccount() *account.Manager {
	return acct
}
