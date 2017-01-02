package main

import (
	"encoding/binary"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/twtiger/gosecco"
	"github.com/twtiger/gosecco/parser"
)

type pathList []string

func (i *pathList) String() string {
	return strings.Join((*i)[:], ",")
}

func (i *pathList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func buildSeccomp(source parser.Source, f *os.File) error {
	settings := gosecco.SeccompSettings{
		DefaultPositiveAction: "allow",
		DefaultNegativeAction: "ENOSYS",
		DefaultPolicyAction:   "ENOSYS",
		ActionOnX32:           "kill",
		ActionOnAuditFailure:  "kill",
	}

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
	var sourcePaths pathList
	flag.Var(&sourcePaths, "rules", "path to file containing rules in gosecco format")
	flag.Parse()

	if len(sourcePaths) < 1 {
		log.Fatal("At least one file must be provided.")
	}

	var combined parser.CombinedSource
	for _, path := range sourcePaths {
		combined.Sources = append(combined.Sources, &parser.FileSource{path})
	}

	if err := buildSeccomp(&combined, os.Stdout); err != nil {
		log.Fatalf("%v", err)
	}
}
