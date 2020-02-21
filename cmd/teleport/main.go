package main

import (
	"github.com/spf13/viper"
	"os"
	"fmt"
	"strings"
	"github.com/muhammadmuhlas/teleport/internal"
	"github.com/sirupsen/logrus"
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
}

func ConfigMap() {
	viper.Set("source.provider", "gitlab")
	viper.Set("source.credential.username", "user@example.com")
	viper.Set("source.credential.password", "password")
	viper.Set("source.credential.access_token", "token")

	viper.Set("target.provider", "bitbucket")
	viper.Set("target.credential.username", "user@example.com")
	viper.Set("target.credential.password", "password")
}

func main() {
	if !GetOrLoadConfig() { os.Exit(0) }
	internal.Log.Println("Getting Repositories from", viper.Get("source.provider"))
	gitlabRepo := internal.CheckGitlab(viper.Get("source.credential.access_token").(string))
	for _, repo := range gitlabRepo {
		if internal.InArray(repo.Namespace.Path, viper.Get("source.whitelist_namespace").([]interface{})) {
			logrus.Println(repo.HTTPURLToRepo)
		}
	}
	internal.Log.Println("Found", len(gitlabRepo), "repositories!")
}

func GetOrLoadConfig() (ready bool) {
	viper.SetConfigName("teleport")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			internal.Log.Warning("Initial config not found!")
			internal.Scanner("Generate initial config? (Y/n) ", func(input string) bool {
				if strings.ToLower(input) == "y" || input == "" {
					ConfigMap()
					viper.SafeWriteConfigAs("teleport.yaml")
					internal.Log.Info("Initial config saved at: teleport.yaml")
					internal.Log.Info("Comeback when teleport.yaml configured")
					return false
				}
				internal.Log.Error("Initial config must be generated.")
				return true
			})
			return false
		} else {
			internal.Log.Error(err)
		}
	}
	return true
}
