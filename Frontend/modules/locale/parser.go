package locale

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type po struct {
	Headers      map[string]string
	Translations map[string]string
}

// A minimal .po file parser.
func parse(s *bufio.Scanner) *po {
	p := &po{
		Headers:      make(map[string]string, 20),
		Translations: make(map[string]string, 500),
	}

	const (
		currentNothing = iota
		currentMsgID
		currentMsgString
	)
	var current byte
	var currentID string
	var currentString string

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if line == "" || line[0] == '#' {
			continue
		}

		switch {
		case strings.HasPrefix(line, "msgid "):
			unq, err := strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(line, "msgid")))
			if err != nil {
				fmt.Println(line)
				fmt.Println(err)
				return nil
			}

			if current != currentNothing && currentID == "" && currentString != "" {
				for _, h := range strings.Split(currentString, "\n") {
					if h == "" {
						continue
					}
					parts := strings.SplitN(h, ": ", 2)
					if len(parts) != 2 {
						continue
					}
					p.Headers[parts[0]] = parts[1]
				}
			} else {
				p.Translations[currentID] = currentString
			}

			currentID = unq
			current = currentMsgID
		case strings.HasPrefix(line, "msgstr "):
			unq, err := strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(line, "msgstr")))
			if err != nil {
				fmt.Println(line)
				fmt.Println(err)
				return nil
			}
			currentString = unq
			current = currentMsgString
		case strings.HasPrefix(line, "msgid_plural "), strings.HasPrefix(line, "msgstr["):
			// We currently don't support plural in the Go parser, at least hold it off 'til it's needed.
			current = currentNothing
		default:
			if current == currentNothing {
				continue
			}

			unq, err := strconv.Unquote(strings.TrimSpace(line))
			if err != nil {
				fmt.Println(line)
				fmt.Println(err)
				return nil
			}
			switch current {
			case currentMsgID:
				currentID += unq
			case currentMsgString:
				currentString += unq
			}
		}
	}

	if current == currentMsgString {
		p.Translations[currentID] = currentString
	}
	return p
}

func parseFile(fName string) (*po, error) {
	f, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	p := parse(s)
	return p, nil
}
