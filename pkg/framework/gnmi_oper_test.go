/*
*Copyright 2019-present Open Networking Foundation
*
*SPDX-License-Identifier: Apache-2.0
 */

package framework

import (
	"testing"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	tg "github.com/stratum/testvectors/proto/target"
)

var (
	TestTarget        = &tg.Target{Address: "localhost:50001"}
	InvalidTestTarget = &tg.Target{Address: "localhost:50012"}
)

func TestProcessGetRequest(t *testing.T) {

	emptyGetReq := &gnmi.GetRequest{Encoding: 2}
	emptyGetResp := &gnmi.GetResponse{}

	validGetReq := &gnmi.GetRequest{
		Path: []*gnmi.Path{
			&gnmi.Path{
				Elem: []*gnmi.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "veth1"}},
					{Name: "state"},
					{Name: "name"},
				},
			},
		},
		Encoding: gnmi.Encoding_PROTO}

	validGetResp := &gnmi.GetResponse{
		Notification: []*gnmi.Notification{
			&gnmi.Notification{
				Timestamp: 1234567890123456789,
				Update: []*gnmi.Update{
					&gnmi.Update{
						Path: &gnmi.Path{
							Elem: []*gnmi.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "veth1"}},
								{Name: "state"},
								{Name: "name"},
							},
						},
						Val: &gnmi.TypedValue{
							Value: &gnmi.TypedValue_StringVal{StringVal: "veth1"},
						},
					},
				},
			},
		},
	}
	type args struct {
		target *tg.Target
		greq   *gnmi.GetRequest
		gresp  *gnmi.GetResponse
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
			InitGNMI(tt.args.target)
			if got := ProcessGetRequest(tt.args.target, tt.args.greq, tt.args.gresp); got != tt.want {
				t.Errorf("ProcessGetRequest() = %v, want %v", got, tt.want)
			}
			TearDownGNMI()
		})
	}
}

func TestProcessSetRequest(t *testing.T) {

	emptySetReq := &gnmi.SetRequest{}
	emptySetResp := &gnmi.SetResponse{Prefix: &gnmi.Path{}}

	validSetReq := &gnmi.SetRequest{
		Update: []*gnmi.Update{
			&gnmi.Update{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_StringVal{StringVal: "GOOD"},
				},
			},
		},
	}

	invalidSetReq := &gnmi.SetRequest{
		Update: []*gnmi.Update{
			&gnmi.Update{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indic"},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_StringVal{StringVal: "GOOD"},
				},
			},
		},
	}

	validSetResp := &gnmi.SetResponse{
		Prefix: &gnmi.Path{},
		Response: []*gnmi.UpdateResult{
			&gnmi.UpdateResult{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Op: gnmi.UpdateResult_UPDATE},
		},
	}

	InitGNMI(TestTarget)
	type args struct {
		target *tg.Target
		sreq   *gnmi.SetRequest
		sresp  *gnmi.SetResponse
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
			if got := ProcessSetRequest(tt.args.target, tt.args.sreq, tt.args.sresp); got != tt.want {
				t.Errorf("ProcessSetRequest() = %v, want %v", got, tt.want)
			}
		})
	}
	TearDownGNMI()
}

