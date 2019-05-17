//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package bootstrap

import (
	"context"
	"encoding/gob"

	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

type ChallengeResponseController interface {
	component.Initer

	Execute(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) (*ChallengePayload, error)
}

type challengeResponseController struct {
	SessionManager SessionManager              `inject:""`
	Cryptography   insolar.CryptographyService `inject:""`
	NodeKeeper     network.NodeKeeper          `inject:""`
	Network        network.HostNetwork         `inject:""`

	options *common.Options
}

type Nonce []byte
type SignedNonce []byte

// NetworkNode                           Discovery NetworkNode
//  1| ------ ChallengeRequest -----> |
//  2| <-- SignedChallengeResponse -- |
//  3| --- SignedChallengeRequest --> |
//  4| <----- ChallengeResponse ----- |
// ------------------------------------

type ChallengeResponseHeader struct {
	Success bool
	Error   string
}

type ChallengeRequest struct {
	SessionID SessionID

	Nonce Nonce
}

type SignedChallengeResponse struct {
	Header  ChallengeResponseHeader
	Payload *SignedChallengePayload
}

type SignedChallengePayload struct {
	SignedNonce       SignedNonce
	XorDiscoveryNonce Nonce
	DiscoveryNonce    Nonce
}

type SignedChallengeRequest struct {
	SessionID SessionID

	SignedDiscoveryNonce SignedNonce
	XorNonce             Nonce
}

type ChallengeResponse struct {
	Header  ChallengeResponseHeader
	Payload *ChallengePayload
}

type ChallengePayload struct {
	// CurrentPulse  insolar.Pulse
	// State         insolar.NetworkState
	AssignShortID insolar.ShortNodeID
}

func init() {
	gob.Register(&ChallengeRequest{})
	gob.Register(&SignedChallengeResponse{})
	gob.Register(&SignedChallengeRequest{})
	gob.Register(&ChallengeResponse{})
}

func (cr *challengeResponseController) processChallenge1(ctx context.Context, request network.Request) (network.Response, error) {
	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.processChallenge1")
	defer span.End()
	data := request.GetData().(*ChallengeRequest)
	xorNonce, err := GenerateNonce()
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, "error generating discovery xor nonce: "+err.Error()), nil
	}
	sign, err := cr.Cryptography.Sign(Xor(data.Nonce, xorNonce))
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, "error signing nonce: "+err.Error()), nil
	}
	discoveryNonce, err := GenerateNonce()
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, "error generating discovery nonce: "+err.Error()), nil
	}
	err = cr.SessionManager.SetDiscoveryNonce(data.SessionID, discoveryNonce)
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, err.Error()), nil
	}
	response := cr.Network.BuildResponse(ctx, request, &SignedChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: true,
		},
		Payload: &SignedChallengePayload{
			SignedNonce:       sign.Bytes(),
			XorDiscoveryNonce: xorNonce,
			DiscoveryNonce:    discoveryNonce,
		},
	})
	return response, nil
}

func (cr *challengeResponseController) buildChallenge1ErrorResponse(ctx context.Context, request network.Request, err string) network.Response {
	log.Warn(err)
	return cr.Network.BuildResponse(ctx, request, &ChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: false,
			Error:   err,
		},
	})
}

func (cr *challengeResponseController) processChallenge2(ctx context.Context, request network.Request) (network.Response, error) {
	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.processChallenge2")
	defer span.End()
	data := request.GetData().(*SignedChallengeRequest)
	cert, discoveryNonce, err := cr.SessionManager.GetChallengeData(data.SessionID)
	if err != nil {
		return cr.buildChallenge2ErrorResponse(ctx, request, err.Error()), nil
	}
	sign := insolar.SignatureFromBytes(data.SignedDiscoveryNonce)
	success := cr.Cryptography.Verify(cert.GetPublicKey(), sign, Xor(data.XorNonce, discoveryNonce))
	if !success {
		return cr.buildChallenge2ErrorResponse(ctx, request, "node %s signature check failed"), nil
	}
	err = cr.SessionManager.ChallengePassed(data.SessionID)
	if err != nil {
		return cr.buildChallenge2ErrorResponse(ctx, request, err.Error()), nil
	}
	response := cr.Network.BuildResponse(ctx, request, &ChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: true,
		},
		Payload: &ChallengePayload{
			AssignShortID: GenerateShortID(cr.NodeKeeper, *cert.GetNodeRef()),
		},
	})
	return response, nil
}

