/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package p4rt

import (
	"bytes"
	"context"
	"time"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	scpb "google.golang.org/genproto/googleapis/rpc/code"

	"github.com/opennetworkinglab/testvectors-runner/pkg/common"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()

var (
	//P4rtClient description
	P4rtClient  v1.P4RuntimeClient
	p4rtContext context.Context
	lock        bool
	//SCV Description
	SCV StreamChannelVar
)

//StreamChannelVar definition
type StreamChannelVar struct {
	sc                                   v1.P4Runtime_StreamChannelClient
	scError                              error
	scCancel                             context.CancelFunc
	masterArbRecvChan, masterArbSendChan chan *v1.MasterArbitrationUpdate
	pktInChan                            chan *v1.PacketIn
	pktOutChan                           chan *v1.PacketOut
	genericStreamMessageChannel          chan *v1.StreamMessageResponse
}

//GetStreamChannel description
func GetStreamChannel(p4rtClient v1.P4RuntimeClient) StreamChannelVar {
	scv := StreamChannelVar{}
	scv.masterArbRecvChan = make(chan *v1.MasterArbitrationUpdate)
	scv.masterArbSendChan = make(chan *v1.MasterArbitrationUpdate)
	scv.pktInChan = make(chan *v1.PacketIn)
	scv.pktOutChan = make(chan *v1.PacketOut)
	scv.genericStreamMessageChannel = make(chan *v1.StreamMessageResponse)
	scContext := context.Background()
	//ctx = context.TODO()
	scContext, scv.scCancel = context.WithCancel(scContext)
	scv.sc, scv.scError = p4rtClient.StreamChannel(scContext)
	if scv.scError != nil {
		log.Errorln(scv.scError)
		log.Fatalln("Unable to get a stream channel")
	}
	go receiveStreamChannel(scv.sc, scv.pktInChan, scv.masterArbRecvChan, scv.genericStreamMessageChannel)
	go sendStreamChannel(scv.sc, scv.pktOutChan, scv.masterArbSendChan)
	return scv
}

//Init starts a P4Runtime client and runs go routines to send and receive stream channel messages from P4Runtime stream channel client
func Init(target *tg.Target) {
	p4rtContext = context.Background()
	P4rtClient, _, _ = common.ConnectP4(p4rtContext, target)
	SCV = GetStreamChannel(P4rtClient)
}

//receiveStreamChannel runs a loop to continuously monitor stream channel client and sorts received messages to appropriate channels
func receiveStreamChannel(sc v1.P4Runtime_StreamChannelClient, pktInChan chan *v1.PacketIn, masterArbitrationRecvChan chan *v1.MasterArbitrationUpdate, genericStreamMessageChannel chan *v1.StreamMessageResponse) {
	for {
		smr, err := sc.Recv()
		if err != nil {
			log.Tracef("Failed to receive a message : %v\n", err)
			//close(waitc)
			return
		}

		switch {
		case smr == nil:
			log.Traceln("Empty message received")
		case smr.GetPacket() != nil:
			log.Traceln("Packet Received")
			pktInChan <- smr.GetPacket()
		case smr.GetArbitration() != nil:
			log.Traceln("Arbitration lock")
			masterArbitrationRecvChan <- smr.GetArbitration()
		default:
			genericStreamMessageChannel <- smr
			log.Traceln("In Process packet in else block")
			log.Tracef("%T\n", smr)
			log.Traceln(smr)
		}
	}
}

//sendStreamChannel runs a loop to continuously monitor pktOut and masterArbitrationReq channels and send messages to stream channel client
func sendStreamChannel(sc v1.P4Runtime_StreamChannelClient, pktOutChan chan *v1.PacketOut, masterArbitrationSendChan chan *v1.MasterArbitrationUpdate) {
	for {
		select {
		case pktOut := <-pktOutChan:
			log.Traceln("In Send Stream Packet Out")
			smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Packet{Packet: pktOut}}
			sendErr := sc.Send(smr)
			if sendErr != nil {
				log.Errorf("send err:%s\n", sendErr)
			}
			log.Traceln("sent packet")
		case masterArbitrationReq := <-masterArbitrationSendChan:
			log.Traceln("In Send Stream Master Arbitration")
			smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Arbitration{Arbitration: masterArbitrationReq}}
			sendErr := sc.Send(smr)
			if sendErr != nil {
				log.Tracef("send err:%s\n", sendErr)
			}
		}
	}
}

