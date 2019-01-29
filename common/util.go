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

package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func Serialize(object interface{}) ([]byte, error) { //모든 stuct 받기 위해 interface{} 타입의 입력으로 선언
	data, err := json.Marshal(object)
	if err != nil {
		panic(fmt.Sprintf("Error encoding : %s", err))
	}
	return data, nil
}


func Deserialize(serializedBytes []byte, object interface{}) error {
	if len(serializedBytes) == 0 {
		return nil
	}
	err := json.Unmarshal(serializedBytes, object)
	if err != nil {
		panic(fmt.Sprintf("Error decoding : %s", err))
	}
	return err
}


func CalculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}