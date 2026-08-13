package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
	"github.com/mythenox/binge-tracker/internal/handler"

	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := app.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	dbQueries := database.New(db)

	programState := &app.State{
		DB:  dbQueries,
		Cfg: &cfg,
	}

	cmds := app.Commands{
		RegisteredCmds: make(map[string]func(*app.State, app.Command) error),
	}
	cmds.Register("init", handler.HandlerInit)
	cmds.Register("play", handler.HandlerPlay)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := flag.Args()[1]
	cmdArgs := flag.Args()[2:]

	err = cmds.Run(programState, app.Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
