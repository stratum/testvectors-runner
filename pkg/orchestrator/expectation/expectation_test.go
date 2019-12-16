/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package expectation

import (
	"testing"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"
	tg "github.com/stratum/testvectors/proto/target"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var TestTarget = &tg.Target{Address: "localhost:50001"}

func TestProcessConfigExpectation(t *testing.T) {
	gnmi.Init(TestTarget)
	defer gnmi.TearDown()
	var (
		emptyConfigExpectation = &tv.ConfigExpectation{}
		validConfigExpectation = &tv.ConfigExpectation{
			GnmiGetRequest: &gpb.GetRequest{
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
				Encoding: gpb.Encoding_PROTO,
			},
			GnmiGetResponse: &gpb.GetResponse{
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
			},
		}
		invalidConfigExpectation = &tv.ConfigExpectation{
			GnmiGetRequest: &gpb.GetRequest{
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
				Encoding: gpb.Encoding_PROTO,
			},
			GnmiGetResponse: &gpb.GetResponse{
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
									Value: &gpb.TypedValue_StringVal{StringVal: "veth"},
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
			want: true,
		},
		{
			name: "Valid Expectation",
			args: args{ce: validConfigExpectation},
			want: true,
		},
		{
			name: "Invalid Expectation",
			args: args{ce: invalidConfigExpectation},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processConfigExpectation(tt.args.ce); got != tt.want {
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
			if got := processControlPlaneExpectation(tt.args.cpe); got != tt.want {
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
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processDataPlaneExpectation(tt.args.dpe); got != tt.want {
				t.Errorf("ProcessDataPlaneExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessTelemetryExpectation(t *testing.T) {
	gnmi.Init(TestTarget)
	p4rt.Init(TestTarget)
	defer gnmi.TearDown()
	defer p4rt.TearDown()
	var (
		telemetryExpectation = &tv.TelemetryExpectation{
			GnmiSubscribeRequest:  &gpb.SubscribeRequest{},
			ActionGroup:           &tv.ActionGroup{},
			GnmiSubscribeResponse: []*gpb.SubscribeResponse{},
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
			if got := processTelemetryExpectation(tt.args.tme); got != tt.want {
				t.Errorf("ProcessTelemetryExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}
