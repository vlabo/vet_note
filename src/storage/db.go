package storage

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/vlabo/vet_note/src/model"
	bolt "go.etcd.io/bbolt"
)

const PatientsBucket = "Patients"

var database *bolt.DB = nil

func InitStorage(dataPath string) error {
	var err error

	// Make sure path exists
	err = os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Open database
	path := dataPath + "/data.db"
	// .Debug("DB", "path", path)
	slog.Info("Storage", "db_path", path)
	database, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}

	return nil
}

func CloseStorage() error {
	err := database.Close()
	database = nil
	return err
}

func AddPatient(patient *model.Patient) error {
	if database == nil {
		return fmt.Errorf("storage not initialized")
	}
	err := database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PatientsBucket))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte(PatientsBucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		patient.Id = uuid.New()
		patientMarshaled, err := cbor.Marshal(patient)
		if err != nil {
			return err
		}
		b.Put(patient.Id[:], patientMarshaled)
		return nil
	})

	return err
}

func GetAllPatients() ([]model.Patient, error) {
	if database == nil {
		return nil, fmt.Errorf("storage not initialized")
	}

	patients := []model.Patient{}

	err := database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PatientsBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", PatientsBucket)
		}

		err := b.ForEach(func(_ []byte, value []byte) error {
			var p model.Patient
			err := cbor.Unmarshal(value, &p)
			patients = append(patients, p)
			return err
		})

		return err
	})

	return patients, err
}
