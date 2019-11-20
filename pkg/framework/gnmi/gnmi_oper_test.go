/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package gnmi_test

import (
	"testing"
	"time"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	tg "github.com/stratum/testvectors/proto/target"
)

var (
	TestTarget        = &tg.Target{Address: "localhost:50001"}
	InvalidTestTarget = &tg.Target{Address: "localhost:50012"}
)

func TestProcessGetRequest(t *testing.T) {

	emptyGetReq := &gpb.GetRequest{Encoding: 2}
	emptyGetResp := &gpb.GetResponse{}

	validGetReq := &gpb.GetRequest{
		Path: []*gpb.Path{
			{
				Elem: []*gpb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "veth1"}},
					{Name: "state"},
					{Name: "name"},
				},
			},
		},
		Encoding: gpb.Encoding_PROTO}

	validGetResp := &gpb.GetResponse{
		Notification: []*gpb.Notification{
			{
				Timestamp: 1234567890123456789,
				Update: []*gpb.Update{
					{
						Path: &gpb.Path{
							Elem: []*gpb.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "veth1"}},
								{Name: "state"},
								{Name: "name"},
							},
						},
						Val: &gpb.TypedValue{
							Value: &gpb.TypedValue_StringVal{StringVal: "veth1"},
						},
					},
				},
			},
		},
	}
	type args struct {
		target *tg.Target
		greq   *gpb.GetRequest
		gresp  *gpb.GetResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Request",
			args: args{
				target: TestTarget,
				greq:   emptyGetReq,
				gresp:  emptyGetResp,
			},
			want: true,
		},
		{
			name: "Valid Request",
			args: args{
				target: TestTarget,
				greq:   validGetReq,
				gresp:  validGetResp,
			},
			want: true,
		},
		{
			name: "Invalid Target",
			args: args{
				target: InvalidTestTarget,
				greq:   validGetReq,
				gresp:  validGetResp,
			},
			want: false,
		},
		{
			name: "Invalid Respnse",
			args: args{
				target: TestTarget,
				greq:   validGetReq,
				gresp:  emptyGetResp,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gnmi.Init(tt.args.target)
			if got := gnmi.ProcessGetRequest(tt.args.greq, tt.args.gresp); got != tt.want {
				t.Errorf("ProcessGetRequest() = %v, want %v", got, tt.want)
			}
			gnmi.TearDown()
		})
	}
}

func TestProcessSetRequest(t *testing.T) {

	emptySetReq := &gpb.SetRequest{}
	emptySetResp := &gpb.SetResponse{Prefix: &gpb.Path{}}

	validSetReq := &gpb.SetRequest{
		Update: []*gpb.Update{
			{
				Path: &gpb.Path{
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_StringVal{StringVal: "GOOD"},
				},
			},
		},
	}

	invalidSetReq := &gpb.SetRequest{
		Update: []*gpb.Update{
			{
				Path: &gpb.Path{
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indic"},
					},
				},
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_StringVal{StringVal: "GOOD"},
				},
			},
		},
	}

	validSetResp := &gpb.SetResponse{
		Prefix: &gpb.Path{},
		Response: []*gpb.UpdateResult{
			{
				Path: &gpb.Path{
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Op: gpb.UpdateResult_UPDATE},
		},
	}

	gnmi.Init(TestTarget)
	type args struct {
		target *tg.Target
		sreq   *gpb.SetRequest
		sresp  *gpb.SetResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Request",
			args: args{
				target: TestTarget,
				sreq:   emptySetReq,
				sresp:  emptySetResp,
			},
			want: true,
		},
		{
			name: "Valid Request",
			args: args{
				target: TestTarget,
				sreq:   validSetReq,
				sresp:  validSetResp,
			},
			want: true,
		},
		{
			name: "Invalid Response",
			args: args{
				target: TestTarget,
				sreq:   validSetReq,
				sresp:  emptySetResp,
			},
			want: false,
		},
		{
			name: "Invalid Request",
			args: args{
				target: TestTarget,
				sreq:   invalidSetReq,
				sresp:  emptySetResp,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gnmi.ProcessSetRequest(tt.args.sreq, tt.args.sresp); got != tt.want {
				t.Errorf("ProcessSetRequest() = %v, want %v", got, tt.want)
			}
		})
	}
	gnmi.TearDown()
}

