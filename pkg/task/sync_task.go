package task

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type SyncTask struct {
	title   string
	srcRepo string
	dstRepo string
	err     error

	privateKeyFile string
}

func NewMigrationTask(srcRepo, dstRepo, privateKeyFile string) *SyncTask {
	dirName := filepath.Base(srcRepo)

	return &SyncTask{
		title:   fmt.Sprintf("syncTask-%s", dirName),
		srcRepo: srcRepo,
		dstRepo: dstRepo,

		privateKeyFile: privateKeyFile,
	}
}

func (t *SyncTask) Run() (err error) {
	// 克隆到临时目录
	srcAuth, repoUrl, err := getAuth(t.srcRepo, t.privateKeyFile)
	if err != nil {
		return
	}
	t.srcRepo = repoUrl

	var dot billy.Filesystem
	dirName := fmt.Sprintf("/tmp/%s", filepath.Base(t.srcRepo))
	defer func(path string) {
		e := os.RemoveAll(path)
		if e != nil {
			logrus.Warnf("Remove dir %s failed: %v", path, e)
		}
	}(dirName)
	_, err = os.Stat(dirName)
	if err == nil {
		err = os.RemoveAll(dirName)
		if err != nil {
			return err
		}
	}
	dot = osfs.New(dirName)

	repo, err := git.Clone(filesystem.NewStorage(dot, cache.NewObjectLRUDefault()), nil, &git.CloneOptions{
		URL:    repoUrl,
		Mirror: true,
		Auth:   srcAuth,
	})
	if err != nil {
		return
	}

	// 推送到目标仓库
	dstAuth, repoUrl, err := getAuth(t.dstRepo, t.privateKeyFile)
	if err != nil {
		return
	}
	t.dstRepo = repoUrl

	err = repo.Push(&git.PushOptions{
		RemoteURL: repoUrl,
		Auth:      dstAuth,
		Force:     true,
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/heads/*",
			"+refs/tags/*:refs/tags/*",
			"+refs/change/*:refs/change/*",
		},
	})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		err = nil
	}
	return
}

func (t *SyncTask) GetTitle() string {
	return t.title
}

func (t *SyncTask) GetSrcRepo() string {
	return t.srcRepo
}

func (t *SyncTask) GetDstRepo() string {
	return t.dstRepo
}

func (t *SyncTask) SetError(err error) {
	t.err = err
}

func (t *SyncTask) GetError() error {
	return t.err
}

func getAuth(repo, privateKeyFile string) (auth transport.AuthMethod, repoUrl string, err error) {
	repoUrl = repo
	if strings.HasPrefix(repo, "git") {
		if privateKeyFile == "" {
			u, err := user.Current()
			if err != nil {
				return auth, repoUrl, err
			}
			privateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir)
		}

		_, err = os.Stat(privateKeyFile)
		if err != nil {
			return auth, repoUrl, err
		}
		auth, err = ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	} else if strings.HasPrefix(repo, "http") {
		if strings.Contains(repo, "@") {
			fields := strings.Split(repo, "@")
			userInfo := strings.Join(fields[0:len(fields)-1], "@")
			repoInfo := fields[len(fields)-1]

			if strings.HasPrefix(userInfo, "http://") {
				userInfo = strings.ReplaceAll(userInfo, "http://", "")
				repoUrl = fmt.Sprintf("%s%s", "http://", repoInfo)
			} else if strings.HasPrefix(userInfo, "https://") {
				userInfo = strings.ReplaceAll(userInfo, "https://", "")
				repoUrl = fmt.Sprintf("%s%s", "https://", repoInfo)
			}

			fields = strings.Split(userInfo, ":")
			username := fields[0]
			password := strings.Join(fields[1:], ":")

			auth = &http.BasicAuth{
				Username: username,
				Password: password,
			}
			logrus.Debugf("http auth: username: %s, password: %s, repoUrl: %s", username, password, repoUrl)
		} else {
			auth = nil
		}
	} else {
		auth = nil
	}
	return
}
