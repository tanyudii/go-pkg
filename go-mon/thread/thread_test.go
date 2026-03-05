package thread

import (
	"errors"
	"sync"
	"testing"
)

func TestNewEventAndGetters(t *testing.T) {
	event := NewEvent("test-event", "test-content", nil)

	if event.GetName() != "test-event" {
		t.Errorf("GetName() = %s; want test-event", event.GetName())
	}

	if event.GetContent() != "test-content" {
		t.Errorf("GetContent() = %v; want test-content", event.GetContent())
	}

	if event.GetErr() != nil {
		t.Errorf("GetErr() = %v; want nil", event.GetErr())
	}
}

func TestNewEventWithError(t *testing.T) {
	expectedErr := errors.New("test error")
	event := NewEvent("error-event", nil, expectedErr)

	if event.GetErr() != expectedErr {
		t.Errorf("GetErr() = %v; want %v", event.GetErr(), expectedErr)
	}
}

func TestGenericGetContent(t *testing.T) {
	event := NewEvent("string-event", "hello world", nil)

	content := GenericGetContent[string](event)
	if content != "hello world" {
		t.Errorf("GenericGetContent[string]() = %s; want hello world", content)
	}
}

func TestNewThreadAndExecute(t *testing.T) {
	expectedContent := "thread result"
	execCount := 0

	thread := NewThread(func() Event {
		execCount++
		return NewEvent("test", expectedContent, nil)
	})

	result := thread.Execute()

	if result.GetName() != "test" {
		t.Errorf("Execute().GetName() = %s; want test", result.GetName())
	}

	if result.GetContent() != expectedContent {
		t.Errorf("Execute().GetContent() = %v; want %v", result.GetContent(), expectedContent)
	}

	if execCount != 1 {
		t.Errorf("Thread executed %d times; want 1", execCount)
	}
}

func TestExecuteThreads(t *testing.T) {
	threads := []Thread{
		NewThread(func() Event {
			return NewEvent("thread1", "result1", nil)
		}),
		NewThread(func() Event {
			return NewEvent("thread2", "result2", nil)
		}),
		NewThread(func() Event {
			return NewEvent("thread3", "result3", nil)
		}),
	}

	results, err := ExecuteThreads(threads)

	if err != nil {
		t.Fatalf("ExecuteThreads() error = %v; want nil", err)
	}

	if len(results) != 3 {
		t.Errorf("ExecuteThreads() results length = %d; want 3", len(results))
	}

	if results["thread1"].GetContent() != "result1" {
		t.Errorf("results[\"thread1\"].GetContent() = %v; want result1", results["thread1"].GetContent())
	}

	if results["thread2"].GetContent() != "result2" {
		t.Errorf("results[\"thread2\"].GetContent() = %v; want result2", results["thread2"].GetContent())
	}

	if results["thread3"].GetContent() != "result3" {
		t.Errorf("results[\"thread3\"].GetContent() = %v; want result3", results["thread3"].GetContent())
	}
}

func TestExecuteThreadsIgnoreError(t *testing.T) {
	expectedErr := errors.New("thread error")

	threads := []Thread{
		NewThread(func() Event {
			return NewEvent("thread1", "result1", nil)
		}),
		NewThread(func() Event {
			return NewEvent("thread2", nil, expectedErr)
		}),
		NewThread(func() Event {
			return NewEvent("thread3", "result3", nil)
		}),
	}

	results, err := ExecuteThreadsIgnoreError(threads)

	if err != nil {
		t.Fatalf("ExecuteThreadsIgnoreError() error = %v; want nil", err)
	}

	if len(results) != 3 {
		t.Errorf("ExecuteThreadsIgnoreError() results length = %d; want 3", len(results))
	}

	// Verify that all threads executed
	if results["thread1"].GetErr() != nil {
		t.Error("thread1 should not have error")
	}

	if results["thread2"].GetErr() != expectedErr {
		t.Errorf("thread2 error = %v; want %v", results["thread2"].GetErr(), expectedErr)
	}

	if results["thread3"].GetErr() != nil {
		t.Error("thread3 should not have error")
	}
}

func TestExecuteThreadsWithError(t *testing.T) {
	expectedErr := errors.New("thread error")

	threads := []Thread{
		NewThread(func() Event {
			return NewEvent("thread1", "result1", nil)
		}),
		NewThread(func() Event {
			return NewEvent("thread2", nil, expectedErr)
		}),
	}

	results, err := ExecuteThreads(threads)

	if err == nil {
		t.Error("ExecuteThreads() error = nil; want error")
	}

	if results != nil {
		t.Error("ExecuteThreads() results should be nil on error")
	}
}

func TestExecuteThreadsConcurrent(t *testing.T) {
	// Test that threads actually run concurrently
	var mu sync.Mutex
	var execOrder []string

	threads := []Thread{
		NewThread(func() Event {
			mu.Lock()
			execOrder = append(execOrder, "thread1")
			mu.Unlock()
			return NewEvent("thread1", nil, nil)
		}),
		NewThread(func() Event {
			mu.Lock()
			execOrder = append(execOrder, "thread2")
			mu.Unlock()
			return NewEvent("thread2", nil, nil)
		}),
		NewThread(func() Event {
			mu.Lock()
			execOrder = append(execOrder, "thread3")
			mu.Unlock()
			return NewEvent("thread3", nil, nil)
		}),
	}

	results, err := ExecuteThreads(threads)

	if err != nil {
		t.Fatalf("ExecuteThreads() error = %v; want nil", err)
	}

	if len(results) != 3 {
		t.Errorf("ExecuteThreads() results length = %d; want 3", len(results))
	}

	if len(execOrder) != 3 {
		t.Errorf("ExecuteThreads() execOrder length = %d; want 3", len(execOrder))
	}
}

func TestExecuteThreadsWithEmptySlice(t *testing.T) {
	threads := []Thread{}

	results, err := ExecuteThreads(threads)

	if err != nil {
		t.Fatalf("ExecuteThreads() error = %v; want nil", err)
	}

	if len(results) != 0 {
		t.Errorf("ExecuteThreads() results length = %d; want 0", len(results))
	}
}
