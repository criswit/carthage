package main // import "github.com/criswit/carthage"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configPath := flag.String("config-path", "", "Path to config")
	sysRootPath := flag.String("sys-root", "/tmp/carthage", "System root. '/' for non-dry run")
	flag.Parse()

	if *configPath == "" {
		log.Fatalf("--config-path is a required flag")
	}
	log.Printf("Starting Carthage with config %q", *configPath)

	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	cfg, err := ReadConfig(b)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	fs := NewFileSet(*sysRootPath)
	_ = fs

	specIndex := map[string]*UnitSpec{}
	unitIndex := map[string]*Unit{}
	for _, spec := range cfg {
		specIndex[spec.Name] = spec
		unitIndex[spec.Name] = &Unit{
			Name: spec.Name,
		}
	}

	root := &Unit{
		Name: "root",
	}
	for name, unit := range unitIndex {
		spec := specIndex[name]
		for _, dname := range spec.Deps {
			unit.AddDependency(unitIndex[dname])
		}
		root.AddDependency(unit)
	}

	plan, err := NewExecutionPlan(root)
	if err != nil {
		log.Fatalf("failed to create execution plan: %v", err)
	}
	plan.String()
	for _, unit := range plan.units {
		fmt.Println(unit.Name)
	}
}
