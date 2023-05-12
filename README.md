# git-syncer

[![go-report](https://goreportcard.com/badge/github.com/DockerContainerService/git-syncer)](https://goreportcard.com/report/github.com/DockerContainerService/git-syncer)
![contributors](https://img.shields.io/github/contributors/DockerContainerService/git-syncer)
![size](https://img.shields.io/github/repo-size/DockerContainerService/git-syncer)
![languages](https://img.shields.io/github/languages/count/DockerContainerService/git-syncer)
![file](https://img.shields.io/github/directory-file-count/DockerContainerService/git-syncer)
![used-by](https://img.shields.io/sourcegraph/rrc/github.com/DockerContainerService/git-syncer)
[![license](https://img.shields.io/github/license/DockerContainerService/git-syncer)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![release](https://img.shields.io/github/v/release/DockerContainerService/git-syncer)](https://github.com/DockerContainerService/git-syncer/releases)
[![download](https://img.shields.io/github/downloads/DockerContainerService/git-syncer/total.svg)](https://api.github.com/repos/DockerContainerService/git-syncer/releases)
[![last-release](https://img.shields.io/github/release-date/DockerContainerService/git-syncer)](https://github.com/DockerContainerService/git-syncer/releases)


## Usage
目标仓库可以不存在，但是存在的仓库必须关闭分支保护或者将分支保护配置为允许强制推送
### Install git-syncer
you can download the latest binary release [here](https://github.com/DockerContainerService/git-syncer/releases)

### Install from source
```bash
go get github.com/DockerContainerService/git-syncer
cd ${GOPATH}/github.com/DockerContainerService/git-syncer
make all
```

### Get usage information
```bash
[root@tencent ~]# ./git-syncer -h
Usage:
  git-syncer [flags]

Flags:
  -c, --config string           config file
  -d, --debug                   enable debug mode
  -h, --help                    help for git-syncer
      --privateKeyFile string   private key file
  -p, --proc int                numbers of worker (default 5)
  -r, --retries int             times to retry failed task (default 3)
```
### Usage example
```json
{
  "git@github.com:DockerContainerService/git-syncer.git": "git@172.17.162.204:demo/git-syncer.git"
}
```
```bash
$ ./git-syncer -c test.json 
[2023-05-11 19:54:55]  INFO Generate task syncTask-git-syncer.git: git@github.com:DockerContainerService/git-syncer.git => git@172.17.162.204:demo/git-syncer.git
[2023-05-11 19:54:55]  INFO Run task syncTask-git-syncer.git: git@github.com:DockerContainerService/git-syncer.git => git@172.17.162.204:demo/git-syncer.git
[2023-05-11 19:55:20]  WARN Retry task syncTask-git-syncer.git, times 1
[2023-05-11 19:55:45]  WARN Retry task syncTask-git-syncer.git, times 2
[2023-05-11 19:56:10]  WARN Retry task syncTask-git-syncer.git, times 3
[2023-05-11 19:56:35] ERROR Run task syncTask-git-syncer.git git@github.com:DockerContainerService/git-syncer.git => git@172.17.162.204:demo/git-syncer.git err: dial tcp 172.17.162.204:22: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
[2023-05-11 19:56:35]  INFO Finished, 0 task succeeded, 1 task failed. failed task list:

syncTask-git-syncer.git: git@github.com:DockerContainerService/git-syncer.git => git@172.17.162.204:demo/git-syncer.git, err: dial tcp 172.17.162.204:22: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=DockerContainerService/git-syncer&type=Date)](https://star-history.com/#DockerContainerService/git-syncer&Date)


