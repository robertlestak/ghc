package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/bgentry/go-netrc/netrc"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

var GITHUB_API = "https://api.github.com/user/repos"
var version string

func main() {
	priv := flag.Bool("p", true, "Private Repository")
	desc := flag.String("d", "", "Description")
	initialize := flag.Bool("i", false, "Initialize")
	initbranch := flag.String("b", "origin", "Initialize Branch")
	org := flag.String("o", "", "Organization")
	clip := flag.Bool("c", false, "Copy clone URL to clipboard")
	clipssh := flag.Bool("s", false, "Copy ssh URL to clipboard")
	checkver := flag.Bool("v", false, "Display current version")
	flag.Parse()
	if *checkver {
		if version == "" {
			version = "latest"
		}
		fmt.Println(version)
		os.Exit(0)
	}
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Respository required")
		return
	}
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}
	mc, merr := netrc.FindMachine(path.Join(usr.HomeDir, ".netrc"), "api.github.com")
	if merr != nil {
		fmt.Println(merr)
		return
	}
	if mc == nil {
		fmt.Println("~/.netrc Required")
		return
	}
	if *org != "" {
		GITHUB_API = "https://api.github.com/orgs/" + *org + "/repos"
	}
	repo := &Repository{
		Name:        args[0],
		Private:     *priv,
		Description: *desc,
	}
	jb, err := json.Marshal(&repo)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := &http.Client{}
	req, err := http.NewRequest("POST", GITHUB_API, bytes.NewBuffer(jb))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(mc.Login, mc.Password)
	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	jd := make(map[string]interface{})
	jerr := json.Unmarshal(b, &jd)
	if jerr != nil {
		fmt.Println(jerr)
		return
	}
	_, er := jd["errors"]
	if er {
		em := jd["errors"]
		je, _ := json.Marshal(&em)
		fmt.Println(string(je))
		return
	}
	url := jd["clone_url"].(string)
	if *clipssh {
		url = jd["ssh_url"].(string)
		clipboard.WriteAll(url)
	} else if *clip {
		url = jd["clone_url"].(string)
		clipboard.WriteAll(url)
	}
	if *initialize {
		icmd := exec.Command("git", "init")
		ierr := icmd.Run()
		if ierr != nil {
			fmt.Println("Init Error:", ierr)
		}
		cmds := []string{"remote", "add", *initbranch, url}
		cmd := exec.Command("git", cmds...)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Remote Add Error:", err)
		}
	}
	fmt.Println(url)
}
