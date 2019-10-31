/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package orchestrator

import (
	"testing"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework"
	tv "github.com/stratum/testvectors/proto/testvector"
)

func TestProcessConfigExpectation(t *testing.T) {
	framework.InitGNMI(TestTarget)
	defer framework.TearDownGNMI()
	var (
		emptyConfigExpectation = &tv.ConfigExpectation{}
		validConfigExpectation = &tv.ConfigExpectation{
			GnmiGetRequest: &gnmi.GetRequest{
				Path: []*gnmi.Path{
					{
						Elem: []*gnmi.PathElem{
							{Name: "interfaces"},
							{Name: "interface", Key: map[string]string{"name": "veth1"}},
							{Name: "state"},
							{Name: "name"},
						},
					},
				},
				Encoding: gnmi.Encoding_PROTO,
			},
			GnmiGetResponse: &gnmi.GetResponse{
				Notification: []*gnmi.Notification{
					{
						Timestamp: 1234567890123456789,
						Update: []*gnmi.Update{
							{
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
			},
		}
	)
	type args struct {
		ce *tv.ConfigExpectation
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Expectation",
			args: args{ce: emptyConfigExpectation},
			want: false,
		},
		{
			name: "Valid Expectation",
			args: args{ce: validConfigExpectation},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessConfigExpectation(tt.args.ce); got != tt.want {
				t.Errorf("ProcessConfigExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessControlPlaneExpectation(t *testing.T) {
	var (
		readExpectation = &tv.ControlPlaneExpectation{
			Expectations: &tv.ControlPlaneExpectation_ReadExpectation_{
				ReadExpectation: &tv.ControlPlaneExpectation_ReadExpectation{},
			},
		}
		pipelineConfigExpectation = &tv.ControlPlaneExpectation{
			Expectations: &tv.ControlPlaneExpectation_PipelineConfigExpectation_{
				PipelineConfigExpectation: &tv.ControlPlaneExpectation_PipelineConfigExpectation{},
			},
		}
		packetInExpectation = &tv.ControlPlaneExpectation{
			Expectations: &tv.ControlPlaneExpectation_PacketInExpectation_{
				PacketInExpectation: &tv.ControlPlaneExpectation_PacketInExpectation{},
			},
		}
	)
	type args struct {
		cpe *tv.ControlPlaneExpectation
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Read Expectation",
			args: args{cpe: readExpectation},
			want: false,
		},
		{
			name: "Empty Pipeline Config Expectation",
			args: args{cpe: pipelineConfigExpectation},
			want: false,
		},
		{
			name: "Empty Packet In Expectation",
			args: args{cpe: packetInExpectation},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessControlPlaneExpectation(tt.args.cpe); got != tt.want {
				t.Errorf("ProcessControlPlaneExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessDataPlaneExpectation(t *testing.T) {
	var (
		dataplaneExpectation = &tv.DataPlaneExpectation{
			Expectations: &tv.DataPlaneExpectation_TrafficExpectation_{
				TrafficExpectation: &tv.DataPlaneExpectation_TrafficExpectation{},
			},
		}
	)
	type args struct {
		dpe *tv.DataPlaneExpectation
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Data Plane Expectation",
			args: args{dpe: dataplaneExpectation},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessDataPlaneExpectation(tt.args.dpe); got != tt.want {
				t.Errorf("ProcessDataPlaneExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessTelemetryExpectation(t *testing.T) {
	framework.InitGNMI(TestTarget)
	framework.Init(TestTarget)
	defer framework.TearDown()
	defer framework.TearDownGNMI()
	var (
		telemetryExpectation = &tv.TelemetryExpectation{
			GnmiSubscribeRequest:  &gnmi.SubscribeRequest{},
			ActionGroup:           &tv.ActionGroup{},
			GnmiSubscribeResponse: []*gnmi.SubscribeResponse{},
		}
	)
	type args struct {
		tme *tv.TelemetryExpectation
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Telemetry Expectation",
			args: args{tme: telemetryExpectation},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessTelemetryExpectation(tt.args.tme); got != tt.want {
				t.Errorf("ProcessTelemetryExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}
