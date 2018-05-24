package testing

import (
	"fmt"
	fake "github.com/huaweicloud/golangsdk/openstack/cce/v3/common"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	th "github.com/huaweicloud/golangsdk/testhelper"
	"net/http"
	"testing"
)

func TestListNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/api/v3/projects/95fb48c35a6c49bca3f99924885d0362/clusters/cec124c2-58f1-11e8-ad73-0255ac101926/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "kind": "List",
    "apiVersion": "v3",
	"items":
	[ 
       {
            "kind": "Host",
            "apiVersion": "v3",
            "metadata": {
                "name": "test-node-1234",
                "uid": "b99acd73-5d7c-11e8-8e76-0255ac101929",
                "creationTimestamp": "2018-05-22 04:58:10.605829241 +0000 UTC",
                "updateTimestamp": "2018-05-22 05:02:26.528912685 +0000 UTC"
            },
            "spec": {
                "flavor": "s1.medium",
                "az": "cn-east-2a",
                "login": {
                    "sshKey": "c2c-keypair",
                    "userPassword": {}
                },
                "rootVolume": {
                    "volumetype": "SATA",
                    "size": 40
                },
                "dataVolumes": [
                    {
                        "volumetype": "SATA",
                        "size": 100
                    }
                ],
                "publicIP": {
                    "eip": {
                        "bandwidth": {}
                    }
                },
                "billingMode": 0
            },
            "status": {
                "phase": "Active",
                "serverId": "41748e56-33d4-46a1-aa57-2c8c29907995",
                "privateIP": "192.168.0.3"
            }
        }
	]
}
		`)
	})

	listNodes := nodes.ListOpts{Name: "test-node-1234"}
	actual, err := nodes.List(fake.ServiceClient(), "cec124c2-58f1-11e8-ad73-0255ac101926").ExtractNode(listNodes)

	if err != nil {
		t.Errorf("Failed to extract nodes: %v", err)
	}

	expected := []nodes.Items{
		{
			Kind:       "Host",
			Apiversion: "v3",
			Metadata: nodes.Metadata{Name: "test-node-1234",
				Uid:               "b99acd73-5d7c-11e8-8e76-0255ac101929",
				CreationTimestamp: "2018-05-22 04:58:10.605829241 +0000 UTC",
				UpdateTimestamp:   "2018-05-22 05:02:26.528912685 +0000 UTC"},
			Spec: nodes.Spec{Az: "cn-east-2a",
				Login:       nodes.Login{SshKey: "c2c-keypair"},
				RootVolume:  nodes.Volume{Size: 40, Volumetype: "SATA"},
				BillingMode: 0,
				DataVolumes: []nodes.Volume{
					{
						Volumetype: "SATA",
						Size:       100,
					}},
				Flavor: "s1.medium",
			},
			Status: nodes.Status{Phase: "Active", ServerID: "41748e56-33d4-46a1-aa57-2c8c29907995", PrivateIP: "192.168.0.3"},
		},
	}

	th.AssertDeepEquals(t, expected, actual)
}

/*
func TestGetV3Node(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/api/v3/projects/95fb48c35a6c49bca3f99924885d0362/clusters/cec124c2-58f1-11e8-ad73-0255ac101926/nodes/cf4bc001-58f1-11e8-ad73-0255ac101926", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, Output)
	})

	actual, err := nodes.Get(fake.ServiceClient(), "cec124c2-58f1-11e8-ad73-0255ac101926", "cf4bc001-58f1-11e8-ad73-0255ac101926").Extract()
	th.AssertNoErr(t, err)
	expected := Expected
	th.AssertDeepEquals(t, expected, actual)

}


func TestCreateV3Node(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/api/v3/projects/95fb48c35a6c49bca3f99924885d0362/clusters/cec124c2-58f1-11e8-ad73-0255ac101926/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")

		th.TestJSONRequest(t, r, `
			{
	  "apiversion": "v3",
	  "kind": "Node",
	  "metadata": {
	    "name": "c2c-hostname"
	  },
	  "spec": {
	    "az": "cn-east-2a",
	    "count": 1,
	    "dataVolumes": [
	      {
	        "size": 100,
	        "volumetype": "SATA"
	      }
	    ],
	    "flavor": "s3.large.2",
	    "login": {
	      "sshKey": "c2c-keypair"
	    },
	    "publicIP": {
	      "eip": {
	        "bandwidth": {
	          "sharetype": "",
	          "size": 0
	        },
	        "iptype": ""
	      }
	    },
	    "rootVolume": {
	      "size": 40,
	      "volumetype": "SATA"
	    }
	  }
	}

`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, Output)
	})
	options := nodes.CreateOpts{Kind: "Node",
		ApiVersion: "v3",
		Metadata:   nodes.CreateMetaData{Name: "c2c-hostname"},
		Spec: nodes.Spec{Flavor: "s3.large.2", Az: "cn-east-2a",
			Login:       nodes.LoginSpec{SshKey: "c2c-keypair"},
			RootVolume:  nodes.VolumeSpec{Size: 40, VolumeType: "SATA"},
			DataVolumes: []nodes.VolumeSpec{{Size: 100, VolumeType: "SATA"}},
			Count:       1},
	}
	actual, err := nodes.Create(fake.ServiceClient(), "cec124c2-58f1-11e8-ad73-0255ac101926", options).Extract()
	th.AssertNoErr(t, err)
	expected := Expected
	th.AssertDeepEquals(t, expected, actual)

}

func TestUpdateV3Node(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/api/v3/projects/95fb48c35a6c49bca3f99924885d0362/clusters/cec124c2-58f1-11e8-ad73-0255ac101926/nodes/cf4bc001-58f1-11e8-ad73-0255ac101926", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "metadata": {
        "name": "c2c-hostname"
    }
}			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, Output)
	})
	options := nodes.UpdateOpts{Metadata: nodes.UpdateMetadata{Name: "c2c-hostname"}}
	actual, err := nodes.Update(fake.ServiceClient(), "cec124c2-58f1-11e8-ad73-0255ac101926", "cf4bc001-58f1-11e8-ad73-0255ac101926", options).Extract()
	th.AssertNoErr(t, err)
	expected := Expected
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/api/v3/projects/95fb48c35a6c49bca3f99924885d0362/clusters/cec124c2-58f1-11e8-ad73-0255ac101926/nodes/cf4bc001-58f1-11e8-ad73-0255ac101926", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
	})

	err := nodes.Delete(fake.ServiceClient(), "cec124c2-58f1-11e8-ad73-0255ac101926", "cf4bc001-58f1-11e8-ad73-0255ac101926").ExtractErr()
	th.AssertNoErr(t, err)

}
*/
