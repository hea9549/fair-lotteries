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

package db

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type DBHandle struct {
	dbName string
	db     *DB
}

type DBProvider struct {
	db        *DB
	mux       sync.Mutex
	dbHandles map[string]*DBHandle
}

func CreateNewDBProvider(levelDbPath string) *DBProvider {
	db := CreateNewDB(levelDbPath)
	db.Open()
	return &DBProvider{db, sync.Mutex{}, make(map[string]*DBHandle)}
}

func (p *DBProvider) Close() {
	p.db.Close()
}

func (p *DBProvider) GetDBHandle(dbName string) *DBHandle {
	p.mux.Lock()
	defer p.mux.Unlock()

	dbHandle := p.dbHandles[dbName]
	if dbHandle == nil {
		dbHandle = &DBHandle{dbName, p.db}
		p.dbHandles[dbName] = dbHandle
	}

	return dbHandle
}

func (h *DBHandle) Get(key []byte) ([]byte, error) {
	return h.db.Get(dbKey(h.dbName, key))
}

func (h *DBHandle) Put(key []byte, value []byte, sync bool) error {
	return h.db.Put(dbKey(h.dbName, key), value, sync)
}

func (h *DBHandle) Delete(key []byte, sync bool) error {
	return h.db.Delete(dbKey(h.dbName, key), sync)
}

func (h *DBHandle) WriteBatch(KVs map[string][]byte, sync bool) error {
	batch := &leveldb.Batch{}
	for k, v := range KVs {
		key := dbKey(h.dbName, []byte(k))
		if v == nil {
			batch.Delete(key)
		} else {
			batch.Put(key, v)
		}
	}
	return h.db.writeBatch(batch, sync)
}

func (h *DBHandle) GetIteratorWithPrefix() KeyValueDBIterator {
	return h.db.GetIteratorWithPrefix([]byte(h.dbName + "_"))
}

func (h *DBHandle) Snapshot() (map[string][]byte, error) {
	return h.db.Snapshot()
}

func dbKey(dbName string, key []byte) []byte {
	dbName = dbName + "_"
	return append([]byte(dbName), key...)
}
