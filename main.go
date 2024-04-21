package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	badger "github.com/dgraph-io/badger/v4"
)

var db *badger.DB = nil

type Patient = struct {
	Id    uuid.UUID
	Name  string
	Owner string
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func getPatientById(id uuid.UUID) (Patient, error) {
	var patient Patient
	bid, err := id.MarshalBinary()
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

func patient(c echo.Context) error {
	l := c.Logger()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		l.Errorf("invalid id: %s", err)
		return c.NoContent(http.StatusOK)
	}

	patient, err := getPatientById(id)
	if err != nil {
		l.Errorf("error getting patient: %v", err)
		return c.NoContent(http.StatusOK)
	}
	return Render(c, http.StatusOK, patientComponent(patient))
}

func patientEdit(c echo.Context) error {
	l := c.Logger()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		id = uuid.New()
	}
	patient, err := getPatientById(id)
	if err != nil {
		l.Errorf("error getting patient: %v", err)
		patient = Patient{Id: id, Name: "", Owner: ""}
	}
	return Render(c, http.StatusOK, page(patientEditComponent(patient)))
}

func patientUpdate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		id = uuid.New()
	}

	name := c.FormValue("name")
	owner := c.FormValue("owner")

	patient := Patient{Id: id, Name: name, Owner: owner}
	db.Update(func(txn *badger.Txn) error {
		value, err := cbor.Marshal(&patient)
		if err != nil {
			return err
		}
		id, _ := id.MarshalBinary()
		entry := badger.NewEntry(id, value)
		err = txn.SetEntry(entry)
		return err
	})
	return Render(c, http.StatusOK, page(patientComponent(patient)))
}

func patientList(c echo.Context) error {

	patients := make([]Patient, 0)
	db.View(func(txn *badger.Txn) error {

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
	return Render(c, http.StatusOK, getPatientsList(patients))

}

func index(c echo.Context) error {
	patients := make([]Patient, 0)
	db.View(func(txn *badger.Txn) error {

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
	return Render(c, http.StatusOK, page(patientsComponent(patients)))
}

func main() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	engine := echo.New()

	engine.Logger.SetLevel(log.DEBUG)
	engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	engine.GET("/", index)
	engine.GET("/patient/new", patientEdit)
	engine.GET("/patient/list", patientList)
	engine.GET("/patient/:id", patient)
	engine.GET("/patient/:id/edit", patientEdit)
	engine.PUT("/patient/:id", patientUpdate)
	engine.Static("/assets", "assets")

	fmt.Println("Listening on :8080")
	engine.Logger.Fatal(engine.Start(":8080"))
}
