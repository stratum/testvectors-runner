/*
*Copyright 2019-present Open Networking Foundation
*
*SPDX-License-Identifier: Apache-2.0
 */

package framework

import (
	"strings"
	"testing"

	config "github.com/abhilashendurthi/p4runtime/proto/p4/config/v1"
	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	tg "github.com/stratum/testvectors/proto/target"
)

func setupTest() {
	Init(TestTarget)
}

func tearDownTest() {
	TearDown()
}

func TestProcessP4PipelineConfigOperation(t *testing.T) {
	log.Infoln(strings.Repeat("*", 100))
	log.Infoln("Start of TestProcessP4PipelineConfigOperation")
	defer log.Infoln("End of TestProcessP4PipelineConfigOperation")
	setupTest()
	defer tearDownTest()
	var (
		deviceID       uint64 = 1
		electionID            = &v1.Uint128{High: 1, Low: 5}
		pipelineCfgReq        = &v1.SetForwardingPipelineConfigRequest{
			DeviceId:   deviceID,
			ElectionId: electionID,
			Action:     v1.SetForwardingPipelineConfigRequest_VERIFY_AND_COMMIT,
			Config: &v1.ForwardingPipelineConfig{
				P4Info: &config.P4Info{
					PkgInfo: &config.PkgInfo{
						Arch: "v1model",
					},
					Tables: []*config.Table{
						&config.Table{
							Preamble: &config.Preamble{
								Id:          33572104,
								Name:        "ingress.l3_fwd.l3_fwd_table",
								Alias:       "l3_fwd_table",
								Annotations: []string{"@switchstack(\"pipeline_stage: L3_LPM\")"},
							},
							MatchFields: []*config.MatchField{
								&config.MatchField{
									Id:       1,
									Name:     "local_metadata.vrf_id",
									Bitwidth: 10,
									Match: &config.MatchField_MatchType_{
										MatchType: config.MatchField_EXACT,
									},
								},
								&config.MatchField{
									Id:       2,
									Name:     "hdr.ipv4_base.dst_addr",
									Bitwidth: 32,
									Match: &config.MatchField_MatchType_{
										MatchType: config.MatchField_LPM,
									},
								},
							},
							ActionRefs: []*config.ActionRef{
								&config.ActionRef{
									Id: 16782370,
								},
								&config.ActionRef{
									Id: 16819938,
								},
								&config.ActionRef{
									Id: 16822646,
								},
							},
							ConstDefaultActionId: 16819938,
							ImplementationId:     285233948,
							Size:                 1024,
						},
						&config.Table{
							Preamble: &config.Preamble{
								Id:          33582129,
								Name:        "ingress.l2_fwd.l2_unicast_table",
								Alias:       "l2_unicast_table",
								Annotations: []string{"@switchstack(\"pipeline_stage: L2\")"},
							},
							MatchFields: []*config.MatchField{
								&config.MatchField{
									Id:       1,
									Name:     "hdr.ethernet.dst_addr",
									Bitwidth: 48,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_EXACT},
								},
							},
							ActionRefs: []*config.ActionRef{
								&config.ActionRef{
									Id: 16838806,
								},
								&config.ActionRef{
									Id:          16800567,
									Scope:       config.ActionRef_DEFAULT_ONLY,
									Annotations: []string{"@defaultonly"},
								},
							},
							Size: 1024,
						},
						&config.Table{
							Preamble: &config.Preamble{
								Id:          33598026,
								Name:        "ingress.punt.punt_table",
								Alias:       "punt_table",
								Annotations: []string{"@switchstack(\"pipeline_stage: INGRESS_ACL\")"},
							},
							MatchFields: []*config.MatchField{
								&config.MatchField{
									Id:       1,
									Name:     "standard_metadata.ingress_port",
									Bitwidth: 9,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       2,
									Name:     "standard_metadata.egress_spec",
									Bitwidth: 9,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       3,
									Name:     "hdr.ethernet.ether_type",
									Bitwidth: 16,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       4,
									Name:     "hdr.ipv4_base.diffserv",
									Bitwidth: 8,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       5,
									Name:     "hdr.ipve_base.ttl",
									Bitwidth: 8,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       6,
									Name:     "hdr.ipv4_base.src_addr",
									Bitwidth: 32,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       7,
									Name:     "hdr.ipv4_base.dst_addr",
									Bitwidth: 32,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       8,
									Name:     "hdr.ipv4_base.protocol",
									Bitwidth: 8,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       9,
									Name:     "local_metadata.icmp_code",
									Bitwidth: 8,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       10,
									Name:     "hdr.vlan_tag[0].vid",
									Bitwidth: 12,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       11,
									Name:     "hdr.vlan_tag[0].pcp",
									Bitwidth: 3,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       12,
									Name:     "local_metadata.class_id",
									Bitwidth: 8,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
								&config.MatchField{
									Id:       13,
									Name:     "local_metadata.vrf_id",
									Bitwidth: 10,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
							},
							ActionRefs: []*config.ActionRef{
								&config.ActionRef{
									Id: 16824483,
								},
								&config.ActionRef{
									Id: 16804007,
								},
								&config.ActionRef{
									Id: 16820507,
								},
								&config.ActionRef{
									Id:          16800567,
									Scope:       config.ActionRef_DEFAULT_ONLY,
									Annotations: []string{"@defaultonly"},
								},
							},
							DirectResourceIds: []uint32{318787623, 352360565},
							Size:              25,
						},
						&config.Table{
							Preamble: &config.Preamble{
								Id:          33573106,
								Name:        "ingress.my_station_table",
								Alias:       "my_station_table",
								Annotations: []string{"@switchstack(\"pipeline_stage: L2\")"},
							},
							MatchFields: []*config.MatchField{
								&config.MatchField{
									Id:       1,
									Name:     "hdr.ethernet.dst_addr",
									Bitwidth: 48,
									Match:    &config.MatchField_MatchType_{MatchType: config.MatchField_TERNARY},
								},
							},
							ActionRefs: []*config.ActionRef{
								&config.ActionRef{
									Id: 16832439,
								},
								&config.ActionRef{
									Id: 16819938,
								},
							},
							Size: 1024,
						},
					},
					Actions: []*config.Action{
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16819938,
								Name:  "nop",
								Alias: "nop",
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16800567,
								Name:  "NoAction",
								Alias: "NoAction",
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16822646,
								Name:  "ingress.l3_fwd.drop",
								Alias: "drop",
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16782370,
								Name:  "ingress.l3_fwd.set_nexthop",
								Alias: "set_nexthop",
							},
							Params: []*config.Action_Param{
								&config.Action_Param{
									Id:       1,
									Name:     "port",
									Bitwidth: 9,
								},
								&config.Action_Param{
									Id:       2,
									Name:     "smac",
									Bitwidth: 48,
								},
								&config.Action_Param{
									Id:       3,
									Name:     "dmac",
									Bitwidth: 48,
								},
								&config.Action_Param{
									Id:       4,
									Name:     "dst_vlan",
									Bitwidth: 12,
								},
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16838806,
								Name:  "ingress.l2_fwd.set_egress_port",
								Alias: "l2_fwd.set_egress_port",
							},
							Params: []*config.Action_Param{
								&config.Action_Param{
									Id:       1,
									Name:     "port",
									Bitwidth: 9,
								},
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16824483,
								Name:  "ingress.punt.set_queue_and_clone_to_cpu",
								Alias: "set_queue_and_clone_to_cpu",
							},
							Params: []*config.Action_Param{
								&config.Action_Param{
									Id:       1,
									Name:     "queue_id",
									Bitwidth: 5,
								},
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16804007,
								Name:  "ingress.punt.set_queue_and_send_to_cpu",
								Alias: "set_queue_and_send_to_cpu",
							},
							Params: []*config.Action_Param{
								&config.Action_Param{
									Id:       1,
									Name:     "queue_id",
									Bitwidth: 5,
								},
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16820507,
								Name:  "ingress.punt.set_egress_port",
								Alias: "punt.set_egress_port",
							},
							Params: []*config.Action_Param{
								&config.Action_Param{
									Id:       1,
									Name:     "port",
									Bitwidth: 9,
								},
							},
						},
						&config.Action{
							Preamble: &config.Preamble{
								Id:    16832439,
								Name:  "ingress.set_l3_admit",
								Alias: "set_l3_admit",
							},
						},
					},
					ActionProfiles: []*config.ActionProfile{
						&config.ActionProfile{
							Preamble: &config.Preamble{
								Id:    285233948,
								Name:  "ingress.l3_fwd.wcmp_action_profile",
								Alias: "wcmp_action_profile",
							},
							TableIds:     []uint32{33572104},
							WithSelector: true,
							Size:         1024,
							MaxGroupSize: 8,
						},
					},
					DirectCounters: []*config.DirectCounter{
						&config.DirectCounter{
							Preamble: &config.Preamble{
								Id:    318787623,
								Name:  "ingress.punt.punt_packet_counter",
								Alias: "punt_packet_counter",
							},
							Spec: &config.CounterSpec{
								Unit: config.CounterSpec_PACKETS,
							},
							DirectTableId: 33598026,
						},
					},
					DirectMeters: []*config.DirectMeter{

						&config.DirectMeter{
							Preamble: &config.Preamble{
								Id:    352360565,
								Name:  "ingress.punt.ingress_port_meter",
								Alias: "ingress_port_meter",
							},
							Spec: &config.MeterSpec{
								Unit: config.MeterSpec_BYTES,
							},
							DirectTableId: 33598026,
						},
					},
					ControllerPacketMetadata: []*config.ControllerPacketMetadata{
						&config.ControllerPacketMetadata{
							Preamble: &config.Preamble{
								Id:          67146229,
								Name:        "packet_in",
								Alias:       "packet_in",
								Annotations: []string{"@controller_header(\"packet_in\")"},
							},
							Metadata: []*config.ControllerPacketMetadata_Metadata{
								&config.ControllerPacketMetadata_Metadata{
									Id:          1,
									Name:        "ingress_physical_port",
									Annotations: []string{"@switchstack(\"field_type: P4_FIELD_TYPE_INGRESS_PORT\")", "@proto_tag(1)"},
									Bitwidth:    9,
								},
								&config.ControllerPacketMetadata_Metadata{
									Id:          2,
									Name:        "padding1",
									Annotations: []string{"@proto_tag(2)"},
									Bitwidth:    7,
								},
								&config.ControllerPacketMetadata_Metadata{
									Id:          3,
									Name:        "ingress_logical_port",
									Annotations: []string{"@proto_tag(3)"},
									Bitwidth:    32,
								},
								&config.ControllerPacketMetadata_Metadata{
									Id:          4,
									Name:        "target_egress_port",
									Annotations: []string{"@switchstack(\"field_type: P4_FIELD_TYPE_EGRESS_PORT\")", "@proto_tag(4)"},
									Bitwidth:    9,
								},
								&config.ControllerPacketMetadata_Metadata{
									Id:          2,
									Name:        "padding2",
									Annotations: []string{"@proto_tag(5)"},
									Bitwidth:    7,
								},
							},
						},
						&config.ControllerPacketMetadata{
							Preamble: &config.Preamble{
								Id:          67121543,
								Name:        "packet_out",
								Alias:       "packet_out",
								Annotations: []string{"@not_extracted_in_egress", "@controller_header(\"packet_out\")"},
							},
							Metadata: []*config.ControllerPacketMetadata_Metadata{
								&config.ControllerPacketMetadata_Metadata{
									Id:          1,
									Name:        "egress_physical_port",
									Annotations: []string{"@switchstack(\"field_type: P4_FIELD_TYPE_EGRESS_PORT\")", "@proto_tag(1)"},
									Bitwidth:    9,
								},
								&config.ControllerPacketMetadata_Metadata{
									Id:          2,
									Name:        "submit_to_ingress",
									Annotations: []string{"@proto_tag(2)"},
									Bitwidth:    1,
								},
								&config.ControllerPacketMetadata_Metadata{
									Id:          3,
									Name:        "padding",
									Annotations: []string{"@proto_tag(3)"},
									Bitwidth:    6,
								},
							},
						},
					},
					TypeInfo: &config.P4TypeInfo{},
				},
				P4DeviceConfig: []byte("{\"header_types\":[{\"name\":\"scalars_0\",\"id\":0,\"fields\":[[\"tmp\",21,false],[\"key_0\",12,false],[\"key_1\",3,false],[\"local_metadata_t.vrf_id\",10,false],[\"local_metadata_t.class_id\",8,false],[\"local_metadata_t.cpu_cos_queue_id\",5,false],[\"local_metadata_t.skip_egress\",1,false],[\"local_metadata_t.egress_spec_at_punt_match\",9,false],[\"local_metadata_t.color\",2,false],[\"local_metadata_t.l4_src_port\",16,false],[\"local_metadata_t.l4_dst_port\",16,false],[\"local_metadata_t.icmp_code\",8,false],[\"local_metadata_t.l3_admit\",1,false],[\"local_metadata_t.dst_vlan\",12,false],[\"_padding_0\",4,false]]},{\"name\":\"standard_metadata\",\"id\":1,\"fields\":[[\"ingress_port\",9,false],[\"egress_spec\",9,false],[\"egress_port\",9,false],[\"clone_spec\",32,false],[\"instance_type\",32,false],[\"drop\",1,false],[\"recirculate_port\",16,false],[\"packet_length\",32,false],[\"enq_timestamp\",32,false],[\"enq_qdepth\",19,false],[\"deq_timedelta\",32,false],[\"deq_qdepth\",19,false],[\"ingress_global_timestamp\",48,false],[\"egress_global_timestamp\",48,false],[\"lf_field_list\",32,false],[\"mcast_grp\",16,false],[\"resubmit_flag\",32,false],[\"egress_rid\",16,false],[\"recirculate_flag\",32,false],[\"checksum_error\",1,false],[\"parser_error\",32,false],[\"priority\",3,false],[\"_padding\",2,false]]},{\"name\":\"ethernet_t\",\"id\":2,\"fields\":[[\"dst_addr\",48,false],[\"src_addr\",48,false],[\"ether_type\",16,false]]},{\"name\":\"ipv4_base_t\",\"id\":3,\"fields\":[[\"version\",4,false],[\"ihl\",4,false],[\"diffserv\",8,false],[\"total_len\",16,false],[\"identification\",16,false],[\"flags\",3,false],[\"frag_offset\",13,false],[\"ttl\",8,false],[\"protocol\",8,false],[\"hdr_checksum\",16,false],[\"src_addr\",32,false],[\"dst_addr\",32,false]]},{\"name\":\"ipv6_base_t\",\"id\":4,\"fields\":[[\"version\",4,false],[\"traffic_class\",8,false],[\"flow_label\",20,false],[\"payload_length\",16,false],[\"next_header\",8,false],[\"hop_limit\",8,false],[\"src_addr\",128,false],[\"dst_addr\",128,false]]},{\"name\":\"icmp_header_t\",\"id\":5,\"fields\":[[\"icmp_type\",8,false],[\"code\",8,false],[\"checksum\",16,false]]},{\"name\":\"tcp_t\",\"id\":6,\"fields\":[[\"src_port\",16,false],[\"dst_port\",16,false],[\"seq_no\",32,false],[\"ack_no\",32,false],[\"data_offset\",4,false],[\"res\",4,false],[\"flags\",8,false],[\"window\",16,false],[\"checksum\",16,false],[\"urgent_ptr\",16,false]]},{\"name\":\"udp_t\",\"id\":7,\"fields\":[[\"src_port\",16,false],[\"dst_port\",16,false],[\"hdr_length\",16,false],[\"checksum\",16,false]]},{\"name\":\"arp_t\",\"id\":8,\"fields\":[[\"hw_type\",16,false],[\"proto_type\",16,false],[\"hw_addr_len\",8,false],[\"proto_addr_len\",8,false],[\"opcode\",16,false],[\"sender_hw_addr\",48,false],[\"sender_proto_addr\",32,false],[\"target_hw_addr\",48,false],[\"target_proto_addr\",32,false]]},{\"name\":\"packet_in_header_t\",\"id\":9,\"fields\":[[\"ingress_physical_port\",9,false],[\"padding1\",7,false],[\"ingress_logical_port\",32,false],[\"target_egress_port\",9,false],[\"padding2\",7,false]]},{\"name\":\"packet_out_header_t\",\"id\":10,\"fields\":[[\"egress_physical_port\",9,false],[\"submit_to_ingress\",1,false],[\"padding\",6,false]]},{\"name\":\"vlan_tag_t\",\"id\":11,\"fields\":[[\"pcp\",3,false],[\"cfi\",1,false],[\"vid\",12,false],[\"ether_type\",16,false]]}],\"headers\":[{\"name\":\"scalars\",\"id\":0,\"header_type\":\"scalars_0\",\"metadata\":true,\"pi_omit\":true},{\"name\":\"standard_metadata\",\"id\":1,\"header_type\":\"standard_metadata\",\"metadata\":true,\"pi_omit\":true},{\"name\":\"ethernet\",\"id\":2,\"header_type\":\"ethernet_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"ipv4_base\",\"id\":3,\"header_type\":\"ipv4_base_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"ipv6_base\",\"id\":4,\"header_type\":\"ipv6_base_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"icmp_header\",\"id\":5,\"header_type\":\"icmp_header_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"tcp\",\"id\":6,\"header_type\":\"tcp_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"udp\",\"id\":7,\"header_type\":\"udp_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"arp\",\"id\":8,\"header_type\":\"arp_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"packet_in\",\"id\":9,\"header_type\":\"packet_in_header_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"packet_out\",\"id\":10,\"header_type\":\"packet_out_header_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"vlan_tag[0]\",\"id\":11,\"header_type\":\"vlan_tag_t\",\"metadata\":false,\"pi_omit\":true},{\"name\":\"vlan_tag[1]\",\"id\":12,\"header_type\":\"vlan_tag_t\",\"metadata\":false,\"pi_omit\":true}],\"header_stacks\":[{\"name\":\"vlan_tag\",\"id\":0,\"header_type\":\"vlan_tag_t\",\"size\":2,\"header_ids\":[11,12]}],\"header_union_types\":[],\"header_unions\":[],\"header_union_stacks\":[],\"field_lists\":[{\"id\":1,\"name\":\"fl\",\"source_info\":{\"filename\":\"max.p4\",\"line\":28,\"column\":11,\"source_fragment\":\"{standard_metadata.ingress_port}\"},\"elements\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"ingress_port\"]}]}],\"errors\":[[\"NoError\",1],[\"PacketTooShort\",2],[\"NoMatch\",3],[\"StackOutOfBounds\",4],[\"HeaderTooShort\",5],[\"ParserTimeout\",6],[\"ParserInvalidArgument\",7]],\"enums\":[],\"parsers\":[{\"name\":\"parser\",\"id\":0,\"init_state\":\"start\",\"parse_states\":[{\"name\":\"start\",\"id\":0,\"parser_ops\":[],\"transitions\":[{\"type\":\"hexstr\",\"value\":\"0x00fd\",\"mask\":null,\"next_state\":\"parse_cpu_header\"},{\"value\":\"default\",\"mask\":null,\"next_state\":\"parse_ethernet\"}],\"transition_key\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"ingress_port\"]}]},{\"name\":\"parse_ethernet\",\"id\":1,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"ethernet\"}],\"op\":\"extract\"}],\"transitions\":[{\"type\":\"hexstr\",\"value\":\"0x8100\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x9100\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x9200\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x9300\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x0800\",\"mask\":null,\"next_state\":\"parse_ipv4\"},{\"type\":\"hexstr\",\"value\":\"0x86dd\",\"mask\":null,\"next_state\":\"parse_ipv6\"},{\"type\":\"hexstr\",\"value\":\"0x0806\",\"mask\":null,\"next_state\":\"parse_arp\"},{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[{\"type\":\"field\",\"value\":[\"ethernet\",\"ether_type\"]}]},{\"name\":\"parse_vlan\",\"id\":2,\"parser_ops\":[{\"parameters\":[{\"type\":\"stack\",\"value\":\"vlan_tag\"}],\"op\":\"extract\"}],\"transitions\":[{\"type\":\"hexstr\",\"value\":\"0x8100\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x9100\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x9200\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x9300\",\"mask\":null,\"next_state\":\"parse_vlan\"},{\"type\":\"hexstr\",\"value\":\"0x0800\",\"mask\":null,\"next_state\":\"parse_ipv4\"},{\"type\":\"hexstr\",\"value\":\"0x86dd\",\"mask\":null,\"next_state\":\"parse_ipv6\"},{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[{\"type\":\"stack_field\",\"value\":[\"vlan_tag\",\"ether_type\"]}]},{\"name\":\"parse_ipv4\",\"id\":3,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"ipv4_base\"}],\"op\":\"extract\"},{\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"tmp\"]},{\"type\":\"expression\",\"value\":{\"type\":\"expression\",\"value\":{\"op\":\"|\",\"left\":{\"type\":\"expression\",\"value\":{\"op\":\"&\",\"left\":{\"type\":\"expression\",\"value\":{\"op\":\"<<\",\"left\":{\"type\":\"expression\",\"value\":{\"op\":\"&\",\"left\":{\"type\":\"field\",\"value\":[\"ipv4_base\",\"frag_offset\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0x1fffff\"}}},\"right\":{\"type\":\"hexstr\",\"value\":\"0x8\"}}},\"right\":{\"type\":\"hexstr\",\"value\":\"0x1fffff\"}}},\"right\":{\"type\":\"expression\",\"value\":{\"op\":\"&\",\"left\":{\"type\":\"expression\",\"value\":{\"op\":\"&\",\"left\":{\"type\":\"field\",\"value\":[\"ipv4_base\",\"protocol\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0x1fffff\"}}},\"right\":{\"type\":\"hexstr\",\"value\":\"0x0001ff\"}}}}}}],\"op\":\"set\"}],\"transitions\":[{\"type\":\"hexstr\",\"value\":\"0x000001\",\"mask\":null,\"next_state\":\"parse_icmp\"},{\"type\":\"hexstr\",\"value\":\"0x000006\",\"mask\":null,\"next_state\":\"parse_tcp\"},{\"type\":\"hexstr\",\"value\":\"0x000011\",\"mask\":null,\"next_state\":\"parse_udp\"},{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[{\"type\":\"field\",\"value\":[\"scalars\",\"tmp\"]}]},{\"name\":\"parse_ipv6\",\"id\":4,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"ipv6_base\"}],\"op\":\"extract\"}],\"transitions\":[{\"type\":\"hexstr\",\"value\":\"0x3a\",\"mask\":null,\"next_state\":\"parse_icmp\"},{\"type\":\"hexstr\",\"value\":\"0x06\",\"mask\":null,\"next_state\":\"parse_tcp\"},{\"type\":\"hexstr\",\"value\":\"0x11\",\"mask\":null,\"next_state\":\"parse_udp\"},{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[{\"type\":\"field\",\"value\":[\"ipv6_base\",\"next_header\"]}]},{\"name\":\"parse_tcp\",\"id\":5,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"tcp\"}],\"op\":\"extract\"},{\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l4_src_port\"]},{\"type\":\"field\",\"value\":[\"tcp\",\"src_port\"]}],\"op\":\"set\"},{\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l4_dst_port\"]},{\"type\":\"field\",\"value\":[\"tcp\",\"dst_port\"]}],\"op\":\"set\"}],\"transitions\":[{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[]},{\"name\":\"parse_udp\",\"id\":6,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"udp\"}],\"op\":\"extract\"},{\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l4_src_port\"]},{\"type\":\"field\",\"value\":[\"udp\",\"src_port\"]}],\"op\":\"set\"},{\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l4_dst_port\"]},{\"type\":\"field\",\"value\":[\"udp\",\"dst_port\"]}],\"op\":\"set\"}],\"transitions\":[{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[]},{\"name\":\"parse_icmp\",\"id\":7,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"icmp_header\"}],\"op\":\"extract\"}],\"transitions\":[{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[]},{\"name\":\"parse_arp\",\"id\":8,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"arp\"}],\"op\":\"extract\"}],\"transitions\":[{\"value\":\"default\",\"mask\":null,\"next_state\":null}],\"transition_key\":[]},{\"name\":\"parse_cpu_header\",\"id\":9,\"parser_ops\":[{\"parameters\":[{\"type\":\"regular\",\"value\":\"packet_out\"}],\"op\":\"extract\"}],\"transitions\":[{\"value\":\"default\",\"mask\":null,\"next_state\":\"parse_ethernet\"}],\"transition_key\":[]}]}],\"parse_vsets\":[],\"deparsers\":[{\"name\":\"deparser\",\"id\":0,\"source_info\":{\"filename\":\"parser.p4\",\"line\":149,\"column\":8,\"source_fragment\":\"pkt_deparser\"},\"order\":[\"packet_in\",\"ethernet\",\"vlan_tag[0]\",\"vlan_tag[1]\",\"ipv4_base\",\"ipv6_base\",\"arp\",\"icmp_header\",\"tcp\",\"udp\"]}],\"meter_arrays\":[{\"name\":\"ingress.punt.ingress_port_meter\",\"id\":0,\"source_info\":{\"filename\":\"max.p4\",\"line\":20,\"column\":40,\"source_fragment\":\"ingress_port_meter\"},\"is_direct\":true,\"rate_count\":2,\"type\":\"bytes\",\"size\":25,\"binding\":\"ingress.punt.punt_table\",\"result_target\":[\"scalars\",\"local_metadata_t.color\"]}],\"counter_arrays\":[{\"name\":\"ingress.punt.punt_packet_counter\",\"id\":0,\"is_direct\":true,\"binding\":\"ingress.punt.punt_table\",\"source_info\":{\"filename\":\"max.p4\",\"line\":22,\"column\":38,\"source_fragment\":\"punt_packet_counter\"}}],\"register_arrays\":[],\"calculations\":[{\"name\":\"calc\",\"id\":0,\"source_info\":{\"filename\":\"ipv4_checksum.p4\",\"line\":29,\"column\":4,\"source_fragment\":\"update_checksum(hdr.ipv4_base.isValid(),...\"},\"algo\":\"csum16\",\"input\":[{\"type\":\"field\",\"value\":[\"ipv4_base\",\"version\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"ihl\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"diffserv\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"total_len\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"identification\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"flags\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"frag_offset\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"ttl\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"protocol\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"src_addr\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"dst_addr\"]}]},{\"name\":\"calc_0\",\"id\":1,\"source_info\":{\"filename\":\"ipv4_checksum.p4\",\"line\":11,\"column\":4,\"source_fragment\":\"verify_checksum(hdr.ipv4_base.isValid(),...\"},\"algo\":\"csum16\",\"input\":[{\"type\":\"field\",\"value\":[\"ipv4_base\",\"version\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"ihl\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"diffserv\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"total_len\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"identification\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"flags\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"frag_offset\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"ttl\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"protocol\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"src_addr\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"dst_addr\"]}]}],\"learn_lists\":[],\"actions\":[{\"name\":\"nop\",\"id\":0,\"runtime_data\":[],\"primitives\":[]},{\"name\":\"nop\",\"id\":1,\"runtime_data\":[],\"primitives\":[]},{\"name\":\"NoAction\",\"id\":2,\"runtime_data\":[],\"primitives\":[]},{\"name\":\"NoAction\",\"id\":3,\"runtime_data\":[],\"primitives\":[]},{\"name\":\"ingress.l3_fwd.drop\",\"id\":4,\"runtime_data\":[],\"primitives\":[{\"op\":\"mark_to_drop\",\"parameters\":[{\"type\":\"header\",\"value\":\"standard_metadata\"}],\"source_info\":{\"filename\":\"max.p4\",\"line\":98,\"column\":18,\"source_fragment\":\"mark_to_drop(standard_metadata)\"}}]},{\"name\":\"ingress.l3_fwd.set_nexthop\",\"id\":5,\"runtime_data\":[{\"name\":\"port\",\"bitwidth\":9},{\"name\":\"smac\",\"bitwidth\":48},{\"name\":\"dmac\",\"bitwidth\":48},{\"name\":\"dst_vlan\",\"bitwidth\":12}],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},{\"type\":\"runtime_data\",\"value\":0}],\"source_info\":{\"filename\":\"max.p4\",\"line\":104,\"column\":6,\"source_fragment\":\"standard_metadata.egress_spec=port\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.dst_vlan\"]},{\"type\":\"runtime_data\",\"value\":3}],\"source_info\":{\"filename\":\"max.p4\",\"line\":105,\"column\":6,\"source_fragment\":\"local_metadata.dst_vlan=dst_vlan\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"ethernet\",\"src_addr\"]},{\"type\":\"runtime_data\",\"value\":1}],\"source_info\":{\"filename\":\"max.p4\",\"line\":106,\"column\":6,\"source_fragment\":\"hdr.ethernet.src_addr=smac\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"ethernet\",\"dst_addr\"]},{\"type\":\"runtime_data\",\"value\":2}],\"source_info\":{\"filename\":\"max.p4\",\"line\":107,\"column\":6,\"source_fragment\":\"hdr.ethernet.dst_addr=dmac\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"ipv4_base\",\"ttl\"]},{\"type\":\"expression\",\"value\":{\"type\":\"expression\",\"value\":{\"op\":\"&\",\"left\":{\"type\":\"expression\",\"value\":{\"op\":\"+\",\"left\":{\"type\":\"field\",\"value\":[\"ipv4_base\",\"ttl\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0xff\"}}},\"right\":{\"type\":\"hexstr\",\"value\":\"0xff\"}}}}],\"source_info\":{\"filename\":\"max.p4\",\"line\":108,\"column\":6,\"source_fragment\":\"hdr.ipv4_base.ttl=hdr.ipv4_base.ttl-1\"}}]},{\"name\":\"ingress.l2_fwd.set_egress_port\",\"id\":6,\"runtime_data\":[{\"name\":\"port\",\"bitwidth\":9}],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},{\"type\":\"runtime_data\",\"value\":0}],\"source_info\":{\"filename\":\"max.p4\",\"line\":144,\"column\":6,\"source_fragment\":\"standard_metadata.egress_spec=port\"}}]},{\"name\":\"ingress.punt.set_queue_and_clone_to_cpu\",\"id\":7,\"runtime_data\":[{\"name\":\"queue_id\",\"bitwidth\":5}],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.cpu_cos_queue_id\"]},{\"type\":\"runtime_data\",\"value\":0}],\"source_info\":{\"filename\":\"max.p4\",\"line\":25,\"column\":4,\"source_fragment\":\"local_metadata.cpu_cos_queue_id=queue_id\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.egress_spec_at_punt_match\"]},{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":26,\"column\":4,\"source_fragment\":\"local_metadata.egress_spec_at_punt_match=standard_metadata.egress_spec\"}},{\"op\":\"clone_ingress_pkt_to_egress\",\"parameters\":[{\"type\":\"hexstr\",\"value\":\"0x000001ff\"},{\"type\":\"hexstr\",\"value\":\"0x1\"}],\"source_info\":{\"filename\":\"max.p4\",\"line\":27,\"column\":4,\"source_fragment\":\"clone3(CloneType.I2E,511,...\"}}]},{\"name\":\"ingress.punt.set_queue_and_send_to_cpu\",\"id\":8,\"runtime_data\":[{\"name\":\"queue_id\",\"bitwidth\":5}],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.cpu_cos_queue_id\"]},{\"type\":\"runtime_data\",\"value\":0}],\"source_info\":{\"filename\":\"max.p4\",\"line\":34,\"column\":4,\"source_fragment\":\"local_metadata.cpu_cos_queue_id=queue_id\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.egress_spec_at_punt_match\"]},{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":35,\"column\":4,\"source_fragment\":\"local_metadata.egress_spec_at_punt_match=standard_metadata.egress_spec\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},{\"type\":\"hexstr\",\"value\":\"0x00fd\"}],\"source_info\":{\"filename\":\"max.p4\",\"line\":36,\"column\":4,\"source_fragment\":\"standard_metadata.egress_spec=0xFD\"}}]},{\"name\":\"ingress.punt.set_egress_port\",\"id\":9,\"runtime_data\":[{\"name\":\"port\",\"bitwidth\":9}],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.egress_spec_at_punt_match\"]},{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":42,\"column\":4,\"source_fragment\":\"local_metadata.egress_spec_at_punt_match=standard_metadata.egress_spec\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},{\"type\":\"runtime_data\",\"value\":0}],\"source_info\":{\"filename\":\"max.p4\",\"line\":43,\"column\":4,\"source_fragment\":\"standard_metadata.egress_spec=port\"}}]},{\"name\":\"ingress.set_l3_admit\",\"id\":10,\"runtime_data\":[],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l3_admit\"]},{\"type\":\"hexstr\",\"value\":\"0x01\"}],\"source_info\":{\"filename\":\"max.p4\",\"line\":167,\"column\":6,\"source_fragment\":\"local_metadata.l3_admit=1w1\"}}]},{\"name\":\"act\",\"id\":11,\"runtime_data\":[],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},{\"type\":\"field\",\"value\":[\"packet_out\",\"egress_physical_port\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":184,\"column\":8,\"source_fragment\":\"standard_metadata.egress_spec=hdr.packet_out.egress_physical_port\"}},{\"op\":\"remove_header\",\"parameters\":[{\"type\":\"header\",\"value\":\"packet_out\"}],\"source_info\":{\"filename\":\"max.p4\",\"line\":185,\"column\":8,\"source_fragment\":\"hdr.packet_out.setInvalid()\"}}]},{\"name\":\"act_0\",\"id\":12,\"runtime_data\":[],\"primitives\":[{\"op\":\"exit\",\"parameters\":[],\"source_info\":{\"filename\":\"max.p4\",\"line\":198,\"column\":6,\"source_fragment\":\"exit\"}}]},{\"name\":\"act_1\",\"id\":13,\"runtime_data\":[],\"primitives\":[{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"key_0\"]},{\"type\":\"field\",\"value\":[\"vlan_tag[0]\",\"vid\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":70,\"column\":6,\"source_fragment\":\"hdr.vlan_tag[0].vid\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"scalars\",\"key_1\"]},{\"type\":\"field\",\"value\":[\"vlan_tag[0]\",\"pcp\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":71,\"column\":6,\"source_fragment\":\"hdr.vlan_tag[0].pcp\"}}]},{\"name\":\"act_2\",\"id\":14,\"runtime_data\":[],\"primitives\":[{\"op\":\"add_header\",\"parameters\":[{\"type\":\"header\",\"value\":\"packet_in\"}],\"source_info\":{\"filename\":\"max.p4\",\"line\":209,\"column\":12,\"source_fragment\":\"hdr.packet_in.setValid()\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"packet_in\",\"ingress_physical_port\"]},{\"type\":\"field\",\"value\":[\"standard_metadata\",\"ingress_port\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":210,\"column\":12,\"source_fragment\":\"hdr.packet_in.ingress_physical_port=standard_metadata.ingress_port\"}},{\"op\":\"assign\",\"parameters\":[{\"type\":\"field\",\"value\":[\"packet_in\",\"target_egress_port\"]},{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.egress_spec_at_punt_match\"]}],\"source_info\":{\"filename\":\"max.p4\",\"line\":211,\"column\":12,\"source_fragment\":\"hdr.packet_in.target_egress_port=local_metadata.egress_spec_at_punt_match\"}},{\"op\":\"exit\",\"parameters\":[],\"source_info\":{\"filename\":\"max.p4\",\"line\":213,\"column\":12,\"source_fragment\":\"exit\"}}]}],\"pipelines\":[{\"name\":\"ingress\",\"id\":0,\"source_info\":{\"filename\":\"max.p4\",\"line\":162,\"column\":8,\"source_fragment\":\"ingress\"},\"init_table\":\"node_2\",\"tables\":[{\"name\":\"tbl_act\",\"id\":0,\"source_info\":{\"filename\":\"max.p4\",\"line\":184,\"column\":38,\"source_fragment\":\"=hdr.packet_out.egress_physical_port;...\"},\"key\":[],\"match_type\":\"exact\",\"type\":\"simple\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[11],\"actions\":[\"act\"],\"base_default_next\":\"node_4\",\"next_tables\":{\"act\":\"node_4\"},\"default_entry\":{\"action_id\":11,\"action_const\":true,\"action_data\":[],\"action_entry_const\":true}},{\"name\":\"ingress.my_station_table\",\"id\":1,\"source_info\":{\"filename\":\"max.p4\",\"line\":171,\"column\":8,\"source_fragment\":\"my_station_table\"},\"key\":[{\"match_type\":\"ternary\",\"name\":\"hdr.ethernet.dst_addr\",\"target\":[\"ethernet\",\"dst_addr\"],\"mask\":null}],\"match_type\":\"ternary\",\"type\":\"simple\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[10,1],\"actions\":[\"ingress.set_l3_admit\",\"nop\"],\"base_default_next\":\"node_6\",\"next_tables\":{\"ingress.set_l3_admit\":\"node_6\",\"nop\":\"node_6\"},\"default_entry\":{\"action_id\":1,\"action_const\":false,\"action_data\":[],\"action_entry_const\":false}},{\"name\":\"ingress.l3_fwd.l3_fwd_table\",\"id\":2,\"source_info\":{\"filename\":\"max.p4\",\"line\":115,\"column\":8,\"source_fragment\":\"l3_fwd_table\"},\"key\":[{\"match_type\":\"exact\",\"name\":\"local_metadata.vrf_id\",\"target\":[\"scalars\",\"local_metadata_t.vrf_id\"],\"mask\":null},{\"match_type\":\"lpm\",\"name\":\"hdr.ipv4_base.dst_addr\",\"target\":[\"ipv4_base\",\"dst_addr\"],\"mask\":null}],\"match_type\":\"lpm\",\"type\":\"indirect_ws\",\"action_profile\":\"ingress.l3_fwd.wcmp_action_profile\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[5,0,4],\"actions\":[\"ingress.l3_fwd.set_nexthop\",\"nop\",\"ingress.l3_fwd.drop\"],\"base_default_next\":\"tbl_act_1\",\"next_tables\":{\"ingress.l3_fwd.set_nexthop\":\"tbl_act_1\",\"nop\":\"tbl_act_1\",\"ingress.l3_fwd.drop\":\"tbl_act_1\"}},{\"name\":\"ingress.l2_fwd.l2_unicast_table\",\"id\":3,\"source_info\":{\"filename\":\"max.p4\",\"line\":148,\"column\":10,\"source_fragment\":\"l2_unicast_table\"},\"key\":[{\"match_type\":\"exact\",\"name\":\"hdr.ethernet.dst_addr\",\"target\":[\"ethernet\",\"dst_addr\"],\"mask\":null}],\"match_type\":\"exact\",\"type\":\"simple\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[6,2],\"actions\":[\"ingress.l2_fwd.set_egress_port\",\"NoAction\"],\"base_default_next\":\"tbl_act_1\",\"next_tables\":{\"ingress.l2_fwd.set_egress_port\":\"tbl_act_1\",\"NoAction\":\"tbl_act_1\"},\"default_entry\":{\"action_id\":2,\"action_const\":false,\"action_data\":[],\"action_entry_const\":false}},{\"name\":\"tbl_act_0\",\"id\":4,\"source_info\":{\"filename\":\"max.p4\",\"line\":198,\"column\":6,\"source_fragment\":\"exit\"},\"key\":[],\"match_type\":\"exact\",\"type\":\"simple\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[12],\"actions\":[\"act_0\"],\"base_default_next\":\"tbl_act_1\",\"next_tables\":{\"act_0\":\"tbl_act_1\"},\"default_entry\":{\"action_id\":12,\"action_const\":true,\"action_data\":[],\"action_entry_const\":true}},{\"name\":\"tbl_act_1\",\"id\":5,\"source_info\":{\"filename\":\"max.p4\",\"line\":70,\"column\":6,\"source_fragment\":\"hdr.vlan_tag[0].vid:ternary;...\"},\"key\":[],\"match_type\":\"exact\",\"type\":\"simple\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[13],\"actions\":[\"act_1\"],\"base_default_next\":\"ingress.punt.punt_table\",\"next_tables\":{\"act_1\":\"ingress.punt.punt_table\"},\"default_entry\":{\"action_id\":13,\"action_const\":true,\"action_data\":[],\"action_entry_const\":true}},{\"name\":\"ingress.punt.punt_table\",\"id\":6,\"source_info\":{\"filename\":\"max.p4\",\"line\":50,\"column\":8,\"source_fragment\":\"punt_table\"},\"key\":[{\"match_type\":\"ternary\",\"name\":\"standard_metadata.ingress_port\",\"target\":[\"standard_metadata\",\"ingress_port\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"standard_metadata.egress_spec\",\"target\":[\"standard_metadata\",\"egress_spec\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.ethernet.ether_type\",\"target\":[\"ethernet\",\"ether_type\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.ipv4_base.diffserv\",\"target\":[\"ipv4_base\",\"diffserv\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.ipv4_base.ttl\",\"target\":[\"ipv4_base\",\"ttl\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.ipv4_base.src_addr\",\"target\":[\"ipv4_base\",\"src_addr\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.ipv4_base.dst_addr\",\"target\":[\"ipv4_base\",\"dst_addr\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.ipv4_base.protocol\",\"target\":[\"ipv4_base\",\"protocol\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"local_metadata.icmp_code\",\"target\":[\"scalars\",\"local_metadata_t.icmp_code\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.vlan_tag[0].vid\",\"target\":[\"scalars\",\"key_0\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"hdr.vlan_tag[0].pcp\",\"target\":[\"scalars\",\"key_1\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"local_metadata.class_id\",\"target\":[\"scalars\",\"local_metadata_t.class_id\"],\"mask\":null},{\"match_type\":\"ternary\",\"name\":\"local_metadata.vrf_id\",\"target\":[\"scalars\",\"local_metadata_t.vrf_id\"],\"mask\":null}],\"match_type\":\"ternary\",\"type\":\"simple\",\"max_size\":25,\"with_counters\":true,\"support_timeout\":false,\"direct_meters\":\"ingress.punt.ingress_port_meter\",\"action_ids\":[7,8,9,3],\"actions\":[\"ingress.punt.set_queue_and_clone_to_cpu\",\"ingress.punt.set_queue_and_send_to_cpu\",\"ingress.punt.set_egress_port\",\"NoAction\"],\"base_default_next\":null,\"next_tables\":{\"ingress.punt.set_queue_and_clone_to_cpu\":null,\"ingress.punt.set_queue_and_send_to_cpu\":null,\"ingress.punt.set_egress_port\":null,\"NoAction\":null},\"default_entry\":{\"action_id\":3,\"action_const\":false,\"action_data\":[],\"action_entry_const\":false}}],\"action_profiles\":[{\"name\":\"ingress.l3_fwd.wcmp_action_profile\",\"id\":0,\"source_info\":{\"filename\":\"max.p4\",\"line\":112,\"column\":55,\"source_fragment\":\"wcmp_action_profile\"},\"max_size\":1024,\"selector\":{\"algo\":\"crc16\",\"input\":[{\"type\":\"field\",\"value\":[\"ipv4_base\",\"src_addr\"]},{\"type\":\"field\",\"value\":[\"ipv4_base\",\"protocol\"]},{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l4_src_port\"]},{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l4_dst_port\"]}]}}],\"conditionals\":[{\"name\":\"node_2\",\"id\":0,\"source_info\":{\"filename\":\"max.p4\",\"line\":183,\"column\":8,\"source_fragment\":\"hdr.packet_out.isValid()\"},\"expression\":{\"type\":\"expression\",\"value\":{\"op\":\"d2b\",\"left\":null,\"right\":{\"type\":\"field\",\"value\":[\"packet_out\",\"$valid$\"]}}},\"true_next\":\"tbl_act\",\"false_next\":\"node_4\"},{\"name\":\"node_4\",\"id\":1,\"source_info\":{\"filename\":\"max.p4\",\"line\":187,\"column\":8,\"source_fragment\":\"standard_metadata.egress_spec==0||...\"},\"expression\":{\"type\":\"expression\",\"value\":{\"op\":\"or\",\"left\":{\"type\":\"expression\",\"value\":{\"op\":\"==\",\"left\":{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0x0000\"}}},\"right\":{\"type\":\"expression\",\"value\":{\"op\":\"==\",\"left\":{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_spec\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0x000d\"}}}}},\"true_next\":\"ingress.my_station_table\",\"false_next\":\"tbl_act_0\"},{\"name\":\"node_6\",\"id\":2,\"source_info\":{\"filename\":\"max.p4\",\"line\":191,\"column\":12,\"source_fragment\":\"local_metadata.l3_admit==1w1\"},\"expression\":{\"type\":\"expression\",\"value\":{\"op\":\"==\",\"left\":{\"type\":\"field\",\"value\":[\"scalars\",\"local_metadata_t.l3_admit\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0x01\"}}},\"true_next\":\"ingress.l3_fwd.l3_fwd_table\",\"false_next\":\"ingress.l2_fwd.l2_unicast_table\"}]},{\"name\":\"egress\",\"id\":1,\"source_info\":{\"filename\":\"max.p4\",\"line\":204,\"column\":8,\"source_fragment\":\"egress\"},\"init_table\":\"node_14\",\"tables\":[{\"name\":\"tbl_act_2\",\"id\":7,\"source_info\":{\"filename\":\"max.p4\",\"line\":209,\"column\":12,\"source_fragment\":\"hdr.packet_in.setValid();...\"},\"key\":[],\"match_type\":\"exact\",\"type\":\"simple\",\"max_size\":1024,\"with_counters\":false,\"support_timeout\":false,\"direct_meters\":null,\"action_ids\":[14],\"actions\":[\"act_2\"],\"base_default_next\":null,\"next_tables\":{\"act_2\":null},\"default_entry\":{\"action_id\":14,\"action_const\":true,\"action_data\":[],\"action_entry_const\":true}}],\"action_profiles\":[],\"conditionals\":[{\"name\":\"node_14\",\"id\":3,\"source_info\":{\"filename\":\"max.p4\",\"line\":208,\"column\":12,\"source_fragment\":\"standard_metadata.egress_port==0xFD\"},\"expression\":{\"type\":\"expression\",\"value\":{\"op\":\"==\",\"left\":{\"type\":\"field\",\"value\":[\"standard_metadata\",\"egress_port\"]},\"right\":{\"type\":\"hexstr\",\"value\":\"0x00fd\"}}},\"false_next\":null,\"true_next\":\"tbl_act_2\"}]}],\"checksums\":[{\"name\":\"cksum\",\"id\":0,\"source_info\":{\"filename\":\"ipv4_checksum.p4\",\"line\":29,\"column\":4,\"source_fragment\":\"update_checksum(hdr.ipv4_base.isValid(),...\"},\"target\":[\"ipv4_base\",\"hdr_checksum\"],\"type\":\"generic\",\"calculation\":\"calc\",\"verify\":false,\"update\":true,\"if_cond\":{\"type\":\"expression\",\"value\":{\"op\":\"d2b\",\"left\":null,\"right\":{\"type\":\"field\",\"value\":[\"ipv4_base\",\"$valid$\"]}}}},{\"name\":\"cksum_0\",\"id\":1,\"source_info\":{\"filename\":\"ipv4_checksum.p4\",\"line\":11,\"column\":4,\"source_fragment\":\"verify_checksum(hdr.ipv4_base.isValid(),...\"},\"target\":[\"ipv4_base\",\"hdr_checksum\"],\"type\":\"generic\",\"calculation\":\"calc_0\",\"verify\":true,\"update\":false,\"if_cond\":{\"type\":\"expression\",\"value\":{\"op\":\"d2b\",\"left\":null,\"right\":{\"type\":\"field\",\"value\":[\"ipv4_base\",\"$valid$\"]}}}}],\"force_arith\":[],\"extern_instances\":[],\"field_aliases\":[[\"queueing_metadata.enq_timestamp\",[\"standard_metadata\",\"enq_timestamp\"]],[\"queueing_metadata.enq_qdepth\",[\"standard_metadata\",\"enq_qdepth\"]],[\"queueing_metadata.deq_timedelta\",[\"standard_metadata\",\"deq_timedelta\"]],[\"queueing_metadata.deq_qdepth\",[\"standard_metadata\",\"deq_qdepth\"]],[\"intrinsic_metadata.ingress_global_timestamp\",[\"standard_metadata\",\"ingress_global_timestamp\"]],[\"intrinsic_metadata.egress_global_timestamp\",[\"standard_metadata\",\"egress_global_timestamp\"]],[\"intrinsic_metadata.lf_field_list\",[\"standard_metadata\",\"lf_field_list\"]],[\"intrinsic_metadata.mcast_grp\",[\"standard_metadata\",\"mcast_grp\"]],[\"intrinsic_metadata.resubmit_flag\",[\"standard_metadata\",\"resubmit_flag\"]],[\"intrinsic_metadata.egress_rid\",[\"standard_metadata\",\"egress_rid\"]],[\"intrinsic_metadata.recirculate_flag\",[\"standard_metadata\",\"recirculate_flag\"]],[\"intrinsic_metadata.priority\",[\"standard_metadata\",\"priority\"]]],\"program\":\"max.p4\",\"__meta__\":{\"version\":[2,18],\"compiler\":\"https://github.com/p4lang/p4c\"}}\n"),
			},
		}
		pipelineCfgResp = &v1.SetForwardingPipelineConfigResponse{}
	)
	type args struct {
		target *tg.Target
		req    *v1.SetForwardingPipelineConfigRequest
		res    *v1.SetForwardingPipelineConfigResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Pipeline Config",
			args: args{
				target: TestTarget,
				req:    pipelineCfgReq,
				res:    pipelineCfgResp,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//fmt.Println(proto.MarshalTextString(tt.args.req))
			if got := ProcessP4PipelineConfigOperation(tt.args.target, tt.args.req, tt.args.res); got != tt.want {
				t.Errorf("ProcessP4PipelineConfigOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessP4WriteRequest(t *testing.T) {
	log.Infoln(strings.Repeat("*", 100))
	log.Infoln("Start of TestProcessP4WriteRequest")
	defer log.Infoln("End of TestProcessP4WriteRequest")
	setupTest()
	defer tearDownTest()
	var (
		deviceID           uint64 = 1
		electionID                = &v1.Uint128{High: 1, Low: 5}
		emptyWriteWithID          = &v1.WriteRequest{DeviceId: deviceID, ElectionId: electionID}
		emptyWriteRequest         = &v1.WriteRequest{}
		emptyWriteResponse        = &v1.WriteResponse{}
	)
	type args struct {
		target *tg.Target
		wreq   *v1.WriteRequest
		wres   *v1.WriteResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Write",
			args: args{
				target: TestTarget,
				wreq:   emptyWriteWithID,
				wres:   emptyWriteResponse,
			},
			want: true,
		},
		{
			name: "Empty Write",
			args: args{
				target: TestTarget,
				wreq:   emptyWriteRequest,
				wres:   emptyWriteResponse,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessP4WriteRequest(tt.args.target, tt.args.wreq, tt.args.wres); got != tt.want {
				t.Errorf("ProcessP4WriteRequest() = %v, want %v", got, tt.want)
			}
		})
	}
	//tearDownTest()
}

func TestProcessPacketIOOperation(t *testing.T) {
	log.Infoln(strings.Repeat("*", 100))
	log.Infoln("Start of TestProcessPacketIOOperation")
	defer log.Infoln("End of TestProcessPacketIOOperation")
	setupTest()
	defer tearDownTest()
	var (
		deviceID           uint64 = 1
		electionID                = &v1.Uint128{High: 1, Low: 5}
		emptyWriteResponse        = &v1.WriteResponse{}
		insertWriteReq            = &v1.WriteRequest{
			DeviceId:   deviceID,
			ElectionId: electionID,
			Updates: []*v1.Update{
				&v1.Update{
					Type: v1.Update_INSERT,
					Entity: &v1.Entity{
						Entity: &v1.Entity_TableEntry{
							TableEntry: &v1.TableEntry{
								TableId:  33598026,
								Priority: 10,
								Match: []*v1.FieldMatch{
									&v1.FieldMatch{
										FieldId: 3,
										FieldMatchType: &v1.FieldMatch_Ternary_{
											Ternary: &v1.FieldMatch_Ternary{
												Value: []byte("\010\000"),
												Mask:  []byte("\377\377"),
											},
										},
									},
								},
								Action: &v1.TableAction{
									Type: &v1.TableAction_Action{
										Action: &v1.Action{
											ActionId: 16804007,
											Params: []*v1.Action_Param{
												&v1.Action_Param{
													ParamId: 1,
													Value:   []byte("\004"),
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
		deleteWriteReq = &v1.WriteRequest{
			DeviceId:   deviceID,
			ElectionId: electionID,
			Updates: []*v1.Update{
				&v1.Update{
					Type: v1.Update_DELETE,
					Entity: &v1.Entity{
						Entity: &v1.Entity_TableEntry{
							TableEntry: &v1.TableEntry{
								TableId:  33598026,
								Priority: 10,
								Match: []*v1.FieldMatch{
									&v1.FieldMatch{
										FieldId: 3,
										FieldMatchType: &v1.FieldMatch_Ternary_{
											Ternary: &v1.FieldMatch_Ternary{
												Value: []byte("\010\000"),
												Mask:  []byte("\377\377"),
											},
										},
									},
								},
								Action: &v1.TableAction{
									Type: &v1.TableAction_Action{
										Action: &v1.Action{
											ActionId: 16804007,
											Params: []*v1.Action_Param{
												&v1.Action_Param{
													ParamId: 1,
													Value:   []byte("\004"),
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
		validPO = &v1.PacketOut{
			Payload: []byte("\x3C\xFD\xFE\xA8\xEA\x31\x00\x00\x00\xC0\x1A\x10\x08\x00\x45\x00\x00\x2E\x00\x01\x00\x00\x40\x00\x66\xCB\x0A\x01\x00\x01\x0A\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b"),
			Metadata: []*v1.PacketMetadata{
				&v1.PacketMetadata{
					MetadataId: 1,
					Value:      []byte("\000\000"),
				},
			},
		}

		invalidPO = &v1.PacketOut{
			Payload: []byte("\x3C\xFD\xFE\xA8\xEA\x31\x00\x00\x00\xC0\x1A\x10\x08\x00\x45\x00\x00\x2E\x00\x01\x00\x00\x40\x00\x66\xCB\x0A\x01\x00\x01\x0A\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b"),
			Metadata: []*v1.PacketMetadata{
				&v1.PacketMetadata{
					MetadataId: 1,
					Value:      []byte("\000\007"),
				},
			},
		}
		validPI = &v1.PacketIn{
			Payload: []byte("\x3C\xFD\xFE\xA8\xEA\x31\x00\x00\x00\xC0\x1A\x10\x08\x00\x45\x00\x00\x2E\x00\x01\x00\x00\x40\x00\x66\xCB\x0A\x01\x00\x01\x0A\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b"),
		}
		invalidPI = &v1.PacketIn{
			Payload: []byte("\x3C\xFD\xFE\xA8\xEA\x31\x00\x00\x00\xC0\x1A\x10\x08\x00\x45\x00\x00\x2E\x00\x01\x00\x00\x40\x00\x66\xCB\x0A\x01\x00\x01\x0A\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29"),
		}
	)
	type args struct {
		target         *tg.Target
		po             *v1.PacketOut
		pi             *v1.PacketIn
		insertWriteReq *v1.WriteRequest
		deleteWriteReq *v1.WriteRequest
		writeResponse  *v1.WriteResponse
	}
	tests := []struct {
		name      string
		args      args
		poWant    bool
		piWant    bool
		writeWant bool
	}{
		{
			name: "Valid Packet",
			args: args{
				target:         TestTarget,
				po:             validPO,
				pi:             validPI,
				insertWriteReq: insertWriteReq,
				deleteWriteReq: deleteWriteReq,
				writeResponse:  emptyWriteResponse,
			},
			poWant:    true,
			piWant:    true,
			writeWant: true,
		},
		{
			name: "Invalid Packet In",
			args: args{
				target:         TestTarget,
				po:             validPO,
				pi:             invalidPI,
				insertWriteReq: insertWriteReq,
				deleteWriteReq: deleteWriteReq,
				writeResponse:  emptyWriteResponse,
			},
			poWant:    true,
			piWant:    false,
			writeWant: true,
		},

		{
			name: "Packet In Timeout",
			args: args{
				target:         TestTarget,
				po:             invalidPO,
				pi:             invalidPI,
				insertWriteReq: insertWriteReq,
				deleteWriteReq: deleteWriteReq,
				writeResponse:  emptyWriteResponse,
			},
			poWant:    true,
			piWant:    false,
			writeWant: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessP4WriteRequest(tt.args.target, tt.args.insertWriteReq, tt.args.writeResponse); got != tt.writeWant {
				t.Errorf("Insert Write ProcessP4WriteRequest() = %v, want %v", got, tt.writeWant)
			}
			if got := ProcessPacketOutOperation(tt.args.target, tt.args.po); got != tt.poWant {
				t.Errorf("ProcessPacketOutOperation() = %v, want %v", got, tt.poWant)
			}
			if got := ProcessPacketIn(tt.args.pi); got != tt.piWant {
				t.Errorf("ProcessPacketIn() = %v, want %v", got, tt.piWant)
			}
			if got := ProcessP4WriteRequest(tt.args.target, tt.args.deleteWriteReq, tt.args.writeResponse); got != tt.writeWant {
				t.Errorf("Delete Write ProcessP4WriteRequest() = %v, want %v", got, tt.writeWant)
			}
		})
	}
}

func TestMasterArbitration(t *testing.T) {
	log.Infoln(strings.Repeat("*", 100))
	log.Infoln("Start of TestMasterArbitration")
	setupTest()
	var (
		deviceID        uint64 = 1
		invalidDeviceID uint64 = 2
		electionID             = &v1.Uint128{High: 1, Low: 5}
		highElectionID         = &v1.Uint128{High: 2, Low: 5}
		scv                    = GetStreamChannel(P4rtClient)
	)
	type args struct {
		scv        StreamChannelVar
		deviceID   uint64
		electionID *v1.Uint128
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid master arbitration",
			args: args{
				scv:        scv,
				deviceID:   deviceID,
				electionID: electionID,
			},
			want: true,
		},
		{
			name: "Higher election",
			args: args{
				scv:        scv,
				deviceID:   deviceID,
				electionID: highElectionID,
			},
			want: true,
		},
		{
			name: "Invalid device ID",
			args: args{
				scv:        scv,
				deviceID:   invalidDeviceID,
				electionID: electionID,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMasterArbitrationLock(tt.args.scv, tt.args.deviceID, tt.args.electionID); got != tt.want {
				t.Errorf("GetMasterArbitrationLock() = %v, want %v", got, tt.want)
			}
		})
	}
	tearDownTest()
	log.Infoln("End of TestMasterArbitration")
}

func TestLowElectionMasterArbitration(t *testing.T) {
	log.Infoln(strings.Repeat("*", 100))
	log.Infoln("Start of TestLowElectionMasterArbitration")
	setupTest()
	var (
		deviceID      uint64 = 1
		electionID           = &v1.Uint128{High: 1, Low: 5}
		lowElectionID        = &v1.Uint128{High: 0, Low: 5}
		scv1                 = GetStreamChannel(P4rtClient)
		scv2                 = GetStreamChannel(P4rtClient)
	)
	type args struct {
		scv1          StreamChannelVar
		scv2          StreamChannelVar
		deviceID      uint64
		electionID    *v1.Uint128
		lowElectionID *v1.Uint128
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Low Election",
			args: args{scv1: scv1, scv2: scv2, deviceID: deviceID, electionID: electionID, lowElectionID: lowElectionID},
		},
	}
	Init(TestTarget)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetMasterArbitrationLock(tt.args.scv1, tt.args.deviceID, tt.args.electionID)
			if got := GetMasterArbitrationLock(tt.args.scv2, tt.args.deviceID, tt.args.lowElectionID); got != tt.want {
				t.Errorf("GetMasterArbitrationLock() = %v, want %v", got, tt.want)
			}
		})
	}
	tearDownTest()
	log.Infoln("End of TestLowElectionMasterArbitration")
}