func TestProcessSubscribeRequest(t *testing.T) {
	resultChan := make(chan bool, 1)
	emptySubReq := &gnmi.SubscribeRequest{}
	emptySubResp := []*gnmi.SubscribeResponse{}
	emptySetReq := &gnmi.SetRequest{}
	emptySetResp := &gnmi.SetResponse{Prefix: &gnmi.Path{}}

	invalidSubReq := &gnmi.SubscribeRequest{
		Request: &gnmi.SubscribeRequest_Subscribe{
			Subscribe: &gnmi.SubscriptionList{
				Subscription: []*gnmi.Subscription{
					&gnmi.Subscription{
						Path: &gnmi.Path{
							Elem: []*gnmi.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "veth3"}},
								{Name: "config"},
								{Name: "health-indator"},
							},
						},
						Mode: gnmi.SubscriptionMode_ON_CHANGE,
					},
				},
				UpdatesOnly: true,
			},
		},
	}

	validSubReq := &gnmi.SubscribeRequest{
		Request: &gnmi.SubscribeRequest_Subscribe{
			Subscribe: &gnmi.SubscriptionList{
				Subscription: []*gnmi.Subscription{
					&gnmi.Subscription{
						Path: &gnmi.Path{
							Elem: []*gnmi.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "veth3"}},
								{Name: "config"},
								{Name: "health-indicator"},
							},
						},
						Mode: gnmi.SubscriptionMode_ON_CHANGE,
					},
				},
				UpdatesOnly: true,
			},
		},
	}
	validSubResp := []*gnmi.SubscribeResponse{
		&gnmi.SubscribeResponse{
			Response: &gnmi.SubscribeResponse_Update{
				Update: &gnmi.Notification{
					Timestamp: 1234567890123456789,
					Update: []*gnmi.Update{
						&gnmi.Update{
							Path: &gnmi.Path{
								Elem: []*gnmi.PathElem{
									{Name: "interfaces"},
									{Name: "interface", Key: map[string]string{"name": "veth3"}},
									{Name: "config"},
									{Name: "health-indicator"},
								},
							},
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_StringVal{StringVal: "GOOD"},
							},
						},
					},
				},
			},
		},
		&gnmi.SubscribeResponse{
			Response: &gnmi.SubscribeResponse_SyncResponse{
				SyncResponse: true,
			},
		},
		&gnmi.SubscribeResponse{
			Response: &gnmi.SubscribeResponse_Update{
				Update: &gnmi.Notification{
					Timestamp: 1234567890123456789,
					Update: []*gnmi.Update{
						&gnmi.Update{
							Path: &gnmi.Path{
								Elem: []*gnmi.PathElem{
									{Name: "interfaces"},
									{Name: "interface", Key: map[string]string{"name": "veth3"}},
									{Name: "config"},
									{Name: "health-indicator"},
								},
							},
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_StringVal{StringVal: "BAD"},
							},
						},
					},
				},
			},
		},
		&gnmi.SubscribeResponse{
			Response: &gnmi.SubscribeResponse_Update{
				Update: &gnmi.Notification{
					Timestamp: 1234567890123456789,
					Update: []*gnmi.Update{
						&gnmi.Update{
							Path: &gnmi.Path{
								Elem: []*gnmi.PathElem{
									{Name: "interfaces"},
									{Name: "interface", Key: map[string]string{"name": "veth3"}},
									{Name: "config"},
									{Name: "health-indicator"},
								},
							},
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_StringVal{StringVal: "GOOD"},
							},
						},
					},
				},
			},
		},
	}
	setReq1 := &gnmi.SetRequest{
		Update: []*gnmi.Update{
			&gnmi.Update{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_StringVal{StringVal: "BAD"},
				},
			},
		},
	}
	setReq2 := &gnmi.SetRequest{
		Update: []*gnmi.Update{
			&gnmi.Update{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_StringVal{StringVal: "GOOD"},
				},
			},
		},
	}
	setResp := &gnmi.SetResponse{
		Prefix: &gnmi.Path{},
		Response: []*gnmi.UpdateResult{
			&gnmi.UpdateResult{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "veth3"}},
						{Name: "config"},
						{Name: "health-indicator"},
					},
				},
				Op: gnmi.UpdateResult_UPDATE},
		},
	}
	InitGNMI(TestTarget)
	type args struct {
		target     *tg.Target
		subreq     *gnmi.SubscribeRequest
		subresp    []*gnmi.SubscribeResponse
		setreq1    *gnmi.SetRequest
		setreq2    *gnmi.SetRequest
		setresp    *gnmi.SetResponse
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
			go ProcessSubscribeRequest(tt.args.target, tt.args.subreq, tt.args.subresp, tt.args.resultChan)
			time.Sleep(2 * time.Millisecond)
			if got := ProcessSetRequest(tt.args.target, tt.args.setreq1, tt.args.setresp) && ProcessSetRequest(tt.args.target, tt.args.setreq2, tt.args.setresp); got != true {
				t.Errorf("ProcessSetRequest() = %v, want %v", got, tt.want)
			}
			if got := <-tt.args.resultChan; got != tt.want {
				t.Errorf("ProcessSubscribeRequest() = %v, want %v", got, tt.want)
			}
		})
	}
	TearDownGNMI()
}

func TestInitGNMI(t *testing.T) {
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
			InitGNMI(tt.args.target)
		})
	}
}