//TearDown closes the stream channel client
func TearDown() {
	log.Traceln("In p4_oper tear down")
	SCV.scCancel()
	if SCV.sc != nil {
		err := SCV.sc.CloseSend()
		if err != nil {
			log.Warnln("Error closing the stream channel:", err)
		}
	}
}

//ProcessP4WriteRequest processes the write request
func ProcessP4WriteRequest(wreq *v1.WriteRequest, wres *v1.WriteResponse) bool {
	log.Traceln("In ProcessP4WriteRequest")
	if wreq == nil {
		return false
	}

	lock = GetMasterArbitrationLock(SCV, wreq.DeviceId, wreq.ElectionId)

	if lock {
		log.Infoln("Sending P4 write request")
		log.Debugf("Write request: %s", wreq)
		resp, err := P4rtClient.Write(p4rtContext, wreq)
		if err != nil {
			log.Errorf("err:%s\n", err)
			return false
		}
		log.Debugf("resp:%s\n", resp)
		return true
	}
	return false
}

//ProcessP4PipelineConfigOperation processes the forwarding request.
func ProcessP4PipelineConfigOperation(req *v1.SetForwardingPipelineConfigRequest, res *v1.SetForwardingPipelineConfigResponse) bool {
	log.Traceln("In ProcessP4PipelineConfigOperation")
	if req == nil {
		return false
	}
	lock = GetMasterArbitrationLock(SCV, req.DeviceId, req.ElectionId)
	if lock {
		log.Infoln("Sending P4 pipeline config")
		log.Tracef("Pipeline config: %s", req)
		resp, err := P4rtClient.SetForwardingPipelineConfig(p4rtContext, req)
		if err != nil {
			log.Errorf("err:%s\n", err)
			return false
		}
		log.Debugf("resp:%s\n", resp)
		return true
	}
	return false
}

//ProcessPacketOutOperation sends packet to stream channel client.
func ProcessPacketOutOperation(po *v1.PacketOut) bool {
	log.Traceln("In ProcessP4 Packet Out")
	var deviceID uint64 = 1
	electionID := &v1.Uint128{High: 1, Low: 5}
	lock = GetMasterArbitrationLock(SCV, deviceID, electionID)
	if lock {
		log.Infoln("Sending packet")
		log.Debugf("Packet info: %s", po)
		SCV.pktOutChan <- po
		return true
	}
	return false
}

//ProcessPacketIn verifies if the packet received is same as expected packet.
func ProcessPacketIn(exp *v1.PacketIn) bool {
	packetMatched := false

	select {
	case ret := <-SCV.pktInChan:
		log.Traceln("In ProcessPacketIn Case PktInChan")
		if bytes.Equal(ret.GetPayload(), exp.GetPayload()) {
			packetMatched = true
			log.Infof("Received packet matches")
			log.Debugf("Packet info: %s", ret)
		} else {
			log.Warningf("Packets don't match\nExpected: % x\nActual  : % x\n", exp.GetPayload(), ret.GetPayload())
		}
		return packetMatched
	case <-time.After(3 * time.Second):
		log.Errorln("Timed out")
	}

	return packetMatched
}

//GetMasterArbitrationLock description
func GetMasterArbitrationLock(scv StreamChannelVar, deviceID uint64, electionID *v1.Uint128) bool {
	log.Traceln("In GetMasterArbitrationLock")
	lockAchieved := false

	arb := &v1.MasterArbitrationUpdate{}
	arb.DeviceId = deviceID
	arb.ElectionId = electionID
	scv.masterArbSendChan <- arb
	select {
	case ret := <-scv.masterArbRecvChan:
		if ret.Status.Code == int32(scpb.Code_OK) {
			log.Traceln("Master lock achieved")
			lockAchieved = true
		} else {
			log.Infoln("Master lock not achieved")
			log.Errorln(ret.Status)
		}
	case <-time.After(3 * time.Second):
		log.Errorln("Timed out")
	}
	return lockAchieved
}
