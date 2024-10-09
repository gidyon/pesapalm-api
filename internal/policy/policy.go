package policy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Options struct {
	Enforcer *casbin.Enforcer
}

type APIServer struct {
	*Options
}

// NewPolicyAPI creates a casbin policy manager singleton
func NewPolicyAPI(ctx context.Context, opt *Options) (_ *APIServer, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to start casbin policy service: %v", err)
		}
	}()

	// Validation
	switch {
	case ctx == nil:
		err = errors.New("missing context")
	case opt == nil:
		err = errors.New("missing options")
	case opt.Enforcer == nil:
		err = errors.New("missing enforcer")
	}
	if err != nil {
		return nil, err
	}

	// Account API
	api := &APIServer{
		Options: opt,
	}

	return api, nil
}

func (api *APIServer) AddPolicies(c *gin.Context) {
	var policies []PolicyRequest

	if err := c.ShouldBindJSON(&policies); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}

	rules := make([][]string, 0, len(policies))

	for _, policy := range policies {
		switch {
		case policy.Policy == "":
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing policy"})
			return
		case policy.Action == "":
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing action"})
			return
		case policy.Subject == "":
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing subject"})
			return
		case policy.Object == "":
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing object"})
			return
		}

		rules = append(rules, []string{policy.Subject, policy.Object, policy.Action})
	}

	ok, err := api.Enforcer.AddPolicies(rules)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to add policy", "details": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "policy added", "success": ok})
}