func (cr *challengeResponseController) buildChallenge2ErrorResponse(ctx context.Context, request network.Request, err string) network.Response {
	log.Warn(err)
	return cr.Network.BuildResponse(ctx, request, &SignedChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: false,
			Error:   err,
		},
	})
}

func (cr *challengeResponseController) Init(ctx context.Context) error {
	cr.Network.RegisterRequestHandler(types.Challenge1, cr.processChallenge1)
	cr.Network.RegisterRequestHandler(types.Challenge2, cr.processChallenge2)
	return nil
}

func (cr *challengeResponseController) sendRequest1(ctx context.Context, discoveryHost *host.Host,
	sessionID SessionID, nonce Nonce) (*SignedChallengePayload, error) {

	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.sendRequest1")
	defer span.End()
	request := cr.Network.NewRequestBuilder().Type(types.Challenge1).Data(&ChallengeRequest{
		SessionID: sessionID, Nonce: nonce}).Build()
	future, err := cr.Network.SendRequestToHost(ctx, request, discoveryHost)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending challenge request")
	}
	response, err := future.WaitResponse(cr.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting response for challenge request")
	}
	data := response.GetData().(*SignedChallengeResponse)
	if !data.Header.Success {
		return nil, errors.Wrap(err, "Discovery node returned error for challenge request: "+data.Header.Error)
	}
	return data.Payload, nil
}

func (cr *challengeResponseController) sendRequest2(ctx context.Context, discoveryHost *host.Host,
	sessionID SessionID, signedDiscoveryNonce SignedNonce, xorNonce Nonce) (*ChallengePayload, error) {

	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.sendRequest2")
	defer span.End()
	request := cr.Network.NewRequestBuilder().Type(types.Challenge2).Data(&SignedChallengeRequest{
		SessionID: sessionID, XorNonce: xorNonce, SignedDiscoveryNonce: signedDiscoveryNonce}).Build()
	future, err := cr.Network.SendRequestToHost(ctx, request, discoveryHost)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending challenge request")
	}
	response, err := future.WaitResponse(cr.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting response for challenge request")
	}
	data := response.GetData().(*ChallengeResponse)
	if !data.Header.Success {
		return nil, errors.Wrap(err, "Discovery node returned error for challenge request: "+data.Header.Error)
	}
	return data.Payload, nil
}

// Execute double challenge response between the node and the discovery node (step 3 of the bootstrap process)
func (cr *challengeResponseController) Execute(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) (*ChallengePayload, error) {
	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.Execute")
	defer span.End()
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, errors.Wrap(err, "error generating nonce")
	}
	inslogger.FromContext(ctx).Debugf("Generated nonce: %s", base58.Encode(nonce))

	data, err := cr.sendRequest1(ctx, discoveryNode.Host, sessionID, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "error executing challenge response (step 1)")
	}

	inslogger.FromContext(ctx).Debugf("Discovery SignedNonce: %s", base58.Encode(data.SignedNonce))
	inslogger.FromContext(ctx).Debugf("Discovery DiscoveryNonce: %s", base58.Encode(data.DiscoveryNonce))
	inslogger.FromContext(ctx).Debugf("Discovery XorDiscoveryNonce: %s", base58.Encode(data.XorDiscoveryNonce))

	sign := insolar.SignatureFromBytes(data.SignedNonce)
	success := cr.Cryptography.Verify(discoveryNode.Node.GetPublicKey(), sign, Xor(nonce, data.XorDiscoveryNonce))
	if !success {
		return nil, errors.New("Error checking signed nonce from discovery node")
	}

	xorNonce, err := GenerateNonce()
	if err != nil {
		return nil, errors.Wrap(err, "error generating xor nonce")
	}
	signedDiscoveryNonce, err := cr.Cryptography.Sign(Xor(xorNonce, data.DiscoveryNonce))
	if err != nil {
		return nil, errors.Wrap(err, "error signing discovery nonce")
	}
	payload, err := cr.sendRequest2(ctx, discoveryNode.Host, sessionID, signedDiscoveryNonce.Bytes(), xorNonce)
	if err != nil {
		return nil, errors.Wrap(err, "error executing challenge response (step 2)")
	}
	return payload, nil
}

func NewChallengeResponseController(options *common.Options) ChallengeResponseController {
	return &challengeResponseController{options: options}
}
