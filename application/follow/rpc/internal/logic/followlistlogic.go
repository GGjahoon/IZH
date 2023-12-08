package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/GGjahoon/IZH/application/follow/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

const userFollowListExpireTime = 3600 * 24 * 2

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
	// todo: add your logic here and delete this line
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	var (
		err              error
		isEnd, isCache   bool
		lastId, cursor   int64
		NeedCountUserIds []int64
		follows          []*model.Follow
		curPage          []*pb.FollowItem
	)

	followedUserIds, _ := l.cacheFollowUserIds(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(followedUserIds) != 0 {
		//命中缓存
		isCache = true
		//查看该列表是否由表明终止符号
		if followedUserIds[len(followedUserIds)-1] == -1 {
			//将终止符取出
			followedUserIds = followedUserIds[:len(followedUserIds)-1]
			isEnd = true
		}
		// 再次判断是否含有有效数据
		if len(followedUserIds) == 0 {
			return &pb.FollowListResponse{}, nil
		}

		// 查询follow表
		follows, err = l.svcCtx.FollowModel.FindByFollowedUserIds(l.ctx, followedUserIds)
		if err != nil {
			l.Logger.Errorf("[FollowList] FollowModel.FindByFollowUserIds error : %v", err)
			return nil, err
		}

		for _, follow := range follows {
			NeedCountUserIds = append(NeedCountUserIds, follow.FollowedUserId)
			curPage = append(curPage, &pb.FollowItem{
				Id:             follow.ID,
				FollowedUserId: follow.FollowedUserId,
				CreatTime:      follow.CreateTime.Unix(),
			})
		}

	} else {
		//未命中缓存，从数据库中查询,最多查询types.CacheMaxFollowCount个记录
		follows, err = l.svcCtx.FollowModel.FindByUserId(l.ctx, in.UserId, types.CacheMaxFollowCount)
		if err != nil {
			l.Logger.Errorf("[FollowList] FollowModel.FindByUserId error : %v ", err)
			return nil, err
		}
		if len(follows) == 0 {
			return &pb.FollowListResponse{}, nil
		}
		var firstPageFollows []*model.Follow
		if len(follows) > int(in.PageSize) {
			firstPageFollows = follows[:in.PageSize]
		} else {
			firstPageFollows = follows
			isEnd = true
		}
		for _, follow := range firstPageFollows {
			NeedCountUserIds = append(NeedCountUserIds, follow.FollowedUserId)
			curPage = append(curPage, &pb.FollowItem{
				Id:             follow.ID,
				FollowedUserId: follow.FollowedUserId,
				CreatTime:      follow.CreateTime.Unix(),
			})
		}

	}

	//去重，添加pageLasatID
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		cursor = pageLast.CreatTime
		if cursor < 0 {
			cursor = 0
		}
		for k, follow := range curPage {
			if follow.CreatTime == in.Cursor && follow.Id == in.Id {
				curPage = curPage[k:]
				break
			}
		}
	}
	// 查询被关注用户的粉丝数
	followCounts, err := l.svcCtx.FollowCountModel.FindByUserIds(l.ctx, NeedCountUserIds)
	if err != nil {
		l.Logger.Errorf("[FollowList] FollowCountModel.FindByUserIds error : %v ", err)
		return nil, err
	}
	userIdFansCount := make(map[int64]int)
	for _, followCount := range followCounts {
		userIdFansCount[followCount.UserId] = followCount.FansCount
	}
	for _, cur := range curPage {
		cur.FansCount = int64(userIdFansCount[cur.FollowedUserId])
	}

	rsp := &pb.FollowListResponse{
		Items:  curPage,
		Cursor: cursor,
		IsEnd:  isEnd,
		LastID: lastId,
	}

	if !isCache {
		//将用户关注的用户id列表异步写入redis
		threading.GoSafe(func() {
			//判断从数据库中查询出的数据是否大于默认分页大小
			if len(follows) < types.CacheMaxFollowCount && len(follows) > 0 {
				//在末尾加上用户id -1 作为结束标志符
				follows = append(follows, &model.Follow{FollowedUserId: -1})
			}
			//以z set数据格式写入redis
			err = l.addCacheFollow(context.Background(), in.UserId, follows)
			if err != nil {
				logx.Errorf("addCacheFollow error : %v", err)
			}
		})
	}

	return rsp, nil
}
func (l *FollowListLogic) cacheFollowUserIds(ctx context.Context,
	userId, cursor, pageSize int64,
) ([]int64, error) {
	key := userFollowKey(userId)

	//判断该key是否存在，如果存在，则为该key续期，避免缓存击穿(续期的逻辑，若未找到，忽略错误即可)
	isExist, err := l.svcCtx.BizRedis.Exists(key)
	if err != nil {
		logx.Errorf("[cacheFollowUserIds] BizRidsExists error : %v", err)
	}
	if isExist {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowListExpireTime)
		if err != nil {
			logx.Errorf("[cacheFollowUserIds] BizReids.ExpireCtx error : %v", err)
		}
	}

	pair, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("[cacheFollowUserIds] BizReids.ZrevrangebyscoreWithScoresAndLimitCtx error : %v", err)
		return nil, err
	}
	var userIds []int64
	for _, p := range pair {
		userId, err := strconv.ParseInt(p.Key, 10, 64)
		if err != nil {
			logx.Errorf("[cacheFollowUserIds] parseint error : %v", err)
			return nil, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func (l *FollowListLogic) addCacheFollow(ctx context.Context, userId int64, follows []*model.Follow) error {
	if len(follows) == 0 {
		return nil
	}
	key := userFollowKey(userId)
	for _, follow := range follows {
		var score int64
		if follow.FollowedUserId == -1 {
			score = 0
		} else {
			score = follow.CreateTime.Unix()
		}

		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(follow.FollowedUserId, 10))
		if err != nil {
			logx.Errorf("[addCacheFollow] error : %v", err)
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowListExpireTime)
}
