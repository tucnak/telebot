package telebot

import (
	"encoding/json"
	"errors"
)

// Checklist is a checklist in a message (Bot API 9.1). As sent, it is a
// Sendable that a bot sends on behalf of a business account, so the send
// options must carry BusinessConnectionID; edit it with Edit.
//
// Bot API: Checklist, InputChecklist.
type Checklist struct {
	// Title of the checklist, 1-255 characters, with its entities or, when
	// sending, its parse mode.
	Title         string          `json:"title"`
	TitleEntities []MessageEntity `json:"title_entities,omitempty"`
	ParseMode     ParseMode       `json:"-"`

	// Tasks, 1-30 when sending.
	Tasks []ChecklistTask `json:"tasks"`

	// (Optional) Users other than the creator can add tasks, and mark tasks
	// as done or not done.
	OthersCanAddTasks        bool `json:"others_can_add_tasks,omitempty"`
	OthersCanMarkTasksAsDone bool `json:"others_can_mark_tasks_as_done,omitempty"`
}

// ChecklistTask is a task in a checklist.
//
// Bot API: ChecklistTask, InputChecklistTask.
type ChecklistTask struct {
	// Identifier of the task; when sending, positive and unique in the
	// checklist.
	ID int `json:"id"`

	// Text of the task, 1-100 characters, with its entities or, when
	// sending, its parse mode.
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	ParseMode    ParseMode       `json:"-"`

	// (Optional) Who completed the task, and when; zero if it isn't done.
	CompletedByUser *User `json:"completed_by_user,omitempty"`
	CompletedByChat *Chat `json:"completed_by_chat,omitempty"`
	CompletionDate  int64 `json:"completion_date,omitempty"`
}

// input returns the checklist as an InputChecklist.
func (c *Checklist) input() (string, error) {
	if c == nil {
		return "", errors.New("telebot: nil checklist")
	}
	w := richWire{"title": c.Title}
	w.set("parse_mode", c.ParseMode)
	if len(c.TitleEntities) > 0 {
		w["title_entities"] = c.TitleEntities
	}
	tasks := make([]richWire, len(c.Tasks))
	for i, t := range c.Tasks {
		tasks[i] = t.input()
	}
	w["tasks"] = tasks
	w.set("others_can_add_tasks", c.OthersCanAddTasks)
	w.set("others_can_mark_tasks_as_done", c.OthersCanMarkTasksAsDone)

	data, err := json.Marshal(w)
	return string(data), err
}

// input returns the task as an InputChecklistTask.
func (t ChecklistTask) input() richWire {
	w := richWire{"id": t.ID, "text": t.Text}
	w.set("parse_mode", t.ParseMode)
	if len(t.TextEntities) > 0 {
		w["text_entities"] = t.TextEntities
	}
	return w
}

// Send delivers the checklist through bot b to recipient, on behalf of the
// business account of opt.BusinessConnectionID.
//
// Bot API: sendChecklist.
func (c *Checklist) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	checklist, err := c.input()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"chat_id":   to.Recipient(),
		"checklist": checklist,
	}
	b.embedSendOptions(params, opt)
	delete(params, "parse_mode")

	data, err := b.Raw("sendChecklist", params)
	if err != nil {
		return nil, err
	}
	return extractMessage(data)
}

// ChecklistTasksDone is a service message about checklist tasks marked as
// done or not done.
//
// Bot API: ChecklistTasksDone.
type ChecklistTasksDone struct {
	// (Optional) The message with the checklist.
	ChecklistMessage *Message `json:"checklist_message,omitempty"`

	MarkedAsDoneTaskIDs    []int `json:"marked_as_done_task_ids,omitempty"`
	MarkedAsNotDoneTaskIDs []int `json:"marked_as_not_done_task_ids,omitempty"`
}

// ChecklistTasksAdded is a service message about tasks added to a
// checklist.
//
// Bot API: ChecklistTasksAdded.
type ChecklistTasksAdded struct {
	// (Optional) The message with the checklist.
	ChecklistMessage *Message `json:"checklist_message,omitempty"`

	Tasks []ChecklistTask `json:"tasks"`
}
