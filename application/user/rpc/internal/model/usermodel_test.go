package model

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/GGjahoon/IZH/pkg/util"
	"github.com/stretchr/testify/require"
)

func InsertUser(t *testing.T) User {
	//create a new user and insert into db
	var randomUser User
	ctx := context.Background()
	randomUser.Username = util.RandomUser()
	randomUser.Mobile = util.RandomMobile()
	randomUser.Avatar = util.RandomAvatar()
	randomUser.CreateTime = time.Now()
	randomUser.UpdateTime = time.Now()
	ret, err := testModel.Insert(ctx, &randomUser)
	require.NoError(t, err)
	fmt.Println(randomUser.Username)

	// find the user inserted just,and compare
	userId, err := ret.LastInsertId()
	fmt.Println(userId)
	require.NoError(t, err)
	user, err := testModel.FindOne(ctx, userId)
	fmt.Println(user)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, randomUser.Username, user.Username)
	require.Equal(t, user.Avatar, randomUser.Avatar)
	require.Equal(t, user.Mobile, randomUser.Mobile)
	require.WithinDuration(t, user.CreateTime, randomUser.CreateTime, time.Second)
	require.WithinDuration(t, user.UpdateTime, randomUser.UpdateTime, time.Second)
	return *user
}

func TestInsertUSer(t *testing.T) {
	InsertUser(t)
}

func TestUpdate(t *testing.T) {
	randomUser := InsertUser(t)

	newUser := randomUser
	newUser.Username = util.RandomUser()
	newUser.Mobile = util.RandomMobile()
	time.Sleep(2 * time.Second)
	ctx := context.Background()
	err := testModel.Update(ctx, &newUser)
	tm := time.Now()
	require.NoError(t, err)
	updatedUser, err := testModel.FindOne(ctx, randomUser.Id)
	require.NoError(t, err)
	require.Equal(t, randomUser.Id, newUser.Id, updatedUser.Id)
	require.Equal(t, updatedUser.Username, newUser.Username)
	require.Equal(t, updatedUser.Mobile, newUser.Mobile)
	require.NotEqual(t, randomUser.Mobile, updatedUser.Mobile)
	require.NotEqual(t, randomUser.Username, updatedUser.Username)
	require.WithinDuration(t, tm, updatedUser.UpdateTime, time.Second)
}

func TestDeleteUser(t *testing.T) {
	randomUser := InsertUser(t)
	ctx := context.Background()
	err := testModel.Delete(ctx, randomUser.Id)
	require.NoError(t, err)
	_, err = testModel.FindOne(ctx, randomUser.Id)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
}

func TestFindUserByMobile(t *testing.T) {
	randomUser := InsertUser(t)
	ctx := context.Background()
	user, err := testModel.FindByMobile(ctx, randomUser.Mobile)
	require.NoError(t, err)
	require.Equal(t, user.Username, randomUser.Username)
	require.Equal(t, randomUser.Mobile, user.Mobile)
	require.Equal(t, randomUser.Avatar, user.Avatar)
}
