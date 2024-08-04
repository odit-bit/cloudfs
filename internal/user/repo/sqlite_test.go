package repo

import (
	"context"
	"testing"

	"github.com/odit-bit/cloudfs/internal/user"
)

func TestInsertAndFindAccount(t *testing.T) {

	db, err := DefaultDB(":memory:", "")
	if err != nil {
		t.Fatal(err)
	}

	adb, err := newAccountDB(db)
	if err != nil {
		t.Fatal(err)
	}

	// Test Insert Account
	testAccount := user.Account{
		Name:         "TestUser",
		HashPassword: "testpassword",
	}

	if err := adb.Insert(context.Background(), &testAccount); err != nil {
		t.Fatal(err)
	}

	// Test Find Account
	foundAccount, err := adb.Find(context.Background(), "TestUser")
	if err != nil {
		t.Fatal(err)
	}

	// Verify the retrieved account matches the inserted account
	if foundAccount.Name != testAccount.Name || string(foundAccount.HashPassword) != string(testAccount.HashPassword) {
		t.Errorf("Expected: %+v, Got: %+v", testAccount, foundAccount)
	}
}
