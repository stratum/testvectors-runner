package framework

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/testutil"

	"github.com/opennetworkinglab/testvectors-runner/pkg/common"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/opennetworkinglab/testvectors/proto/target"
)

var log = logger.NewLogger()

var (
	gnmiContext context.Context
	gnmiClient  gnmi.GNMIClient
	gnmiError   error
	gnmiCancel  context.CancelFunc
)

//InitGNMI starts a gNMI client connection to switch under test
func InitGNMI(target *tg.Target) {
	gnmiContext = context.Background()
	gnmiClient, gnmiCancel, gnmiError = common.Connect(gnmiContext, target)
	if gnmiError != nil {
		log.Errorln(gnmiError)
		log.Fatalln("Unable to get a gnmi client")
	}
}

//TearDownGNMI closes the gNMI connection
func TearDownGNMI() {
	log.Traceln("In gnmi_oper tear down")
	gnmiCancel()
}

//ProcessGetRequest sends a request to SUT and gives the response
func ProcessGetRequest(target *tg.Target, greq *gnmi.GetRequest, gresp *gnmi.GetResponse) bool {
	log.Infoln("Sending get request")
	resp, err := gnmiClient.Get(gnmiContext, greq)
	if err != nil {
		log.Errorln(err)
		return false
	}
	isEqual := testutil.GetResponseEqual(resp, gresp, testutil.IgnoreTimestamp{})
	if !isEqual {
		log.Warningf("Get responses are unequal\nExpected: %s\nActual  : %s\n", gresp, resp)
	} else {
		log.Infoln("Get responses are equal")
		log.Debugf("Get response: %s\n", resp)
	}
	return isEqual
}

//ProcessSetRequest sends a set request to SUT and gives the response
func ProcessSetRequest(target *tg.Target, sreq *gnmi.SetRequest, sresp *gnmi.SetResponse) bool {
	log.Traceln("In ProcessSetRequest")
	log.Infoln("Sending set request")
	resp, err := gnmiClient.Set(gnmiContext, sreq)
	if err != nil {
		log.Errorln(err)
		return false
	}
	//FIXME
	//resetting timestamp as a work around to ignore timestamp during comparison
	origRespTimestamp := resp.Timestamp
	origSRespTimestamp := sresp.Timestamp
	resp.Timestamp = 0
	sresp.Timestamp = 0
	isEqual := proto.Equal(resp, sresp)
	if !isEqual {
		log.Warningf("Set responses are unequal\nexpected: %s\nactual: %s\n", sresp, resp)
	} else {
		log.Infoln("Set responses are equal")
		log.Debugf("Set response: %s\n", resp)
	}

	//Reset timestamp to original value
	resp.Timestamp = origRespTimestamp
	sresp.Timestamp = origSRespTimestamp
	return isEqual
}

//ProcessSubscribeRequest opens a subscription channel and verifies the response
func ProcessSubscribeRequest(target *tg.Target, sreq *gnmi.SubscribeRequest, sresp []*gnmi.SubscribeResponse, resultChan chan bool) bool {

	subcl, err := gnmiClient.Subscribe(gnmiContext)
	if err != nil {
		log.Infoln(err)
		return false
	}
	defer subcl.CloseSend()
	result := true
	waitc := make(chan struct{})
	log.Tracef("Length of expected result: %d\n\n", len(sresp))
	go func() {
		for _, exp := range sresp {
			in, err := subcl.Recv()
			if err != nil {
				log.Fatalf("Failed to receive a subscription response : %v", err)
				result = result && false
			}
			if in.GetUpdate() != nil && testutil.NotificationSetEqual([]*gnmi.Notification{exp.GetUpdate()}, []*gnmi.Notification{in.GetUpdate()}, testutil.IgnoreTimestamp{}) {
				log.Traceln("In GetUpdate condition")
				log.Infoln("Subscription responses are equal")
				log.Debugf("Subscription response: %s\n", in)
				result = result && true
			} else if testutil.SubscribeResponseEqual(exp, in) {
				//continue
				log.Infoln("Subscription responses are equal")
				log.Debugf("Subscription response: %s\n", in)
				result = result && true
			} else {
				log.Warnf("Subscription responses are unequal expected:\n%s \nactual:\n%s", exp, in)
				result = result && false
			}
		}
		close(waitc)
		return
	}()
	log.Infoln("Sending subscription request")
	err = subcl.Send(sreq)
	if err != nil {
		log.Errorln(err)
		result = false
	}
	<-waitc
	resultChan <- result
	return true
}