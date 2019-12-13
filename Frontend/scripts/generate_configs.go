// +build ignore

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/thehowl/conf"
)

type simplePage struct {
	Handler, Template, TitleBar, KyutGrill string
	MinPrivilegesRaw                       uint64
}

type noTemplate struct {
	Handler, TitleBar, KyutGrill string
	MinPrivileges                uint64
}

var simplePages = [...]simplePage{{"/", "homepage.html", "Home Page", "homepage2.jpg", 0}, {"/login", "login.html", "Log in", "login2.jpg", 0}, {"/settings/avatar", "settings/avatar.html", "Change avatar", "settings2.jpg", 2}, {"/dev/tokens", "dev/tokens.html", "Your API tokens", "dev.jpg", 2}, {"/beatmaps/rank_request", "beatmaps/rank_request.html", "Request beatmap ranking", "request_beatmap_ranking.jpg", 2}, {"/donate", "support.html", "Support Ripple", "donate2.png", 0}, {"/doc", "doc.html", "Documentation", "documentation.jpg", 0}, {"/doc/:id", "doc_content.html", "View document", "documentation.jpg", 0}, {"/help", "help.html", "Contact support", "help.jpg", 0}, {"/leaderboard", "leaderboard.html", "Leaderboard", "leaderboard2.jpg", 0}, {"/friends", "friends.html", "Friends", "", 2}, {"/changelog", "changelog.html", "Changelog", "changelog.jpg", 0}, {"/team", "team.html", "Team", "", 0}, {"/pwreset", "pwreset.html", "Reset password", "", 0}, {"/about", "about.html", "About", "", 0}}

func main() {
	for _, p := range simplePages {
		fmt.Print("=> ", p.Handler+" ... ")
		noTemplateP := noTemplate{
			Handler:       p.Handler,
			TitleBar:      p.TitleBar,
			KyutGrill:     p.KyutGrill,
			MinPrivileges: p.MinPrivilegesRaw,
		}
		d := []byte("{{/*###\n")
		confData, err := conf.ExportRaw(&noTemplateP)
		if err != nil {
			panic(err)
		}
		d = append(d, confData...)
		fileData, err := ioutil.ReadFile("templates/" + p.Template)
		if err != nil {
			panic(err)
		}
		d = append(d, []byte("*/}}\n")...)
		d = append(d, fileData...)
		err = ioutil.WriteFile("templates/"+p.Template, d, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("ok.")
	}
}
