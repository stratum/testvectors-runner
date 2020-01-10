/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package action

import (
	"testing"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	tg "github.com/stratum/testvectors/proto/target"
	tv "github.com/stratum/testvectors/proto/testvector"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"
)

var (
	TestTarget   = &tg.Target{Address: "localhost:50001"}
	emptyAction  = &tv.Action{}
	configAction = &tv.Action{
		Actions: &tv.Action_ConfigOperation{
			ConfigOperation: &tv.ConfigOperation{},
		},
	}
	controlPlaneWriteAction = &tv.Action{
		Actions: &tv.Action_ControlPlaneOperation{
			ControlPlaneOperation: &tv.ControlPlaneOperation{
				Operations: &tv.ControlPlaneOperation_WriteOperation_{
					WriteOperation: &tv.ControlPlaneOperation_WriteOperation{},
				},
			},
		},
	}
	controlPlanePipelineConfigAction = &tv.Action{
		Actions: &tv.Action_ControlPlaneOperation{
			ControlPlaneOperation: &tv.ControlPlaneOperation{
				Operations: &tv.ControlPlaneOperation_PipelineConfigOperation_{
					PipelineConfigOperation: &tv.ControlPlaneOperation_PipelineConfigOperation{},
				},
			},
		},
	}
	controlPlanePacketOutAction = &tv.Action{
		Actions: &tv.Action_ControlPlaneOperation{
			ControlPlaneOperation: &tv.ControlPlaneOperation{
				Operations: &tv.ControlPlaneOperation_PacketOutOperation_{
					PacketOutOperation: &tv.ControlPlaneOperation_PacketOutOperation{
						P4PacketOut: &v1.PacketOut{},
					},
				},
			},
		},
	}
	dataPlaneAction = &tv.Action{
		Actions: &tv.Action_DataPlaneStimulus{
			DataPlaneStimulus: &tv.DataPlaneStimulus{
				Stimuli: &tv.DataPlaneStimulus_TrafficStimulus_{
					TrafficStimulus: &tv.DataPlaneStimulus_TrafficStimulus{},
				},
			},
		},
	}
	alarmAction = &tv.Action{
		Actions: &tv.Action_AlarmStimulus{
			AlarmStimulus: &tv.AlarmStimulus{},
		},
	}
	portAction = &tv.Action{
		Actions: &tv.Action_PortStimulus{
			PortStimulus: &tv.PortStimulus{},
		},
	}
	managementAction = &tv.Action{
		Actions: &tv.Action_ManagementOperation{
			ManagementOperation: &tv.ManagementOperation{},
		},
	}
)

func TestProcessSequentialActionGroup(t *testing.T) {
	var (
		emptySag = &tv.SequentialActionGroup{}
		validSag = &tv.SequentialActionGroup{
			Actions: []*tv.Action{
				emptyAction,
				emptyAction,
			},
		}
	)
	type args struct {
		sag *tv.SequentialActionGroup
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Sequential Action Group Test",
			args: args{sag: emptySag},
			want: true,
		},
		{
			name: "Valid Sequential Action Group Test",
			args: args{sag: validSag},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processSequentialActionGroup(tt.args.sag); got != tt.want {
				t.Errorf("ProcessSequentialActionGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessParallelActionGroup(t *testing.T) {
	var (
		emptyPag = &tv.ParallelActionGroup{}
		validPag = &tv.ParallelActionGroup{
			Actions: []*tv.Action{
				emptyAction,
				emptyAction,
			},
		}
	)
	type args struct {
		pag *tv.ParallelActionGroup
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Parallel Action Group Test",
			args: args{pag: emptyPag},
			want: true,
		},
		{
			name: "Valid Parallel Action Group Test",
			args: args{pag: validPag},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processParallelActionGroup(tt.args.pag); got != tt.want {
				t.Errorf("ProcessParallelActionGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessRandomizedActionGroup(t *testing.T) {

	var (
		emptyRag = &tv.RandomizedActionGroup{}
	)
	type args struct {
		rag *tv.RandomizedActionGroup
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Random Action Group Test",
			args: args{rag: emptyRag},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processRandomizedActionGroup(tt.args.rag); got != tt.want {
				t.Errorf("ProcessRandomizedActionGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessAction(t *testing.T) {
	gnmi.Init(TestTarget)
	portMap := make(map[string]string)
	portMap["1"] = "veth0"
	portMap["1"] = "veth2"
	dpMode := "direct"
	p4rt.Init(TestTarget, dpMode, portMap)
	defer p4rt.TearDown()
	defer gnmi.TearDown()

	type args struct {
		action *tv.Action
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Config Action",
			args: args{action: configAction},
			want: true,
		},
		{
			name: "Control Plane Write Action",
			args: args{action: controlPlaneWriteAction},
			want: false,
		},
		{
			name: "Control Plane Pipeline Config Action",
			args: args{action: controlPlanePipelineConfigAction},
			want: false,
		},
		{
			name: "Control Plane Packet Out Action",
			args: args{action: controlPlanePacketOutAction},
			want: true,
		},
		{
			name: "Data Plane Traffic Action",
			args: args{action: dataPlaneAction},
			want: false,
		},
		{
			name: "Management Action",
			args: args{action: managementAction},
			want: false,
		},
		{
			name: "Port Action",
			args: args{action: portAction},
			want: false,
		},
		{
			name: "Alarm Action",
			args: args{action: alarmAction},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processAction(tt.args.action); got != tt.want {
				t.Errorf("ProcessAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
