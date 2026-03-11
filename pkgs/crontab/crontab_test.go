package crontab

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		enableSeconds bool
		wantSeconds   bool
	}{
		{
			name:          "with seconds enabled",
			enableSeconds: true,
			wantSeconds:   true,
		},
		{
			name:          "with seconds disabled",
			enableSeconds: false,
			wantSeconds:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &CronOpts{
				EnableSeconds: tt.enableSeconds,
				Jobs:          []Job{},
			}
			ctx := context.Background()

			ct := New(ctx, opts)

			if ct.Ctx != ctx {
				t.Errorf("New() ctx = %v, want %v", ct.Ctx, ctx)
			}
			if ct.WGroup == nil {
				t.Error("New() WGroup should not be nil")
			}
		})
	}
}

func TestNew_WithJobs(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: true,
		Jobs: []Job{
			{Name: "test-job", Schedule: "*/5 * * * * *", Disable: false},
		},
	}
	ctx := context.Background()

	ct := New(ctx, opts)

	if len(ct.CronServer.Entries()) != 0 {
		t.Errorf("New() should not add jobs, got %d entries", len(ct.CronServer.Entries()))
	}
}

func TestCrontab_StartCron(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()
	ct := New(ctx, opts)

	ct.StartCron()
	time.Sleep(10 * time.Millisecond)

	if len(ct.CronServer.Entries()) != 0 {
		t.Errorf("StartCron() should not add entries, got %d", len(ct.CronServer.Entries()))
	}
	ct.CronServer.Stop()
}

func TestCrontab_StopCron(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()
	ct := New(ctx, opts)

	ct.StartCron()
	time.Sleep(10 * time.Millisecond)

	done := make(chan struct{})
	go func() {
		ct.StopCron()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("StopCron() timed out")
	}
}

func TestCrontab_StopCron_WaitForTasks(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()
	ct := New(ctx, opts)

	ct.StartCron()

	var taskStarted, taskDone bool
	var mu sync.Mutex

	_, err := ct.CronServer.AddFunc("@every 1s", func() {
		mu.Lock()
		taskStarted = true
		mu.Unlock()

		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		taskDone = true
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Failed to add cron job: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	started := taskStarted
	mu.Unlock()

	if !started {
		t.Skip("Task not started yet, skipping wait test")
	}

	done := make(chan struct{})
	go func() {
		ct.StopCron()
		close(done)
	}()

	select {
	case <-done:
		mu.Lock()
		if !taskDone {
			t.Error("StopCron() returned before task completed")
		}
		mu.Unlock()
	case <-time.After(2 * time.Second):
		t.Error("StopCron() timed out waiting for task")
	}
}

func TestCronOpts_Job(t *testing.T) {
	job := Job{
		Name:     "test-job",
		Schedule: "*/5 * * * *",
		Disable:  false,
	}

	if job.Name != "test-job" {
		t.Errorf("Job.Name = %v, want 'test-job'", job.Name)
	}
	if job.Schedule != "*/5 * * * *" {
		t.Errorf("Job.Schedule = %v, want '*/5 * * * *'", job.Schedule)
	}
	if job.Disable != false {
		t.Errorf("Job.Disable = %v, want false", job.Disable)
	}
}

func TestCronOpts_JobDisabled(t *testing.T) {
	job := Job{
		Name:     "disabled-job",
		Schedule: "*/10 * * * *",
		Disable:  true,
	}

	if job.Disable != true {
		t.Errorf("Job.Disable = %v, want true", job.Disable)
	}
}

func TestCrontab_StartStopMultipleTimes(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		ct := New(ctx, opts)
		ct.StartCron()
		time.Sleep(10 * time.Millisecond)
		ct.StopCron()
	}
}

func TestCronOpts_Validate(t *testing.T) {
	tests := []struct {
		name         string
		opts         CronOpts
		wantJobCount int
	}{
		{
			name: "empty jobs",
			opts: CronOpts{
				EnableSeconds: true,
				Jobs:          []Job{},
			},
			wantJobCount: 0,
		},
		{
			name: "single job",
			opts: CronOpts{
				EnableSeconds: true,
				Jobs: []Job{
					{Name: "job1", Schedule: "* * * * *", Disable: false},
				},
			},
			wantJobCount: 1,
		},
		{
			name: "multiple jobs",
			opts: CronOpts{
				EnableSeconds: false,
				Jobs: []Job{
					{Name: "job1", Schedule: "* * * * *", Disable: false},
					{Name: "job2", Schedule: "*/5 * * * *", Disable: false},
					{Name: "job3", Schedule: "*/10 * * * *", Disable: true},
				},
			},
			wantJobCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.opts.Jobs) != tt.wantJobCount {
				t.Errorf("Jobs count = %v, want %v", len(tt.opts.Jobs), tt.wantJobCount)
			}
		})
	}
}

func TestCrontab_ConcurrentAccess(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()
	ct := New(ctx, opts)

	ct.StartCron()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ct.CronServer.Entries()
		}()
	}

	wg.Wait()
	ct.StopCron()
}

func TestCrontab_InvalidSchedule(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()
	ct := New(ctx, opts)

	ct.StartCron()

	_, err := ct.CronServer.AddFunc("invalid schedule", func() {})
	if err == nil {
		t.Error("AddFunc() with invalid schedule should return error")
	}

	ct.StopCron()
}

func TestCrontab_RemoveEntry(t *testing.T) {
	opts := &CronOpts{
		EnableSeconds: false,
		Jobs:          []Job{},
	}
	ctx := context.Background()
	ct := New(ctx, opts)

	ct.StartCron()

	entryID, err := ct.CronServer.AddFunc("@every 1h", func() {})
	if err != nil {
		t.Fatalf("Failed to add cron job: %v", err)
	}

	if len(ct.CronServer.Entries()) != 1 {
		t.Errorf("After AddFunc() entries = %v, want 1", len(ct.CronServer.Entries()))
	}

	ct.CronServer.Remove(entryID)

	if len(ct.CronServer.Entries()) != 0 {
		t.Errorf("After Remove() entries = %v, want 0", len(ct.CronServer.Entries()))
	}

	ct.StopCron()
}
