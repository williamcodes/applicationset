package applicationsets

import (
	"context"
	"time"

	"github.com/argoproj-labs/applicationset/test/e2e/fixture/applicationsets/utils"
	argov1alpha1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/argoproj/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// this implements the "then" part of given/when/then
type Consequences struct {
	context *Context
	actions *Actions
}

func (c *Consequences) Expect(e Expectation) *Consequences {
	// this invocation makes sure this func is not reported as the cause of the failure - we are a "test helper"
	c.context.t.Helper()
	var message string
	var state state
	timeout := time.Duration(30) * time.Second
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(3 * time.Second) {
		state, message = e(c)
		switch state {
		case succeeded:
			log.Infof("expectation succeeded: %s", message)
			return c
		case failed:
			c.context.t.Fatalf("failed expectation: %s", message)
			return c
		}
		log.Infof("expectation pending: %s", message)
	}
	c.context.t.Fatal("timeout waiting for: " + message)
	return c
}

func (c *Consequences) And(block func()) *Consequences {
	c.context.t.Helper()
	block()
	return c
}

func (c *Consequences) Given() *Context {
	return c.context
}

func (c *Consequences) When() *Actions {
	return c.actions
}

func (c *Consequences) app(name string) *argov1alpha1.Application {
	apps := c.apps()

	for index, app := range apps {
		if app.Name == name {
			return &apps[index]
		}
	}

	return nil
}

func (c *Consequences) apps() []argov1alpha1.Application {

	fixtureClient := utils.GetE2EFixtureK8sClient()
	list, err := fixtureClient.AppClientset.ArgoprojV1alpha1().Applications(utils.ArgoCDNamespace).List(context.Background(), v1.ListOptions{})
	errors.CheckError(err)

	if list == nil {
		list = &argov1alpha1.ApplicationList{}
	}

	return list.Items
}
