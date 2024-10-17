package nsxt

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

const nsxRemovedVersion = "9.0.0"

type resourceFunc func() *schema.Resource

func commonVersionCheck(m interface{}) bool {
	initNSXVersion(getPolicyConnector(m))
	return !util.NsxVersionHigherOrEqual(nsxRemovedVersion)
}

func genericWrapper(originalFunc interface{}, name string, fail bool, action string) interface{} {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		if commonVersionCheck(m) {
			switch f := originalFunc.(type) {
			case schema.ReadContextFunc:
				return f(ctx, d, m)
			case schema.CreateContextFunc:
				return f(ctx, d, m)
			case schema.UpdateContextFunc:
				return f(ctx, d, m)
			case schema.DeleteContextFunc:
				return f(ctx, d, m)
			}
		}
		if fail {
			log.Printf("[INFO] Failing %s for resource %s: removed from NSX %s", action, name, nsxRemovedVersion)
			return diag.FromErr(mpResourceRemovedError(name))
		}
		log.Printf("[INFO] Skipping %s for resource %s: removed from NSX %s", action, name, nsxRemovedVersion)
		return nil
	}
}

func readWrapper(originalFunc schema.ReadContextFunc, name string, fail bool) schema.ReadContextFunc {
	return genericWrapper(originalFunc, name, fail, "read").(schema.ReadContextFunc)
}

func deleteWrapper(originalFunc schema.DeleteContextFunc, name string) schema.DeleteContextFunc {
	return genericWrapper(originalFunc, name, false, "delete").(schema.DeleteContextFunc)
}

func createWrapper(originalFunc schema.CreateContextFunc, name string) schema.CreateContextFunc {
	return genericWrapper(originalFunc, name, false, "create").(schema.CreateContextFunc)
}

func updateWrapper(originalFunc schema.UpdateContextFunc, name string) schema.UpdateContextFunc {
	return genericWrapper(originalFunc, name, false, "update").(schema.UpdateContextFunc)
}

func importerWrapper(originalImporter *schema.ResourceImporter, name string) *schema.ResourceImporter {
	wrappedFunc := func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if commonVersionCheck(m) {
			return originalImporter.StateContext(ctx, d, m)
		}
		log.Printf("[INFO] Failing import for resource %s: removed from NSX %s", name, nsxRemovedVersion)
		return nil, mpResourceRemovedError(name)
	}
	originalImporter.StateContext = wrappedFunc
	return originalImporter
}

func removedResourceWrapper(realResourceFunc resourceFunc, name string) *schema.Resource {
	resource := realResourceFunc()
	resource.ReadContext = readWrapper(resource.ReadContext, name, false)
	if resource.DeleteContext != nil {
		resource.DeleteContext = deleteWrapper(resource.DeleteContext, name)
	}
	if resource.CreateContext != nil {
		resource.CreateContext = createWrapper(resource.CreateContext, name)
	}
	if resource.UpdateContext != nil {
		resource.UpdateContext = updateWrapper(resource.UpdateContext, name)
	}
	if resource.Importer != nil {
		resource.Importer = importerWrapper(resource.Importer, name)
	}
	return resource
}

func removedDataSourceWrapper(realDataSourceFunc resourceFunc, name string) *schema.Resource {
	resource := realDataSourceFunc()
	resource.ReadContext = readWrapper(resource.ReadContext, name, true)
	return resource
}
