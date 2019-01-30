package main

import (
	"bufio"
	"fmt"
	gcolor "github.com/gookit/color"
	"github.com/kanbara/extend-switch-roles-tint/config"
	"github.com/kanbara/extend-switch-roles-tint/strategy"
	"github.com/muesli/gamut"
	"gopkg.in/alecthomas/kingpin.v2"
	"image/color"
	"os"
)


var (
	app = kingpin.New("extend-switch-roles-tint", "Generate colours for your AWS Extend Switch Roles plugin")
	debug = kingpin.Flag("debug", "debug mode").Hidden().Default("false").Bool()
	show = app.Flag("show", "Show the colours (requires a true colour capable term like iTerm2, konsole, PuTTY").
		Default("false").Bool()
	generator = app.Flag("generator", "The type of palette to use").
		HintOptions("Pastel").
		HintOptions("Warm").String()
	outfile = app.Flag("outfile", "Output file").Short('o').Default("extend-roles.conf").String()

	fromfile = app.Command("from-file", "Generate colours and edit a config automatically")
	infile = fromfile.Arg("infile", "Input file").Required().ExistingFile()
	globs = fromfile.Flag("globs", "Wildcard globs to use to separate accounts by name").Short('g').Required().Strings()
    shade = fromfile.Flag("shade", "Use shading for all entries that match a glob, instead of the same colour for each").
    	Default("false").Bool()
	fromnum = app.Command("from-number", "Generate a certain number of colours")
	numcolours = fromnum.Arg("number", "Number of colours to generate").Required().Int()
)

func showColours (colours []color.Color) {
	for _, colour := range colours {
		foreR, foreG, foreB, _ := gamut.Contrast(colour).RGBA()
		r, g, b, _ := colour.RGBA()
		c := gcolor.NewRGBStyle(gcolor.RGB(uint8(foreR), uint8(foreG), uint8(foreB)), gcolor.RGB(uint8(r), uint8(g), uint8(b)))
		c.Printf("colour")
		fmt.Printf(" ")
	}
	fmt.Println()
}

func main() {

	kingpin.CommandLine.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case fromfile.FullCommand():
		file, err := os.Open(*infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var data []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}

		conf, err := config.Parse(data, *debug)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		strat := strategy.SplitStrategy{}
		newColours, newConf, err := strat.Replace(conf, *globs, *debug)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *show {
			showColours(newColours)
		}

		of, err := os.Create(*outfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		b := bufio.NewWriter(of)
		_, err = b.WriteString(newConf.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		b.Flush()

		fmt.Printf("Wrote new conf with %d new colours to %v. Enjoy!", len(newColours), *outfile)


	case fromnum.FullCommand():
		fmt.Printf("not done\n")
	}
}