package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

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

func loadFuel(f *shuttles.ShuttleFactory, fname string, index int) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if err := f.SupplyLane(scanner.Text(), index); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	f.Stop()
	return nil
}

func main() {
	var lanes int
	var fuel []string
	var help bool
	var out string

	shuttleFlagSet := flag.NewFlagSet("shuttles", flag.ExitOnError)
	shuttleFlagSet.IntVar(&lanes, "p", 2, "parallel processes to run")
	shuttleFlagSet.Func("f", "specify file(s) to read from, will be [^SHT0^,^SHT1^...] in the provided command", func(s string) error {
		fuel = append(fuel, s)
		return nil
	})
	shuttleFlagSet.StringVar(&out, "o", "", "file to write (json) results to")
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

	factory := shuttles.NewShuttleFactory(shuttleFlagSet.Args(), lanes)
	for index, fname := range fuel {
		go func(n string, i int) {
			if err := loadFuel(factory, n, i); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}(fname, index)
	}

	if err := factory.Start(context.Background(), lanes); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outputs := factory.GetShuttleOutputs()
	for _, v := range outputs {
		fmt.Println("===== BEGIN OF SHUTTLE", v.ID, "OUTPUT =====")
		fmt.Println("=====", v.Arguments, "===>", v.Injected, "=====")
		fmt.Println(v.Output)
		fmt.Println("===== END OF SHUTTLE", v.ID, "OUTPUT =====")
		fmt.Println("")
	}

	if out != "" {
		if err := ioutil.WriteFile(out, []byte(out), 0644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
