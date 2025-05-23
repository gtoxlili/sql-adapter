// Copyright 2020 by Blank-Xu. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqladaptertest

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/casbin/casbin/v2"
	. "github.com/gtoxlili/sql-adapter"
)

const (
	testRbacModelFile  = "../testdata/rbac_model.conf"
	testRbacPolicyFile = "../testdata/rbac_policy.csv"
)

var testDefaultPolicy = [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}}

func TestAdapter(t *testing.T) {
	for driverName, db := range testDBs {
		t.Logf("adapter test for [%s] start", driverName)

		testTableName(t, db, driverName)
		testSaveLoad(t, db, driverName, "sqladapter_test_save_load")
		testAutoSave(t, db, driverName, "sqladapter_test_auto_save")
		testFilteredPolicy(t, db, driverName, "sqladapter_test_filtered_policy")
		testUpdatePolicy(t, db, driverName, "sqladapter_test_update_policy")
		testUpdatePolicies(t, db, driverName, "sqladapter_test_update_policies")
		testUpdateFilteredPolicies(t, db, driverName, "sqladapter_test_update_filtered_policies")

		t.Logf("adapter test for [%s] finished", driverName)
	}
}

func initPolicy(t *testing.T, db *sql.DB, driverName, tableName string) {
	// Because the DB is empty at first,
	// so we need to load the policy from the file adapter (.CSV) first.
	e, err := casbin.NewEnforcer(testRbacModelFile, testRbacPolicyFile)
	if err != nil {
		t.Fatal("casbin NewEnforcer failed, err: ", err)
	}

	// create sqladapter
	a, err := NewAdapter(db, driverName, tableName)
	if err != nil {
		t.Fatal("sqladapter NewAdapter failed, err: ", err)
	}

	// This is a trick to save the current policy to the DB.
	// We can't call e.SavePolicy() because the adapter in the enforcer is still the file adapter.
	// The current policy means the policy in the Casbin enforcer (aka in memory).
	err = a.SavePolicy(e.GetModel())
	if err != nil {
		t.Fatal("sqladapter SavePolicy failed, err: ", err)
	}

	// clear current policy.
	e.ClearPolicy()
	policies, err := e.GetPolicy()
	validateNilError(t, err)
	validatePolicies(t, policies, [][]string{})

	// load policy from DB.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		t.Fatal("sqladapter LoadPolicy failed, err: ", err)
	}

	policies, err = e.GetPolicy()
	validateNilError(t, err)
	validatePolicies(t, policies, testDefaultPolicy)
}

func testTableName(t *testing.T, db *sql.DB, driverName string) {
	tests := []struct {
		name      string
		tableName string
	}{
		{"01_empty_name", ""},
		{"02_test_name", "test_name"},
	}

	for _, tt := range tests {
		t.Run("TableName_"+tt.name, func(t *testing.T) {
			if _, err := NewAdapter(db, driverName, tt.tableName); err != nil {
				t.Errorf("test [%s] in data source [%s] failed, err: %v", tt.name, driverName, err)
			}
		})
	}
}

func testSaveLoad(t *testing.T, db *sql.DB, driverName, tableName string) {
	t.Run("Save_Load", func(t *testing.T) {
		initPolicy(t, db, driverName, tableName)

		a, _ := NewAdapter(db, driverName, tableName)
		e, _ := casbin.NewEnforcer(testRbacModelFile, a)
		policies, err := e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, testDefaultPolicy)
	})
}

