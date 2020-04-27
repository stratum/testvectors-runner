/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package tests

import (
	"testing"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	"github.com/golang/protobuf/proto"

	"github.com/stratum/testvectors-runner/pkg/framework/dataplane"
	"github.com/stratum/testvectors-runner/pkg/framework/p4rt"
	"github.com/stratum/testvectors-runner/pkg/test/setup"
	"github.com/stratum/testvectors-runner/pkg/test/teardown"
	"github.com/stretchr/testify/assert"
)

var (
	payload       = "\x3c\xfd\xfe\xa8\xea\x31\x00\x00\x00\xc0\x1a\x10\x08\x00\x45\x00\x00\x2e\x00\x01\x00\x00\x40\x00\x66\xcb\x0a\x01\x00\x01\x0a\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b"
	pktOutToPort1 = `
payload: "\x3c\xfd\xfe\xa8\xea\x31\x00\x00\x00\xc0\x1a\x10\x08\x00\x45\x00\x00\x2e\x00\x01\x00\x00\x40\x00\x66\xcb\x0a\x01\x00\x01\x0a\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b"
metadata: <
  metadata_id: 1
  value: "\000\001"
>`
	pktOutToPort0 = `
payload: "\x3c\xfd\xfe\xa8\xea\x31\x00\x00\x00\xc0\x1a\x10\x08\x00\x45\x00\x00\x2e\x00\x01\x00\x00\x40\x00\x66\xcb\x0a\x01\x00\x01\x0a\x02\x00\x01\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b"
metadata: <
  metadata_id: 1
  value: "\000\000"
>`
	writeRequest = `
device_id: 1
election_id: <
  low: 4
>
updates: <
  type: INSERT
  entity: <
    table_entry: <
      table_id: 33598026
      priority: 10
      match: <
        field_id: 3
        ternary: <
          value: "\010\000"
          mask: "\377\377"
        >
      >
      action: <
        action: <
          action_id: 16820507
          params: <
            param_id: 1
            value: "\000\002"
          >
        >
      >
    >
  >
>`
	deleteWriteRequest = `
device_id: 1
election_id: <
  low: 4
>
updates: <
  type: DELETE
  entity: <
    table_entry: <
      table_id: 33598026
      priority: 10
      match: <
        field_id: 3
        ternary: <
          value: "\010\000"
          mask: "\377\377"
        >
      >
      action: <
        action: <
          action_id: 16820507
          params: <
            param_id: 1
            value: "\000\002"
          >
        >
      >
    >
  >
>`
)

// PktIoOutDirectToDataPlaneTest sends packets directly out of a physical port. It Skips the ingress pipeline and any processing.
func (st Test) PktIoOutDirectToDataPlaneTest(t *testing.T) {
	// Start packet capturing
	setup.TestCase()

	// Build packet-out
	pktOut := &v1.PacketOut{}
	if err := proto.UnmarshalText(pktOutToPort1, pktOut); err != nil {
		log.Fatalf("Error parsing proto message of type %T\n%s", pktOut, err)
	}
	// Send packet-out
	result := p4rt.ProcessPacketOutOperation(pktOut)
	assert.True(t, result, "PacketOut operation failed")

	// Check if we received packets from data plane port 1
	result = dataplane.ProcessTrafficExpectation([][]byte{[]byte(payload)}, []uint32{1})
	assert.True(t, result, "Packet not received on port 1")
	// Check if we received no packets from data plane port 2
	result = dataplane.ProcessTrafficExpectation([][]byte{}, []uint32{2})
	assert.True(t, result, "Unexpected packet received on port 2")

	// Stop packet capturing
	teardown.TestCase()
}

// PktIoOutToIngressPipelineACLRedirectToPortTest sends packets out through the ingress pipeline and redirect it to a port via an ACL rule.
func (st Test) PktIoOutToIngressPipelineACLRedirectToPortTest(t *testing.T) {
	// Start packet capturing
	setup.TestCase()

	// Build write request
	request := &v1.WriteRequest{}
	if err := proto.UnmarshalText(writeRequest, request); err != nil {
		log.Fatalf("Error parsing proto message of type %T\n%s", request, err)
	}

	// Insert table entry
	result := p4rt.ProcessP4WriteRequest(request, nil)
	assert.True(t, result, "Write request failed")

	// Build packet-out
	pktOut := &v1.PacketOut{}
	if err := proto.UnmarshalText(pktOutToPort0, pktOut); err != nil {
		log.Fatalf("Error parsing proto message of type %T\n%s", pktOut, err)
	}
	// Send packet-out
	result = p4rt.ProcessPacketOutOperation(pktOut)
	assert.True(t, result, "PacketOut operation failed")

	// Check if we received packets from data plane port 2
	result = dataplane.ProcessTrafficExpectation([][]byte{[]byte(payload)}, []uint32{2})
	assert.True(t, result, "Packet not received on port 1")
	// Check if we received no packets from data plane port 1
	result = dataplane.ProcessTrafficExpectation([][]byte{}, []uint32{1})
	assert.True(t, result, "Unexpected packet received on port 2")

	// Build delete write request
	request = &v1.WriteRequest{}
	if err := proto.UnmarshalText(deleteWriteRequest, request); err != nil {
		log.Fatalf("Error parsing proto message of type %T\n%s", request, err)
	}

	// Delete table entry
	result = p4rt.ProcessP4WriteRequest(request, nil)
	assert.True(t, result, "Write request failed")
	// Stop packet capturing
	teardown.TestCase()
}
