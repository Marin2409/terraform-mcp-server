package tools

import (
	tfeTools "github.com/hashicorp/terraform-mcp-server/pkg/mcp-official/tools/tfe"
	"github.com/hashicorp/terraform-mcp-server/pkg/toolsets"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	log "github.com/sirupsen/logrus"
)

func RegisterTools(svr *mcp.Server, logger *log.Logger, enabledToolsets []string) {
	if toolsets.IsToolEnabled("list_workspaces", enabledToolsets) {
		mcp.AddTool(svr, tfeTools.ListWorkpsacesTool(), tfeTools.ListWorkspacesFunc)
	}

	if toolsets.IsToolEnabled("list_terraform_orgs", enabledToolsets) {
		mcp.AddTool(svr, tfeTools.ListTerraformOrganizationsTool(), tfeTools.ListTerraformOrganizationsFunc)
	}

	if toolsets.IsToolEnabled("read_workspace_tags", enabledToolsets) {
		mcp.AddTool(svr, tfeTools.ReadWorkspaceTagsTool(), tfeTools.ReadWorkspaceTagsFunc)
	}

	if toolsets.IsToolEnabled("list_workspace_policy_sets", enabledToolsets) {
		mcp.AddTool(svr, tfeTools.ListWorkspacePolicySetsTool(), tfeTools.ListWorkspacePolicySetsFunc)
	}

	if toolsets.IsToolEnabled("attach_policy_set_to_workspaces", enabledToolsets) {
		mcp.AddTool(svr, tfeTools.AttachPolicySetToWorkspacesTool(), tfeTools.AttachPolicySetToWorkspacesFunc)
	}
}
