package orm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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
		InitSpan(db, "create")
	}); err != nil {
		return err
	}

	if err := db.Callback().Query().Before("gorm:queryBefore").Register("gorm:queryBefore:metric:trace", func(d *gorm.DB) {
		InitSpan(db, "query")
	}); err != nil {
		return err
	}

	if err := db.Callback().Update().Before("gorm:updateBefore").Register("gorm:updateBefore:metric:trace", func(d *gorm.DB) {
		InitSpan(db, "update")
	}); err != nil {
		return err
	}

	if err := db.Callback().Delete().Before("gorm:deleteBefore").Register("gorm:deleteBefore:metric:trace", func(d *gorm.DB) {
		InitSpan(db, "delete")
	}); err != nil {
		return err
	}

	if err := db.Callback().Row().Before("gorm:rowBefore").Register("gorm:rowBefore:metric:trace", func(d *gorm.DB) {
		InitSpan(db, "row")
	}); err != nil {
		return err
	}

	if err := db.Callback().Raw().Before("gorm:rawBefore").Register("gorm:rawBefore:metric:trace", func(d *gorm.DB) {
		InitSpan(db, "raw")
	}); err != nil {
		return err
	}

	// After Callbakc
	if err := db.Callback().Create().After("gorm:createAfter").Register("gorm:createAfter:metric:trace", func(d *gorm.DB) {
		MetricAndGetSpan(db, "create")
	}); err != nil {
		return err
	}

	if err := db.Callback().Query().After("gorm:queryAfter").Register("gorm:queryAfter:metric:trace", func(d *gorm.DB) {
		MetricAndGetSpan(db, "query")
	}); err != nil {
		return err
	}

	if err := db.Callback().Update().After("gorm:updateAfter").Register("gorm:updateAfter:metric:trace", func(d *gorm.DB) {
		MetricAndGetSpan(db, "update")
	}); err != nil {
		return err
	}

	if err := db.Callback().Delete().After("gorm:deleteAfter").Register("gorm:deleteAfter:metric:trace", func(d *gorm.DB) {
		MetricAndGetSpan(db, "delete")
	}); err != nil {
		return err
	}

	if err := db.Callback().Row().After("gorm:rowAfter").Register("gorm:rowAfter:metric:trace", func(d *gorm.DB) {
		MetricAndGetSpan(db, "row")
	}); err != nil {
		return err
	}

	if err := db.Callback().Row().After("gorm:rawAfter").Register("gorm:rawAfter:metric:trace", func(d *gorm.DB) {
		MetricAndGetSpan(db, "raw")
	}); err != nil {
		return err
	}

	return nil
}

func InitSpan(db *gorm.DB, method string) {
	startTime := time.Now().Unix()

	db.InstanceSet(fmt.Sprintf("gorm:%s_start_time", method), startTime)

	ctx := db.Statement.Context
	tracer := trace.TracerFromContext(ctx)
	_, span := tracer.Start(ctx, fmt.Sprintf("gorm:%s", method), oteltrace.WithSpanKind(oteltrace.SpanKindClient))

	db.InstanceSet(fmt.Sprintf("gorm:%s_span", method), span)
}

func MetricAndGetSpan(db *gorm.DB, method string) {
	//统计耗时及错误
	startTime, ok := db.InstanceGet(fmt.Sprintf("gorm:%s_start_time", method))
	if !ok {
		return
	}
	stSecond := startTime.(int64)
	st := time.Unix(stSecond, 0)
	//统计该方法耗时
	metricClientReqDur.Observe(time.Since(st).Microseconds(), db.Statement.Table, method)
	//统计该方法是否出现错误
	metricClientReqErrTotal.Inc(db.Statement.Table, method, strconv.FormatBool(db.Statement.Error != nil))

	//span相关操作
	v, ok := db.InstanceGet(fmt.Sprintf("gorm:%s_span", method))
	if !ok {
		return
	}
	span := v.(oteltrace.Span)

	if db.Statement.Error != nil {
		//如果关于db的方法执行错误
		span.RecordError(db.Statement.Error)
	}

	span.SetAttributes(
		semconv.DBSQLTableKey.String(db.Statement.Table),
		semconv.DBStatementKey.String(db.Statement.SQL.String()),
	)
	span.End()
}
