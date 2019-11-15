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

// connect opens a new gRPC connection to the target speciifed by the
// ConnectionArgs. It returns the p4runtime Client connection, and a function
// which can be called to close the connection. If an error is encountered
// during opening the connection, it is returned.
func connect(tg *tvb.Target) connection {
	if tg.Address == "" {
		return connection{connError: errors.New("an address must be specified")}
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, tg.Address, grpc.WithInsecure())
	if err != nil {
		return connection{connError: fmt.Errorf("cannot dial target %s, %v", tg.Address, err)}
	}
	return connection{ctx: ctx, client: v1.NewP4RuntimeClient(conn), cancel: func() { conn.Close() }}
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
		log.Errorln(scv.scError)
		log.Fatalln("Unable to get a stream channel")
	}
	go receiveStreamChannel(scv.sc, scv.pktInChan, scv.masterArbRecvChan, scv.genericStreamMessageChannel)
	go sendStreamChannel(scv.sc, scv.pktOutChan, scv.masterArbSendChan)
	return scv
}

//GetMasterArbitrationLock description
func getMasterArbitrationLock(scv streamChannel, deviceID uint64, electionID *v1.Uint128) bool {
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
