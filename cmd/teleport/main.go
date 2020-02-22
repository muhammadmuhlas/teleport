package main

import (
	"fmt"
	. "github.com/muhammadmuhlas/teleport/internal"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

func init() {
	fmt.Println(`
   __       __                      __ 
  / /____  / /__  ____  ____  _____/ /_
 / __/ _ \/ / _ \/ __ \/ __ \/ ___/ __/
/ /_/  __/ /  __/ /_/ / /_/ / /  / /_  
\__/\___/_/\___/ .___/\____/_/   \__/  
              /_/
==================================================
Teleport your repositories to another git provider
==================================================
`)
	NewConfig("teleport", "yaml")
}
func main() {
	Log.Println("Initializing Config")
	if !loadConfig() {
		os.Exit(0)
	}

	Log.Println("Analyzing Source...")
	SourceReporter(Config.GetString("source.provider"))
	TargetReporter(Config.GetString("target.provider"))

	Log.Println("Initiate Teleport Sequence...")
	t := time.Now()
	for _, repo := range MarkToMigrate {
		BeginTeleport(&repo)
	}
	Log.Printf("%d Repositories Teleported at %.2f seconds!", len(MarkToMigrate), time.Now().Sub(t).Seconds())
}

func loadConfig() (ready bool) {
	if err := Config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Log.Warning("Initial config not found!")
			Scanner("Generate initial config? (Y/n) ", func(input string) bool {
				if strings.ToLower(input) == "y" || input == "" {
					ConfigMap()
					Config.SafeWriteConfigAs("teleport.yaml")
					Log.Info("Initial config saved at: teleport.yaml")
					Log.Info("Comeback when teleport.yaml configured")
					return false
				}
				Log.Error("Initial config must be generated.")
				return true
			})
			return false
		} else {
			Log.Error(err)
		}
	}
	return true
}
