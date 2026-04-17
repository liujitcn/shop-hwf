package biz

import (
	"context"
	"encoding/json"
	"errors"
	"gitee.com/liujit/shop/server/api/admin"
	"gitee.com/liujit/shop/server/api/common"
	"gitee.com/liujit/shop/server/lib/data"
	"gitee.com/liujit/shop/server/lib/data/models"
	"gitee.com/liujit/shop/server/lib/utils/str"
	"gitee.com/liujit/shop/server/lib/utils/timeutil"
	"gitee.com/liujit/shop/server/lib/utils/trans"
	"gitee.com/liujit/shop/server/pkg/service/admin/job"
	"gitee.com/liujit/shop/server/pkg/service/admin/task"
	"github.com/robfig/cron/v3"
	"go.newcapec.cn/ncttools/nmskit/log"
)

type BaseJobCase struct {
	data.BaseJobRepo
	cron *cron.Cron
	task map[string]task.TaskExec
}

// NewBaseJobCase new a BaseJob use case.
func NewBaseJobCase(
	baseJobRepo data.BaseJobRepo,
	task map[string]task.TaskExec,
) *BaseJobCase {
	return &BaseJobCase{
		BaseJobRepo: baseJobRepo,
		cron:        cron.New(cron.WithSeconds()),
		task:        task,
	}
}

func (c *BaseJobCase) GetFromID(ctx context.Context, id int64) (*models.BaseJob, error) {
	return c.Find(ctx, &data.BaseJobCondition{Id: id})
}

func (c *BaseJobCase) Page(ctx context.Context, req *admin.PageBaseJobRequest) (*admin.PageBaseJobResponse, error) {
	condition := &data.BaseJobCondition{
		Status:       int32(req.GetStatus()),
		Name:         req.GetName(),
		InvokeTarget: req.GetInvokeTarget(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseJob, 0)
	for _, item := range page {
		list = append(list, &admin.BaseJob{
			Id:             item.ID,
			Name:           item.Name,
			CronExpression: item.CronExpression,
			InvokeTarget:   item.InvokeTarget,
			Args:           c.ConvertToBaseJobArgs(item.Args),
			Status:         common.Status(item.Status),
			EntryId:        item.EntryID,
			CreatedAt:      timeutil.TimeToTimeString(item.CreatedAt),
			UpdatedAt:      timeutil.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageBaseJobResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseJobCase) ConvertToProto(item *models.BaseJob) *admin.BaseJobForm {
	res := &admin.BaseJobForm{
		Id:             item.ID,
		Name:           item.Name,
		CronExpression: item.CronExpression,
		InvokeTarget:   item.InvokeTarget,
		Args:           c.ConvertToBaseJobArgs(item.Args),
		Status:         trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *BaseJobCase) ConvertToModel(item *admin.BaseJobForm) *models.BaseJob {
	res := &models.BaseJob{
		ID:             item.GetId(),
		Name:           item.GetName(),
		CronExpression: item.GetCronExpression(),
		InvokeTarget:   item.GetInvokeTarget(),
		Args:           str.ConvertAnyToJsonString(item.GetArgs()),
		Status:         int32(item.GetStatus()),
	}
	return res
}

func (c *BaseJobCase) ConvertToBaseJobArgs(argsStr string) []*admin.BaseJobArgs {
	args := make([]*admin.BaseJobArgs, 0)
	err := json.Unmarshal([]byte(argsStr), &args)
	if err != nil {
		log.Error(err)
	}
	return args
}

// Start 开始任务
func (c *BaseJobCase) Start(ctx context.Context, baseJob *models.BaseJob) error {
	invokeTarget, ok := c.task[baseJob.InvokeTarget]
	if !ok {
		return errors.New("invokeTarget not exist")
	}

	argsMap := make(map[string]string)
	args := c.ConvertToBaseJobArgs(baseJob.Args)
	for _, item := range args {
		argsMap[item.Key] = item.Value
	}

	execJob := &job.ExecJob{
		JobId:        baseJob.ID,
		Args:         argsMap,
		InvokeTarget: invokeTarget,
	}

	// 添加任务
	entryId, err := c.cron.AddJob(baseJob.CronExpression, execJob)
	if err != nil {
		log.Errorf("cron add error: %v", err)
		return err
	}
	// 更新
	baseJob.EntryID = int32(entryId)
	err = c.UpdateByID(ctx, baseJob)
	if err != nil {
		return err
	}
	return nil
}

// Stop 停止任务
func (c *BaseJobCase) Stop(entryID int32) chan bool {
	ch := make(chan bool)
	go func() {
		c.cron.Remove(cron.EntryID(entryID))
		log.Infof("remove success %d", entryID)
		ch <- true
	}()
	return ch
}

// Exec 执行任务
func (c *BaseJobCase) Exec(ctx context.Context, baseJob *models.BaseJob) error {
	if invokeTarget, ok := c.task[baseJob.InvokeTarget]; ok {
		argsMap := make(map[string]string)
		args := c.ConvertToBaseJobArgs(baseJob.Args)
		for _, item := range args {
			argsMap[item.Key] = item.Value
		}
		execJob := &job.ExecJob{
			JobId:        baseJob.ID,
			Args:         argsMap,
			InvokeTarget: invokeTarget,
		}
		execJob.Run()
		if execJob.Status == common.BaseJobLogStatus_FAIL {
			return errors.New(execJob.ErrMsg)
		}
	} else {
		return errors.New("invokeTarget not exist")
	}
	return nil
}

// Init 初始化定时任务
func (c *BaseJobCase) Init(ctx context.Context) error {
	// 清空job启动时返回的id
	err := c.CleanEntryId(ctx, 0)
	if err != nil {
		return err
	}
	// 查询全部定时任务
	baseJobList := make([]*models.BaseJob, 0)
	baseJobList, err = c.FindAll(ctx, &data.BaseJobCondition{
		Status: int32(common.Status_ENABLE),
	})
	if err != nil {
		return err
	}
	for _, item := range baseJobList {
		err = c.Start(ctx, item)
		if err != nil {
			return err
		}
	}

	c.cron.Start()

	// 优雅退出监听 context
	go func() {
		<-ctx.Done()
		log.Infof("context cancel received, stopping cron")
		c.cron.Stop()
	}()
	return nil
}
