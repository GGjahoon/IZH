package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/GGjahoon/IZH/application/article/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	prefixArticles = "biz#articles#%d#%d"
	articlesExpire = 3600 * 24 * 2
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *pb.ArticlesRequest) (*pb.ArticlesResponse, error) {
	// valid request paramaters
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}

	var (
		sortField       string
		sortLikeNum     int64
		sortPublishTime string
	)
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
		sortLikeNum = in.Cursor
	} else {
		sortField = "publish_time"
		sortPublishTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	}
	var (
		err            error
		isCache, isEnd bool
		lastId, cursor int64
		curPage        []*pb.ArticleItem
		articles       []*model.Article
	)
	// step1 : get the article list from redis
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.PageSize, in.Cursor, in.SortType)
	if len(articleIds) > 0 {
		//赋值表示已被缓存命中
		isCache = true
		//判断文章列表是否结束
		if articleIds[len(articleIds)-1] == -1 {
			isEnd = true
		}
		articles, err = l.articlesById(l.ctx, articleIds)
		if err != nil {
			return nil, err
		}
		// var cmpFunc func(a, b *model.Article) int
		// if sortField == "like_num" {
		// 	cmpFunc = func(a, b *model.Article) int {
		// 		return cmp.Compare(b.LikeNum, a.LikeNum)
		// 	}
		// } else {
		// 	cmpFunc = func(a, b *model.Article) int {
		// 		return cmp.Compare(b.PublishTime.Unix(), a.PublishTime.Unix())
		// 	}
		// }
		// slices.SortFunc(articles, cmpFunc)

		for _, article := range articles {
			curPage = append(curPage, &pb.ArticleItem{
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				LikeCount:    article.LikeNum,
				CommentCount: article.CommentNum,
				PublishTime:  article.PublishTime.Unix(),
			})
		}
	} else {
		//缓存内没有该用户的文章列表，直接从数据库查询,默认查询20页
		articles, err = l.svcCtx.ArticleModel.ArticlesByUserId(l.ctx, in.UserId, sortLikeNum,
			sortPublishTime, sortField, types.DefaultLimit)
		if err != nil {
			l.Logger.Errorf("ArticlesByUserId userId:%d sortField:%s error:%v", in.UserId, sortField, err)
			return nil, err
		}
		var firstPageArticles []*model.Article
		if len(articles) > int(in.PageSize) {
			firstPageArticles = articles[:int(in.PageSize)]
		} else {
			firstPageArticles = articles
			isEnd = true
		}
		for _, article := range firstPageArticles {
			curPage = append(curPage, &pb.ArticleItem{
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				LikeCount:    article.LikeNum,
				CommentCount: article.CommentNum,
				PublishTime:  article.PublishTime.Unix(),
			})
		}
	}

	//修改返回的参数：cursor(curpage中最后一个article的cursor返回，以便下次使用)
	//articleId:curpage中最后一个文章的id
	//去重操作
	if len(curPage) > 0 {
		//修改返回的参数：cursor(curpage中最后一个article的cursor返回，以便下次使用)
		//articleId:curpage中最后一个文章的id
		//加上去重的操作
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		if in.SortType == types.SortPublishTime {
			cursor = pageLast.PublishTime
		} else {
			cursor = pageLast.LikeCount
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, article := range curPage {
			if in.SortType == types.SortPublishTime {
				if article.PublishTime == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			} else {
				if article.LikeCount == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			}
		}

	}
	//构建rsp
	rsp := &pb.ArticlesResponse{
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
		Articles:  curPage,
	}
	//仅当未从redis中命中时走入该分支
	//判断isCache，若为false，则未被写入到redis中,将当前的文章列表以sortset的形式写入到redis中（异步写入）
	if !isCache {
		threading.GoSafe(func() {
			//若能从db中查询到articles且长度小于limit：200，则在articles切片末尾添加上articleID为-1的article，以表示结束
			if len(articles) > 0 && len(articles) < types.DefaultLimit {
				articles = append(articles, &model.Article{Id: -1})
			}
			err = l.addCacheArticles(context.Background(), articles, in.UserId, in.SortType)
			if err != nil {
				l.Logger.Errorf("addCacheArticles error: %v", err)
			}
		})
	}

	return rsp, nil
}

func articlesKey(userId int64, sortType int32) string {
	return fmt.Sprintf(prefixArticles, userId, sortType)
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context,
	userId, pageSize, cursor int64,
	sortType int32,
) ([]int64, error) {
	key := articlesKey(userId, sortType)
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		l.Logger.Errorf("ZrevrangebyscoreWithScoresAndLimitCtx key :%s err : %v", key, err)
		return nil, err
	}
	var ids []int64
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			l.Logger.Errorf("ParseInt key:%s err:%v", pair.Key, err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (l *ArticlesLogic) articlesById(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](
		func(source chan<- int64) {
			for _, articleId := range articleIds {
				if articleId == -1 {
					continue
				}
				source <- articleId
			}
		}, func(id int64, writer mr.Writer[*model.Article], cancel func(error)) {
			article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, id)
			if err != nil {
				cancel(err)
				return
			}
			writer.Write(article)
		}, func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) {
			var articles []*model.Article
			for article := range pipe {
				articles = append(articles, article)
			}
			writer.Write(articles)
		})
	if err != nil {
		return nil, err
	}
	return articles, nil
}
func (l *ArticlesLogic) addCacheArticles(ctx context.Context,
	articles []*model.Article,
	userId int64,
	sortType int32,
) error {
	if len(articles) == 0 {
		return nil
	}
	key := articlesKey(userId, sortType)
	for _, article := range articles {
		var score int64
		if sortType == types.SortLikeCount {
			score = article.LikeNum
		} else {
			score = article.PublishTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(article.Id)))
		if err != nil {
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
}
