package framework

import (
	"bytes"
	"context"
	"time"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	scpb "google.golang.org/genproto/googleapis/rpc/code"

	"github.com/opennetworkinglab/testvectors-runner/pkg/common"
	tg "github.com/opennetworkinglab/testvectors/proto/target"
)

var (
	p4rtClient                                           v1.P4RuntimeClient
	p4rtContext, scContext                               context.Context
	sc                                                   v1.P4Runtime_StreamChannelClient
	p4rtError                                            error
	lock                                                 bool
	p4rtCancel                                           context.CancelFunc
	pktInChan                                            = make(chan *v1.PacketIn)
	pktOutChan                                           = make(chan *v1.PacketOut)
	masterArbitrationRecvChan, masterArbitrationSendChan = make(chan *v1.MasterArbitrationUpdate), make(chan *v1.MasterArbitrationUpdate)
	genericStreamMessageChannel                          = make(chan *v1.StreamMessageResponse)
)

//Init starts a P4Runtime client and runs go routines to send and receive stream channel messages from P4Runtime stream channel client
func Init(target *tg.Target) {
	p4rtContext = context.Background()
	p4rtClient, _, _ = common.ConnectP4(p4rtContext, target)
	scContext = context.Background()
	//ctx = context.TODO()
	scContext, p4rtCancel = context.WithCancel(scContext)
	sc, p4rtError = p4rtClient.StreamChannel(scContext)
	if p4rtError != nil {
		log.Errorln(p4rtError)
		log.Fatalln("Unable to get a stream channel")
	}
	go receiveStreamChannel(sc, pktInChan, masterArbitrationRecvChan)
	go sendStreamChannel(sc, pktOutChan, masterArbitrationSendChan)

}

//receiveStreamChannel runs a loop to continuously monitor stream channel client and sorts received messages to appropriate channels
func receiveStreamChannel(sc v1.P4Runtime_StreamChannelClient, pktInChan chan *v1.PacketIn, masterArbitrationRecvChan chan *v1.MasterArbitrationUpdate) {
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
	p4rtCancel()
	if sc != nil {
		sc.CloseSend()
	}
}

//ProcessP4WriteRequest processes the write request
func ProcessP4WriteRequest(target *tg.Target, wreq *v1.WriteRequest, wres *v1.WriteResponse) bool {
	log.Traceln("In ProcessP4WriteRequest")
	//FIXME: How to obtain master lock with different election id
	if !lock {
		lock = GetMasterArbitrationLock(sc, wreq.DeviceId, wreq.ElectionId)
	}

	if lock {
		log.Infoln("Sending P4 write request")
		log.Debugf("Write request: %s", wreq)
		resp, err := p4rtClient.Write(p4rtContext, wreq)
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
func ProcessP4PipelineConfigOperation(target *tg.Target, req *v1.SetForwardingPipelineConfigRequest, res *v1.SetForwardingPipelineConfigResponse) bool {
	log.Traceln("In ProcessP4PipelineConfigOperation")
	//FIXME: How to obtain master lock with different election id
	if !lock {
		lock = GetMasterArbitrationLock(sc, req.DeviceId, req.ElectionId)
	}
	if lock {
		log.Infoln("Sending P4 pipeline config")
		log.Tracef("Pipeline config: %s", req)
		resp, err := p4rtClient.SetForwardingPipelineConfig(p4rtContext, req)
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
func ProcessPacketOutOperation(target *tg.Target, po *v1.PacketOut) bool {
	log.Traceln("In ProcessP4 Packet Out")
	//smr := &v1.StreamMessageRequest{Update: &v1.StreamMessageRequest_Packet{Packet: po}}

	var deviceID uint64 = 1
	electionID := &v1.Uint128{High: 1, Low: 5}
	//FIXME: How to obtain master lock with different election id
	if !lock {
		lock = GetMasterArbitrationLock(sc, deviceID, electionID)
	}
	if lock {
		log.Infoln("Sending packet")
		log.Debugf("Packet info: %s", po)
		pktOutChan <- po
		return true
	}
	return false
}

//ProcessPacketIn verifies if the packet received is same as expected packet.
func ProcessPacketIn(exp *v1.PacketIn) bool {
	packetMatched := false

	select {
	case ret := <-pktInChan:
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

//GetMasterArbitrationLock sends a master arbitration request to stream channel client. Returns true if returned status code is Code_OK.
func GetMasterArbitrationLock(sc v1.P4Runtime_StreamChannelClient, deviceID uint64, electionID *v1.Uint128) bool {
	log.Traceln("In GetMasterArbitrationLock")
	lockAchieved := false

	arb := &v1.MasterArbitrationUpdate{}
	arb.DeviceId = deviceID
	arb.ElectionId = electionID
	masterArbitrationSendChan <- arb
	select {
	case ret := <-masterArbitrationRecvChan:
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
