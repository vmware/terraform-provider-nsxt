package unittests

import (
	"testing"

	// "../mocks/mocks"
	"github.com/vmware/terraform-provider-nsxt/mocks"

	"github.com/golang/mock/gomock"
)
func TestSampleResourceRead(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // mock SDK
    _ = mocks.NewMockSampleService(ctrl)
    // mockSdk := mocks.NewMockSampleService(ctrl)
    // mockSdk.EXPECT().
    //     Get("123").
    //     Return(&samplesdk.Sample{ID: "123", Name: "mocked"}, nil)
    
    // // fake client
    // client := infra.Client{Sdk: mockSdk}
    
    // // fake resource data
    // r := resourceSample()
    // d := r.TestResourceData()
    // d.SetId("123")
    
    // diags := sampleresourceRead(context.Background(), d, client)

    // require.Empty(t, diags)
    // require.Equal(t, "mocked", d.Get("name"))
}