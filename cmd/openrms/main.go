package main

import (
	"flag"
	"github.com/goioc/di"
	"github.com/qvistgaard/openrms/internal/config"
	"github.com/qvistgaard/openrms/internal/postprocess"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

var wg sync.WaitGroup

func init() {

	di.RegisterBeanFactory("config", di.Singleton, func() (interface{}, error) {
		flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
		file, err := ioutil.ReadFile(*flagConfig)
		if err != nil {
			return nil, err
		}

		c := make(map[string]interface{})
		err = yaml.Unmarshal(file, c)
		if err != nil {
			return nil, err
		}
		return c, nil
	})

	log.Info(c)
	os.Exit(0)
}

func main() {
	flagConfig := flag.String("config", "config.yaml", "OpenRMS Config file")
	file, err := ioutil.ReadFile(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Make configurable
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)

	implement, err := config.CreateImplementFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	rules, err := config.CreateRaceRulesFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	repository, err := config.CreateCarRepositoryFromConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	course, err := config.CreateCourseFromConfig(file, rules)
	if err != nil {
		log.Fatal(err)
	}

	processors, err := config.CreatePostProcessors(file)
	if err != nil {
		log.Fatal(err)
	}
	postProcess := postprocess.CreatePostProcess(processors)

	wg.Add(1)
	go eventloop(implement)

	wg.Add(1)
	go processEvents(implement, postProcess, repository, course, rules)

	wg.Wait()
}
