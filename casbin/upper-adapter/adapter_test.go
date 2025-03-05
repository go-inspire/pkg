package upperadapter

import (
	"fmt"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"log"
	"strings"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
)

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	t.Helper()
	myRes, _ := e.GetPolicy()
	log.Print("Policy: ", myRes)

	m := make(map[string]bool, len(res))
	for _, value := range res {
		key := strings.Join(value, ",")
		m[key] = true
	}

	for _, value := range myRes {
		key := strings.Join(value, ",")
		if !m[key] {
			t.Error("Policy: ", myRes, ", supposed to be ", res)
			break
		}
	}
}

func initPolicy(t *testing.T, a *Adapter) {
	// Because the db is empty at first,
	// so we need to load the policy from the file adapter (.CSV) first.
	e, err := casbin.NewEnforcer("internal/rbac_model.conf", "internal/rbac_policy.csv")

	// This is a trick to save the current policy to the db.
	// We can't call e.SavePolicy() because the adapter in the enforcer is still the file adapter.
	// The current policy means the policy in the Casbin enforcer (aka in memory).
	err = a.SavePolicy(e.GetModel())
	if err != nil {
		panic(err)
	}

	// Clear the current policy.
	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// Load the policy from db.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		panic(err)
	}
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testSaveLoad(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf", a)
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testAutoSave(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf", a)

	// AutoSave is enabled by default.
	// Now we disable it.
	e.EnableAutoSave(false)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
	// it doesn't affect the policy in the storage.
	_, err = e.AddPolicy("alice", "data1", "write")
	logErr("AddPolicy")
	// Reload the policy from the storage to see the effect.
	err = e.LoadPolicy()
	logErr("LoadPolicy")
	// This is still the original policy.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Now we enable the AutoSave.
	e.EnableAutoSave(true)

	// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
	// but also affects the policy in the storage.
	_, err = e.AddPolicy("alice", "data1", "write")
	logErr("AddPolicy2")
	// Reload the policy from the storage to see the effect.
	err = e.LoadPolicy()
	logErr("LoadPolicy2")
	// The policy has a new rule: {"alice", "data1", "write"}.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"alice", "data1", "write"}})

	// Remove the added rule.
	_, err = e.RemovePolicy("alice", "data1", "write")
	logErr("RemovePolicy")
	err = e.LoadPolicy()
	logErr("LoadPolicy3")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Remove "data2_admin" related policy rules via a filter.
	// Two rules: {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"} are deleted.
	_, err = e.RemoveFilteredPolicy(0, "data2_admin")
	logErr("RemoveFilteredPolicy")
	err = e.LoadPolicy()
	logErr("LoadPolicy4")

	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
}

func testFilteredPolicy(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf")
	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	// Load only alice's policies
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"alice"}})
	logErr("LoadFilteredPolicy")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}})

	// Load only bob's policies
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"bob"}})
	logErr("LoadFilteredPolicy2")
	testGetPolicy(t, e, [][]string{{"bob", "data2", "write"}})

	// Load policies for data2_admin
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"data2_admin"}})
	logErr("LoadFilteredPolicy3")
	testGetPolicy(t, e, [][]string{{"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Load policies for alice and bob
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"alice", "bob"}})
	logErr("LoadFilteredPolicy4")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
}

func testRemovePolicies(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	err = a.AddPolicies("p", "p", [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}, {"max", "data1", "delete"}})
	logErr("AddPolicies")

	// Load policies for max
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"max"}})
	logErr("LoadFilteredPolicy")

	testGetPolicy(t, e, [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}, {"max", "data1", "delete"}})

	// Remove policies
	err = a.RemovePolicies("p", "p", [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}})
	logErr("RemovePolicies")

	// Reload policies for max
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"max"}})
	logErr("LoadFilteredPolicy2")

	testGetPolicy(t, e, [][]string{{"max", "data1", "delete"}})
}

