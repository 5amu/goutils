package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/5amu/goutils/shuttles"
)

func usage() {
	fmt.Println("Usage: shuttles [-h|-p|-f...] <ELSE...>")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("    -h    show help and exit")
	fmt.Println("    -p    parallel processes to run")
	fmt.Println("    -f    specify file(s) to read from, will be")
	fmt.Println("          [^SHT0^,^SHT1^...] in the provided command")
	fmt.Println("")
	fmt.Println("POSITIONAL:")
	fmt.Println("    whatever you'll put after recognized flags will")
	fmt.Println("    be parsed as a command to encapsulate and launched")
	fmt.Println("    as a shuttle in the shuttle factory :)")
	fmt.Println("")
}

func loadFuel(f shuttles.ShuttleFactory, fname string, index int) error {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	lines := strings.Split("\n", string(data))
	for _, line := range lines {
		if err := f.SupplyLane(line, index); err != nil {
			return err
		}
	}
	f.Stop()
	return nil
}

func main() {
	var lanes int
	var fuel []string
	var help bool

	shuttleFlagSet := flag.NewFlagSet("shuttles", flag.ExitOnError)
	shuttleFlagSet.IntVar(&lanes, "p", 2, "parallel processes to run")
	shuttleFlagSet.Func("f", "specify file(s) to read from, will be [^SHT0^,^SHT1^...] in the provided command", func(s string) error {
		fuel = append(fuel, s)
		return nil
	})
	shuttleFlagSet.BoolVar(&help, "h", false, "show help and exit")

	if len(os.Args) < 2 {
		usage()
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	shuttleFlagSet.Parse(os.Args[1:])

	if help || os.Args[1] == "help" {
		usage()
		os.Exit(0)
	}

	if len(fuel) == 0 {
		usage()
		fmt.Println("missing fuel: -f flag")
		os.Exit(1)
	}

	factory := shuttles.NewShuttleFactory(strings.Join(shuttleFlagSet.Args(), " "), lanes)
	for index, fname := range fuel {
		go func(n string, i int) {
			if err := loadFuel(*factory, n, i); err != nil {
				panic(err)
			}
		}(fname, index)
	}

	if err := factory.Start(context.Background(), lanes); err != nil {
		panic(err)
	}
}
