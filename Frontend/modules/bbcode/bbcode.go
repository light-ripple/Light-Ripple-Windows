// Package bbcode implements BBCode compiling for Hanayo.
package bbcode

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/frustra/bbcode"
	"github.com/microcosm-cc/bluemonday"
)

var bbcodeCompiler = func() bbcode.Compiler {
	compiler := bbcode.NewCompiler(true, true)
	compiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ul"
		style := node.GetOpeningTag().Value
		switch style {
		case "a":
			out.Attrs["style"] = "list-style-type: lower-alpha;"
		case "A":
			out.Attrs["style"] = "list-style-type: upper-alpha;"
		case "i":
			out.Attrs["style"] = "list-style-type: lower-roman;"
		case "I":
			out.Attrs["style"] = "list-style-type: upper-roman;"
		case "1":
			out.Attrs["style"] = "list-style-type: decimal;"
		default:
			out.Attrs["style"] = "list-style-type: disc;"
		}

		if len(node.Children) == 0 {
			out.AppendChild(bbcode.NewHTMLTag(""))
		} else {
			node.Info = []*bbcode.HTMLTag{out, out}
			tags := node.Info.([]*bbcode.HTMLTag)
			for _, child := range node.Children {
				curr := tags[1]
				curr.AppendChild(node.Compiler.CompileTree(child))
			}
			if len(tags[1].Children) > 0 {
				last := tags[1].Children[len(tags[1].Children)-1]
				if len(last.Children) > 0 && last.Children[len(last.Children)-1].Name == "br" {
					last.Children[len(last.Children)-1] = bbcode.NewHTMLTag("")
				}
			} else {
				tags[1].AppendChild(bbcode.NewHTMLTag(""))
			}
		}
		return out, false
	})

	compiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		parent := node.Parent
		for parent != nil {
			if parent.ID == bbcode.OPENING_TAG && parent.GetOpeningTag().Name == "list" {
				out := bbcode.NewHTMLTag("")
				out.Name = "li"
				tags := parent.Info.([]*bbcode.HTMLTag)
				if len(tags[1].Children) > 0 {
					last := tags[1].Children[len(tags[1].Children)-1]
					if len(last.Children) > 0 && last.Children[len(last.Children)-1].Name == "br" {
						last.Children[len(last.Children)-1] = bbcode.NewHTMLTag("")
					}
				} else {
					tags[1].AppendChild(bbcode.NewHTMLTag(""))
				}
				tags[1] = out
				tags[0].AppendChild(out)

				if len(parent.Children) == 0 {
					out.AppendChild(bbcode.NewHTMLTag(""))
				} else {
					for _, child := range node.Children {
						curr := tags[1]
						curr.AppendChild(node.Compiler.CompileTree(child))
					}
				}
				if node.ClosingTag != nil {
					tag := bbcode.NewHTMLTag(node.ClosingTag.Raw)
					bbcode.InsertNewlines(tag)
					out.AppendChild(tag)
				}
				return nil, false
			}
			parent = parent.Parent
		}
		return bbcode.DefaultTagCompiler(node)
	})

	compiler.SetTag("youtube", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		var youtubeID string

		content := bbcode.CompileText(node)
		youtubeLink, err := url.Parse(content)
		if err != nil {
			youtubeID = content
		} else {
			youtubeID = youtubeLink.Query().Get("v")
			if youtubeID == "" {
				youtubeID = content
			}
		}

		tag := bbcode.NewHTMLTag("")
		tag.Name = "iframe"
		tag.Attrs = map[string]string{
			"style":           "width: 100%; max-height: 100%;",
			"src":             "https://www.youtube.com/embed/" + youtubeID,
			"frameborder":     "0",
			"allowfullscreen": "",
		}
		tag.AppendChild(nil)

		container := bbcode.NewHTMLTag("")
		container.Name = "div"
		container.Attrs["class"] = "youtube video container"
		container.AppendChild(tag)

		return container, false
	})

	compiler.SetTag("left", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		out.Attrs["style"] = "text-align: left;"
		return out, true
	})
	compiler.SetTag("right", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		out.Attrs["style"] = "text-align: right;"
		return out, true
	})

	compiler.SetTag("container", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		args := node.GetOpeningTag().Args
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		out.Attrs["style"] = ""
		out.Attrs["class"] = ""
		if _, err := strconv.Atoi(args["width"]); err == nil {
			out.Attrs["style"] += "width: " + args["width"] + "px;"
		}
		if args["compact"] != "" {
			out.Attrs["class"] += "compact-container "
		}
		if args["center"] != "" {
			out.Attrs["style"] += "margin: 0 auto;"
		}
		return out, true
	})

	compiler.SetTag("hr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		out.Attrs["class"] = "ui divider"
		out.AppendChild(nil)
		return out, false
	})

	compiler.SetTag("email", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "a"
		val := node.GetOpeningTag().Value
		if val == "" {
			out.Attrs["href"] = "mailto:" + bbcode.CompileText(node)
			out.AppendChild(bbcode.NewHTMLTag(bbcode.CompileText(node)))
			return out, false
		}
		out.Attrs["href"] = "mailto:" + val
		return out, true
	})

	compiler.SetTag("size", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "span"
		if size, err := strconv.Atoi(node.GetOpeningTag().Value); err == nil && size > 0 {
			if size > 15 {
				size = 15
			}
			out.Attrs["style"] = fmt.Sprintf("font-size: %dpt; line-height: %[1]dpt;", size*6)
		}
		return out, true
	})

	return compiler
}()

var emojis = []string{
	"peppy",
	"barney",
	"akerino",
	"foka",
	"kappy",
	"creepypeppy",
	"peppyfiero",
	"djpeppy",
	"kappa",
}
var emojiReplacer = func() *strings.Replacer {
	var list []string
	for _, e := range emojis {
		list = append(list, ":"+e+":", "[img=/static/emotes/"+e+".png]:"+e+":[/img]")
	}
	return strings.NewReplacer(list...)
}()

// Compile takes some BBCode and converts it into safe HTML output.
func Compile(s string) string {
	s = emojiReplacer.Replace(s)
	s = strings.TrimSpace(s)
	return mondaySanitise(bbcodeCompiler.Compile(s))
}

var policy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("style", "class").Globally()
	p.AllowElements("iframe")
	p.AllowAttrs("style", "src", "frameborder", "allowfullscreen").OnElements("iframe")
	return p
}()

func mondaySanitise(source string) string {
	return policy.Sanitize(source)
}
