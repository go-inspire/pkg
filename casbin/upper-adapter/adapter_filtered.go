package upperadapter

import (
	"github.com/casbin/casbin/v2/model"
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
	"reflect"
)

// Filter  defines the filtering cols for a FilteredAdapter's policy.
// Empty values are ignored, but all others must match the filter.
type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

// LoadFilteredPolicy  load policy cols that match the filter.
// filterPtr must be a pointer.
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	if filter == nil {
		return a.LoadPolicy(model)
	}

	filterValue, ok := filter.(*Filter)
	if !ok {
		return errors.Errorf("filter must type of *upperadapter.Filter but invalid filter type %s", reflect.TypeOf(filter).String())
	}

	lines, err := a.selectWhereIn(filterValue)
	if err != nil {
		return err
	}

	for _, line := range lines {
		a.loadPolicyLine(line, model)
	}

	a.filtered = true

	return nil
}

// IsFiltered  returns true if the loaded policy cols has been filtered.
func (a *Adapter) IsFiltered() bool {
	return a.filtered
}

// selectWhereIn  select eligible data by filter from the table.
func (a *Adapter) selectWhereIn(filter *Filter) (lines []*CasbinRule, err error) {
	cond := db.And()
	for _, col := range [maxParamLength]struct {
		name string
		arg  []string
	}{
		{"p_type", filter.PType},
		{"v0", filter.V0},
		{"v1", filter.V1},
		{"v2", filter.V2},
		{"v3", filter.V3},
		{"v4", filter.V4},
		{"v5", filter.V5},
	} {
		l := len(col.arg)
		if l == 0 {
			continue
		}

		if l == 1 {
			cond = cond.And(db.Cond{col.name: col.arg[0]})
		} else {
			cond = cond.And(db.Cond{col.name: db.In(ToInterfaceSlice(col.arg)...)})
		}
	}

	err = a.cols.Find(cond).All(&lines)
	return lines, err
}

func ToInterfaceSlice[T interface{}](slice []T) []interface{} {
	res := make([]interface{}, len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res
}
