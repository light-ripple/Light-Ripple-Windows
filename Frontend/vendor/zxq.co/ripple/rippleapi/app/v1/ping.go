package v1

import (
	"math/rand"
	"time"

	"zxq.co/ripple/rippleapi/common"
)

var rn = rand.New(rand.NewSource(time.Now().UnixNano()))

var kaomojis = [...]string{
	"Σ(ノ°▽°)ノ",
	"( ƅ°ਉ°)ƅ",
	"ヽ(　･∀･)ﾉ",
	"˭̡̞(◞⁎˃ᆺ˂)◞*✰",
	"(p^-^)p",
	"(ﾉ^∇^)ﾉﾟ",
	"ヽ(〃･ω･)ﾉ",
	"(۶* ‘ꆚ’)۶”",
	"（。＞ω＜）。",
	"（ﾉ｡≧◇≦）ﾉ",
	"ヾ(｡･ω･)ｼ",
	"(ﾉ･д･)ﾉ",
	".+:｡(ﾉ･ω･)ﾉﾞ",
	"Σ(*ﾉ´>ω<｡`)ﾉ",
	"ヾ（〃＾∇＾）ﾉ♪",
	"＼（＠￣∇￣＠）／",
	"＼(^▽^＠)ノ",
	"ヾ(@^▽^@)ノ",
	"(((＼（＠v＠）／)))",
	"＼(*T▽T*)／",
	"＼（＾▽＾）／",
	"＼（Ｔ∇Ｔ）／",
	"ヽ( ★ω★)ノ",
	"ヽ(；▽；)ノ",
	"ヾ(。◕ฺ∀◕ฺ)ノ",
	"ヾ(＠† ▽ †＠）ノ",
	"ヾ(＠^∇^＠)ノ",
	"ヾ(＠^▽^＠)ﾉ",
	"ヾ（＠＾▽＾＠）ノ",
	"ヾ(＠゜▽゜＠）ノ",
	"(.=^・ェ・^=)",
	"((≡^⚲͜^≡))",
	"(^･o･^)ﾉ”",
	"(^._.^)ﾉ",
	"(^人^)",
	"(=；ェ；=)",
	"(=｀ω´=)",
	"(=｀ェ´=)",
	"（=´∇｀=）",
	"(=^･^=)",
	"(=^･ｪ･^=)",
	"(=^‥^=)",
	"(=ＴェＴ=)",
	"(=ｘェｘ=)",
	"＼(=^‥^)/’`",
	"~(=^‥^)/",
	"└(=^‥^=)┐",
	"ヾ(=ﾟ･ﾟ=)ﾉ",
	"ヽ(=^･ω･^=)丿",
	"d(=^･ω･^=)b",
	"o(^・x・^)o",
	"V(=^･ω･^=)v",
	"(⁎˃ᆺ˂)",
	"(,,^・⋏・^,,)",
}

var randomSentences = [...]string{
	"Proudly sponsored by Kirotuso!",
	"The brace is on fire!",
	"deverupa ga daisuki!",
	"It works!!!!",
	"Feelin' groovy!",
	"sudo rm -rf /",
	"Hi! I'm Flowey! Flowey the flower!",
	"Ripple devs are actually cats",
	"Support Howl's fund for buying a power supply for his SSD!",
	"Superman dies",
	"PP when?",
	"RWC hype",
	"I'd just like to interject for a moment.",
	"Running on an apple pie!",
	":thinking:",
	"The total entropy of an isolated system can only increase over time",
	"Where are my testicles, Summer?",
	"Why don't you ask the smartest people in the universe? Oh yeah, you can't. They blew up.",
}

func surpriseMe() string {
	n := int(time.Now().UnixNano())
	return randomSentences[n%len(randomSentences)] + " " + kaomojis[n%len(kaomojis)]
}

type pingResponse struct {
	common.ResponseBase
	ID              int                   `json:"user_id"`
	Privileges      common.Privileges     `json:"privileges"`
	UserPrivileges  common.UserPrivileges `json:"user_privileges"`
	PrivilegesS     string                `json:"privileges_string"`
	UserPrivilegesS string                `json:"user_privileges_string"`
}

// PingGET is a message to check with the API that we are logged in, and know what are our privileges.
func PingGET(md common.MethodData) common.CodeMessager {
	var r pingResponse
	r.Code = 200

	if md.ID() == 0 {
		r.Message = "You have not given us a token, so we don't know who you are! But you can still login with POST /tokens " + kaomojis[rn.Intn(len(kaomojis))]
	} else {
		r.Message = surpriseMe()
	}

	r.ID = md.ID()
	r.Privileges = md.User.TokenPrivileges
	r.UserPrivileges = md.User.UserPrivileges
	r.PrivilegesS = md.User.TokenPrivileges.String()
	r.UserPrivilegesS = md.User.UserPrivileges.String()

	return r
}

type surpriseMeResponse struct {
	common.ResponseBase
	Cats [100]string `json:"cats"`
}

// SurpriseMeGET generates cute cats.
//
// ... Yes.
func SurpriseMeGET(md common.MethodData) common.CodeMessager {
	var r surpriseMeResponse
	r.Code = 200
	for i := 0; i < 100; i++ {
		r.Cats[i] = surpriseMe()
	}
	return r
}