func TestProcessSubscribeRequest(t *testing.T) {
	resultChan := make(chan bool)
	emptySubReq := &gpb.SubscribeRequest{}
	emptySubResp := []*gpb.SubscribeResponse{}
	emptySetReq := &gpb.SetRequest{}
	emptySetResp := &gpb.SetResponse{Prefix: &gpb.Path{}}

	invalidSubReq := &gpb.SubscribeRequest{
		Request: &gpb.SubscribeRequest_Subscribe{
			Subscribe: &gpb.SubscriptionList{
				Subscription: []*gpb.Subscription{
					{
						Path: &gpb.Path{
							Elem: []*gpb.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "veth3"}},
								{Name: "config"},
								{Name: "health-indator"},
							},
						},
						Mode: gpb.SubscriptionMode_ON_CHANGE,
					},
				},
				UpdatesOnly: true,
			},
		},
	}

	validSubReq := &gpb.SubscribeRequest{
		Request: &gpb.SubscribeRequest_Subscribe{
			Subscribe: &gpb.SubscriptionList{
				Subscription: []*gpb.Subscription{
					{
						Path: &gpb.Path{
							Elem: []*gpb.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "veth3"}},
								{Name: "config"},
								{Name: "health-indicator"},
							},
						},
						Mode: gpb.SubscriptionMode_ON_CHANGE,
					},
				},
				UpdatesOnly: true,
			},
		},
	}
	validSubResp := []*gpb.SubscribeResponse{
		{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Timestamp: 1234567890123456789,
					Update: []*gpb.Update{
						{
							Path: &gpb.Path{
								Elem: []*gpb.PathElem{
									{Name: "interfaces"},
									{Name: "interface", Key: map[string]string{"name": "veth3"}},
									{Name: "config"},
									{Name: "health-indicator"},
								},
							},
							Val: &gpb.TypedValue{
								Value: &gpb.TypedValue_StringVal{StringVal: "GOOD"},
							},
						},
					},
				},
			},
		},
		{
			Response: &gpb.SubscribeResponse_SyncResponse{
				SyncResponse: true,
			},
		},
		{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Timestamp: 1234567890123456789,
					Update: []*gpb.Update{
						{
							Path: &gpb.Path{
								Elem: []*gpb.PathElem{
									{Name: "interfaces"},
									{Name: "interface", Key: map[string]string{"name": "veth3"}},
									{Name: "config"},
									{Name: "health-indicator"},
								},
							},
							Val: &gpb.TypedValue{
								Value: &gpb.TypedValue_StringVal{StringVal: "BAD"},
							},
						},
					},
				},
			},
		},
		{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Timestamp: 1234567890123456789,
					Update: []*gpb.Update{
						{
							Path: &gpb.Path{
								Elem: []*gpb.PathElem{
									{Name: "interfaces"},
									{Name: "interface", Key: map[string]string{"name": "veth3"}},
									{Name: "config"},
									{Name: "health-indicator"},
								},
							},
							Val: &gpb.TypedValue{
								Value: &gpb.TypedValue_StringVal{StringVal: "GOOD"},
							},
						},
					},
				},
			},
		},
	}
	setReq1 := &gpb.SetRequest{
		Update: []*gpb.Update{
			{
				Path: &gpb.Path{
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_StringVal{StringVal: "BAD"},
				},
			},
		},
	}
	setReq2 := &gpb.SetRequest{
		Update: []*gpb.Update{
			{
				Path: &gpb.Path{
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_StringVal{StringVal: "GOOD"},
				},
			},
		},
	}
	setResp := &gpb.SetResponse{
		Prefix: &gpb.Path{},
		Response: []*gpb.UpdateResult{
			{
				Path: &gpb.Path{
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Op: gpb.UpdateResult_UPDATE},
		},
	}
	gnmi.Init(TestTarget)
	type args struct {
		target     *tg.Target
		subreq     *gpb.SubscribeRequest
		subresp    []*gpb.SubscribeResponse
		setreq1    *gpb.SetRequest
		setreq2    *gpb.SetRequest
		setresp    *gpb.SetResponse
		resultChan chan bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Sub",
			args: args{
				target:     TestTarget,
				subreq:     emptySubReq,
				subresp:    emptySubResp,
				setreq1:    emptySetReq,
				setreq2:    emptySetReq,
				setresp:    emptySetResp,
				resultChan: resultChan,
			},
			want: true,
		},
		{
			name: "Invalid Sub",
			args: args{
				target:     TestTarget,
				subreq:     invalidSubReq,
				subresp:    emptySubResp,
				setreq1:    emptySetReq,
				setreq2:    emptySetReq,
				setresp:    emptySetResp,
				resultChan: resultChan,
			},
			want: true,
		},
		{
			name: "Valid Sub",
			args: args{
				target:     TestTarget,
				subreq:     validSubReq,
				subresp:    validSubResp,
				setreq1:    setReq1,
				setreq2:    setReq2,
				setresp:    setResp,
				resultChan: resultChan,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstRespChan := make(chan struct{})
			go gnmi.ProcessSubscribeRequest(tt.args.subreq, tt.args.subresp, firstRespChan, tt.args.resultChan)

			select {
			case <-firstRespChan:
				if got := gnmi.ProcessSetRequest(tt.args.setreq1, tt.args.setresp) && gnmi.ProcessSetRequest(tt.args.setreq2, tt.args.setresp); got != true {
					t.Errorf("ProcessSetRequest() = %v, want %v", got, tt.want)
				}
				if got := <-tt.args.resultChan; got != tt.want {
					t.Errorf("ProcessSubscribeRequest() = %v, want %v", got, tt.want)
				}
			case <-time.After(gnmi.SubTimeout):
				t.Errorf("ProcessSubscribeRequest() = %v, want %v", false, tt.want)
			}

		})
	}
	gnmi.TearDown()
}

func TestInit(t *testing.T) {
	type args struct {
		target *tg.Target
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid Target",
			args: args{target: TestTarget},
		},
		{
			name: "Invalid Target",
			args: args{target: InvalidTestTarget},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gnmi.Init(tt.args.target)
		})
	}
}
