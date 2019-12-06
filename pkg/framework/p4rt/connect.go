/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package p4rt

import (
	"context"
	"errors"
	"fmt"
	"time"

	scpb "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"

	tvb "github.com/stratum/testvectors/proto/target"
)

//CtxTimeout for contexts
const CtxTimeout = 3 * time.Second

//connect starts a gRPC connection to the target specified.
//It returns connection struct with P4Runtime client, close function
//If an error is encountered during opening the connection, it is returned.
func connect(tg *tvb.Target) connection {
	log.Debug("In p4_oper connect")
	if tg.Address == "" {
		return connection{connError: errors.New("an address must be specified")}
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, tg.Address, grpc.WithInsecure())
	if err != nil {
		return connection{connError: fmt.Errorf("cannot dial target %s, %v", tg.Address, err)}
	}
	return connection{ctx: ctx, client: v1.NewP4RuntimeClient(conn), cancel: func() { conn.Close() }}
}

//recvStreamChannel runs a loop to continuously monitor stream channel client and sorts received messages to appropriate channels
func recvStreamChannel(sc v1.P4Runtime_StreamChannelClient, pktInChan chan *v1.PacketIn, masterArbitrationRecvChan chan *v1.MasterArbitrationUpdate, genericStreamMessageChannel chan *v1.StreamMessageResponse) {
	for {
		if sc == nil {
			log.Debugf("Stream channel is nil or closed")
			return
		}
		smr, err := sc.Recv()
		if err != nil {
			log.Debugf("Failed to receive a message : %v\n", err)
			return
		}

		switch {
		case smr == nil:
			log.Debug("Empty message received")
		case smr.GetPacket() != nil:
			log.Debug("Packet Received")
			pktInChan <- smr.GetPacket()
		case smr.GetArbitration() != nil:
			log.Debug("Arbitration lock")
			masterArbitrationRecvChan <- smr.GetArbitration()
		default:
			genericStreamMessageChannel <- smr
			log.Debug("In Process packet in else block")
			log.Debugf("%T\n", smr)
			log.Debug(smr)
		}
	}
}

//sendStreamChannel runs a loop to continuously monitor pktOut and masterArbitrationReq channels and send messages to stream channel client
func sendStreamChannel(sc v1.P4Runtime_StreamChannelClient, pktOutChan chan *v1.PacketOut, masterArbitrationSendChan chan *v1.MasterArbitrationUpdate) {
	for {
		select {
		case pktOut := <-pktOutChan:
			log.Debug("In Send Stream Packet Out")
			smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Packet{Packet: pktOut}}
			sendErr := sc.Send(smr)
			if sendErr != nil {
				log.Errorf("send err:%s\n", sendErr)
			}
			log.Debug("sent packet")
		case masterArbitrationReq := <-masterArbitrationSendChan:
			log.Debug("In Send Stream Master Arbitration")
			smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Arbitration{Arbitration: masterArbitrationReq}}
			sendErr := sc.Send(smr)
			if sendErr != nil {
				log.Debugf("send err:%s\n", sendErr)
			}
		}
	}
}

//GetStreamChannel description
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
	go recvStreamChannel(scv.sc, scv.pktInChan, scv.masterArbRecvChan, scv.genericStreamMessageChannel)
	go sendStreamChannel(scv.sc, scv.pktOutChan, scv.masterArbSendChan)
	return scv
}

//GetMasterArbitrationLock description
func getMasterArbitrationLock(scv streamChannel, deviceID uint64, electionID *v1.Uint128) bool {
	lockAchieved := false

	arb := &v1.MasterArbitrationUpdate{}
	arb.DeviceId = deviceID
	arb.ElectionId = electionID
	scv.masterArbSendChan <- arb
	select {
	case ret := <-scv.masterArbRecvChan:
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
