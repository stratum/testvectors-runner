/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package p4rt

import (
	"context"
	"time"

	scpb "google.golang.org/genproto/googleapis/rpc/code"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
)

//CtxTimeout for contexts
const CtxTimeout = 3 * time.Second

//streamChannel struct stores stream channel client, cancel function and channels to receive stream messages
type streamChannel struct {
	sc                                   v1.P4Runtime_StreamChannelClient
	scError                              error
	cancel                               context.CancelFunc
	masterArbRecvChan, masterArbSendChan chan *v1.MasterArbitrationUpdate
	pktInChan                            chan *v1.PacketIn
	pktOutChan                           chan *v1.PacketOut
	genericStreamMessageChannel          chan *v1.StreamMessageResponse
}

//Close P4Runtime_StreamChannelClient
func (s streamChannel) Close() {
	s.cancel()
	if s.sc != nil {
		err := s.sc.CloseSend()
		if err != nil {
			log.Warn("Error closing the stream channel:", err)
		}
	}
}

//GetStreamChannel gets a new P4Runtime stream channel client, starts Recv() and Send() goroutines
func getStreamChannel(p4rtClient v1.P4RuntimeClient) streamChannel {
	scv := streamChannel{}
	scv.masterArbRecvChan = make(chan *v1.MasterArbitrationUpdate)
	scv.masterArbSendChan = make(chan *v1.MasterArbitrationUpdate)
	scv.pktInChan = make(chan *v1.PacketIn)
	scv.pktOutChan = make(chan *v1.PacketOut)
	scv.genericStreamMessageChannel = make(chan *v1.StreamMessageResponse)
	scContext := context.Background()
	scContext, scv.cancel = context.WithCancel(scContext)
	scv.sc, scv.scError = p4rtClient.StreamChannel(scContext)
	if scv.scError != nil {
		log.Error(scv.scError)
		log.Fatal("Unable to get a stream channel")
	}
	go scv.Recv()
	go scv.Send()
	return scv
}

//Recv runs a loop to continuously monitor stream channel client and sorts received messages to appropriate channels
//This method is called as go routine.
func (s streamChannel) Recv() {
	for {
		if s.sc == nil {
			log.Debugf("Stream channel is nil or closed")
			return
		}
		smr, err := s.sc.Recv()
		if err != nil {
			log.Debugf("Failed to receive a message : %v\n", err)
			return
		}

		switch {
		case smr == nil:
			log.Debug("Empty message received")
		case smr.GetPacket() != nil:
			log.Debug("Packet Received")
			s.pktInChan <- smr.GetPacket()
		case smr.GetArbitration() != nil:
			log.Debug("Arbitration lock")
			s.masterArbRecvChan <- smr.GetArbitration()
		default:
			s.genericStreamMessageChannel <- smr
			log.Debug("In Process packet in else block")
			log.Debugf("%T\n", smr)
			log.Debug(smr)
		}
	}
}

//Send runs a loop to continuously monitor pktOut and masterArbitrationReq channels and send messages to stream channel client
//This method is called as go routine.
func (s streamChannel) Send() {
	for {
		select {
		case pktOut := <-s.pktOutChan:
			log.Debug("In Send Stream Packet Out")
			smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Packet{Packet: pktOut}}
			sendErr := s.sc.Send(smr)
			if sendErr != nil {
				log.Errorf("send err:%s\n", sendErr)
			}
			log.Debug("sent packet")
		case masterArbitrationReq := <-s.masterArbSendChan:
			log.Debug("In Send Stream Master Arbitration")
			smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Arbitration{Arbitration: masterArbitrationReq}}
			sendErr := s.sc.Send(smr)
			if sendErr != nil {
				log.Debugf("send err:%s\n", sendErr)
			}
		}
	}
}

//getMasterArbitrationLock sends master arbitration request to stream channel with provided deviceID and electionID
//returns true if master lock is achieved, false in case of error or timeout
func (s streamChannel) getMasterArbitrationLock(deviceID uint64, electionID *v1.Uint128) bool {
	lockAchieved := false

	arb := &v1.MasterArbitrationUpdate{}
	arb.DeviceId = deviceID
	arb.ElectionId = electionID
	s.masterArbSendChan <- arb
	select {
	case ret := <-s.masterArbRecvChan:
		if ret.Status.Code == int32(scpb.Code_OK) {
			log.Debug("Master lock achieved")
			lockAchieved = true
		} else {
			log.Errorf("Error getting master lock: %v", ret.Status)
		}
	case <-time.After(CtxTimeout):
		log.Error("Timed out waiting for master arbitration response")
	}
	return lockAchieved
}
