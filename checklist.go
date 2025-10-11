package telebot

// Checklist contains information about a checklist.
// Can be sent only in business chats.
type Checklist struct {
	Title                   string          `json:"title"`
	Entities                []MessageEntity `json:"title_entities,omitempty"`
	Tasks                   []ChecklistTask `json:"tasks"`
	OtherCanAddTasks        bool            `json:"other_can_add_tasks,omitempty"`
	OtherCanMarkTasksAsDone bool            `json:"other_can_mark_tasks_as_done,omitempty"`
}

// ChecklistTask describes a task in a checklist
type ChecklistTask struct {
	ID                 int64           `json:"id"`
	Text               string          `json:"text"`
	Entities           []MessageEntity `json:"title_entities,omitempty"`
	CompletedByUser    *User           `json:"completed_by_user,omitempty"`
	CompletionUnixDate int64           `json:"completion_date,omitempty"`
}

// ChecklistTasksDone describes a service message about checklist tasks marked as done or not done.
type ChecklistTasksDone struct {
	ChecklistMessage *Message `json:"checklist_message,omitempty"`
	DoneTasks        []int64  `json:"marked_as_done_task_ids,omitempty"`
	NotDoneTasks     []int64  `json:"marked_as_not_done_task_ids,omitempty"`
}

// ChecklistTasksAdded describes a service message about tasks added to a checklist.
type ChecklistTasksAdded struct {
	Message *Message        `json:"checklist_message,omitempty"`
	Tasks   []ChecklistTask `json:"tasks"`
}
