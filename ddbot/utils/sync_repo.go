package utils

import (
	"bytes"
	"ddbot/models"
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"os/exec"
	"runtime/debug"

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

func SyncRepo(config *models.DDEnv) {
	if CloneRepoCheck() {
		baseScriptsPath := config.RepoBaseDir
		if CheckDirOrFileIsExist(baseScriptsPath) {
			fmt.Printf("脚本仓库目录已存在，执行pull")
			repoPull(baseScriptsPath)
		} else {
			fmt.Printf("脚本仓库目录不存在，执行clone")
			repoClone(repoUrl, baseScriptsPath)
		}
	} else {
		fmt.Printf("为了避免程序内置的用户名密码被滥用，所以会有使用场景检查，当前环境不符合使用要求。")
	}
}

func repoClone(url string, directory string) {

	// Clone the given repository to the given directory
	fmt.Printf("git clone %s to %s", url, directory)

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
			fmt.Printf("已经是最新代码，暂无更新。")
		} else if errPull.Error() == "authentication required" {
			fmt.Printf("用户密码登陆失败，更新失败。")
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
		fmt.Printf("还原本地修改（新增文件不受影响）防止更新冲突.....\n%v", outInfo.String())
	}
}
