package router_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtt28/shortener/api/router"
	_ "github.com/xtt28/shortener/api/test/init"
)

func TestInitRouter(t *testing.T) {
	initRouterDiscardReturnValue := func() {
		router.InitRouter()
	}
	assert.NotPanics(t, initRouterDiscardReturnValue)
}
