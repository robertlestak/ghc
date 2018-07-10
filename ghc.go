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
	"os/user"
	"path"
	"strings"
	"time"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

var GITHUB_API = "https://api.github.com/user/repos"

func addTeam(tid string, owner string, repo string, perm string) error {
	var e error
	c := &http.Client{}
	if perm != "pull" && perm != "push" && perm != "admin" {
		perm = ""
	}
	GITHUB_API = "https://api.github.com/teams/:team_id/repos/:owner/:repo"
	GITHUB_API = strings.Replace(GITHUB_API, ":team_id", tid, 1)
	GITHUB_API = strings.Replace(GITHUB_API, ":owner", owner, 1)
	GITHUB_API = strings.Replace(GITHUB_API, ":repo", repo, 1)
	req, err := http.NewRequest("PUT", GITHUB_API, nil)
	req.Header.Set("Accept", "application/vnd.github.hellcat-preview+json")
	if perm != "" {
		q := req.URL.Query()
		q.Add("permission", perm)
		req.URL.RawQuery = q.Encode()
	}
	if err != nil {
		return err
	}
	mc := getmachine()
	req.SetBasicAuth(mc.Login, mc.Password)
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, berr := ioutil.ReadAll(res.Body)
	if berr != nil {
		return berr
	}
	fmt.Println(string(b))
	return e
}

func createTeams(repo string, teams []string, perm string) error {
	var e error
	owner, err := getUsername()
	if err != nil {
		return err
	}
	for _, v := range teams {
		err := addTeam(v, owner, repo, perm)
		if err != nil {
			fmt.Println(err)
		}
		if len(teams) > 1 {
			time.Sleep(time.Millisecond * 500)
		}
	}
	return e
}

func getmachine() *netrc.Machine {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mc, merr := netrc.FindMachine(path.Join(usr.HomeDir, ".netrc"), "api.github.com")
	if merr != nil {
		fmt.Println(merr)
		os.Exit(1)
	}
	if mc == nil {
		fmt.Println("~/.netrc Required")
		os.Exit(1)
	}
	return mc
}

func getUsername() (string, error) {
	var u string
	var e error
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return u, err
	}
	mc := getmachine()
	c := &http.Client{}
	req.SetBasicAuth(mc.Login, mc.Password)
	res, err := c.Do(req)
	if err != nil {
		return u, err
	}
	defer res.Body.Close()
	b, berr := ioutil.ReadAll(res.Body)
	if berr != nil {
		return u, berr
	}
	jd := make(map[string]interface{})
	jerr := json.Unmarshal(b, &jd)
	if jerr != nil {
		return u, jerr
	}
	_, er := jd["login"]
	if er {
		u = jd["login"].(string)
	}
	return u, e
}

func trimStringSlice(s []string) []string {
	var ns []string
	for _, v := range s {
		if strings.TrimSpace(v) != "" {
			ns = append(ns, strings.TrimSpace(v))
		}
	}
	return ns
}

func main() {
	priv := flag.Bool("p", true, "Private Repository")
	desc := flag.String("d", "", "Description")
	org := flag.String("o", "", "Organization")
	clip := flag.Bool("c", false, "Copy clone URL to clipboard")
	clipssh := flag.Bool("s", false, "Copy ssh URL to clipboard")
	teamsString := flag.String("t", "", "Teams (comma separated)")
	perm := flag.String("u", "", "Team Permissions (pull, push, admin)")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Respository required")
		return
	}
	var teams []string
	if *teamsString != "" {
		teams = strings.Split(*teamsString, ",")
		teams = trimStringSlice(teams)
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
	mc := getmachine()
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
	if len(teams) > 0 {
		terr := createTeams(repo.Name, teams, *perm)
		if terr != nil {
			fmt.Println(err)
		}
	}
	url := jd["clone_url"].(string)
	if *clipssh {
		url = jd["ssh_url"].(string)
		clipboard.WriteAll(url)
	} else if *clip {
		url = jd["clone_url"].(string)
		clipboard.WriteAll(url)
	}
	fmt.Println(url)
}
