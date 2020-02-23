package internal

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"time"
)

func BeginTeleport(migration *Migration) {
	t := time.Now()
	Log.Info("Teleporting ", migration.Name)
	// Clone
	r, err := git.PlainClone("teleport-tmp/"+migration.Name, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "muhlas_ekrut",
			Password: Config.GetString("source.credential.access_token"),
		},
		URL: migration.GitSource,
	})

	// Open If cloned
	if err != nil {
		if Config.GetBool("source.reuse_tmp") {
			r, err = git.PlainOpen("teleport-tmp/" + migration.Name)
			if err != nil {
				Log.Panicln(err)
			}
		} else {
			Log.Panicln(err, "please do remove dictionary teleport-tmp or allow reuse_tmp")
		}
	}

	// Initialize Worktree
	w, err := r.Worktree()
	if err != nil {
		Log.Panicln(err)
	}

	// Pull Changes
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: Config.GetString("source.credential.username"),
			Password: Config.GetString("source.credential.access_token"),
		},
		Force: true,
	})

	// Create Target Repository and Add Remotes
	targetRepo := CreateOrGetBitbucketRepository(Config.GetString("target.namespace"), migration.Name, "true", "no_public_forks")
	for _, cloneLink := range targetRepo.Links["clone"].([]interface{}) {
		if linkKey := cloneLink.(map[string]interface{})["name"]; linkKey == "https" {
			linkType := cloneLink.(map[string]interface{})["href"]
			_, err = r.CreateRemote(&config.RemoteConfig{
				Name: "target",
				URLs: []string{linkType.(string)},
			})
		}
	}

	// Push to target Repositories each branch
	for _, v := range Config.GetStringSlice("source.branch") {
		// Checkout
		remoteRef, err := r.Reference(plumbing.ReferenceName("refs/remotes/origin/"+v), true)
		if err != nil {
			Log.Error(migration.Name, "@", v, " is ", err)
			continue
		}
		newRef := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+v), remoteRef.Hash())
		if err := r.Storer.SetReference(newRef); err != nil {
			Log.Error(migration.Name, "@", v, " is ", err)
		}
		if err := w.Checkout(&git.CheckoutOptions{Branch: newRef.Name(), Create: false}); err != nil {
			Log.Warning(migration.Name, "@", v, " is ", err)
		}

		// Push
		err = r.Push(&git.PushOptions{
			RemoteName: "target",
			Auth: &http.BasicAuth{
				Username: Config.GetString("target.credential.username"),
				Password: Config.GetString("target.credential.password"),
			},
		})
		if err != nil && err.Error() != "already up-to-date" {
			Log.Warning(migration.Name, "@", v, " is ", err)
		} else {
			Log.Info(migration.Name, "@", v, " has been teleported!")
		}
	}

	// Push All Tags To Target
	tagrefs, err := r.Tags()
	if err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if obj, err := r.TagObject(t.Hash()); err == nil {
			err = r.Push(&git.PushOptions{
				RemoteName: "target",
				RefSpecs:   []config.RefSpec{config.RefSpec(t.Name() + ":" + t.Name())},
				Auth: &http.BasicAuth{
					Username: Config.GetString("target.credential.username"),
					Password: Config.GetString("target.credential.password"),
				},
			})
			if err != nil && err.Error() != "already up-to-date" {
				Log.Warning(migration.Name, "@", obj.Name, " is ", err)
			} else {
				Log.Info(migration.Name, "@", obj.Name, " has been teleported!")
			}
		}
		return nil
	}); err != nil {
		Log.Println(err)
	}

	Log.Printf("%s Teleported at %.2f seconds!", migration.Name, time.Now().Sub(t).Seconds())
}
