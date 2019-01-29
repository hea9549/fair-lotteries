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

package txpool

import (
	"errors"
	"sync"

	"github.com/hea9549/fair-lotteries/core"
	"github.com/hea9549/fair-lotteries/log"
)

var ErrExistTx = errors.New("transaction is already in txpool")

type Repository struct {
	mux         sync.RWMutex
	unCommitTxs map[string]core.Transaction
}

func NewRepository() *Repository {
	return &Repository{
		mux:         sync.RWMutex{},
		unCommitTxs: make(map[string]core.Transaction),
	}
}

func (r *Repository) AddTransaction(tx core.Transaction) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.unCommitTxs[tx.ID]; ok {
		return ErrExistTx
	}

	r.unCommitTxs[tx.ID] = tx
	return nil
}

func (r *Repository) GetAllUnCommitTransaction() []core.Transaction {
	r.mux.Lock()
	defer r.mux.Unlock()

	retSlice := make([]core.Transaction, 0)

	for _, v := range r.unCommitTxs {
		retSlice = append(retSlice, v)
	}

	return retSlice
}

func (r *Repository) RemoveCommitTransaction(txs []core.Transaction) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, oneTx := range txs {
		if _, ok := r.unCommitTxs[oneTx.ID]; ok {
			delete(r.unCommitTxs, oneTx.ID)
		} else {
			log.Warn(&log.Fields{"txId": oneTx.ID}, "Fail to remove commit transaction in txpool. \n There is no transaction in pool")
		}
	}
}
