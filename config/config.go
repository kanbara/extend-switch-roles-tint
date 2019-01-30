package config

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"image/color"
	"reflect"
	"regexp"
	"strings"
)

type Config struct {
	Entries  []Entry
}

type Entry struct {
	AccountName string
	NotColourData []string
	Colour color.Color
}

func (c* Config) String() string {
	text := ""
	for _, e := range c.Entries {
		text += e.String() + "\n"
	}

	return text
}

func (e* Entry) String() string {
	accountNameStr := fmt.Sprintf("[%v]\n", e.AccountName)
	otherData := ""
	for _, str := range e.NotColourData {
		otherData += str + "\n"
	}

	colourStr := ""

	if e.Colour != nil {
		colorfulColour := e.Colour.(colorful.Color)
		hex := strings.ToLower(colorfulColour.Hex()[1:])
		colourStr = fmt.Sprintf("color = %v\n", hex)
	} else {
		colourStr = fmt.Sprintf("colour = %v\n", e.Colour)
	}

	return accountNameStr + otherData + colourStr
}

func parseAccountName(pattern *regexp.Regexp, accountString string) string {

	matches := pattern.FindSubmatch([]byte(accountString))
	accountName := matches[1]
	return string(accountName)
}

func parseColour(pattern *regexp.Regexp, colourString string) color.Color {
	matches := pattern.FindSubmatch([]byte(colourString))
	colour := matches[1]

	// looks like #FFCC06, our input is like ffcc06
	hexCol := fmt.Sprintf("#%v", strings.ToUpper(string(colour)))
	return gamut.Hex(hexCol)
}

func Parse(data []string, debug bool) (Config, error) {

	// match something like [profile account], capturing account, with optional profile flag
	profileRegex := regexp.MustCompile(`^\[(?:profile\s*)?(.+)]$`)
	colourRegex := regexp.MustCompile(`^(?:color)\s?=\s?(\S{6})$`)

	config := Config{}

	currentEntry := Entry{}
	for _, line := range data {
		// we have 3 cases to consider: new profile, profile data, and colour
		switch {
		// we found a profile
		case profileRegex.Match([]byte(strings.TrimSpace(line))):

			if !reflect.DeepEqual(Entry{}, currentEntry) {
				if debug {
					fmt.Printf("writing entry %v\n", currentEntry)
				}
				config.Entries = append(config.Entries, currentEntry)
			}

			if debug {
				fmt.Printf("Found profile %v\n", line)
			}
			currentEntry = Entry{}
			currentEntry.AccountName = parseAccountName(profileRegex, line)
		// we found a colour
		case colourRegex.Match([]byte(strings.TrimSpace(line))):
			if debug {
				fmt.Printf("Found colour %v\n", line)
			}
			colour := parseColour(colourRegex, line)

			currentEntry.Colour = colour

		// don't add newlines ?
		case strings.TrimSpace(line) == "":
			continue
		// any other parts that aren't important
		default:
			if debug {
				fmt.Printf("Found random other stuff %v\n", line)
			}
			currentEntry.NotColourData = append(currentEntry.NotColourData, line)
		}
	}

	if !reflect.DeepEqual(Entry{}, currentEntry) {
		config.Entries = append(config.Entries, currentEntry)
	}

	return config, nil
}