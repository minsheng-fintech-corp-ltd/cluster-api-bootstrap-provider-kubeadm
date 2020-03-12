package controllers

import (
	"strings"

	"sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/ignition"
)

func getCommandsDropins(preKubeadmCommand []string, postKubeadminCommand []string) []ignition.Dropin {
	if len(preKubeadmCommand) == 0 && len(postKubeadminCommand) == 0 {
		return []ignition.Dropin{}
	}
	builder := strings.Builder{}
	builder.WriteString("[Service]\n")
	for _, command := range preKubeadmCommand {
		builder.WriteString("ExecStartPre=" + command + "\n")
	}
	for _, command := range postKubeadminCommand {
		builder.WriteString("ExecStartPost=" + command + "\n")
	}
	return []ignition.Dropin{
		{
			Name:    "10-commands.conf",
			Content: builder.String(),
		},
	}
}
