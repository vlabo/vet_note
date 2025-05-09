// Code generated by BobGen sqlite v0.31.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/dm"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/bob/dialect/sqlite/um"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
)

// Patient is an object representing the database table.
type Patient struct {
	ID          int32                     `db:"id,pk" `
	Type        null.Val[string]          `db:"type" `
	Name        null.Val[string]          `db:"name" `
	Gender      null.Val[string]          `db:"gender" `
	Age         null.Val[decimal.Decimal] `db:"age" `
	ChipID      null.Val[string]          `db:"chip_id" `
	Weight      null.Val[float32]         `db:"weight" `
	Castrated   null.Val[decimal.Decimal] `db:"castrated" `
	Note        null.Val[string]          `db:"note" `
	Owner       null.Val[string]          `db:"owner" `
	OwnerPhone  null.Val[string]          `db:"owner_phone" `
	Folder      null.Val[decimal.Decimal] `db:"folder" `
	IndexFolder null.Val[decimal.Decimal] `db:"index_folder" `
	CreatedAt   null.Val[time.Time]       `db:"created_at" `
	UpdatedAt   null.Val[time.Time]       `db:"updated_at" `
	DeletedAt   null.Val[time.Time]       `db:"deleted_at" `

	R patientR `db:"-" `
}

// PatientSlice is an alias for a slice of pointers to Patient.
// This should almost always be used instead of []*Patient.
type PatientSlice []*Patient

// Patients contains methods to work with the patients table
var Patients = sqlite.NewTablex[*Patient, PatientSlice, *PatientSetter]("", "patients")

// PatientsQuery is a query on the patients table
type PatientsQuery = *sqlite.ViewQuery[*Patient, PatientSlice]

// patientR is where relationships are stored.
type patientR struct {
	Procedures ProcedureSlice // fk_procedures_0
}

type patientColumnNames struct {
	ID          string
	Type        string
	Name        string
	Gender      string
	Age         string
	ChipID      string
	Weight      string
	Castrated   string
	Note        string
	Owner       string
	OwnerPhone  string
	Folder      string
	IndexFolder string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}

var PatientColumns = buildPatientColumns("patients")

type patientColumns struct {
	tableAlias  string
	ID          sqlite.Expression
	Type        sqlite.Expression
	Name        sqlite.Expression
	Gender      sqlite.Expression
	Age         sqlite.Expression
	ChipID      sqlite.Expression
	Weight      sqlite.Expression
	Castrated   sqlite.Expression
	Note        sqlite.Expression
	Owner       sqlite.Expression
	OwnerPhone  sqlite.Expression
	Folder      sqlite.Expression
	IndexFolder sqlite.Expression
	CreatedAt   sqlite.Expression
	UpdatedAt   sqlite.Expression
	DeletedAt   sqlite.Expression
}

func (c patientColumns) Alias() string {
	return c.tableAlias
}

func (patientColumns) AliasedAs(alias string) patientColumns {
	return buildPatientColumns(alias)
}

func buildPatientColumns(alias string) patientColumns {
	return patientColumns{
		tableAlias:  alias,
		ID:          sqlite.Quote(alias, "id"),
		Type:        sqlite.Quote(alias, "type"),
		Name:        sqlite.Quote(alias, "name"),
		Gender:      sqlite.Quote(alias, "gender"),
		Age:         sqlite.Quote(alias, "age"),
		ChipID:      sqlite.Quote(alias, "chip_id"),
		Weight:      sqlite.Quote(alias, "weight"),
		Castrated:   sqlite.Quote(alias, "castrated"),
		Note:        sqlite.Quote(alias, "note"),
		Owner:       sqlite.Quote(alias, "owner"),
		OwnerPhone:  sqlite.Quote(alias, "owner_phone"),
		Folder:      sqlite.Quote(alias, "folder"),
		IndexFolder: sqlite.Quote(alias, "index_folder"),
		CreatedAt:   sqlite.Quote(alias, "created_at"),
		UpdatedAt:   sqlite.Quote(alias, "updated_at"),
		DeletedAt:   sqlite.Quote(alias, "deleted_at"),
	}
}

