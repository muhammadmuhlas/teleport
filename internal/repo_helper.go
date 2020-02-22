package internal

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/muhammadmuhlas/teleport/internal/payload"
	"github.com/parnurzeal/gorequest"
)

func GetAllGitlabRepositories() (gr payload.GitlabRepositories) {
	request := gorequest.New()
	if _, _, err := request.
		Get("https://gitlab.com/api/v4/projects/?simple=yes&visibility=private&per_page=100000&page=1").
		Set("PRIVATE-TOKEN", Config.GetString("source.credential.access_token")).
		EndStruct(&gr); err != nil {
		Log.Panic(err)
	}
	return
}

func CheckBitbucketRepository(owner, repoSlug string) bool {
	bbCon := bitbucket.NewBasicAuth(Config.GetString("target.credential.username"), Config.GetString("target.credential.password"))
	opt := &bitbucket.RepositoryOptions{
		Owner:    owner,
		RepoSlug: repoSlug,
	}

	_, err := bbCon.Repositories.Repository.Get(opt)
	if err != nil {
		Log.Debug(err)
		return false
	}
	return true
}

func CreateBitbucketRepository(owner, repoSlug, isPrivate, forkPolicy string) *bitbucket.Repository {
	bbCon := bitbucket.NewBasicAuth(Config.GetString("target.credential.username"), Config.GetString("target.credential.password"))
	opt := &bitbucket.RepositoryOptions{
		Owner:      owner,
		RepoSlug:   repoSlug,
		IsPrivate:  isPrivate,
		ForkPolicy: forkPolicy,
	}

	r, err := bbCon.Repositories.Repository.Create(opt)
	if err != nil {
		Log.Panic(err)
	}
	return r
}
