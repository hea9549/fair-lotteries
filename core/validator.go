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
	"errors"
	"github.com/hea9549/fair-lotteries/common"
	"github.com/hea9549/fair-lotteries/log"
	"reflect"
	"time"
)

var ErrInsufficientFields = errors.New("previous seal or transaction list seal is not set")
var ErrEmptyTxList = errors.New("empty TxList")

type Validator struct {}

func (t *Validator) ValidateBlock(seal []byte, comparisonBlock Block) (bool, error) {

	comparisonSeal, err := t.BuildBlockSeal(comparisonBlock.Timestamp, comparisonBlock.PrevSeal, comparisonBlock.TxSeal)

	if err != nil {
		log.Error(nil,"[blockchain] error while build seal")
		return false, err
	}
	return bytes.Compare(seal, comparisonSeal) == 0, nil
}

// ValidateTransaction 함수는 주어진 Transaction이 이 txSeal에 올바로 있는지를 확인한다.
func (t *Validator) ValidateTransaction(txSeal [][]byte, transaction Transaction) (bool, error) {
	hash, err := transaction.CalculateSeal()
	if err != nil {
		return false, err
	}

	index := -1
	for i, h := range txSeal {
		if bytes.Compare(h, hash) == 0 {
			index = i
		}
	}

	if index == -1 {
		return false, nil
	}

	var siblingIndex, parentIndex int
	for index > 0 {
		var isLeft bool
		if index%2 == 0 {
			siblingIndex = index - 1
			parentIndex = (index - 1) / 2
			isLeft = false
		} else {
			siblingIndex = index + 1
			parentIndex = index / 2
			isLeft = true
		}

		var parentHash []byte
		if isLeft {
			parentHash = calculateIntermediateNodeHash(txSeal[index], txSeal[siblingIndex])
		} else {
			parentHash = calculateIntermediateNodeHash(txSeal[siblingIndex], txSeal[index])
		}

		if bytes.Compare(parentHash, txSeal[parentIndex]) != 0 {
			return false, nil
		}

		index = parentIndex
	}

	return true, nil
}
func (t *Validator) ValidateTxSeal(txSeal [][]byte, txList []Transaction) (bool, error) {

	if isEmpty(txList) {
		return true, nil
	}

	leafNodeList, err := convertToLeafNodeList(txList)
	if err != nil {
		return false, err
	}

	tree, err := buildTree(leafNodeList, leafNodeList)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(txSeal, tree), nil
}


func convertToLeafNodeList(txList []Transaction) ([][]byte, error) {
	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, err := tx.CalculateSeal()
		if err != nil {
			return nil, err
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	return leafNodeList, nil
}
func isEmpty(txList []Transaction) bool {
	if len(txList) == 0 {
		return true
	}
	return false
}


func (t *Validator) BuildBlockSeal(timeStamp time.Time, prevSeal []byte, txSeal [][]byte) ([]byte, error) {
	timestamp, err := timeStamp.MarshalText()
	if err != nil {
		return nil, err
	}

	if prevSeal == nil || txSeal == nil {
		return nil, ErrInsufficientFields
	}

	var rootHash []byte
	if len(txSeal) == 0 {
		rootHash = make([]byte, 0)
	} else {
		rootHash = txSeal[0]
	}
	combined := append(prevSeal, rootHash...)
	combined = append(combined, timestamp...)

	seal := common.CalculateHash(combined)
	return seal, nil
}

func (t *Validator) BuildTxSeal(txList []Transaction) ([][]byte, error) {
	if len(txList) == 0 {
		return nil, ErrEmptyTxList
	}

	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, err := tx.CalculateSeal()
		if err != nil {
			return nil, err
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	// leafNodeList의 개수는 짝수개로 맞춤. (홀수 일 경우 마지막 Tx를 중복 저장.)
	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	tree, err := buildTree(leafNodeList, leafNodeList)
	if err != nil {
		return nil, err
	}

	// Validator 는 MerkleTree의 루트노드(tree[0])를 Proof로 간주함
	return tree, nil
}


func buildTree(nodeList [][]byte, fullNodeList [][]byte) ([][]byte, error) {
	intermediateNodeList := make([][]byte, 0)
	for i := 0; i < len(nodeList); i += 2 {
		leftIndex, rightIndex := i, i+1

		if i+1 == len(nodeList) {
			rightIndex = i
		}

		leftNode, rightNode := nodeList[leftIndex], nodeList[rightIndex]

		intermediateNode := calculateIntermediateNodeHash(leftNode, rightNode)

		intermediateNodeList = append(intermediateNodeList, intermediateNode)

		if len(nodeList) == 2 {
			return append(intermediateNodeList, fullNodeList...), nil
		}
	}

	newFullNodeList := append(intermediateNodeList, fullNodeList...)

	return buildTree(intermediateNodeList, newFullNodeList)
}


func calculateIntermediateNodeHash(leftHash []byte, rightHash []byte) []byte {
	combinedHash := append(leftHash, rightHash...)

	return common.CalculateHash(combinedHash)
}