type patientWhere[Q sqlite.Filterable] struct {
	ID          sqlite.WhereMod[Q, int32]
	Type        sqlite.WhereNullMod[Q, string]
	Name        sqlite.WhereNullMod[Q, string]
	Gender      sqlite.WhereNullMod[Q, string]
	Age         sqlite.WhereNullMod[Q, decimal.Decimal]
	ChipID      sqlite.WhereNullMod[Q, string]
	Weight      sqlite.WhereNullMod[Q, float32]
	Castrated   sqlite.WhereNullMod[Q, decimal.Decimal]
	Note        sqlite.WhereNullMod[Q, string]
	Owner       sqlite.WhereNullMod[Q, string]
	OwnerPhone  sqlite.WhereNullMod[Q, string]
	Folder      sqlite.WhereNullMod[Q, decimal.Decimal]
	IndexFolder sqlite.WhereNullMod[Q, decimal.Decimal]
	CreatedAt   sqlite.WhereNullMod[Q, time.Time]
	UpdatedAt   sqlite.WhereNullMod[Q, time.Time]
	DeletedAt   sqlite.WhereNullMod[Q, time.Time]
}

func (patientWhere[Q]) AliasedAs(alias string) patientWhere[Q] {
	return buildPatientWhere[Q](buildPatientColumns(alias))
}

func buildPatientWhere[Q sqlite.Filterable](cols patientColumns) patientWhere[Q] {
	return patientWhere[Q]{
		ID:          sqlite.Where[Q, int32](cols.ID),
		Type:        sqlite.WhereNull[Q, string](cols.Type),
		Name:        sqlite.WhereNull[Q, string](cols.Name),
		Gender:      sqlite.WhereNull[Q, string](cols.Gender),
		Age:         sqlite.WhereNull[Q, decimal.Decimal](cols.Age),
		ChipID:      sqlite.WhereNull[Q, string](cols.ChipID),
		Weight:      sqlite.WhereNull[Q, float32](cols.Weight),
		Castrated:   sqlite.WhereNull[Q, decimal.Decimal](cols.Castrated),
		Note:        sqlite.WhereNull[Q, string](cols.Note),
		Owner:       sqlite.WhereNull[Q, string](cols.Owner),
		OwnerPhone:  sqlite.WhereNull[Q, string](cols.OwnerPhone),
		Folder:      sqlite.WhereNull[Q, decimal.Decimal](cols.Folder),
		IndexFolder: sqlite.WhereNull[Q, decimal.Decimal](cols.IndexFolder),
		CreatedAt:   sqlite.WhereNull[Q, time.Time](cols.CreatedAt),
		UpdatedAt:   sqlite.WhereNull[Q, time.Time](cols.UpdatedAt),
		DeletedAt:   sqlite.WhereNull[Q, time.Time](cols.DeletedAt),
	}
}

var PatientErrors = &patientErrors{
	ErrUniquePkMainPatients: &UniqueConstraintError{s: "pk_main_patients"},
}

type patientErrors struct {
	ErrUniquePkMainPatients *UniqueConstraintError
}

// PatientSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type PatientSetter struct {
	ID          omit.Val[int32]               `db:"id,pk" `
	Type        omitnull.Val[string]          `db:"type" `
	Name        omitnull.Val[string]          `db:"name" `
	Gender      omitnull.Val[string]          `db:"gender" `
	Age         omitnull.Val[decimal.Decimal] `db:"age" `
	ChipID      omitnull.Val[string]          `db:"chip_id" `
	Weight      omitnull.Val[float32]         `db:"weight" `
	Castrated   omitnull.Val[decimal.Decimal] `db:"castrated" `
	Note        omitnull.Val[string]          `db:"note" `
	Owner       omitnull.Val[string]          `db:"owner" `
	OwnerPhone  omitnull.Val[string]          `db:"owner_phone" `
	Folder      omitnull.Val[decimal.Decimal] `db:"folder" `
	IndexFolder omitnull.Val[decimal.Decimal] `db:"index_folder" `
	CreatedAt   omitnull.Val[time.Time]       `db:"created_at" `
	UpdatedAt   omitnull.Val[time.Time]       `db:"updated_at" `
	DeletedAt   omitnull.Val[time.Time]       `db:"deleted_at" `
}

