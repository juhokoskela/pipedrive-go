package v2

import (
	"context"
	"strconv"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Task struct {
	ID               TaskID     `json:"id"`
	Title            string     `json:"title,omitempty"`
	CreatorID        UserID     `json:"creator_id,omitempty"`
	Description      *string    `json:"description,omitempty"`
	ProjectID        ProjectID  `json:"project_id,omitempty"`
	IsDone           bool       `json:"is_done,omitempty"`
	IsMilestone      bool       `json:"is_milestone,omitempty"`
	DueDate          *string    `json:"due_date,omitempty"`
	StartDate        *string    `json:"start_date,omitempty"`
	ParentTaskID     *TaskID    `json:"parent_task_id,omitempty"`
	AssigneeIDs      []UserID   `json:"assignee_ids,omitempty"`
	Priority         *int       `json:"priority,omitempty"`
	AddTime          *time.Time `json:"add_time,omitempty"`
	UpdateTime       *time.Time `json:"update_time,omitempty"`
	MarkedAsDoneTime *time.Time `json:"marked_as_done_time,omitempty"`
}

type TaskDeleteResult struct {
	ID TaskID `json:"id"`
}

type TasksService struct{ client *Client }

type ListTasksOption interface{ applyListTasks(*listTasksOptions) }
type CreateTaskOption interface{ applyCreateTask(*createTaskOptions) }
type UpdateTaskOption interface{ applyUpdateTask(*updateTaskOptions) }

type TaskOption interface {
	CreateTaskOption
	UpdateTaskOption
}

type TaskRequestOption interface {
	ListTasksOption
	CreateTaskOption
	UpdateTaskOption
	taskRequestOptions() []pipedrive.RequestOption
}

type listTasksOptions struct {
	params         genv2.GetTasksParams
	requestOptions []pipedrive.RequestOption
}
type createTaskOptions struct {
	payload        taskPayload
	requestOptions []pipedrive.RequestOption
}
type updateTaskOptions struct {
	payload        taskPayload
	requestOptions []pipedrive.RequestOption
}

type nullableTaskValue[T any] struct {
	set   bool
	value *T
}

type taskPayload struct {
	title          *string
	projectID      *ProjectID
	parentTaskID   nullableTaskValue[TaskID]
	description    nullableTaskValue[string]
	done           *bool
	milestone      *bool
	dueDate        nullableTaskValue[string]
	startDate      nullableTaskValue[string]
	assigneeID     nullableTaskValue[UserID]
	assigneeIDs    []UserID
	assigneeIDsSet bool
	priority       nullableTaskValue[int]
}

type listTasksOptionFunc func(*listTasksOptions)

func (f listTasksOptionFunc) applyListTasks(cfg *listTasksOptions) { f(cfg) }

type taskFieldOption func(*taskPayload)

func (f taskFieldOption) applyCreateTask(cfg *createTaskOptions) { f(&cfg.payload) }
func (f taskFieldOption) applyUpdateTask(cfg *updateTaskOptions) { f(&cfg.payload) }

type taskRequestOption struct{ options []pipedrive.RequestOption }

func (o taskRequestOption) taskRequestOptions() []pipedrive.RequestOption { return o.options }
func (o taskRequestOption) applyListTasks(cfg *listTasksOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o taskRequestOption) applyCreateTask(cfg *createTaskOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o taskRequestOption) applyUpdateTask(cfg *updateTaskOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}

func WithTaskRequestOptions(opts ...pipedrive.RequestOption) TaskRequestOption {
	return taskRequestOption{options: opts}
}
func WithTasksPageSize(limit int) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) {
		if limit > 0 {
			cfg.params.Limit = &limit
		}
	})
}
func WithTasksCursor(cursor string) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) {
		if cursor != "" {
			cfg.params.Cursor = &cursor
		}
	})
}
func WithTasksDone(done bool) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) { cfg.params.IsDone = &done })
}
func WithTasksMilestone(milestone bool) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) { cfg.params.IsMilestone = &milestone })
}
func WithTasksAssigneeID(id UserID) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) { value := int(id); cfg.params.AssigneeId = &value })
}
func WithTasksProjectID(id ProjectID) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) { value := int(id); cfg.params.ProjectId = &value })
}
func WithTasksParentTaskID(id TaskID) ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) {
		value := strconv.FormatInt(int64(id), 10)
		cfg.params.ParentTaskId = &value
	})
}
func WithTasksRootOnly() ListTasksOption {
	return listTasksOptionFunc(func(cfg *listTasksOptions) { value := "null"; cfg.params.ParentTaskId = &value })
}

