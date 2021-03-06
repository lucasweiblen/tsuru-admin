// Copyright 2014 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tsuru/tsuru/app"
	"github.com/tsuru/tsuru/cmd"
	"launchpad.net/gnuflag"
)

type planCreate struct {
	memory     int64
	swap       int64
	cpushare   int
	setDefault bool
	router     string
	fs         *gnuflag.FlagSet
}

func (c *planCreate) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("plan-create", gnuflag.ExitOnError)
		memory := "Amount of available memory for units in bytes."
		c.fs.Int64Var(&c.memory, "memory", 0, memory)
		c.fs.Int64Var(&c.memory, "m", 0, memory)
		swap := "Amount of available swap space for units in bytes."
		c.fs.Int64Var(&c.swap, "swap", 0, swap)
		c.fs.Int64Var(&c.swap, "s", 0, swap)
		cpushare := "Relative cpu share each unit will have available."
		c.fs.IntVar(&c.cpushare, "cpushare", 0, cpushare)
		c.fs.IntVar(&c.cpushare, "c", 0, cpushare)
		setDefault := "Set plan as default, this will remove the default flag from any other plan."
		c.fs.BoolVar(&c.setDefault, "default", false, setDefault)
		c.fs.BoolVar(&c.setDefault, "d", false, setDefault)
		router := "The name of the router used by this plan."
		c.fs.StringVar(&c.router, "router", "", router)
		c.fs.StringVar(&c.router, "r", "", router)

	}
	return c.fs
}

func (c *planCreate) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "plan-create",
		Usage:   "plan-create <name> -c cpushare [-m memory] [-s swap] [-r router] [--default]",
		Desc:    `Creates a new plan for being used when creating apps.`,
		MinArgs: 1,
	}
}

func (c *planCreate) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL("/plans")
	if err != nil {
		return err
	}
	plan := app.Plan{
		Name:     context.Args[0],
		Memory:   c.memory,
		Swap:     c.swap,
		CpuShare: c.cpushare,
		Default:  c.setDefault,
		Router:   c.router,
	}
	planData, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	body := strings.NewReader(string(planData))
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	if err != nil {
		fmt.Fprintf(context.Stdout, "Failed to create plan!\n")
		return err
	}
	fmt.Fprintf(context.Stdout, "Plan successfully created!\n")
	return nil
}

type planRemove struct{}

func (c *planRemove) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "plan-remove",
		Usage:   "plan-remove <name>",
		Desc:    `Removes a plan from the database.`,
		MinArgs: 1,
	}
}

func (c *planRemove) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL("/plans/" + context.Args[0])
	if err != nil {
		return err
	}
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	if err != nil {
		fmt.Fprintf(context.Stdout, "Failed to remove plan!\n")
		return err
	}
	fmt.Fprintf(context.Stdout, "Plan successfully removed!\n")
	return nil
}
