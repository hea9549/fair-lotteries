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

package blockchain

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hea9549/fair-lotteries/common"
	"github.com/hea9549/fair-lotteries/common/db"
	"github.com/hea9549/fair-lotteries/core"
)

var ErrPrevSealMismatch = errors.New("PrevSeal value mismatch")
var ErrSealValidation = errors.New("seal validation failed")
var ErrTxSealValidation = errors.New("txSeal validation failed")
var ErrNoValidator = errors.New("validator not defined")

const (
	blockSealDB   = "block_seal"
	blockHeightDB = "block_height"
	transactionDB = "transaction"
	utilDB        = "util"
	lastBlockKey  = "last_block"
)

type Repository struct {
	DBProvider *db.DBProvider
	mux        *sync.RWMutex
	validator  *core.Validator
}

func NewRepository(dbPath string) (*Repository, error) {
	validator := new(core.Validator)
	dbProvider := db.CreateNewDBProvider(dbPath)

	return &Repository{
		mux:        &sync.RWMutex{},
		DBProvider: dbProvider,
		validator:  validator,
	}, nil
}

func (y *Repository) AddBlock(block core.Block) error {
	serializedBlock, err := common.Serialize(block)
	if err != nil {
		return err
	}

	err = y.validateBlock(block)
	if err != nil {
		return err
	}

	utilDB := y.DBProvider.GetDBHandle(utilDB)
	blockSealDB := y.DBProvider.GetDBHandle(blockSealDB)
	blockHeightDB := y.DBProvider.GetDBHandle(blockHeightDB)
	transactionDB := y.DBProvider.GetDBHandle(transactionDB)
	err = blockSealDB.Put(block.Seal, serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockHeightDB.Put([]byte(fmt.Sprint(block.Height)), block.Seal, true)
	if err != nil {
		return err
	}

	err = utilDB.Put([]byte(lastBlockKey), serializedBlock, true)
	if err != nil {
		return err
	}

	for _, tx := range block.GetTxList() {
		serializedTX, err := common.Serialize(tx)
		if err != nil {
			return err
		}

		err = transactionDB.Put([]byte(tx.ID), serializedTX, true)
		if err != nil {
			return err
		}

		err = utilDB.Put([]byte(tx.ID), block.Seal, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (y *Repository) GetBlockByHeight(block core.Block, height uint64) error {
	blockHeightDB := y.DBProvider.GetDBHandle(blockHeightDB)

	blockSeal, err := blockHeightDB.Get([]byte(fmt.Sprint(height)))
	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

func (y *Repository) GetBlockBySeal(block core.Block, seal []byte) error {
	blockSealDB := y.DBProvider.GetDBHandle(blockSealDB)

	serializedBlock, err := blockSealDB.Get(seal)
	if err != nil {
		return err
	}

	err = common.Deserialize(serializedBlock, block)

	return err
}

func (y *Repository) GetBlockByTxID(block core.Block, txID string) error {
	utilDB := y.DBProvider.GetDBHandle(utilDB)

	blockSeal, err := utilDB.Get([]byte(txID))

	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

func (y *Repository) GetLastBlock(block core.Block) error {
	utilDB := y.DBProvider.GetDBHandle(utilDB)

	serializedBlock, err := utilDB.Get([]byte(lastBlockKey))
	if serializedBlock == nil || err != nil {
		return err
	}

	err = common.Deserialize(serializedBlock, block)

	return err
}

func (y *Repository) GetTransactionByTxID(transaction core.Transaction, txID string) error {
	transactionDB := y.DBProvider.GetDBHandle(transactionDB)

	serializedTX, err := transactionDB.Get([]byte(txID))
	if err != nil {
		return err
	}

	err = common.Deserialize(serializedTX, transaction)

	return err
}

func (y *Repository) validateBlock(block core.Block) error {
	if y.validator == nil {
		return ErrNoValidator
	}

	utilDB := y.DBProvider.GetDBHandle(utilDB)

	lastBlockByte, err := utilDB.Get([]byte(lastBlockKey))
	if err != nil {
		return err
	}
	if lastBlockByte != nil && !block.IsPrev(lastBlockByte) {
		return ErrPrevSealMismatch
	}

	// Validate the Seal of the new block using the validator
	result, err := y.validator.ValidateBlock(block.Seal, block)
	if err != nil {
		return err
	}

	if !result {
		return ErrSealValidation
	}

	// Validate the TxSeal of the new block using the validator
	result, err = y.validator.ValidateTxSeal(block.TxSeal, block.GetTxList())
	if err != nil {
		return err
	}

	if !result {
		return ErrTxSealValidation
	}

	return nil
}
