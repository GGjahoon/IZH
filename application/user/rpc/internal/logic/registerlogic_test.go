package logic_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/code"
	mockuserModel "github.com/GGjahoon/IZH/application/user/rpc/internal/mock"
	model "github.com/GGjahoon/IZH/application/user/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogics(t *testing.T) {
	user := randomUser(t)
	testCases := []struct {
		name          string
		request       *service.RegisterRequest
		buildStubs    func(store *mockuserModel.MockUserModel)
		checkResponse func(t *testing.T, rsp *service.RegisterResponse, err error)
	}{
		{
			name: "ok",
			request: &service.RegisterRequest{
				Username: user.Username,
				Mobile:   user.Mobile,
				Avatar:   user.Avatar,
			},
			buildStubs: func(store *mockuserModel.MockUserModel) {
				u := model.User{
					Username: user.Username,
					Mobile:   user.Mobile,
					Avatar:   user.Avatar,
				}
				result := NewResult(int(user.Id), true)
				store.EXPECT().Insert(gomock.Any(), gomock.Eq(&u)).Times(1).Return(result, nil)
			},
			checkResponse: func(t *testing.T, rsp *service.RegisterResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, user.Id, rsp.UserId)
			},
		},
		{
			name: "internal_error",
			request: &service.RegisterRequest{
				Username: user.Username,
				Mobile:   user.Mobile,
				Avatar:   user.Avatar,
			},
			buildStubs: func(store *mockuserModel.MockUserModel) {
				store.EXPECT().Insert(gomock.Any(), gomock.Any()).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rsp *service.RegisterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, err, sql.ErrConnDone)
			},
		},
		{
			name: "UserNameEmpty",
			request: &service.RegisterRequest{
				Username: "",
				Mobile:   user.Mobile,
				Avatar:   user.Avatar,
			},
			buildStubs: func(store *mockuserModel.MockUserModel) {
				store.EXPECT().Insert(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, rsp *service.RegisterResponse, err error) {
				require.Error(t, err)
				require.Empty(t, rsp)
				require.Equal(t, err, code.RegisterNameEmpty)
			},
		},
		{
			name: "Can not get id",
			request: &service.RegisterRequest{
				Username: user.Username,
				Mobile:   user.Mobile,
				Avatar:   user.Avatar,
			},
			buildStubs: func(store *mockuserModel.MockUserModel) {
				u := model.User{
					Username: user.Username,
					Mobile:   user.Mobile,
					Avatar:   user.Avatar,
				}
				result := NewResult(int(user.Id), false)
				store.EXPECT().Insert(gomock.Any(), gomock.Eq(&u)).Times(1).Return(result, nil)
			},
			checkResponse: func(t *testing.T, rsp *service.RegisterResponse, err error) {
				require.Error(t, err)
				require.Empty(t, rsp)
				require.Equal(t, err, errors.New("cannot get user id"))
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			crtl := gomock.NewController(t)
			defer crtl.Finish()

			store := mockuserModel.NewMockUserModel(crtl)
			// build stubs
			testCase.buildStubs(store)
			testServer := NewTestServer(t, store)
			rsp, err := testServer.Register(context.Background(), testCase.request)
			testCase.checkResponse(t, rsp, err)
		})
	}
}

type Result struct {
	Id int
	X  bool
}

func NewResult(id int, x bool) sql.Result {
	return &Result{
		Id: id,
		X:  x,
	}
}

func (result *Result) LastInsertId() (int64, error) {
	if result.X == true {
		return int64(result.Id), nil
	} else {
		return -1, errors.New("cannot get user id")
	}
}

func (result *Result) RowsAffected() (int64, error) {
	return 0, nil
}
