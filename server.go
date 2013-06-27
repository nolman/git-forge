package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"text/template"
)

type Status struct {
	Success bool
}

type ReceiveHookConfig struct {
	BasePath string
	AppPath string
	AppName string
	DatabaseUrl string
}

func handler(w http.ResponseWriter, r *http.Request) {
	hook_config := ReceiveHookConfig{
		BasePath: os.Getenv("REPO_DIRECTORY"),
		AppPath: path.Join(os.Getenv("REPO_DIRECTORY"), r.FormValue("application")),
		AppName: r.FormValue("application"),
		DatabaseUrl: r.FormValue("database_url"),
	}
	var success bool
	success = createBareRepo(hook_config)
	createDeployHook(hook_config)
	data,_ := json.Marshal(Status{success})
	fmt.Fprintf(w, string(data))
}

func createBareRepo(hook_config ReceiveHookConfig) bool{
	if !execAndWait("mkdir", "-p", hook_config.AppPath){ return false }
	if !execAndWait("git", "init", hook_config.AppPath, "--bare"){ return false }
	return true
}

func createDeployHook(hook_config ReceiveHookConfig) bool{
	post_receive_hook := path.Join(hook_config.AppPath, "hooks", "post-receive")
	t, _ := template.ParseFiles("tmpl/post-receive.sh")
	var content bytes.Buffer
	t.Execute(&content, hook_config)
	ioutil.WriteFile(post_receive_hook, content.Bytes(), 0755)
	return true
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func execAndWait(command string, args ...string) bool{
	program, _ := exec.LookPath(command)
	fmt.Println("Running %v %v", program, args)
	cmd := exec.Command(program, args...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("command failed:" + err.Error())
		return false;
	}
	return true;
}

func main() {
fmt.Println("STARTING")
	if os.Getenv("REPO_DIRECTORY") == "" {
		panic("ENV REPO_DIRECTORY was undefined unable to load git-forge")
	}
	http.HandleFunc("/git-forge", handler)
	http.ListenAndServe(":8080", nil)
}