func testAutoSave(t *testing.T, db *sql.DB, driverName, tableName string) {
	const testName = "AutoSave_"

	var err error
	t.Run(testName+"01_not_EnableAutoSave", func(t *testing.T) {
		initPolicy(t, db, driverName, tableName)

		a, _ := NewAdapter(db, driverName, tableName)
		e, _ := casbin.NewEnforcer(testRbacModelFile, a)

		// AutoSave is enabled by default.
		// Now we disable it.
		e.EnableAutoSave(false)

		// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
		// it doesn't affect the policy in the storage.
		if _, err = e.AddPolicy("alice", "data1", "write"); err != nil {
			t.Errorf("%s test failed, err: %v", "AddPolicy", err)
		}
		// Reload the policy from the storage to see the effect.
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err := e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, testDefaultPolicy)

		if _, err = e.AddPolicies([][]string{{"alice_1", "data_1", "read_1"}, {"bob_1", "data_1", "write_1"}}); err != nil {
			t.Errorf("%s test failed, err: %v", "AddPolicies", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err = e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, testDefaultPolicy)
	})

	t.Run(testName+"02_EnableAutoSave", func(t *testing.T) {
		initPolicy(t, db, driverName, tableName)

		a, _ := NewAdapter(db, driverName, tableName)
		e, _ := casbin.NewEnforcer(testRbacModelFile, a)

		// Now we enable the AutoSave.
		e.EnableAutoSave(true)

		// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
		// but also affects the policy in the storage.
		if _, err = e.AddPolicy("alice", "data1", "write"); err != nil {
			t.Errorf("%s test failed, err: %v", "AddPolicy", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		// The policy has a new rule: {"alice", "data1", "write"}.
		policies, err := e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, append(testDefaultPolicy, []string{"alice", "data1", "write"}))

		if _, err = e.AddPolicies([][]string{{"alice_2", "data_2", "read_2"}, {"bob_2", "data_2", "write_2"}}); err != nil {
			t.Errorf("%s test failed, err: %v", "AddPolicies", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err = e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies,
			append(testDefaultPolicy, []string{"alice", "data1", "write"}, []string{"alice_2", "data_2", "read_2"}, []string{"bob_2", "data_2", "write_2"}))

		if _, err = e.RemovePolicies([][]string{{"alice_2", "data_2", "read_2"}, {"bob_2", "data_2", "write_2"}}); err != nil {
			t.Errorf("%s test failed, err: %v", "RemovePolicies", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err = e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, append(testDefaultPolicy, []string{"alice", "data1", "write"}))

		if _, err = e.RemovePolicy("alice", "data1", "write"); err != nil {
			t.Errorf("%s test failed, err: %v", "RemovePolicy", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err = e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, testDefaultPolicy)

		// Remove "data2_admin" related policy rules via a filter.
		// Two rules: {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"} are deleted.
		if _, err = e.RemoveFilteredPolicy(0, "data2_admin"); err != nil {
			t.Errorf("%s test failed, err: %v", "RemoveFilteredPolicy", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err = e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
	})
}

func testFilteredPolicy(t *testing.T, db *sql.DB, driverName, tableName string) {
	initPolicy(t, db, driverName, tableName)

	var err error
	a, _ := NewAdapter(db, driverName, tableName)
	e, _ := casbin.NewEnforcer(testRbacModelFile, a)
	e.SetAdapter(a)

	tests := []struct {
		name         string
		addPolicy    []interface{}
		filterPolicy *Filter
		expectPolicy [][]string
	}{
		{
			name:         "01_filter_alice",
			filterPolicy: &Filter{V0: []string{"alice"}},
			expectPolicy: [][]string{{"alice", "data1", "read"}},
		},
		{
			name:         "02_filter_bob",
			filterPolicy: &Filter{V0: []string{"bob"}},
			expectPolicy: [][]string{{"bob", "data2", "write"}},
		},
		{
			name:         "03_filter_data2_admin",
			filterPolicy: &Filter{V0: []string{"data2_admin"}},
			expectPolicy: [][]string{{"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}},
		},
		{
			name:         "04_filter_alice_bob",
			filterPolicy: &Filter{V0: []string{"alice", "bob"}},
			expectPolicy: [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}},
		},
		{
			name:      "05_filter_AddPolicy",
			addPolicy: []interface{}{"bob", "data1", "write"},
			filterPolicy: &Filter{
				PType: []string{"p"},
				V0:    []string{"bob", "data2_admin"},
				V1:    []string{"data1", "data2"},
				V2:    []string{"write"},
			},
			expectPolicy: [][]string{{"bob", "data1", "write"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "write"}},
		},
	}
	for _, tt := range tests {
		t.Run("FilteredPolicy_"+tt.name, func(t *testing.T) {
			if len(tt.addPolicy) > 0 {
				if _, err = e.AddPolicy(tt.addPolicy...); err != nil {
					t.Errorf("%s AddPolicy test failed, err: %v", tt.name, err)
				}
			}
			if err = e.LoadFilteredPolicy(tt.filterPolicy); err != nil {
				t.Errorf("%s LoadFilteredPolicy test failed, err: %v", tt.name, err)
			}
			policies, err := e.GetPolicy()
			validateNilError(t, err)
			validatePolicies(t, policies, tt.expectPolicy)
		})
	}
}

func testUpdatePolicy(t *testing.T, db *sql.DB, driverName, tableName string) {
	var err error
	t.Run("UpdatePolicy", func(t *testing.T) {
		initPolicy(t, db, driverName, tableName)

		a, _ := NewAdapter(db, driverName, tableName)
		e, _ := casbin.NewEnforcer(testRbacModelFile, a)

		e.EnableAutoSave(true)
		if _, err = e.UpdatePolicy([]string{"alice", "data1", "read"}, []string{"alice", "data1", "write"}); err != nil {
			t.Errorf("%s test failed, err: %v", "UpdatePolicy", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err := e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, [][]string{{"alice", "data1", "write"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
	})
}

func testUpdatePolicies(t *testing.T, db *sql.DB, driverName, tableName string) {
	var err error
	t.Run("UpdatePolicies", func(t *testing.T) {
		initPolicy(t, db, driverName, tableName)

		a, _ := NewAdapter(db, driverName, tableName)
		e, _ := casbin.NewEnforcer(testRbacModelFile, a)

		e.EnableAutoSave(true)
		if _, err = e.UpdatePolicies([][]string{{"alice", "data1", "write"}, {"bob", "data2", "write"}}, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "read"}}); err != nil {
			t.Errorf("%s test failed, err: %v", "UpdatePolicies", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err := e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "read"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
	})
}

func testUpdateFilteredPolicies(t *testing.T, db *sql.DB, driverName, tableName string) {
	var err error
	t.Run("UpdateFilteredPolicies", func(t *testing.T) {
		initPolicy(t, db, driverName, tableName)

		a, _ := NewAdapter(db, driverName, tableName)
		e, _ := casbin.NewEnforcer(testRbacModelFile, a)

		e.EnableAutoSave(true)
		if _, err = e.UpdateFilteredPolicies([][]string{{"alice", "data1", "write"}}, 0, "alice", "data1", "read"); err != nil {
			t.Errorf("%s test failed, err: %v", "UpdateFilteredPolicies", err)
		}
		if _, err = e.UpdateFilteredPolicies([][]string{{"bob", "data2", "read"}}, 0, "bob", "data2", "write"); err != nil {
			t.Errorf("%s test failed, err: %v", "UpdateFilteredPolicies", err)
		}
		if err = e.LoadPolicy(); err != nil {
			t.Errorf("%s test failed, err: %v", "LoadPolicy", err)
		}
		policies, err := e.GetPolicy()
		validateNilError(t, err)
		validatePolicies(t, policies, [][]string{{"alice", "data1", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"bob", "data2", "read"}})
	})
}

func validatePolicies(t *testing.T, getPolicy, wantPolicy [][]string) {
	t.Helper()

	if len(wantPolicy) != len(getPolicy) {
		t.Error("get policy: \n", getPolicy, "supposed to be: \n", wantPolicy)
		return
	}

	m := make(map[string]struct{}, len(getPolicy))
	for _, record := range getPolicy {
		key := strings.Join(record, ",")
		m[key] = struct{}{}
	}

	for _, record := range wantPolicy {
		key := strings.Join(record, ",")
		if _, ok := m[key]; !ok {
			t.Error("get policy: \n", getPolicy, "supposed to be: \n", wantPolicy)
			break
		}
	}
}

func validateNilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
}
