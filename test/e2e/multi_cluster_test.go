//go:build multicluster
// +build multicluster

package e2e

import (
	"testing"

	"github.com/argoproj/argo-workflows/v3/test/e2e/fixtures"
	"github.com/stretchr/testify/suite"
)

type MultiClusterSuite struct {
	fixtures.E2ESuite
}

func (s *MultiClusterSuite) TestMultiCluster() {
	s.Given().
		Workflow(`
metadata:
  generateName: multi-cluster-
spec:
  entrypoint: main
  artifactRepositoryRef:
    key: empty
  templates:
    - name: main
      cluster: other
      namespace: default
      container:
        image: argoproj/argosay:v2
`).
		When().
		SubmitWorkflow().
		WaitForWorkflow(fixtures.ToBeSucceeded)
}

func (s *MultiClusterSuite) TestDisallowedLocalNamespace() {
	s.Given().
		Workflow(`
metadata:
  generateName: multi-cluster-
spec:
  entrypoint: main
  artifactRepositoryRef:
    key: empty
  templates:
    - name: main
      namespace: default
      container:
        image: argoproj/argosay:v2
`).
		When().
		SubmitWorkflow().
		WaitForWorkflow(fixtures.ToBeErrored).
		Then().
		ExpectWorkflow(fixtures.StatusMessageContains(`namespace "argo" is forbidden from creating resources in cluster "" namespace "default"`))
}

func (s *MultiClusterSuite) TestDisallowedCluster() {
	s.Given().
		Workflow(`
metadata:
  generateName: multi-cluster-
spec:
  entrypoint: main
  artifactRepositoryRef:
    key: empty
  templates:
    - name: main
      cluster: other
      container:
        image: argoproj/argosay:v2
`).
		When().
		SubmitWorkflow().
		WaitForWorkflow(fixtures.ToBeErrored).
		Then().
		ExpectWorkflow(fixtures.StatusMessageContains(`namespace "argo" is forbidden from creating resources in cluster "other" namespace "argo"`))
}

func TestMultiClusterSuite(t *testing.T) {
	suite.Run(t, new(MultiClusterSuite))
}
