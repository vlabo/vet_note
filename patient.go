package main

import (
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
)

type Examination struct {
	Date        time.Time
	Type        string
	Description string
}
type Patient = struct {
	Id           uuid.UUID
	Name         string
	Owner        string
	Examinations []Examination
}

func getPatientById(db *badger.DB, id uuid.UUID) (*Patient, error) {
	var patient *Patient
	bid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(bid)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			err = cbor.Unmarshal(val, &patient)
			return err
		})
	})

	return patient, err
}

func updatePatent(patient *Patient) error {
	return db.Update(func(txn *badger.Txn) error {
		value, err := cbor.Marshal(&patient)
		if err != nil {
			return err
		}
		id, _ := patient.Id.MarshalBinary()
		entry := badger.NewEntry(id, value)
		err = txn.SetEntry(entry)
		return err
	})
}

func getPatientList(db *badger.DB) []Patient {
	patients := make([]Patient, 0)
	_ = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var patient Patient
				_ = cbor.Unmarshal(v, &patient)

				patients = append(patients, patient)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return patients
}

type Patients []Patient

func (p Patients) String(i int) string {
	return p[i].Name + " " + p[i].Owner
}

func (p Patients) Len() int {
	return len(p)
}
