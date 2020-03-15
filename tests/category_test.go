// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build integration

package tests

import (
	"testing"
)

func TestCreateCategory(t *testing.T) {
	categoryName := "My category"
	client := createClient(t)
	category, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	if category.ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, category.ID)
	}

	if category.UserID <= 0 {
		t.Fatalf(`Invalid userID, got "%v"`, category.UserID)
	}

	if category.Title != categoryName {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, category.Title, categoryName)
	}
}

func TestCreateCategoryWithEmptyTitle(t *testing.T) {
	client := createClient(t)
	_, err := client.CreateCategory("")
	if err == nil {
		t.Fatal(`The category title should be mandatory`)
	}
}

func TestCannotCreateDuplicatedCategory(t *testing.T) {
	client := createClient(t)

	categoryName := "My category"
	_, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateCategory(categoryName)
	if err == nil {
		t.Fatal(`Duplicated categories should not be allowed`)
	}
}

func TestUpdateCategory(t *testing.T) {
	categoryName := "My category"
	client := createClient(t)
	category, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	categoryName = "Updated category"
	category, err = client.UpdateCategory(category.ID, categoryName)
	if err != nil {
		t.Fatal(err)
	}

	if category.ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, category.ID)
	}

	if category.UserID <= 0 {
		t.Fatalf(`Invalid userID, got "%v"`, category.UserID)
	}

	if category.Title != categoryName {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, category.Title, categoryName)
	}
}

func TestListCategories(t *testing.T) {
	categoryName := "My category"
	client := createClient(t)

	_, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	categories, err := client.Categories()
	if err != nil {
		t.Fatal(err)
	}

	if len(categories) != 2 {
		t.Fatalf(`Invalid number of categories, got "%v" instead of "%v"`, len(categories), 2)
	}

	if categories[0].ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, categories[0].ID)
	}

	if categories[0].UserID <= 0 {
		t.Fatalf(`Invalid userID, got "%v"`, categories[0].UserID)
	}

	if categories[0].Title != "All" {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, categories[0].Title, "All")
	}

	if categories[1].ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, categories[0].ID)
	}

	if categories[1].UserID <= 0 {
		t.Fatalf(`Invalid userID, got "%v"`, categories[1].UserID)
	}

	if categories[1].Title != categoryName {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, categories[1].Title, categoryName)
	}
}

func TestDeleteCategory(t *testing.T) {
	client := createClient(t)

	category, err := client.CreateCategory("My category")
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteCategory(category.ID)
	if err != nil {
		t.Fatal(`Removing a category should not raise any error`)
	}
}

func TestCannotDeleteCategoryOfAnotherUser(t *testing.T) {
	client := createClient(t)
	categories, err := client.Categories()
	if err != nil {
		t.Fatal(err)
	}

	client = createClient(t)
	err = client.DeleteCategory(categories[0].ID)
	if err == nil {
		t.Fatal(`Removing a category that belongs to another user should be forbidden`)
	}
}
