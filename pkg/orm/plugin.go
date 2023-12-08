package orm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func NewCustomPlugin() gorm.Plugin {
	return &CustomePlugin{}
}

type CustomePlugin struct {
}

func (p *CustomePlugin) Name() string {
	return "CustomPlugin"
}

// 之后需要使用db的几个方法：create query update delete row raw 等方法前后添加上需要运行的插件函数即可。
func (p *CustomePlugin) Initialize(db *gorm.DB) error {
	//Befor Callback
	if err := db.Callback().Create().Before("gorm:createBefore").Register("gorm:createBefore:metric:trace", func(db *gorm.DB) {
		//初始化metrics，记录操作开始时间
		makeMetrics(db, "create")
		//使用db的ctx初始化tracer（span）
		makeTrace(db, "create")
	}); err != nil {
		return err
	}

	if err := db.Callback().Query().Before("gorm:queryBefore").Register("gorm:queryBefore:metric:trace", func(db *gorm.DB) {
		//初始化metrics，记录操作开始时间
		makeMetrics(db, "query")
		//使用db的ctx初始化tracer（span）
		makeTrace(db, "query")
	}); err != nil {
		return err
	}

	if err := db.Callback().Update().Before("gorm:updateBefore").Register("gorm:updateBefore:metric:trace", func(db *gorm.DB) {
		//初始化metrics，记录操作开始时间
		makeMetrics(db, "update")
		//使用db的ctx初始化tracer（span）
		makeTrace(db, "update")
	}); err != nil {
		return err
	}

	if err := db.Callback().Delete().Before("gorm:deleteBefore").Register("gorm:deleteBefore:metric:trace", func(db *gorm.DB) {
		//初始化metrics，记录操作开始时间
		makeMetrics(db, "delete")
		//使用db的ctx初始化tracer（span）
		makeTrace(db, "delete")
	}); err != nil {
		return err
	}

	if err := db.Callback().Row().Before("gorm:rowBefore").Register("gorm:rowBefore:metric:trace", func(db *gorm.DB) {
		//初始化metrics，记录操作开始时间
		makeMetrics(db, "row")
		//使用db的ctx初始化tracer（span）
		makeTrace(db, "row")
	}); err != nil {
		return err
	}

	if err := db.Callback().Raw().Before("gorm:rawBefore").Register("gorm:rawBefore:metric:trace", func(db *gorm.DB) {
		//初始化metrics，记录操作开始时间
		makeMetrics(db, "raw")
		//使用db的ctx初始化tracer（span）
		makeTrace(db, "raw")
	}); err != nil {
		return err
	}

	// After Callbakc
	if err := db.Callback().Create().After("gorm:createAfter").Register("gorm:createAfter:metric:trace", func(db *gorm.DB) {
		//获取db操作开始时间，使用metric记录 耗时、对象表、方法,供prometheus监控
		getMetrics(db, "create")
		//获取db操作开始前初始化好的span，获取后断言为接口，若发生错误则记录， 并记录操作的对象Table及sql语句
		getTrace(db, "create")
	}); err != nil {
		return err
	}

	if err := db.Callback().Query().After("gorm:queryAfter").Register("gorm:queryAfter:metric:trace", func(db *gorm.DB) {
		getMetrics(db, "query")
		getTrace(db, "query")
	}); err != nil {
		return err
	}

	if err := db.Callback().Update().After("gorm:updateAfter").Register("gorm:updateAfter:metric:trace", func(db *gorm.DB) {
		getMetrics(db, "update")
		getTrace(db, "update")
	}); err != nil {
		return err
	}

	if err := db.Callback().Delete().After("gorm:deleteAfter").Register("gorm:deleteAfter:metric:trace", func(db *gorm.DB) {
		getMetrics(db, "delete")
		getTrace(db, "delete")
	}); err != nil {
		return err
	}

	if err := db.Callback().Row().After("gorm:rowAfter").Register("gorm:rowAfter:metric:trace", func(db *gorm.DB) {
		getMetrics(db, "row")
		getTrace(db, "row")
	}); err != nil {
		return err
	}

	if err := db.Callback().Raw().After("gorm:rawAfter").Register("gorm:rawAfter:metric:trace", func(db *gorm.DB) {
		getMetrics(db, "raw")
		getTrace(db, "raw")
	}); err != nil {
		return err
	}

	return nil
}

func makeMetrics(db *gorm.DB, method string) {
	startTime := time.Now().Unix()
	//设置操作开始时间
	db.InstanceSet(fmt.Sprintf("gorm:%s_start)time", method), startTime)
}
func getMetrics(db *gorm.DB, method string) {
	// 获取操作开始时间
	startTime, ok := db.InstanceGet(fmt.Sprintf("gorm:%s_start)time", method))
	if !ok {
		return
	}
	stSecond := startTime.(int64)
	st := time.Unix(stSecond, 0)
	//metric中记录耗时、操作的table，方法
	metricClientReqDur.Observe(time.Since(st).Milliseconds(), db.Statement.Table, method)
	//记录table，method，是否发生错误
	metricClientReqErrTotal.Inc(db.Statement.Table, method, strconv.FormatBool(db.Statement.Error != nil))
}

func makeTrace(db *gorm.DB, method string) {
	ctx := db.Statement.Context
	//以db的上下文初始化tracer
	tracer := trace.TracerFromContext(ctx)
	// 初始化span
	_, span := tracer.Start(ctx, fmt.Sprintf("gorm:%s", method), oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	// 将span设置为k-v对
	db.InstanceSet(fmt.Sprintf("gorm:%s_span", method), span)
}
func getTrace(db *gorm.DB, method string) {
	// 以k获取span
	v, ok := db.InstanceGet(fmt.Sprintf("gorm:%s_span", method))
	if !ok {
		return
	}
	// 将获取的v断言
	span := v.(oteltrace.Span)
	//如果发生错误，span记录错误
	if db.Statement.Error != nil {
		span.RecordError(db.Statement.Error)
	}
	// span必记录的：所操作的table名称、执行的sql语句
	span.SetAttributes(
		semconv.DBSQLTableKey.String(db.Statement.Table),
		semconv.DBStatementKey.String(db.Statement.SQL.String()),
	)
	span.End()
}
