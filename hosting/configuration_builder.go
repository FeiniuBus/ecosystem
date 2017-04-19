package hosting

import (
	"path"

	"sort"

	"github.com/go-ini/ini"
)

type ConfigurationBuilder struct {
	BaseDir string
	Config  *Configuration
	sources []*Configuration
}

func (this *ConfigurationBuilder) AddIniFile(fileName string) *ConfigurationBuilder {
	file, err := ini.Load(this.getFilePath(fileName))
	if err != nil {
		panic("File load error：" + err.Error())
	}
	this.sources = append(this.sources, NewConfigurationFromFile(file))
	return this
}

func (this *ConfigurationBuilder) Build() *Configuration {
	var allSections SectionPointerSlice
	var targetSections []*Section
	len := len(this.sources)
	for i, c := range this.sources {
		for _, section := range c.Sections {
			section.Level = len - i
			allSections = append(allSections, section)
		}
	}
	sort.Sort(allSections)
	for i, section := range allSections {
		if i > 0 {
			merge := section.Merge(allSections[i-1])
			if merge != nil {
				targetSections = append(targetSections, merge)
			}
		}
	}
	this.Config = NewConfiguration(targetSections)
	return this.Config
}

func (this *ConfigurationBuilder) SetBasePath(path string) *ConfigurationBuilder {
	this.BaseDir = path
	return this
}

func NewConfigurationBuilder() *ConfigurationBuilder {
	builder := new(ConfigurationBuilder)
	builder.sources = make([]*Configuration, 0, 3)
	return builder
}

func (this *ConfigurationBuilder) getFilePath(fileName string) string {
	return path.Join(this.BaseDir, fileName)
}
