package host

import (
	"errors"

	configapilatest "github.com/openshift/origin/pkg/cmd/server/api/latest"
	configvalidation "github.com/openshift/origin/pkg/cmd/server/api/validation"
	"github.com/openshift/origin/pkg/diagnostics/types"
)

// NodeConfigCheck is a Diagnostic to check that the node config file is valid
type NodeConfigCheck struct {
	NodeConfigFile string
}

const NodeConfigCheckName = "NodeConfigCheck"

func (d NodeConfigCheck) Name() string {
	return NodeConfigCheckName
}

func (d NodeConfigCheck) Description() string {
	return "Check the node config file"
}
func (d NodeConfigCheck) CanRun() (bool, error) {
	if len(d.NodeConfigFile) == 0 {
		return false, errors.New("must have node config file")
	}

	return true, nil
}
func (d NodeConfigCheck) Check() types.DiagnosticResult {
	r := types.NewDiagnosticResult(NodeConfigCheckName)
	r.Debugf("DH1001", "Looking for node config file at '%s'", d.NodeConfigFile)
	nodeConfig, err := configapilatest.ReadAndResolveNodeConfig(d.NodeConfigFile)
	if err != nil {
		r.Errorf("DH1002", err, "Could not read node config file '%s':\n(%T) %[2]v", d.NodeConfigFile, err)
		return r
	}

	r.Infof("DH1003", "Found a node config file: %[1]s", d.NodeConfigFile)

	for _, err := range configvalidation.ValidateNodeConfig(nodeConfig) {
		r.Errorf("DH1004", err, "Validation of node config file '%s' failed:\n(%T) %[2]v", d.NodeConfigFile, err)
	}
	return r
}
