// Package discordoauth provides constant for using OAuth2 to access Discord.
package discordoauth

import "golang.org/x/oauth2"

// Endpoint is Discord's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://discordapp.com/api/oauth2/authorize",
	TokenURL: "https://discordapp.com/api/oauth2/token",
}
