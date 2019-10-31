/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package orchestrator

import (
	"testing"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework"
	tg "github.com/stratum/testvectors/proto/target"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var (
	TestTarget  = &tg.Target{Address: "localhost:50001"}
	EmptyTarget *tg.Target
)

func TestProcessTestVector(t *testing.T) {
	framework.InitGNMI(TestTarget)
	defer framework.TearDownGNMI()
	var (
		emptyTestVector = &tv.TestVector{}

		allScenariosTV = &tv.TestVector{
			TestCases: []*tv.TestCase{
				{
					ActionGroups: []*tv.ActionGroup{
						{
							ActionGroup: &tv.ActionGroup_SequentialActionGroup{
								SequentialActionGroup: &tv.SequentialActionGroup{},
							},
							ActionGroupId: "ag1",
						},
						{
							ActionGroup: &tv.ActionGroup_ParallelActionGroup{
								ParallelActionGroup: &tv.ParallelActionGroup{},
							},
							ActionGroupId: "ag2",
						},
						{
							ActionGroup: &tv.ActionGroup_RandomizedActionGroup{
								RandomizedActionGroup: &tv.RandomizedActionGroup{},
							},
							ActionGroupId: "ag3",
						},
						{},
					},
					TestCaseId: "tc1",
				},
				{
					Expectations: []*tv.Expectation{
						{
							Expectations: &tv.Expectation_ConfigExpectation{
								ConfigExpectation: &tv.ConfigExpectation{},
							},
							ExpectationId: "e1",
						},
						{
							Expectations: &tv.Expectation_ControlPlaneExpectation{
								ControlPlaneExpectation: &tv.ControlPlaneExpectation{},
							},
							ExpectationId: "e2",
						},
						{
							Expectations: &tv.Expectation_TelemetryExpectation{
								TelemetryExpectation: &tv.TelemetryExpectation{},
							},
							ExpectationId: "e3",
						},
						{
							Expectations: &tv.Expectation_DataPlaneExpectation{
								DataPlaneExpectation: &tv.DataPlaneExpectation{},
							},
							ExpectationId: "e4",
						},
						{},
					},
					TestCaseId: "tc2",
				},
			},
		}
	)

	type args struct {
		tv1 *tv.TestVector
		tg1 *tg.Target
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Target",
			args: args{tv1: emptyTestVector, tg1: EmptyTarget},
			want: false,
		},
		{
			name: "Empty Test Vector",
			args: args{tv1: emptyTestVector, tg1: TestTarget},
			want: true,
		},
		{
			name: "All Scenarios Test Vector",
			args: args{tv1: allScenariosTV, tg1: TestTarget},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessTestVector(tt.args.tv1, tt.args.tg1); got != tt.want {
				t.Errorf("ProcessTestVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
validTestVector = &tv.TestVector{
			TestCases: []*tv.TestCase{
				&tv.TestCase{
					ActionGroups: []*tv.ActionGroup{
						&tv.ActionGroup{
							ActionGroup: &tv.ActionGroup_SequentialActionGroup{
								SequentialActionGroup: &tv.SequentialActionGroup{
									Actions: []*tv.Action{
										&tv.Action{
											Actions: &tv.Action_ConfigOperation{
												ConfigOperation: &tv.ConfigOperation{
													GnmiSetRequest: &gnmi.SetRequest{
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
													GnmiSetResponse: &gnmi.SetResponse{
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
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

	invalidTestVector = &tv.TestVector{
		TestCases: []*tv.TestCase{
			&tv.TestCase{
				ActionGroups: []*tv.ActionGroup{
					&tv.ActionGroup{
						ActionGroup: &tv.ActionGroup_SequentialActionGroup{
							SequentialActionGroup: &tv.SequentialActionGroup{
								Actions: []*tv.Action{
									&tv.Action{
										Actions: &tv.Action_ConfigOperation{
											ConfigOperation: &tv.ConfigOperation{
												GnmiSetRequest: &gnmi.SetRequest{
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
												},
												GnmiSetResponse: &gnmi.SetResponse{
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
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	{
		name: "Invalid Test Vector",
		args: args{tv1: invalidTestVector, tg1: TestTarget},
		want: false,
	},*/

func TestProcessTestCase(t *testing.T) {
	var (
		tc1 = &tv.TestCase{}
	)
	type args struct {
		tc *tv.TestCase
		tg *tg.Target
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Target",
			args: args{tc: tc1, tg: EmptyTarget},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessTestCase(tt.args.tc, tt.args.tg); got != tt.want {
				t.Errorf("ProcessTestCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
