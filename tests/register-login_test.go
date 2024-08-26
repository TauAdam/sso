package tests

import (
	"github.com/TauAdam/sso/tests/suite"
	"testing"
)

const (
	secret = "secret-string"
	appID  = 1

	passwordLength = 8
)

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, suite := suite.New(t)

}
