package main

import (
	"encoding/binary"
	"flag"
	"log"
	"os"

	"github.com/twtiger/gosecco"
	"github.com/twtiger/gosecco/parser"
)

func buildSeccomp(path string, f *os.File) error {
	settings := gosecco.SeccompSettings{
		DefaultPositiveAction: "allow",
		DefaultNegativeAction: "ENOSYS",
		DefaultPolicyAction:   "ENOSYS",
		ActionOnX32:           "kill",
		ActionOnAuditFailure:  "kill",
	}

	source := &parser.FileSource{path}
	bpf, err := gosecco.PrepareSource(source, settings)
	if err != nil {
		return err
	}

	for _, rule := range bpf {
		if err := binary.Write(f, binary.LittleEndian, rule); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var rulePath string
	flag.StringVar(&rulePath, "rules", "", "path to file containing rules in gosecco format")
	flag.Parse()

	if len(rulePath) < 1 {
		log.Fatal("The path to a file containing rules must be provided")
	}

	if err := buildSeccomp(rulePath, os.Stdout); err != nil {
		log.Fatalf("%v", err)
	}
}
