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
	"github.com/hea9549/fair-lotteries/common"
	"github.com/hea9549/fair-lotteries/log"
	"time"
)

type Transaction struct {
	ID        string
	Type      string
	Timestamp time.Time
	Function  string
	Args      []string
	Signature []byte
}

func (t *Transaction) GetContent() ([]byte, error) {
	serialized, err := common.Serialize(t)
	if err != nil {
		return nil, err
	}

	return serialized, nil
}

func (t *Transaction) CalculateSeal() ([]byte, error) {
	serializedTx, err := common.Serialize(t)
	if err != nil {
		log.Error(nil,"[blockchain] error while serialize transaction")
		return nil, err
	}

	return common.CalculateHash(serializedTx), nil
}


