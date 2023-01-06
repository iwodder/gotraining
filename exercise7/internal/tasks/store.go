package tasks

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

type dbFn func(tx *bolt.Tx) error
type BoltRepository struct {
	db *bolt.DB
}

func NewRepository(db string) (*BoltRepository, error) {
	conn, err := getDb(db)
	if err != nil {
		return nil, err
	}
	return &BoltRepository{
		db: conn,
	}, nil
}

func (b *BoltRepository) Store(t Task) Task {
	if updateFn, err := serializeTask(&t); err == nil {
		err = b.db.Update(updateFn)
	}
	return t
}

func (b *BoltRepository) LoadAll() []Task {
	var ret []Task
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))

		var buff bytes.Buffer
		return b.ForEach(func(k, v []byte) error {
			var t Task
			buff.Write(v)
			dec := gob.NewDecoder(&buff)
			err := dec.Decode(&t)
			if err != nil {
				return err
			}
			ret = append(ret, t)
			return nil
		})
	})
	if err != nil {
		return nil
	}
	return ret
}

func (b *BoltRepository) Update(task Task) Task {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(&task)
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, task.ID)
		b.Put(key, buff.Bytes())
		return nil
	})
	if err != nil {
		return task
	}
	return task
}

func (b *BoltRepository) LoadTask(name string) (*Task, error) {
	var ret *Task
	err := b.db.View(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		tx.Bucket([]byte("tasks")).ForEach(func(k, v []byte) error {
			var t Task
			buff.Write(v)
			dec := gob.NewDecoder(&buff)
			err := dec.Decode(&t)
			if err != nil {
				return err
			}
			if t.Name == name {
				ret = &t

			}
			return nil
		})
		if ret == nil {
			return errors.New(fmt.Sprintf("Unable to find task with name %s", name))
		}
		return nil
	})
	return ret, err
}

func getDb(db string) (*bolt.DB, error) {
	open, err := bolt.Open(db, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := open.Update(createTaskBucket()); err != nil {
		return nil, err
	}
	return open, nil
}

func serializeTask(t *Task) (dbFn, error) {
	return func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		t.ID, _ = b.NextSequence()

		var buff bytes.Buffer
		enc := gob.NewEncoder(&buff)
		err := enc.Encode(t)

		if err != nil {
			return err
		}

		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, t.ID)
		return b.Put(key, buff.Bytes())
	}, nil
}

func createTaskBucket() dbFn {
	return func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		if err != nil {
			return err
		}
		return nil
	}
}
