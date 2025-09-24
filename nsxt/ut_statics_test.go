package nsxt

var (
	emptyBody     = ""
	nsxtVersion = `
{
    "node_version": "9.1.0",
    "product_version": "9.1.0"
}
`
	serviceWithId = `
{
  "resource_type": "Service",
  "description": "My HTTP",
  "id": "foo",
  "display_name": "foo",
  "path": "/infra/services/my-http",
  "parent_path": "/infra/services/my-http",
  "relative_path": "my-http",
  "service_entries": [
      {
          "resource_type": "L4PortSetServiceEntry",
          "id": "MyHttpEntry",
          "display_name": "MyHttpEntry",
          "path": "/infra/services/my-http/service-entries/MyHttpEntry",
          "parent_path": "/infra/services/my-http",
          "relative_path": "MyHttpEntry",
          "destination_ports": [
              "8080"
          ],
          "l4_protocol": "TCP",
          "_create_user": "admin",
          "_create_time": 1517310677617,
          "_last_modified_user": "admin",
          "_last_modified_time": 1517310677617,
          "_system_owned": false,
          "_protection": "NOT_PROTECTED",
          "_revision": 0
      }
  ],
  "_create_user": "admin",
  "_create_time": 1517310677604,
  "_last_modified_user": "admin",
  "_last_modified_time": 1517310677604,
  "_system_owned": false,
  "_protection": "NOT_PROTECTED",
  "_revision": 0
}
`
	segmentWithId = `
{
  "resource_type": "Segment",
  "id": "web-tier",
  "display_name": "web-tier",
  "path": "/infra/tier-1s/cgw/segments/web-tier",
  "parent_path": "/infra/tier-1s/cgw",
  "relative_path": "web-tier",
  "subnets": [
    {
      "gateway_address": "40.1.1.1/16",
      "dhcp_ranges": [
        "40.1.2.0/24"
      ]
    }
  ],
  "connectivity_path": "/infra/tier-1s/mgw",
  "_create_user": "admin",
  "_create_time": 1516668961954,
  "_last_modified_user": "admin",
  "_last_modified_time": 1516668961954,
  "_system_owned": false,
  "_protection": "NOT_PROTECTED",
  "_revision": 0
}
`


)