func WithTaskTitle(title string) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.title = &title })
}
func WithTaskProjectID(id ProjectID) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.projectID = &id })
}
func WithTaskParentTaskID(id TaskID) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.parentTaskID = nullableTaskValue[TaskID]{set: true, value: &id} })
}
func ClearTaskParentTaskID() TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.parentTaskID = nullableTaskValue[TaskID]{set: true} })
}
func WithTaskDescription(description string) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.description = nullableTaskValue[string]{set: true, value: &description} })
}
func ClearTaskDescription() TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.description = nullableTaskValue[string]{set: true} })
}
func WithTaskDone(done bool) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.done = &done })
}
func WithTaskMilestone(milestone bool) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.milestone = &milestone })
}
func WithTaskDueDate(date string) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.dueDate = nullableTaskValue[string]{set: true, value: &date} })
}
func ClearTaskDueDate() TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.dueDate = nullableTaskValue[string]{set: true} })
}
func WithTaskStartDate(date string) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.startDate = nullableTaskValue[string]{set: true, value: &date} })
}
func ClearTaskStartDate() TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.startDate = nullableTaskValue[string]{set: true} })
}
func WithTaskAssigneeID(id UserID) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.assigneeID = nullableTaskValue[UserID]{set: true, value: &id} })
}
func ClearTaskAssigneeID() TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.assigneeID = nullableTaskValue[UserID]{set: true} })
}
func WithTaskAssigneeIDs(ids ...UserID) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.assigneeIDsSet = true; p.assigneeIDs = append([]UserID{}, ids...) })
}
func WithTaskPriority(priority int) TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.priority = nullableTaskValue[int]{set: true, value: &priority} })
}
func ClearTaskPriority() TaskOption {
	return taskFieldOption(func(p *taskPayload) { p.priority = nullableTaskValue[int]{set: true} })
}

func newListTasksOptions(opts []ListTasksOption) listTasksOptions {
	var cfg listTasksOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyListTasks(&cfg)
		}
	}
	return cfg
}
func newCreateTaskOptions(opts []CreateTaskOption) createTaskOptions {
	var cfg createTaskOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyCreateTask(&cfg)
		}
	}
	return cfg
}
func newUpdateTaskOptions(opts []UpdateTaskOption) updateTaskOptions {
	var cfg updateTaskOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyUpdateTask(&cfg)
		}
	}
	return cfg
}
func taskRequestOptionValues(opts []TaskRequestOption) []pipedrive.RequestOption {
	var out []pipedrive.RequestOption
	for _, opt := range opts {
		if opt != nil {
			out = append(out, opt.taskRequestOptions()...)
		}
	}
	return out
}

func (p taskPayload) body() map[string]interface{} {
	body := map[string]interface{}{}
	if p.title != nil {
		body["title"] = *p.title
	}
	if p.projectID != nil {
		body["project_id"] = int(*p.projectID)
	}
	if p.parentTaskID.set {
		body["parent_task_id"] = p.parentTaskID.value
	}
	if p.description.set {
		body["description"] = p.description.value
	}
	if p.done != nil {
		if *p.done {
			body["done"] = 1
		} else {
			body["done"] = 0
		}
	}
	if p.milestone != nil {
		if *p.milestone {
			body["milestone"] = 1
		} else {
			body["milestone"] = 0
		}
	}
	if p.dueDate.set {
		body["due_date"] = p.dueDate.value
	}
	if p.startDate.set {
		body["start_date"] = p.startDate.value
	}
	if p.assigneeID.set {
		body["assignee_id"] = p.assigneeID.value
	}
	if p.assigneeIDsSet {
		body["assignee_ids"] = p.assigneeIDs
	}
	if p.priority.set {
		body["priority"] = p.priority.value
	}
	return body
}

func (s *TasksService) List(ctx context.Context, opts ...ListTasksOption) ([]Task, *string, error) {
	cfg := newListTasksOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}
func (s *TasksService) ListPager(opts ...ListTasksOption) *pipedrive.CursorPager[Task] {
	cfg := newListTasksOptions(opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Task, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}
func (s *TasksService) ForEach(ctx context.Context, fn func(Task) error, opts ...ListTasksOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *TasksService) Get(ctx context.Context, id TaskID, opts ...TaskRequestOption) (*Task, error) {
	if err := validateID(id, "task id"); err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, taskRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetTaskWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Task](resp.HTTPResponse, resp.Body, "task")
}

func (s *TasksService) Create(ctx context.Context, opts ...CreateTaskOption) (*Task, error) {
	cfg := newCreateTaskOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.AddTaskWithBodyWithResponse(ctx, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Task](resp.HTTPResponse, resp.Body, "task")
}

func (s *TasksService) Update(ctx context.Context, id TaskID, opts ...UpdateTaskOption) (*Task, error) {
	if err := validateID(id, "task id"); err != nil {
		return nil, err
	}
	cfg := newUpdateTaskOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.UpdateTaskWithBodyWithResponse(ctx, int(id), "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Task](resp.HTTPResponse, resp.Body, "task")
}

func (s *TasksService) Delete(ctx context.Context, id TaskID, opts ...TaskRequestOption) (*TaskDeleteResult, error) {
	if err := validateID(id, "task id"); err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, taskRequestOptionValues(opts)...)
	resp, err := s.client.gen.DeleteTaskWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[TaskDeleteResult](resp.HTTPResponse, resp.Body, "task delete")
}

func (s *TasksService) list(ctx context.Context, params genv2.GetTasksParams, requestOptions []pipedrive.RequestOption) ([]Task, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.GetTasksWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	return decodeV2List[Task](resp.HTTPResponse, resp.Body)
}
