package hosting

import (
	"errors"
	"log"
	"path"
	"sort"

	"github.com/go-ini/ini"
)

type ConfigurationBuilder struct {
	Config  *Configuration
	baseDir string
	sources []*Configuration
}

func (this *ConfigurationBuilder) AddIniFile(fileName string) *ConfigurationBuilder {
	file, err := ini.Load(this.getFilePath(fileName))
	if err != nil {
		strError := "configuration file load error：" + err.Error()
		log.Fatal(strError)
		panic(strError)
	}
	this.sources = append(this.sources, NewConfigurationFromFile(file))
	return this
}

func (this *ConfigurationBuilder) AddIniFileOptional(fileName string) *ConfigurationBuilder {
	file, err := ini.Load(this.getFilePath(fileName))
	if err == nil {
		this.sources = append(this.sources, NewConfigurationFromFile(file))
	}
	return this
}

func (this *ConfigurationBuilder) Build() (*Configuration, error) {
	l := len(this.sources)
	if l <= 0 {
		return nil, errors.New("none config file to build")
	}
	var allSections SectionPointerSlice
	for i, c := range this.sources {
		for _, section := range c.Sections {
			section.Level = l - i
			allSections = append(allSections, section)
		}
	}
	sort.Sort(allSections)
	targetSectionsMap := make(map[string]*Section)
	for i, section := range allSections {
		if i == 0 {
			targetSectionsMap[section.Name] = section
		} else {
			mergeSection := allSections[i-1].Merge(section)
			targetSectionsMap[mergeSection.Name] = mergeSection
		}
	}
	targetSections := make([]*Section, len(targetSectionsMap))
	i := 0
	for _, section := range targetSectionsMap {
		targetSections[i] = section
		i++
	}
	this.Config = NewConfiguration(targetSections)
	return this.Config, nil
}

func (this *ConfigurationBuilder) BuildToObject(v interface{}) error {
	cfg, err := this.Build()
	if err != nil {
		return err
	}
	return cfg.Object(v)
}

func (this *ConfigurationBuilder) SetBasePath(path string) *ConfigurationBuilder {
	this.baseDir = path
	return this
}

func NewConfigurationBuilder() *ConfigurationBuilder {
	builder := new(ConfigurationBuilder)
	builder.sources = make([]*Configuration, 0, 3)
	return builder
}

func (this *ConfigurationBuilder) getFilePath(fileName string) string {
	return path.Join(this.baseDir, fileName)
}
