//nolint:revive
package segments

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/segments"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type SegmentDiscoveryProfileBindingMapClientContext utl.ClientContext

func NewSegmentDiscoveryProfileBindingMapsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *SegmentDiscoveryProfileBindingMapClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewSegmentDiscoveryProfileBindingMapsClient(connector)

	case utl.Multitenancy:
		client = client1.NewSegmentDiscoveryProfileBindingMapsClient(connector)

	default:
		return nil
	}
	return &SegmentDiscoveryProfileBindingMapClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c SegmentDiscoveryProfileBindingMapClientContext) Get(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) (model0.SegmentDiscoveryProfileBindingMap, error) {
	var obj model0.SegmentDiscoveryProfileBindingMap
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.SegmentDiscoveryProfileBindingMapsClient)
		obj, err = client.Get(infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Multitenancy:
		client := c.Client.(client1.SegmentDiscoveryProfileBindingMapsClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c SegmentDiscoveryProfileBindingMapClientContext) Delete(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.SegmentDiscoveryProfileBindingMapsClient)
		err = client.Delete(infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam)

	case utl.Multitenancy:
		client := c.Client.(client1.SegmentDiscoveryProfileBindingMapsClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c SegmentDiscoveryProfileBindingMapClientContext) Patch(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model0.SegmentDiscoveryProfileBindingMap) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.SegmentDiscoveryProfileBindingMapsClient)
		err = client.Patch(infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam, segmentDiscoveryProfileBindingMapParam)

	case utl.Multitenancy:
		client := c.Client.(client1.SegmentDiscoveryProfileBindingMapsClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam, segmentDiscoveryProfileBindingMapParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c SegmentDiscoveryProfileBindingMapClientContext) Update(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model0.SegmentDiscoveryProfileBindingMap) (model0.SegmentDiscoveryProfileBindingMap, error) {
	var err error
	var obj model0.SegmentDiscoveryProfileBindingMap

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.SegmentDiscoveryProfileBindingMapsClient)
		obj, err = client.Update(infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam, segmentDiscoveryProfileBindingMapParam)

	case utl.Multitenancy:
		client := c.Client.(client1.SegmentDiscoveryProfileBindingMapsClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, infraSegmentIdParam, segmentDiscoveryProfileBindingMapIdParam, segmentDiscoveryProfileBindingMapParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c SegmentDiscoveryProfileBindingMapClientContext) List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.SegmentDiscoveryProfileBindingMapListResult, error) {
	var err error
	var obj model0.SegmentDiscoveryProfileBindingMapListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.SegmentDiscoveryProfileBindingMapsClient)
		obj, err = client.List(infraSegmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Multitenancy:
		client := c.Client.(client1.SegmentDiscoveryProfileBindingMapsClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, infraSegmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}