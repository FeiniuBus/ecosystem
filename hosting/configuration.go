package hosting

import (
	"bytes"
	"encoding/json"

	"regexp"

	"github.com/go-ini/ini"
)

type Configuration struct {
	Sections []*Section
}

func (cfg *Configuration) Section(name string) *Section {
	for _, section := range cfg.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

func (cfg *Configuration) Remove(name string) {
	len := len(cfg.Sections)
	for i, section := range cfg.Sections {
		if section.Name == name {
			if i >= len-1 {
				cfg.Sections = cfg.Sections[:i]
			} else {
				cfg.Sections = append(cfg.Sections[:i], cfg.Sections[i+1:]...)
			}
		}
	}
}

func (cfg *Configuration) Object(v interface{}) error {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, section := range cfg.Sections {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\"")
		buffer.WriteString(section.Name)
		buffer.WriteString("\":{")
		for j, item := range section.Items {
			if j > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString("\"")
			buffer.WriteString(item.Key)
			buffer.WriteString("\":")
			oriStr := item.Value.(string)
			if ok, _ := regexp.MatchString("^(-|\\+)?\\d+(\\.\\d+)?$", oriStr); ok {
				buffer.WriteString(oriStr)
			} else if oriStr == "true" || oriStr == "false" {
				buffer.WriteString(oriStr)
			} else {
				buffer.WriteString("\"")
				buffer.WriteString(oriStr)
				buffer.WriteString("\"")
			}
			j++
		}
		buffer.WriteString("}")
	}
	buffer.WriteString("}")
	err := json.Unmarshal(buffer.Bytes(), v)
	return err
}

type Section struct {
	Name  string
	Level int
	Items []*Item
}

func (s *Section) Item(key string) interface{} {
	for _, item := range s.Items {
		if item.Key == key {
			return item
		}
	}
	return nil
}

func (s *Section) Remove(key string) {
	len := len(s.Items)
	for i, item := range s.Items {
		if item.Key == key {
			if i >= len-1 {
				s.Items = s.Items[:i]
			} else {
				s.Items = append(s.Items[:i], s.Items[i+1:]...)
			}
		}
	}
}

func (s *Section) Merge(other *Section) *Section {
	if s.Name == other.Name {
		target := &Section{
			Level: 0,
			Name:  s.Name,
		}
		section := make(map[string]*Item)
		for _, v := range s.Items {
			section[v.Key] = v
		}
		for _, v := range other.Items {
			section[v.Key] = v
		}
		for _, item := range section {
			target.Items = append(target.Items, item)
		}
		return target
	}
	return other
}

type Item struct {
	Key   string
	Value interface{}
}

type SectionPointerSlice []*Section

func (s SectionPointerSlice) Len() int {
	return len(s)
}
func (s SectionPointerSlice) Less(i, j int) bool {
	if s[i].Name == s[j].Name {
		return s[i].Level > s[j].Level
	}
	return s[i].Name < s[j].Name
}
func (s SectionPointerSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ItemPointerSlice []*Item

func (s ItemPointerSlice) Len() int {
	return len(s)
}
func (s ItemPointerSlice) Less(i, j int) bool {
	return s[i].Key < s[j].Key
}
func (s ItemPointerSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewConfiguration(sections []*Section) *Configuration {
	c := new(Configuration)
	c.Sections = sections
	return c
}
func NewConfigurationFromFile(file *ini.File) *Configuration {
	c := new(Configuration)
	sections := file.Sections()
	c.Sections = make([]*Section, len(sections))
	for i, section := range sections {
		keys := section.Keys()
		s := new(Section)
		s.Name = section.Name()
		s.Items = make([]*Item, len(keys))
		for i, key := range keys {
			item := new(Item)
			item.Key = key.Name()
			item.Value = key.Value()
			s.Items[i] = item
		}
		c.Sections[i] = s
	}
	return c
}