func testAddPolicies(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	err = a.AddPolicies("p", "p", [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}})
	logErr("AddPolicies")

	// Load policies for max
	err = e.LoadFilteredPolicy(&Filter{V0: []string{"max"}})
	logErr("LoadFilteredPolicy")

	testGetPolicy(t, e, [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}})
}

func testUpdatePolicies(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	err = a.UpdatePolicy("p", "p", []string{"bob", "data2", "write"}, []string{"alice", "data2", "write"})
	logErr("UpdatePolicy")

	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"alice", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	err = a.UpdatePolicies("p", "p", [][]string{{"alice", "data1", "read"}, {"alice", "data2", "write"}}, [][]string{{"bob", "data1", "read"}, {"bob", "data2", "write"}})
	logErr("UpdatePolicies")

	testGetPolicy(t, e, [][]string{{"bob", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testUpdateFilteredPolicies(t *testing.T, a *Adapter) {
	// Initialize some policy in db.
	initPolicy(t, a)
	// Note: you don't need to look at the above code
	// if you already have a working db with policy inside.

	// Now the db has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	e, _ := casbin.NewEnforcer("internal/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	e.UpdateFilteredPolicies([][]string{{"alice", "data1", "write"}}, 0, "alice", "data1", "read")
	e.UpdateFilteredPolicies([][]string{{"bob", "data2", "read"}}, 0, "bob", "data2", "write")
	e.LoadPolicy()
	testGetPolicyWithoutOrder(t, e, [][]string{{"alice", "data1", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"bob", "data2", "read"}})
}

func testGetPolicyWithoutOrder(t *testing.T, e *casbin.Enforcer, res [][]string) {
	myRes, _ := e.GetPolicy()
	log.Print("Policy: ", myRes)

	if !arrayEqualsWithoutOrder(myRes, res) {
		t.Error("Policy: ", myRes, ", supposed to be ", res)
	}
}

func arrayEqualsWithoutOrder(a [][]string, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}

	mapA := make(map[int]string)
	mapB := make(map[int]string)
	order := make(map[int]struct{})
	l := len(a)

	for i := 0; i < l; i++ {
		mapA[i] = util.ArrayToString(a[i])
		mapB[i] = util.ArrayToString(b[i])
	}

	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			if _, ok := order[j]; ok {
				if j == l-1 {
					return false
				} else {
					continue
				}
			}
			if mapA[i] == mapB[j] {
				order[j] = struct{}{}
				break
			} else if j == l-1 {
				return false
			}
		}
	}
	return true
}

func TestAdapters(t *testing.T) {
	db.LC().SetLevel(db.LogLevelDebug)

	//conn, err := mysql.ParseURL("jquant:jquant123@tcp(172.16.2.154:3306)/jquant")
	conn, err := mysql.ParseURL("dev:dev123@tcp(localhost:3306)/dev")
	session, err := mysql.Open(conn)
	if err != nil {
		log.Fatalf("session.Open(): %q\n", err)
	}

	defer func(session db.Session) {
		_ = session.Close()
	}(session)

	schema := fmt.Sprintf(createTableSQL, defaultTableName, defaultTableName)
	_, _ = session.SQL().Exec(schema)

	a, err := NewAdapter(session, "casbin_rule")

	t.Run("testSaveLoad", func(t *testing.T) {
		testSaveLoad(t, a)
	})
	t.Run("testAutoSave", func(t *testing.T) {
		testAutoSave(t, a)
	})
	t.Run("testFilteredPolicy", func(t *testing.T) {
		testFilteredPolicy(t, a)
	})
	t.Run("testAddPolicies", func(t *testing.T) {
		testAddPolicies(t, a)
	})
	t.Run("testRemovePolicies", func(t *testing.T) {
		testRemovePolicies(t, a)
	})
	t.Run("testUpdatePolicies", func(t *testing.T) {
		testUpdatePolicies(t, a)
	})
	t.Run("testUpdateFilteredPolicies", func(t *testing.T) {
		testUpdateFilteredPolicies(t, a)
	})
}
