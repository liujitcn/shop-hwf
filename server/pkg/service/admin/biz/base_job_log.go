package biz

import (
	"context"
	"encoding/json"
	"gitee.com/liujit/shop/server/api/admin"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/timeutil"
	queueData "go.newcapec.cn/ncttools/nmskit-bootstrap/queue/data"
	"go.newcapec.cn/ncttools/nmskit/log"
	"time"
)

type BaseJobLogCase struct {
	data.BaseJobLogRepo
}

// NewBaseJobLogCase new a BaseJobLog use case.
func NewBaseJobLogCase(baseJobLogRepo data.BaseJobLogRepo) *BaseJobLogCase {
	return &BaseJobLogCase{
		BaseJobLogRepo: baseJobLogRepo,
	}
}

func (c *BaseJobLogCase) GetFromID(ctx context.Context, id int64) (*models.BaseJobLog, error) {
	return c.Find(ctx, &data.BaseJobLogCondition{
		Id: id,
	})
}

func (c *BaseJobLogCase) Page(ctx context.Context, req *admin.PageBaseJobLogRequest) (*admin.PageBaseJobLogResponse, error) {
	executeTime := req.GetExecuteTime()
	var startTime, endTime *time.Time
	if len(executeTime) == 2 {
		startTime = timeutil.StringTimeToTime(executeTime[0])
		endTime = timeutil.StringTimeToTime(executeTime[1])
		if endTime != nil {
			t := endTime.AddDate(0, 0, 1)
			endTime = &t
		}
	}
	condition := &data.BaseJobLogCondition{
		JobId:            req.GetJobId(),
		Status:           int32(req.GetStatus()),
		ExecuteStartTime: startTime,
		ExecuteEndTime:   endTime,
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseJobLog, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PageBaseJobLogResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseJobLogCase) ConvertToProto(item *models.BaseJobLog) *admin.BaseJobLog {
	processTime := time.Duration(item.ProcessTime) * time.Millisecond
	return &admin.BaseJobLog{
		Id:          item.ID,
		JobId:       item.JobID,
		Input:       item.Input,
		Output:      item.Output,
		Error:       item.Error,
		Status:      common.BaseJobLogStatus(item.Status),
		ProcessTime: processTime.String(),
		ExecuteTime: timeutil.TimeToTimeString(item.ExecuteTime),
	}
}

func (c *BaseJobLogCase) SaveJobLog(message queueData.Message) error {
	rb, err := json.Marshal(message.Values)
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var m map[string]*models.BaseJobLog
	err = json.Unmarshal(rb, &m)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	if v, ok := m["data"]; ok {
		err = c.Create(context.TODO(), v)
		if err != nil {
			return err
		}
	}
	return nil
}
