package logic_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	mockuserModel "github.com/GGjahoon/IZH/application/user/rpc/internal/mock"
	model "github.com/GGjahoon/IZH/application/user/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/GGjahoon/IZH/pkg/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestFindById(t *testing.T) {
	user := randomUser(t)

	testCases := []struct {
		name           string
		req            service.FindByIdRequest
		buildStubs     func(model *mockuserModel.MockUserModel)
		checkeResponse func(t *testing.T, rsp *service.FindByIdResponse, err error)
	}{
		{
			name: "ok",
			req: service.FindByIdRequest{
				UserId: user.Id,
			},
			buildStubs: func(userModel *mockuserModel.MockUserModel) {
				userModel.EXPECT().FindOne(gomock.Any(), gomock.Eq(user.Id)).Times(1).Return(&user, nil)
			},
			checkeResponse: func(t *testing.T, rsp *service.FindByIdResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, rsp.UserId, user.Id)
				require.Equal(t, rsp.Username, user.Username)
				require.Equal(t, rsp.Avatar, user.Avatar)
			},
		},
		{
			name: "Not Found",
			req: service.FindByIdRequest{
				UserId: user.Id,
			},
			buildStubs: func(userModel *mockuserModel.MockUserModel) {
				userModel.EXPECT().FindOne(gomock.Any(), gomock.Eq(user.Id)).Times(1).Return(nil, sqlx.ErrNotFound)
			},
			checkeResponse: func(t *testing.T, rsp *service.FindByIdResponse, err error) {
				myError := errors.New("the user is not found in db")
				require.Error(t, err)
				require.Equal(t, err, myError)
			},
		},
		{
			name: "DB internal Error",
			req: service.FindByIdRequest{
				UserId: user.Id,
			},
			buildStubs: func(userModel *mockuserModel.MockUserModel) {
				userModel.EXPECT().FindOne(gomock.Any(), gomock.Eq(user.Id)).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkeResponse: func(t *testing.T, rsp *service.FindByIdResponse, err error) {
				myError := errors.New("db internal error")
				require.Error(t, err)
				require.Equal(t, err, myError)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//create the crtl of mock model
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userModel := mockuserModel.NewMockUserModel(ctrl)
			//build stubs
			tc.buildStubs(userModel)

			//create a new server
			testServer := NewTestServer(t, userModel)

			//call the findbyid method
			rsp, err := testServer.FindById(context.Background(), &tc.req)

			//check the response of testServer
			tc.checkeResponse(t, rsp, err)

		})
	}
}
func randomUser(t *testing.T) model.User {
	return model.User{
		Id:       util.RandoInt(1, 1000),
		Username: util.RandomUser(),
		Mobile:   util.RandomMobile(),
		Avatar:   util.RandomAvatar(),
	}
}
