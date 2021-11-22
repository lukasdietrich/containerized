package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/lukasdietrich/dockerized/updater/internal/fetcher"
)

func init() {
	// set authentication globally
	ssh.DefaultAuthBuilder = func(user string) (ssh.AuthMethod, error) {
		publicKeysCallback := ssh.PublicKeysCallback{
			User:     user,
			Callback: loadDeployKey(),
			HostKeyCallbackHelper: ssh.HostKeyCallbackHelper{
				HostKeyCallback: acceptGithubHostKeys,
			},
		}

		return &publicKeysCallback, nil
	}
}

type Meta struct {
	repo     *git.Repository
	worktree *git.Worktree
}

func Open() (*Meta, error) {
	log.Printf("Cloning repository into memory.")

	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:               "git@github.com:lukasdietrich/dockerized.git",
		RemoteName:        "origin",
		ReferenceName:     plumbing.ReferenceName("refs/heads/master"),
		SingleBranch:      true,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Tags:              git.NoTags,
	})

	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return &Meta{repo: repo, worktree: worktree}, nil
}

func (m *Meta) UpdateVersions() error {
	if err := m.worktree.Pull(&git.PullOptions{}); err != nil {
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}
	}

	folderSlice, err := m.worktree.Filesystem.ReadDir(m.worktree.Filesystem.Root())
	if err != nil {
		return err
	}

	var anyUpdated bool

	for _, folder := range folderSlice {
		if folder.IsDir() {
			foldername := folder.Name()

			updated, err := m.updateApplication(foldername)
			if err != nil {
				return err
			}

			anyUpdated = anyUpdated || updated
		}
	}

	if anyUpdated {
		log.Print("At least one update was committed. Pushing to repository.")
		return m.repo.Push(&git.PushOptions{})
	}

	return nil
}

type appSpec struct {
	DisplayName string `json:"displayName"`
	Origin      string `json:"origin"`
}

func (m *Meta) updateApplication(foldername string) (bool, error) {
	log.Printf("Start of %q", foldername)
	defer log.Printf("End of %q", foldername)

	fs, err := m.worktree.Filesystem.Chroot(foldername)
	if err != nil {
		return false, err
	}

	app, err := readAppSpec(fs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Print(".dockerized.json does not exist. skipping folder.")
			return false, nil
		}

		return false, err
	}

	currentVersion, err := readVersion(fs)
	if err != nil {
		return false, err
	}

	latestRelease, err := fetcher.Latest(app.Origin)
	if err != nil {
		return false, err
	}

	log.Printf("Current version = %q, latest version = %q (published at %q).",
		currentVersion, latestRelease.Version, latestRelease.PublishedAt)

	if currentVersion != latestRelease.Version {
		log.Print("Latest version differs from current version.")

		writeVersion(fs, latestRelease.Version)
		return true, m.commitUpdate(app, latestRelease)
	}

	return false, nil
}

func (m *Meta) commitUpdate(app *appSpec, release *fetcher.Release) error {
	if err := m.worktree.AddWithOptions(&git.AddOptions{}); err != nil {
		return err
	}

	message := fmt.Sprintf("Update %q to version %q", app.DisplayName, release.Version)
	log.Printf("Commit message: %q.", message)

	hash, err := m.worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  os.Getenv("DOCKERIZED_GIT_AUTHOR_NAME"),
			Email: os.Getenv("DOCKERIZED_GIT_AUTHOR_EMAIL"),
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	log.Printf("Commit hash: %q.", hash)
	return nil
}

func readAppSpec(fs billy.Basic) (*appSpec, error) {
	f, err := fs.Open(".dockerized.json")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var app appSpec

	if err := json.NewDecoder(f).Decode(&app); err != nil {
		return nil, err
	}

	return &app, nil
}

func readVersion(fs billy.Basic) (string, error) {
	b, err := util.ReadFile(fs, "VERSION")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func writeVersion(fs billy.Basic, version string) error {
	return util.WriteFile(fs, "VERSION", []byte(version), 0644)
}
