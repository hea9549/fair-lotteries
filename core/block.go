/* 
 * Copyright 2019 hea9549
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"bytes"
	"github.com/hea9549/fair-lotteries/common"
	"github.com/hea9549/fair-lotteries/log"
	"time"
)

type Block struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []*Transaction
	TxSeal    [][]byte
	Timestamp time.Time
}

func (block *Block) PutTx(transaction Transaction) error {
	if block.TxList == nil {
		block.TxList = make([]*Transaction, 0)
	}
	block.TxList = append(block.TxList, &transaction)

	return nil
}

func (block *Block) GetTxList() []Transaction {
	txList := make([]Transaction, 0)
	for _, tx := range block.TxList {
		txList = append(txList, *tx)
	}
	return txList
}

// check to ready to publish.
// todo : is check all?
func (block *Block) IsReadyToPublish() bool {
	return block.Seal != nil
}

func (block *Block) IsPrev(serializedPrevBlock []byte) bool {
	prevBlock := &Block{}
	err :=common.Deserialize(serializedPrevBlock,prevBlock)

	if err != nil{
		log.Error(nil,"[blockchain] error while check prev. cant' deserialize prev block")
		return false
	}

	return bytes.Compare(prevBlock.Seal, block.Seal) == 0
}

