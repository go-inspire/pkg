package upperadapter

import "github.com/upper/db/v4"

// AddPolicies  add multiple policy cols to the storage. 实现 persist.BatchAdapter  接口
func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	args := make([][]interface{}, 0, 8)

	for _, rule := range rules {
		arg := a.genArgs(ptype, rule)
		args = append(args, arg)
	}

	return a.insertRules(args)
}

// RemovePolicies  remove policy cols. 实现 persist.BatchAdapter  接口
func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) (err error) {
	return a.db.Tx(func(tx db.Session) error {
		for _, rule := range rules {
			cond := a.genCond(ptype, rule)
			if err = tx.Collection(a.tableName).Find(cond).Delete(); err != nil {
				return err
			}
		}
		return nil
	})
}
