package strategy

import (
	"fmt"
	"github.com/kanbara/extend-switch-roles-tint/config"
	"github.com/muesli/gamut"
	"github.com/ryanuber/go-glob"
	"image/color"
)

const defaultBin = ""

type Strategy interface {
	 Generate(conf config.Config, conditions []string) ([]color.Color, error)
	 GetBins(conf config.Config, conditions []string, debug bool) map[string][]config.Entry
}

type Bins map[string]Bin
type Bin []config.Entry
type SplitStrategy struct {}

func (s *SplitStrategy) GetBins(conf config.Config, globs []string, debug bool) (Bins, int) {
	bins := Bins{}
	numBins := 0

	for _, entry := range conf.Entries {
		globbed := false
		for _, g := range globs {
			if glob.Glob(g, entry.AccountName) {
				bins[g] = append(bins[g], entry)
				if debug {
					fmt.Printf("name %v matches condition %v\n", entry.AccountName, g)
				}
				globbed = true
				break
			}
		}
		if !globbed {
			bins[defaultBin] = append(bins[defaultBin], entry)
		}
	}

	for _, bin := range bins {
		if len(bin) != 0 {
			numBins+=1
		}
	}

	return bins, numBins
}

func (s *SplitStrategy) Replace(conf config.Config, globs []string, debug bool) ([]color.Color, config.Config, error) {

	bins, numBins := s.GetBins(conf, globs, debug)

	// check if glob fails, do we still generate a colour? add test for this
	// e.g. "*.dev" "*.int" with no "*.int" matches, we should know that and also know which globs to write

	if debug {
		for cond, bin := range bins {
			fmt.Printf("%v:\n", cond)
			for _, entry := range bin {
				fmt.Printf("\t%v\n", entry)
			}
		}

		fmt.Printf("Num bins: %v\n", numBins)
	}

	colours, err := gamut.Generate(numBins, gamut.PastelGenerator{})
	if err != nil {
		return nil, config.Config{}, err
	}

	curColor := 0
	for _, bin := range bins {
		newColour := colours[curColor]
		for _, entry := range bin {
			for i, origEntry := range conf.Entries {
				if entry.AccountName == origEntry.AccountName {
					origEntry.Colour = newColour
					conf.Entries[i] = origEntry
					}
			}
		}
		curColor += 1
	}
	
	return colours, conf, nil
}