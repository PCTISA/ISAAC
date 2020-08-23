package command

import (
	"fmt"

	"github.com/PCTISA/ISAAC/log"
	"github.com/PCTISA/ISAAC/multiplexer"
	"github.com/patrickmn/go-cache"
)

// Example is a command
type Example struct {
	Command  string
	HelpText string

	/* Optional rate-limiting settings */
	RateLimitMax int
	RateLimitDB  *cache.Cache

	/* Anything that needs passed to a given command can be added here */
	Logger *log.Logs // For example, the Logger
}

// Init is called by the multiplexer before the bot starts to initialize any
// variables the command needs.
func (c Example) Init(m *multiplexer.Mux) {
	// Nothing to init
}

// Handle is called by the multiplexer whenever a user triggers the command.
func (c Example) Handle(ctx *multiplexer.Context) {
	ctx.ChannelSend("Congradulations! You've run your first command")

	c.Logger.CmdErr(
		ctx,
		fmt.Errorf("this is an example command error"),
		"Command errors print to the console and to chat.",
	)
}

// HandleHelp is not called by the multiplexer. It is used by the
// `!help` command (if included) to provide a bigger description of the
// command's functionality.
func (c Example) HandleHelp(ctx *multiplexer.Context) {
	ctx.ChannelSend("Much bigger/more detailed command description")
}

// Settings is called by the multiplexer on startup to process any settings
// associated with that command.
func (c Example) Settings() *multiplexer.CommandSettings {
	return &multiplexer.CommandSettings{
		Command:  c.Command,
		HelpText: c.HelpText,

		/* Optional rate limiting */
		RateLimitMax: c.RateLimitMax,
		RateLimitDB:  c.RateLimitDB,
	}
}