func (s PatientSetter) SetColumns() []string {
	vals := make([]string, 0, 16)
	if !s.ID.IsUnset() {
		vals = append(vals, "id")
	}

	if !s.Type.IsUnset() {
		vals = append(vals, "type")
	}

	if !s.Name.IsUnset() {
		vals = append(vals, "name")
	}

	if !s.Gender.IsUnset() {
		vals = append(vals, "gender")
	}

	if !s.Age.IsUnset() {
		vals = append(vals, "age")
	}

	if !s.ChipID.IsUnset() {
		vals = append(vals, "chip_id")
	}

	if !s.Weight.IsUnset() {
		vals = append(vals, "weight")
	}

	if !s.Castrated.IsUnset() {
		vals = append(vals, "castrated")
	}

	if !s.Note.IsUnset() {
		vals = append(vals, "note")
	}

	if !s.Owner.IsUnset() {
		vals = append(vals, "owner")
	}

	if !s.OwnerPhone.IsUnset() {
		vals = append(vals, "owner_phone")
	}

	if !s.Folder.IsUnset() {
		vals = append(vals, "folder")
	}

	if !s.IndexFolder.IsUnset() {
		vals = append(vals, "index_folder")
	}

	if !s.CreatedAt.IsUnset() {
		vals = append(vals, "created_at")
	}

	if !s.UpdatedAt.IsUnset() {
		vals = append(vals, "updated_at")
	}

	if !s.DeletedAt.IsUnset() {
		vals = append(vals, "deleted_at")
	}

	return vals
}

func (s PatientSetter) Overwrite(t *Patient) {
	if !s.ID.IsUnset() {
		t.ID, _ = s.ID.Get()
	}
	if !s.Type.IsUnset() {
		t.Type, _ = s.Type.GetNull()
	}
	if !s.Name.IsUnset() {
		t.Name, _ = s.Name.GetNull()
	}
	if !s.Gender.IsUnset() {
		t.Gender, _ = s.Gender.GetNull()
	}
	if !s.Age.IsUnset() {
		t.Age, _ = s.Age.GetNull()
	}
	if !s.ChipID.IsUnset() {
		t.ChipID, _ = s.ChipID.GetNull()
	}
	if !s.Weight.IsUnset() {
		t.Weight, _ = s.Weight.GetNull()
	}
	if !s.Castrated.IsUnset() {
		t.Castrated, _ = s.Castrated.GetNull()
	}
	if !s.Note.IsUnset() {
		t.Note, _ = s.Note.GetNull()
	}
	if !s.Owner.IsUnset() {
		t.Owner, _ = s.Owner.GetNull()
	}
	if !s.OwnerPhone.IsUnset() {
		t.OwnerPhone, _ = s.OwnerPhone.GetNull()
	}
	if !s.Folder.IsUnset() {
		t.Folder, _ = s.Folder.GetNull()
	}
	if !s.IndexFolder.IsUnset() {
		t.IndexFolder, _ = s.IndexFolder.GetNull()
	}
	if !s.CreatedAt.IsUnset() {
		t.CreatedAt, _ = s.CreatedAt.GetNull()
	}
	if !s.UpdatedAt.IsUnset() {
		t.UpdatedAt, _ = s.UpdatedAt.GetNull()
	}
	if !s.DeletedAt.IsUnset() {
		t.DeletedAt, _ = s.DeletedAt.GetNull()
	}
}

func (s *PatientSetter) Apply(q *dialect.InsertQuery) {
	q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
		return Patients.BeforeInsertHooks.RunHooks(ctx, exec, s)
	})

	if len(q.Table.Columns) == 0 {
		q.Table.Columns = s.SetColumns()
	}

	q.AppendValues(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		vals := make([]bob.Expression, 0, 16)
		if !s.ID.IsUnset() {
			vals = append(vals, sqlite.Arg(s.ID))
		}

		if !s.Type.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Type))
		}

		if !s.Name.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Name))
		}

		if !s.Gender.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Gender))
		}

		if !s.Age.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Age))
		}

		if !s.ChipID.IsUnset() {
			vals = append(vals, sqlite.Arg(s.ChipID))
		}

		if !s.Weight.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Weight))
		}

		if !s.Castrated.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Castrated))
		}

		if !s.Note.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Note))
		}

		if !s.Owner.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Owner))
		}

		if !s.OwnerPhone.IsUnset() {
			vals = append(vals, sqlite.Arg(s.OwnerPhone))
		}

		if !s.Folder.IsUnset() {
			vals = append(vals, sqlite.Arg(s.Folder))
		}

		if !s.IndexFolder.IsUnset() {
			vals = append(vals, sqlite.Arg(s.IndexFolder))
		}

		if !s.CreatedAt.IsUnset() {
			vals = append(vals, sqlite.Arg(s.CreatedAt))
		}

		if !s.UpdatedAt.IsUnset() {
			vals = append(vals, sqlite.Arg(s.UpdatedAt))
		}

		if !s.DeletedAt.IsUnset() {
			vals = append(vals, sqlite.Arg(s.DeletedAt))
		}

		return bob.ExpressSlice(ctx, w, d, start, vals, "", ", ", "")
	}))
}

