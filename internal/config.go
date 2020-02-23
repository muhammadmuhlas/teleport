package internal

import "github.com/spf13/viper"

var Config *viper.Viper

func NewConfig(filename, extension string) {
	viper.SetConfigName(filename)
	viper.SetConfigType(extension)
	viper.AddConfigPath(".")
	Config = viper.GetViper()
}

func ConfigMap() {
	viper.Set("source.provider", "gitlab")
	viper.Set("source.credential.username", "user@example.com")
	viper.Set("source.credential.password", "password")
	viper.Set("source.credential.access_token", "token")
	viper.Set("source.whitelist_namespace", []string{"team_a", "team_b"})
	viper.Set("source.blacklist_repositories", []string{"repo_a", "repo-b"})
	viper.Set("source.branch", []string{"master", "staging", "dev"})
	viper.Set("source.reuse_tmp", true)

	viper.Set("target.provider", "bitbucket")
	viper.Set("target.namespace", "team")
	viper.Set("target.credential.username", "user@example.com")
	viper.Set("target.credential.password", "password")
	viper.Set("target.allow_update", true)
}
