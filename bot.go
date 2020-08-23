package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PCTISA/ISAAC/command"
	"github.com/PCTISA/ISAAC/config"
	"github.com/PCTISA/ISAAC/log"
	"github.com/PCTISA/ISAAC/multiplexer"

	"github.com/bwmarrin/discordgo"
	goenv "github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type environment struct {
	Token     string `env:"BOT_TOKEN"`
	Debug     bool   `env:"DEBUG" envDefault:"false"`
	DataDir   string `env:"DATA_DIR" envDefault:"data/"`
	ConfigURL string `env:"CONFIG_URL"`
}

var (
	env  = environment{}
	cfg  *config.BotConfig
	logs *log.Logs

	prefix = "!"
)

func init() {
	/* Parse enviorment variables */
	if err := goenv.Parse(&env); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/* Check if URL is being specified */
	path := env.DataDir + "config.json"
	if len(env.ConfigURL) > 0 {
		path = env.ConfigURL
	}

	/* Parse config */
	var err error
	cfg, err = config.Get(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/* Define logging setup */
	logs = log.New(env.Debug)
}

func main() {
	/* Initialize DiscordGo */
	logs.Primary.Info("Starting Bot...")
	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		logs.Primary.WithError(err).Error("Problem starting bot")
	}
	logs.Primary.Info("Bot started")

	/* Initialize Mux */
	mux, err := multiplexer.New(prefix)
	if err != nil {
		logs.Primary.WithError(err).Fatalf("Unable to create multixplexer")
	}

	/* Use the logging middleware with the multiplexer */
	mux.UseMiddleware(logs.MuxMiddleware)

	/* Set Permissions */
	mux.SetPermissions(cfg.Permissions)

	/* Setup Errors */
	mux.SetErrors(&multiplexer.ErrorTexts{
		CommandNotFound: "Command not found.",
		NoPermissions:   "You do not have permissions to execute that command.",
		RateLimited:     "You've used this command too many times, wait a bit and try again.",
	})

	/* === Register all the things === */

	/* Register the commands with the multiplexer*/
	mux.Register(
		command.Gatekeeper{
			Command:  "role",
			HelpText: "Manage your access to roles, and their related channels",
			Logger:   logs,
		},
		command.JPEG{
			Command:  "jpeg",
			HelpText: "More JPEG for the last image. 'nuff said",
			Logger:   logs,
		},
		command.Help{
			Command:  "help",
			HelpText: "Displays help  information regarding the bot's commands",
			Logger:   logs,
		},
	)

	for k := range cfg.SimpleCommands {
		mux.RegisterSimple(multiplexer.SimpleCommand{
			Command:  k,
			Content:  cfg.SimpleCommands[k],
			HelpText: "This is a simple command",
		})
	}

	/* Configure multiplexer options */
	mux.SetOptions(&multiplexer.Options{
		IgnoreDMs:        true,
		IgnoreBots:       true,
		IgnoreNonDefault: true,
		IgnoreEmpty:      true,
	})

	/* Initialize the commands */
	mux.Initialize()

	/* === End Register === */

	/* Handle commands and start DiscordGo */
	dg.AddHandler(mux.Handle)

	err = dg.Open()
	if err != nil {
		logs.Primary.WithError(err).Error(
			"Problem opening websocket connection.",
		)
		return
	}

	/* Set a fun status message */
	idle := 0
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		IdleSince: &idle,
		Game: &discordgo.Game{
			Name: "with your network",
			Type: discordgo.GameTypeGame,
		},
		Status: "online",
	})

	defer dg.Close()

	/* Wait for interrupt */
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
