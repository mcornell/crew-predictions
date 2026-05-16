package tasks

import (
	"context"
	"fmt"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CloudTasksConfig captures the GCP routing config needed to enqueue a poll
// task. All fields are required when the real enqueuer is in use.
type CloudTasksConfig struct {
	ProjectID    string // e.g. "crew-predictions"
	Location     string // e.g. "us-east5"
	QueueID      string // e.g. "match-polling"
	TargetURL    string // e.g. "https://crew-predictions-3qbonlobra-ul.a.run.app/admin/poll-scores"
	AdminKey     string // value of X-Admin-Key on the dispatched HTTP request
}

// CloudTasksEnqueuer is the production Enqueuer implementation. It schedules
// a one-shot HTTP POST to TargetURL?matchID=<id> at the given run time. Task
// names are deterministic (matchID + unix-seconds) so duplicate enqueues
// within Cloud Tasks' ~1h dedup window are no-ops.
type CloudTasksEnqueuer struct {
	client *cloudtasks.Client
	cfg    CloudTasksConfig
}

// NewCloudTasksEnqueuer constructs a real Cloud Tasks-backed Enqueuer. The
// client uses Application Default Credentials — in Cloud Run, that's the
// runtime service account. The caller is responsible for closing the
// underlying client via Close() when shutting down.
func NewCloudTasksEnqueuer(ctx context.Context, cfg CloudTasksConfig) (*CloudTasksEnqueuer, error) {
	if cfg.ProjectID == "" || cfg.Location == "" || cfg.QueueID == "" || cfg.TargetURL == "" {
		return nil, fmt.Errorf("cloud tasks config: ProjectID/Location/QueueID/TargetURL all required, got %+v", cfg)
	}
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloud tasks: new client: %w", err)
	}
	return &CloudTasksEnqueuer{client: client, cfg: cfg}, nil
}

// Close releases the underlying gRPC client. Safe to call multiple times.
func (e *CloudTasksEnqueuer) Close() error {
	if e == nil || e.client == nil {
		return nil
	}
	return e.client.Close()
}

func (e *CloudTasksEnqueuer) EnqueuePoll(ctx context.Context, matchID string, runAt time.Time) error {
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", e.cfg.ProjectID, e.cfg.Location, e.cfg.QueueID)
	taskName := fmt.Sprintf("%s/tasks/poll-%s-%d", queuePath, sanitizeMatchID(matchID), runAt.UTC().Unix())
	url := e.cfg.TargetURL + "?matchID=" + matchID
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			Name:        taskName,
			ScheduleTime: timestamppb.New(runAt.UTC()),
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					Headers: map[string]string{
						"X-Admin-Key": e.cfg.AdminKey,
					},
				},
			},
		},
	}
	if _, err := e.client.CreateTask(ctx, req); err != nil {
		// ALREADY_EXISTS within the dedup window is the deterministic-name
		// idempotency path; treat it as success.
		if isAlreadyExists(err) {
			return nil
		}
		return fmt.Errorf("cloud tasks: create task %s: %w", taskName, err)
	}
	return nil
}

// sanitizeMatchID makes a match ID safe to embed in a Cloud Tasks task name.
// Cloud Tasks names must match [A-Za-z0-9_-]+. ESPN match IDs are numeric
// but defensively replace anything outside the safe set.
func sanitizeMatchID(id string) string {
	var sb strings.Builder
	for _, r := range id {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	return sb.String()
}

// isAlreadyExists returns true for the gRPC ALREADY_EXISTS code, which
// Cloud Tasks returns when a task with the same deterministic name was
// recently enqueued (within its ~1h dedup window).
func isAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	// status.FromError is the canonical way, but the cloud tasks SDK wraps
	// errors deeply. A string match on the gRPC code is acceptable here —
	// the error message format is stable and ALREADY_EXISTS is rare enough
	// that a string check won't cause false positives in practice.
	return strings.Contains(err.Error(), "AlreadyExists") || strings.Contains(err.Error(), "code = AlreadyExists")
}
