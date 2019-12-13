package pagemappings

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

var pageMappings = map[int]interface{}{
	1:  "/",
	2:  "/login",
	3:  "/register",
	4:  "/", // user cp (deleted)
	5:  "/settings/avatar",
	6:  "/settings",
	7:  "/settings/password",
	8:  "/settings/userpage",
	9:  "/", // some deleted pages
	10: "/",
	11: "/",
	12: "/",
	13: "/leaderboard",
	14: "/doc",
	15: "/", // ripple v1 documentation stuff. this stuff is so old I won't even bother.
	16: func(u url.URL) string {
		return fmt.Sprintf("/doc/%s", u.Query().Get("id"))
	},
	17: "/changelog",
	18: "/pwreset",
	19: func(u url.URL) string {
		return fmt.Sprintf("/pwreset/continue?k=%s", u.Query().Get("k"))
	},
	20: "/", // Beta keys
	21: "/about",
	22: "/",          // reports
	23: "/doc/rules", // rules
	24: "/",          // my report
	25: "/",          // report
	26: "/friends",
	27: "https://status.ripple.moe",
	28: "/", // user lookup
	29: "/2fa_gateway",
	30: "/settings/2fa",
	31: "/beatmaps/rank_request",
	32: "/dev/applications",
	33: "/dev/applications", // Theorically, this should be something like /dev/applications/<id>, but no-one ever used that page so who gives a fuck.
	34: "/donate",
	35: "/team",
	36: "/irc",
	37: "/beatmaps",
	38: "/register/verify",
	39: "/register/welcome",
	40: "/settings/discord",
	41: "/register", // elmo
}

// CheckRedirect checks if the request is to be redirected to another page.
// This is to avoid broken links because of the old website.
func CheckRedirect(c *gin.Context) {
	p := c.Request.URL.Path
	if p != "/" && p != "/index.php" {
		c.Next()
		return
	}

	if c.Query("u") != "" {
		c.Redirect(302, "/u/"+c.Query("u"))
		c.Abort()
		return
	} else if i, _ := strconv.Atoi(c.Query("p")); i != 0 {
		mapped := pageMappings[i]
		if mapped == nil {
			u := c.Request.URL
			u.Host = "old.ripple.moe"
			c.Redirect(302, u.String())
			return
		}
		if str, ok := mapped.(string); ok {
			c.Redirect(302, str)
			c.Abort()
			return
		}
		if f, ok := mapped.(func(url.URL) string); ok {
			c.Redirect(302, f(*c.Request.URL))
			c.Abort()
			return
		}
	}
}
