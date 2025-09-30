package main

type Route struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Type    string `yaml:"type"`
	Backend string `yaml:"backend"`
}
