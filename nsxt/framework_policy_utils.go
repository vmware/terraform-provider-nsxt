/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	tf_api "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
)

type policyContextModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

func NewUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func GetProjectIDFromContext(ctx context.Context, context types.Object, diag diag.Diagnostics) string {
	var contextConfig policyContextModel
	if context.IsNull() {
		return ""
	}
	diag.Append(context.As(ctx, &contextConfig, basetypes.ObjectAsOptions{})...)
	if diag.HasError() {
		return ""
	}
	return contextConfig.ProjectID.ValueString()
}

func GetSessionContext(ctx context.Context, client interface{}, context types.Object, diag diag.Diagnostics) tf_api.SessionContext {
	var clientType tf_api.ClientType
	projectID := GetProjectIDFromContext(ctx, context, diag)
	if projectID != "" {
		clientType = tf_api.Multitenancy
	} else if IsPolicyGlobalManager(client) {
		clientType = tf_api.Global
	} else {
		clientType = tf_api.Local
	}
	return tf_api.SessionContext{ProjectID: projectID, ClientType: clientType}
}

func GetOrGenerateID2(ctx context.Context, client interface{}, nsxID types.String, context types.Object, diag diag.Diagnostics, presenceChecker func(tf_api.SessionContext, string, client.Connector, diag.Diagnostics) bool) string {
	connector := GetPolicyConnector(client)

	if nsxID.IsNull() {
		return NewUUID()
	}
	id := nsxID.ValueString()
	if id == "" {
		return NewUUID()
	}

	exists := presenceChecker(GetSessionContext(ctx, client, context, diag), id, connector, diag)
	if diag.HasError() {
		return ""
	}

	if exists {
		diag.AddError("Failed to add resource", fmt.Sprintf("resource with id %s already exists", id))
	}

	return id
}

func NsxtPolicyPathResourceImporter(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	err := NsxtPolicyPathResourceImporterHelper(ctx, request, response)
	if errors.Is(err, ErrNotAPolicyPath) {
		return
	} else if err != nil {
		response.Diagnostics.AddError("Resource import failed", err.Error())
		return
	}
	// return
}

func NsxtPolicyPathResourceImporterHelper(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) error {
	importID := request.ID
	if isPolicyPath(importID) {
		pathSegs := strings.Split(importID, "/")
		if strings.Contains(pathSegs[1], "infra") {
			response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("nsx_id"), pathSegs[len(pathSegs)-1])...)
		} else if pathSegs[1] == "orgs" && pathSegs[3] == "projects" {
			if len(pathSegs) < 5 {
				return fmt.Errorf("invalid policy multitenancy path %s", importID)
			}
			// pathSegs[2] should contain the organization. Once we support multiple organization, it should be
			// assigned into the context as well
			response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("nsx_id"), pathSegs[len(pathSegs)-1])...)
			response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("context").AtName("project_id"), pathSegs[4])...)

		}
		return nil
	}
	return ErrNotAPolicyPath
}

// func getFrameworkPolicyTagsFromSchema(d *schema.ResourceData) []model.Tag {
// 	tags, _ := getCustomizedPolicyTagsFromSchema(d, "tag")
// 	return tags
// }