func (s PatientSetter) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return um.Set(s.Expressions()...)
}

func (s PatientSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 16)

	if !s.ID.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "id")...),
			sqlite.Arg(s.ID),
		}})
	}

	if !s.Type.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "type")...),
			sqlite.Arg(s.Type),
		}})
	}

	if !s.Name.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "name")...),
			sqlite.Arg(s.Name),
		}})
	}

	if !s.Gender.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "gender")...),
			sqlite.Arg(s.Gender),
		}})
	}

	if !s.Age.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "age")...),
			sqlite.Arg(s.Age),
		}})
	}

	if !s.ChipID.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "chip_id")...),
			sqlite.Arg(s.ChipID),
		}})
	}

	if !s.Weight.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "weight")...),
			sqlite.Arg(s.Weight),
		}})
	}

	if !s.Castrated.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "castrated")...),
			sqlite.Arg(s.Castrated),
		}})
	}

	if !s.Note.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "note")...),
			sqlite.Arg(s.Note),
		}})
	}

	if !s.Owner.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "owner")...),
			sqlite.Arg(s.Owner),
		}})
	}

	if !s.OwnerPhone.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "owner_phone")...),
			sqlite.Arg(s.OwnerPhone),
		}})
	}

	if !s.Folder.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "folder")...),
			sqlite.Arg(s.Folder),
		}})
	}

	if !s.IndexFolder.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "index_folder")...),
			sqlite.Arg(s.IndexFolder),
		}})
	}

	if !s.CreatedAt.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "created_at")...),
			sqlite.Arg(s.CreatedAt),
		}})
	}

	if !s.UpdatedAt.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "updated_at")...),
			sqlite.Arg(s.UpdatedAt),
		}})
	}

	if !s.DeletedAt.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "deleted_at")...),
			sqlite.Arg(s.DeletedAt),
		}})
	}

	return exprs
}

// FindPatient retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindPatient(ctx context.Context, exec bob.Executor, IDPK int32, cols ...string) (*Patient, error) {
	if len(cols) == 0 {
		return Patients.Query(
			SelectWhere.Patients.ID.EQ(IDPK),
		).One(ctx, exec)
	}

	return Patients.Query(
		SelectWhere.Patients.ID.EQ(IDPK),
		sm.Columns(Patients.Columns().Only(cols...)),
	).One(ctx, exec)
}

// PatientExists checks the presence of a single record by primary key
func PatientExists(ctx context.Context, exec bob.Executor, IDPK int32) (bool, error) {
	return Patients.Query(
		SelectWhere.Patients.ID.EQ(IDPK),
	).Exists(ctx, exec)
}

// AfterQueryHook is called after Patient is retrieved from the database
func (o *Patient) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = Patients.AfterSelectHooks.RunHooks(ctx, exec, PatientSlice{o})
	case bob.QueryTypeInsert:
		ctx, err = Patients.AfterInsertHooks.RunHooks(ctx, exec, PatientSlice{o})
	case bob.QueryTypeUpdate:
		ctx, err = Patients.AfterUpdateHooks.RunHooks(ctx, exec, PatientSlice{o})
	case bob.QueryTypeDelete:
		ctx, err = Patients.AfterDeleteHooks.RunHooks(ctx, exec, PatientSlice{o})
	}

	return err
}

// PrimaryKeyVals returns the primary key values of the Patient
func (o *Patient) PrimaryKeyVals() bob.Expression {
	return sqlite.Arg(o.ID)
}

func (o *Patient) pkEQ() dialect.Expression {
	return sqlite.Quote("patients", "id").EQ(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		return o.PrimaryKeyVals().WriteSQL(ctx, w, d, start)
	}))
}

