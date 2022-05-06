package db

import (
	"context"
	"testing"
)

// These tests require you to be connected to the database
// It is to check if our sqlc methods work
// which are otherwize difficult to mock

func TestCreateUser(t *testing.T) {
	user := CreateUserParams{
		Firstname: "Savez",
		Lastname:  "Siddiqui",
		Phone:     "7408963464",
		Addr:      "some random address",
	}

	err := testQueries.CreateUser(context.Background(), user)

	if err != nil {
		t.Error("Failed with error", err)
	}
}

func TestGetUsers(t *testing.T) {
	_, err := testQueries.ListUser(context.Background())

	if err != nil {
		t.Error("Failed with error", err)
	}
}

func TestDeactivateUser(t *testing.T) {
	err := testQueries.DeactivateUser(context.Background(), "7408963464")

	if err != nil {
		t.Error("Failed with error", err)
	}
}

func TestDeleteUser(t *testing.T) {
	err := testQueries.DeleteUser(context.Background())

	if err != nil {
		t.Error("Failed with error", err)
	}
}
