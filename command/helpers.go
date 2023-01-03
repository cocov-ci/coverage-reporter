package command

import (
	"github.com/urfave/cli/v2"
)

func getToken(ctx *cli.Context) string { return ctx.String("token") }