// Update uses an executor to update the Patient
func (o *Patient) Update(ctx context.Context, exec bob.Executor, s *PatientSetter) error {
	v, err := Patients.Update(s.UpdateMod(), um.Where(o.pkEQ())).One(ctx, exec)
	if err != nil {
		return err
	}

	o.R = v.R
	*o = *v

	return nil
}

// Delete deletes a single Patient record with an executor
func (o *Patient) Delete(ctx context.Context, exec bob.Executor) error {
	_, err := Patients.Delete(dm.Where(o.pkEQ())).Exec(ctx, exec)
	return err
}

// Reload refreshes the Patient using the executor
func (o *Patient) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Patients.Query(
		SelectWhere.Patients.ID.EQ(o.ID),
	).One(ctx, exec)
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

// AfterQueryHook is called after PatientSlice is retrieved from the database
func (o PatientSlice) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = Patients.AfterSelectHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeInsert:
		ctx, err = Patients.AfterInsertHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeUpdate:
		ctx, err = Patients.AfterUpdateHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeDelete:
		ctx, err = Patients.AfterDeleteHooks.RunHooks(ctx, exec, o)
	}

	return err
}

func (o PatientSlice) pkIN() dialect.Expression {
	if len(o) == 0 {
		return sqlite.Raw("NULL")
	}

	return sqlite.Quote("patients", "id").In(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		pkPairs := make([]bob.Expression, len(o))
		for i, row := range o {
			pkPairs[i] = row.PrimaryKeyVals()
		}
		return bob.ExpressSlice(ctx, w, d, start, pkPairs, "", ", ", "")
	}))
}

// copyMatchingRows finds models in the given slice that have the same primary key
// then it first copies the existing relationships from the old model to the new model
// and then replaces the old model in the slice with the new model
func (o PatientSlice) copyMatchingRows(from ...*Patient) {
	for i, old := range o {
		for _, new := range from {
			if new.ID != old.ID {
				continue
			}
			new.R = old.R
			o[i] = new
			break
		}
	}
}

// UpdateMod modifies an update query with "WHERE primary_key IN (o...)"
func (o PatientSlice) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return bob.ModFunc[*dialect.UpdateQuery](func(q *dialect.UpdateQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return Patients.BeforeUpdateHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *Patient:
				o.copyMatchingRows(retrieved)
			case []*Patient:
				o.copyMatchingRows(retrieved...)
			case PatientSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a Patient or a slice of Patient
				// then run the AfterUpdateHooks on the slice
				_, err = Patients.AfterUpdateHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

// DeleteMod modifies an delete query with "WHERE primary_key IN (o...)"
func (o PatientSlice) DeleteMod() bob.Mod[*dialect.DeleteQuery] {
	return bob.ModFunc[*dialect.DeleteQuery](func(q *dialect.DeleteQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return Patients.BeforeDeleteHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *Patient:
				o.copyMatchingRows(retrieved)
			case []*Patient:
				o.copyMatchingRows(retrieved...)
			case PatientSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a Patient or a slice of Patient
				// then run the AfterDeleteHooks on the slice
				_, err = Patients.AfterDeleteHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

func (o PatientSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals PatientSetter) error {
	if len(o) == 0 {
		return nil
	}

	_, err := Patients.Update(vals.UpdateMod(), o.UpdateMod()).All(ctx, exec)
	return err
}

func (o PatientSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	_, err := Patients.Delete(o.DeleteMod()).Exec(ctx, exec)
	return err
}

func (o PatientSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	o2, err := Patients.Query(sm.Where(o.pkIN())).All(ctx, exec)
	if err != nil {
		return err
	}

	o.copyMatchingRows(o2...)

	return nil
}

type patientJoins[Q dialect.Joinable] struct {
	typ        string
	Procedures func(context.Context) modAs[Q, procedureColumns]
}

func (j patientJoins[Q]) aliasedAs(alias string) patientJoins[Q] {
	return buildPatientJoins[Q](buildPatientColumns(alias), j.typ)
}

func buildPatientJoins[Q dialect.Joinable](cols patientColumns, typ string) patientJoins[Q] {
	return patientJoins[Q]{
		typ:        typ,
		Procedures: patientsJoinProcedures[Q](cols, typ),
	}
}

func patientsJoinProcedures[Q dialect.Joinable](from patientColumns, typ string) func(context.Context) modAs[Q, procedureColumns] {
	return func(ctx context.Context) modAs[Q, procedureColumns] {
		return modAs[Q, procedureColumns]{
			c: ProcedureColumns,
			f: func(to procedureColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Procedures.Name().As(to.Alias())).On(
						to.PatientID.EQ(from.ID),
					))
				}

				return mods
			},
		}
	}
}

// Procedures starts a query for related objects on procedures
func (o *Patient) Procedures(mods ...bob.Mod[*dialect.SelectQuery]) ProceduresQuery {
	return Procedures.Query(append(mods,
		sm.Where(ProcedureColumns.PatientID.EQ(sqlite.Arg(o.ID))),
	)...)
}

func (os PatientSlice) Procedures(mods ...bob.Mod[*dialect.SelectQuery]) ProceduresQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = sqlite.ArgGroup(o.ID)
	}

	return Procedures.Query(append(mods,
		sm.Where(sqlite.Group(ProcedureColumns.PatientID).In(PKArgs...)),
	)...)
}

