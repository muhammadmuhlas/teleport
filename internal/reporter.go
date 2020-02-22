package internal

import (
	"github.com/muhammadmuhlas/teleport/internal/payload"
)

type Migration struct {
	Name      string
	GitSource string
	GitTarget string
}

var SourceBag interface{}
var MarkToMigrate []Migration

func SourceReporter(source string) {
	Log.Println("Fetching repositories from", source)
	switch source {
	case "gitlab":
		SourceBag = GetAllGitlabRepositories()
		for _, repo := range SourceBag.(payload.GitlabRepositories) {
			if InArray(repo.Namespace.Path, Config.GetStringSlice("source.whitelist_namespace")) {
				Log.Println(repo.HTTPURLToRepo)
			}
		}
		Log.Println("Found", len(SourceBag.(payload.GitlabRepositories)), "repositories!")
	default:
		Log.Panic("Source not Supported!")
	}
}

func TargetReporter(target string) {
	if SourceBag == nil {
		Log.Panic("SourceBag is nil")
	}
	Log.Println("Fetching repositories from", target)
	switch target {
	case "bitbucket":
		for _, repo := range SourceBag.(payload.GitlabRepositories) {
			bbUrl := BitbucketHTTPURLBuilder(Config.GetString("target.credential.username"), Config.GetString("target.namespace"), repo.Path)
			if CheckBitbucketRepository(Config.GetString("target.namespace"), repo.Path) {
				Log.Warning("Skipping ", repo.Path, " (target already have this repository)")
				continue
			}
			Log.Info("Teleport Available, ", repo.HTTPURLToRepo, " => ", bbUrl)
			if !InArray(repo.Path, Config.GetStringSlice("source.blacklist_repositories")) {
				MarkToMigrate = append(MarkToMigrate, Migration{
					Name:      repo.Path,
					GitSource: repo.HTTPURLToRepo,
					GitTarget: bbUrl,
				})
				continue
			}
			Log.Warning("Skipping ", repo.Path, " (repository on blacklist)")
		}
	default:
		Log.Panic("Target not Supported!")
	}
}
