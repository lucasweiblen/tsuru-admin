// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/tsuru/tsuru/app"
	"github.com/tsuru/tsuru/cmd"
	"github.com/tsuru/tsuru/cmd/cmdtest"
	"gopkg.in/check.v1"
)

func (s *S) TestPlanCreateInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "plan-create",
		Usage:   "plan-create <name> -c cpushare [-m memory] [-s swap] [-r router] [--default]",
		Desc:    "Creates a new plan for being used when creating apps.",
		MinArgs: 1,
	}
	c.Assert((&planCreate{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestPlanCreate(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Args:   []string{"myplan"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "", Status: http.StatusCreated},
		CondFunc: func(req *http.Request) bool {
			var plan app.Plan
			err := json.NewDecoder(req.Body).Decode(&plan)
			c.Assert(err, check.IsNil)
			expected := app.Plan{
				Name:     "myplan",
				Memory:   0,
				Swap:     0,
				CpuShare: 100,
				Default:  false,
				Router:   "",
			}
			c.Assert(plan, check.DeepEquals, expected)
			return req.URL.Path == "/plans" && req.Method == "POST"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	command := planCreate{}
	command.Flags().Parse(true, []string{"-c", "100"})
	err := command.Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, "Plan successfully created!\n")
}

func (s *S) TestPlanCreateFlags(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Args:   []string{"myplan"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "", Status: http.StatusCreated},
		CondFunc: func(req *http.Request) bool {
			var plan app.Plan
			err := json.NewDecoder(req.Body).Decode(&plan)
			c.Assert(err, check.IsNil)
			expected := app.Plan{
				Name:     "myplan",
				Memory:   1024,
				Swap:     512,
				CpuShare: 100,
				Default:  true,
				Router:   "myrouter",
			}
			c.Assert(plan, check.DeepEquals, expected)
			return req.URL.Path == "/plans" && req.Method == "POST"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	command := planCreate{}
	command.Flags().Parse(true, []string{"-c", "100", "-m", "1024", "-s", "512", "-d", "-r", "myrouter"})
	err := command.Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, "Plan successfully created!\n")
}

func (s *S) TestPlanCreateError(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Args:   []string{"myplan"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "", Status: http.StatusConflict},
		CondFunc: func(req *http.Request) bool {
			return req.URL.Path == "/plans" && req.Method == "POST"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	command := planCreate{}
	err := command.Run(&context, client)
	c.Assert(err, check.NotNil)
	c.Assert(stdout.String(), check.Equals, "Failed to create plan!\n")
}

func (s *S) TestPlanRemoveInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "plan-remove",
		Usage:   "plan-remove <name>",
		Desc:    "Removes a plan from the database.",
		MinArgs: 1,
	}
	c.Assert((&planRemove{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestPlanRemove(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Args:   []string{"myplan"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "", Status: http.StatusOK},
		CondFunc: func(req *http.Request) bool {
			return req.URL.Path == "/plans/myplan" && req.Method == "DELETE"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	command := planRemove{}
	err := command.Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, "Plan successfully removed!\n")
}

func (s *S) TestPlanRemoveError(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Args:   []string{"myplan"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "", Status: http.StatusInternalServerError},
		CondFunc: func(req *http.Request) bool {
			return req.URL.Path == "/plans/myplan" && req.Method == "DELETE"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	command := planRemove{}
	err := command.Run(&context, client)
	c.Assert(err, check.NotNil)
	c.Assert(stdout.String(), check.Equals, "Failed to remove plan!\n")
}
