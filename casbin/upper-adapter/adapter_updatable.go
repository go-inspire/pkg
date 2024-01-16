package upperadapter

import (
	"github.com/pkg/errors"
	"github.com/upper/db/v4"
	"strconv"
)

// UpdatePolicy update a policy rule from storage.
// This is part of the Auto-Save feature.
func (a *Adapter) UpdatePolicy(sec, ptype string, oldRule, newRule []string) error {
	return a.updatePolicy(a.db, sec, ptype, oldRule, newRule)
}

func (a *Adapter) updatePolicy(session db.Session, sec, ptype string, oldRule, newRule []string) error {
	cond := a.genCond(ptype, oldRule)
	args := a.genUpdateArgs(ptype, newRule)
	return session.Collection(a.tableName).Find(cond).Update(args)
}

// UpdatePolicies updates policy cols to storage.
func (a *Adapter) UpdatePolicies(sec, ptype string, oldRules, newRules [][]string) (err error) {
	if len(oldRules) != len(newRules) {
		return errors.New("old cols size not equal to new cols size")
	}

	return a.db.Tx(func(tx db.Session) error {
		for idx := range oldRules {
			if err = a.updatePolicy(a.db, sec, ptype, oldRules[idx], newRules[idx]); err != nil {
				return err
			}
		}
		return nil
	})

}

// UpdateFilteredPolicies deletes old cols and adds new cols.
func (a *Adapter) UpdateFilteredPolicies(sec, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) (oldPolicies [][]string, err error) {
	l := fieldIndex + len(fieldValues)
	cond := db.And(db.Cond{"p_type": ptype})
	for idx := 0; idx < 6; idx++ {
		if fieldIndex <= idx && idx < l {
			value := fieldValues[idx-fieldIndex]
			if value != "" {
				cond = cond.And(db.Cond{"v" + strconv.Itoa(idx): value})
			}
		}
	}

	var oldRules []*CasbinRule
	res := a.cols.Find(cond)
	if err = res.All(&oldRules); err != nil {
		err = errors.Wrap(err, "load old policies")
		return
	}

	if err = res.Delete(); err != nil {
		err = errors.Wrap(err, "delete old policies")
		return
	}

	newPolicies := make([][]interface{}, 0, len(newRules))
	for _, rule := range newRules {
		arg := a.genArgs(ptype, rule)
		newPolicies = append(newPolicies, arg)
	}

	if err = a.insertRules(newPolicies); err != nil {
		err = errors.Wrap(err, "insert new policies")
		return
	}

	oldPolicies = make([][]string, 0, len(oldRules))
	for _, rule := range oldRules {
		oldRule := []string{rule.PType, rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5}
		oldPolicy := make([]string, 0, len(oldRule))
		for _, val := range oldRule {
			if val == "" {
				break
			}
			oldPolicy = append(oldPolicy, val)
		}
		oldPolicies = append(oldPolicies, oldPolicy)
	}

	return

}
