package internal

import (
	"github.com/muhammadmuhlas/teleport/internal/payload"
	"github.com/parnurzeal/gorequest"
)

func CheckGitlab(token string) (gr payload.GitlabRepositories) {
	request := gorequest.New()
	if _, _, err := request.
		Get("https://gitlab.com/api/v4/projects/?simple=yes&visibility=private&per_page=100000&page=1").
		Set("PRIVATE-TOKEN", token).
		EndStruct(&gr); err != nil {
		Log.Panic(err)
	}
	return
}
