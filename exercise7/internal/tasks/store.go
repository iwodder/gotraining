package tasks

import (
	"bytes"
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

func (b *BoltRepository) Store(t Task) (Task, error) {
	if updateFn, err := serializeTask(&t); err == nil {
		err = b.db.Update(updateFn)
		if err != nil {
			return Task{}, err
		}
	}

	return t, nil
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

		b.Put([]byte(task.Name), buff.Bytes())
		return nil
	})
	if err != nil {
		return task
	}
	return task
}

func (b *BoltRepository) Load(name string) (*Task, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		data = b.Get([]byte(name))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New(fmt.Sprintf("No task with name %s", name))
	}
	var t Task
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (b *BoltRepository) Delete(task Task) {
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		return b.Delete([]byte(task.Name))
	})
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

		var buff bytes.Buffer
		enc := gob.NewEncoder(&buff)
		err := enc.Encode(t)

		if err != nil {
			return err
		}

		return b.Put([]byte(t.Name), buff.Bytes())
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