func (o *Patient) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "Procedures":
		rels, ok := retrieved.(ProcedureSlice)
		if !ok {
			return fmt.Errorf("patient cannot load %T as %q", retrieved, name)
		}

		o.R.Procedures = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.Patient = o
			}
		}
		return nil
	default:
		return fmt.Errorf("patient has no relationship %q", name)
	}
}

func ThenLoadPatientProcedures(queryMods ...bob.Mod[*dialect.SelectQuery]) sqlite.Loader {
	return sqlite.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadPatientProcedures(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load PatientProcedures", retrieved)
		}

		err := loader.LoadPatientProcedures(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadPatientProcedures loads the patient's Procedures into the .R struct
func (o *Patient) LoadPatientProcedures(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.Procedures = nil

	related, err := o.Procedures(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.Patient = o
	}

	o.R.Procedures = related
	return nil
}

// LoadPatientProcedures loads the patient's Procedures into the .R struct
func (os PatientSlice) LoadPatientProcedures(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	procedures, err := os.Procedures(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.Procedures = nil
	}

	for _, o := range os {
		for _, rel := range procedures {
			if o.ID != rel.PatientID.GetOrZero() {
				continue
			}

			rel.R.Patient = o

			o.R.Procedures = append(o.R.Procedures, rel)
		}
	}

	return nil
}

func insertPatientProcedures0(ctx context.Context, exec bob.Executor, procedures1 []*ProcedureSetter, patient0 *Patient) (ProcedureSlice, error) {
	for i := range procedures1 {
		procedures1[i].PatientID = omitnull.From(patient0.ID)
	}

	ret, err := Procedures.Insert(bob.ToMods(procedures1...)).All(ctx, exec)
	if err != nil {
		return ret, fmt.Errorf("insertPatientProcedures0: %w", err)
	}

	return ret, nil
}

func attachPatientProcedures0(ctx context.Context, exec bob.Executor, count int, procedures1 ProcedureSlice, patient0 *Patient) (ProcedureSlice, error) {
	setter := &ProcedureSetter{
		PatientID: omitnull.From(patient0.ID),
	}

	err := procedures1.UpdateAll(ctx, exec, *setter)
	if err != nil {
		return nil, fmt.Errorf("attachPatientProcedures0: %w", err)
	}

	return procedures1, nil
}

func (patient0 *Patient) InsertProcedures(ctx context.Context, exec bob.Executor, related ...*ProcedureSetter) error {
	if len(related) == 0 {
		return nil
	}

	var err error

	procedures1, err := insertPatientProcedures0(ctx, exec, related, patient0)
	if err != nil {
		return err
	}

	patient0.R.Procedures = append(patient0.R.Procedures, procedures1...)

	for _, rel := range procedures1 {
		rel.R.Patient = patient0
	}
	return nil
}

func (patient0 *Patient) AttachProcedures(ctx context.Context, exec bob.Executor, related ...*Procedure) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	procedures1 := ProcedureSlice(related)

	_, err = attachPatientProcedures0(ctx, exec, len(related), procedures1, patient0)
	if err != nil {
		return err
	}

	patient0.R.Procedures = append(patient0.R.Procedures, procedures1...)

	for _, rel := range related {
		rel.R.Patient = patient0
	}

	return nil
}
