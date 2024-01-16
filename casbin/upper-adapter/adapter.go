package upperadapter

import (
	"bytes"
	"context"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
	"strconv"
)

// defaultTableName  if tableName == "", the Adapter will use this default table name.
const defaultTableName = "casbin_rule"

const createTableSQL = `
CREATE TABLE IF NOT EXISTS %s(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL,
    INDEX idx_%s (p_type,v0,v1)
) ;
`

// maxParamLength  .
const maxParamLength = 7

// CasbinRule  defines the casbin rule model.
// It used for save or load policy lines from sqlx connected database.
type CasbinRule struct {
	PType string `db:"p_type"`
	V0    string `db:"v0"`
	V1    string `db:"v1"`
	V2    string `db:"v2"`
	V3    string `db:"v3"`
	V4    string `db:"v4"`
	V5    string `db:"v5"`
}

// Adapter  define the sqlx upper-adapter for Casbin.
// It can load policy lines or save policy lines from sqlx connected database.
type Adapter struct {
	ctx       context.Context
	tableName string
	db        db.Session
	cols      db.Collection

	filtered bool
}

// NewAdapter  the constructor for Adapter.
// db should connected to database and controlled by user.
// If tableName == "", the Adapter will automatically create a table named "casbin_rule".
func NewAdapter(db db.Session, tableName string) (*Adapter, error) {
	return NewAdapterContext(context.Background(), db, tableName)
}

// NewAdapterContext  the constructor for Adapter.
// db should connected to database and controlled by user.
// If tableName == "", the Adapter will automatically create a table named "casbin_rule".
func NewAdapterContext(ctx context.Context, db db.Session, tableName string) (*Adapter, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	// check db connecting
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	if tableName == "" {
		tableName = defaultTableName
	}

	schema := fmt.Sprintf(createTableSQL, tableName, tableName)
	_, err = db.SQL().Exec(schema)

	adapter := Adapter{
		ctx:       ctx,
		tableName: tableName,
		db:        db,
		cols:      db.Collection(tableName),
		filtered:  true,
	}

	return &adapter, nil
}

// truncate  clear the table.
func (a *Adapter) truncate() error {
	return a.cols.Truncate()
}

// deleteAll  clear the table.
func (a *Adapter) deleteAll() error {
	return a.cols.Find().Delete()
}

// isTableExist  check the table exists.
func (a *Adapter) isTableExist() bool {
	exists, err := a.cols.Exists()
	return exists && err == nil
}

// truncateAndInsertRules  clear table and insert new rows.
func (a *Adapter) truncateAndInsertRules(rules [][]interface{}) error {
	if err := a.truncate(); err != nil {
		return err
	}
	return a.insertRules(rules)
}

// deleteAllAndInsertRules  clear table and insert new rows.
func (a *Adapter) deleteAllAndInsertRules(rules [][]interface{}) error {
	if err := a.deleteAll(); err != nil {
		return err
	}

	return a.insertRules(rules)
}

func (a *Adapter) insertRules(rules [][]interface{}) error {
	//PType string `db:"p_type"`
	//V0    string `db:"v0"`
	//V1    string `db:"v1"`
	//V2    string `db:"v2"`
	//V3    string `db:"v3"`
	//V4    string `db:"v4"`
	//V5    string `db:"v5"`
	return a.db.Tx(func(tx db.Session) error {
		batch := tx.SQL().InsertInto(a.tableName).
			Columns("p_type", "v0", "v1", "v2", "v3", "v4", "v5").
			Batch(len(rules))
		for _, rule := range rules {
			batch.Values(rule...)
		}
		batch.Done()
		return batch.Wait()
	})

}

// LoadPolicy  load all policy cols from the storage.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []*CasbinRule
	err := a.cols.Find().All(&lines)
	if err != nil {
		return err
	}

	for _, line := range lines {
		a.loadPolicyLine(line, model)
	}

	return nil
}

// SavePolicy  save policy cols to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	args := make([][]interface{}, 0, 64)

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			arg := a.genArgs(ptype, rule)
			args = append(args, arg)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			arg := a.genArgs(ptype, rule)
			args = append(args, arg)
		}
	}

	return a.deleteAllAndInsertRules(args)
}

// AddPolicy  add one policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	args := a.genArgs(ptype, rule)

	return a.insertRules([][]interface{}{args})
}

// RemovePolicy  remove policy cols from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	cond := a.genCond(ptype, rule)
	return a.cols.Find(cond).Delete()

}

// RemoveFilteredPolicy  remove policy cols that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(_ string, ptype string, fieldIndex int, fieldValues ...string) error {
	cond := db.And(db.Cond{"p_type": ptype})

	l := fieldIndex + len(fieldValues)
	for idx := 0; idx < 6; idx++ {
		if fieldIndex <= idx && idx < l {
			value := fieldValues[idx-fieldIndex]
			if value != "" {
				cond = cond.And(db.Cond{"v" + strconv.Itoa(idx): value})
			}
		}
	}
	return a.cols.Find(cond).Delete()
}

// loadPolicyLine  load a policy line to model.
func (*Adapter) loadPolicyLine(line *CasbinRule, model model.Model) {
	if line == nil {
		return
	}

	var lineBuf bytes.Buffer

	lineBuf.Grow(64)
	lineBuf.WriteString(line.PType)

	args := [6]string{line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}
	for _, arg := range args {
		if arg != "" {
			lineBuf.WriteByte(',')
			lineBuf.WriteString(arg)
		}
	}

	_ = persist.LoadPolicyLine(lineBuf.String(), model)
}

// genArgs  generate args from ptype and rule.
func (*Adapter) genArgs(ptype string, rule []string) []interface{} {
	l := len(rule)

	args := make([]interface{}, maxParamLength)
	args[0] = ptype

	for idx := 0; idx < l; idx++ {
		args[idx+1] = rule[idx]
	}

	for idx := l + 1; idx < maxParamLength; idx++ {
		args[idx] = ""
	}

	return args
}

// genUpdateArgs  generate args from ptype and rule.
func (*Adapter) genUpdateArgs(ptype string, rule []string) map[string]interface{} {
	l := len(rule)

	args := make(map[string]interface{}, maxParamLength)
	args["p_type"] = ptype

	for idx := 0; idx < l; idx++ {
		args["v"+strconv.Itoa(idx)] = rule[idx]
	}

	for idx := l + 1; idx < maxParamLength-1; idx++ {
		args["v"+strconv.Itoa(idx)] = ""
	}

	return args
}

// genCond  generate cond from ptype and rule.
func (*Adapter) genCond(ptype string, rule []string) *db.AndExpr {
	cond := db.And(db.Cond{"p_type": ptype})

	for idx, arg := range rule {
		if arg != "" {
			cond = cond.And(db.Cond{"v" + strconv.Itoa(idx): arg})
		}
	}

	return cond
}
