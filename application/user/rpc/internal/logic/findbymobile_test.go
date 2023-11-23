package logic_test

import (
	"context"
	"testing"

	mockuserModel "github.com/GGjahoon/IZH/application/user/rpc/internal/mock"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/GGjahoon/IZH/pkg/encrypt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestFindByMobile(t *testing.T) {
	user := randomUser(t)
	var err error
	user.Mobile, err = encrypt.EncMobile(user.Mobile)
	require.NoError(t, err)
	testCases := []struct {
		name           string
		req            service.FindByMobileRequest
		buildStubs     func(userModel *mockuserModel.MockUserModel)
		checkeResponse func(t *testing.T, rsp *service.FindByMobileResponse, err error)
	}{
		{
			name: "ok",
			req: service.FindByMobileRequest{
				Mobile: user.Mobile,
			},
			buildStubs: func(userModel *mockuserModel.MockUserModel) {
				userModel.EXPECT().FindByMobile(gomock.Any(), gomock.Eq(user.Mobile)).Times(1).Return(&user, nil)
			},
			checkeResponse: func(t *testing.T, rsp *service.FindByMobileResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, rsp.UserId, user.Id)
				require.Equal(t, rsp.Username, user.Username)
				require.Equal(t, rsp.Avatar, user.Avatar)
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
			rsp, err := testServer.FindByMobile(context.Background(), &tc.req)

			//check the response of testServer
			tc.checkeResponse(t, rsp, err)

		})
	}
}
