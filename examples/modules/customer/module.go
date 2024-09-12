package main

import (
	"github.com/trinitytechnology/ebrick/module"
)

type CustomerModule struct {
}

// Install implements plugin.Plugin.
func (p *CustomerModule) Initialize(opt *module.Options) error {
	// Init Tables
	return nil
}

func (p *CustomerModule) Name() string {
	return "Customer Management"
}

func (p *CustomerModule) Version() string {
	return "1.0.0"
}

func (p *CustomerModule) Description() string {
	return "Customer Management"
}

func (p *CustomerModule) Id() string {
	return "Customer"
}

var Module CustomerModule
