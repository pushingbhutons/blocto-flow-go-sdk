/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/portto/blocto-flow-go-sdk"
	"github.com/portto/blocto-flow-go-sdk/client"
	"github.com/portto/blocto-flow-go-sdk/crypto"
	"github.com/portto/blocto-flow-go-sdk/examples"
	"github.com/portto/blocto-flow-go-sdk/templates"
)

func main() {
	AddAccountKeyDemo()
}

func AddAccountKeyDemo() {
	ctx := context.Background()

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	acctAddr, acctKey, acctSigner := examples.RandomAccount(flowClient)

	// Create the new key to add to your account
	myPrivateKey := examples.RandomPrivateKey()
	myAcctKey := flow.NewAccountKey().
		FromPrivateKey(myPrivateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	// Create a Cadence script that will add another key to our account.
	addKeyScript, err := templates.AddAccountKey(myAcctKey)
	examples.Handle(err)

	// Create a transaction to execute the script.
	// The transaction is signed by our account key so it has permission to add keys.
	addKeyTx := flow.NewTransaction().
		SetScript(addKeyScript).
		SetProposalKey(acctAddr, acctKey.ID, acctKey.SequenceNumber).
		SetPayer(acctAddr).
		// This defines which accounts are accessed by this transaction
		AddAuthorizer(acctAddr)

	// Sign the transaction with the new account.
	err = addKeyTx.SignEnvelope(acctAddr, acctKey.ID, acctSigner)
	examples.Handle(err)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, *addKeyTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, addKeyTx.ID())

	fmt.Println("Public key added to account!")
}
