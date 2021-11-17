package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"ddbot/models"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	//repoUrl 非公开脚本仓库地址(编译时传入)
	repoUrl = "compile_repo_url"
	// github 非公开仓库用户名(编译时传入)
	gitUsername = "compile_git_username"
	// github 非公开仓库访问token(编译时传入)
	gitToken = "compile_git_token"
)

// SyncRepo
// @description	根据传入仓库信息配置，同步仓库
// @auth	@iouAkira
// @param1  config *models.DDEnv
func SyncRepo(config *models.DDEnv) {
	if CloneRepoCheck() {
		baseScriptsPath := config.RepoBaseDir
		if CheckDirOrFileIsExist(baseScriptsPath) {
			fmt.Printf("脚本仓库目录已存在，执行pull\n")
			repoPull(baseScriptsPath)
		} else {
			fmt.Printf("脚本仓库目录不存在，执行clone\n")
			repoClone(repoUrl, baseScriptsPath)
		}
	} else {
		fmt.Printf("为了避免程序内置的用户名密码被滥用，所以会有使用场景检查，当前环境不符合使用要求。\n")
	}
}

// repoClone
// @description	根据传入仓库信息配置，clone仓库
// @auth	@iouAkira
// @param1     url string
// @param2     directory string
func repoClone(url string, directory string) {
	// Clone the given repository to the given directory
	fmt.Printf("git clone %s to %s\n", url, directory)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: gitUsername, // yes, this can be anything except an empty string
			Password: gitToken,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}

// pullRepo
// @description	根据传入仓库信息配置，更新仓库
// @auth	@iouAkira
// @param1     repoPath string
func repoPull(path string) {
	//对异常状态进行补货并输出到缓冲区
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover %v\n", err)
			debug.PrintStack()
		}
	}()
	//resetHard(path) //还原本地修改操作放到shell_default_scripts.sh里面
	// We instantiate a new repository targeting the given path (the .git folder)
	r, errP := git.PlainOpen(path)
	CheckIfError(errP)

	// Get the working directory for the repository
	w, errW := r.Worktree()
	CheckIfError(errW)

	//Pull the latest changes from the origin remote and merge into the current branch
	errPull := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      true,
		Auth: &http.BasicAuth{
			Username: gitUsername,
			Password: gitToken,
		}})
	if errPull != nil {
		if errPull.Error() == "already up-to-date" {
			fmt.Printf("已经是最新代码，暂无更新。\n")
		} else if errPull.Error() == "authentication required" {
			fmt.Printf("用户密码登陆失败，更新失败。\n")
		} else {
			fmt.Printf(errPull.Error())
		}
	} else {
		CheckIfError(errPull)
		// 获取最后一次提交的信息。
		ref, errH := r.Head()
		CheckIfError(errH)

		commit, errC := r.CommitObject(ref.Hash())
		CheckIfError(errC)
		fmt.Printf("%v", commit)
	}
}

// resetHard
// @description	根据传入仓库流经还原本地修改，防止更新仓库冲突
// @auth	@iouAkira
// @param     repoPath string
func resetHard(path string) {
	//var execResult string
	var cmdArguments []string
	resetCmd := []string{"git", "-C", path, "reset", "--hard"}

	for i, v := range resetCmd {
		if i >= 1 {
			cmdArguments = append(cmdArguments, v)
		}
	}
	command := exec.Command(resetCmd[0], cmdArguments...)
	outInfo := bytes.Buffer{}
	command.Stdout = &outInfo
	err := command.Start()
	if err != nil {
		fmt.Printf(err.Error())
	}
	if err = command.Wait(); err != nil {
		fmt.Printf(err.Error())
	} else {
		//fmt.Println(command.ProcessState.Pid())
		//fmt.Println(command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
		fmt.Printf("还原本地修改（新增文件不受影响）防止更新冲突.....\n%v\n", outInfo.String())
	}
}
