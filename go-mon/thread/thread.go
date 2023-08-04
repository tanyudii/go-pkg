package thread

type Event interface {
	GetName() string
	GetErr() error
	GetContent() any
}

type event struct {
	name    string
	content any
	err     error
}

func (e *event) GetName() string {
	return e.name
}

func (e *event) GetErr() error {
	return e.err
}

func (e *event) GetContent() any {
	return e.content
}

func GenericGetContent[T any](e Event) T {
	return e.GetContent().(T)
}

func NewEvent(name string, content interface{}, err error) Event {
	return &event{
		name:    name,
		content: content,
		err:     err,
	}
}

type Thread interface {
	Execute() Event
}

type thread struct {
	exec func() Event
}

func (t *thread) Execute() Event {
	return t.exec()
}

func NewThread(exec func() Event) Thread {
	return &thread{exec}
}

func executeThreads(threads []Thread, ignoreError bool) (map[string]Event, error) {
	mapEventByEventName := make(map[string]Event)
	totalThread := len(threads)
	chanEvent := make(chan Event, totalThread)
	for _, t := range threads {
		go func(thread Thread) {
			chanEvent <- thread.Execute()
		}(t)
	}
	for i := 0; i < totalThread; i++ {
		e := <-chanEvent
		if !ignoreError && e.GetErr() != nil {
			return nil, e.GetErr()
		}
		mapEventByEventName[e.GetName()] = e
	}
	return mapEventByEventName, nil
}

func ExecuteThreads(threads []Thread) (map[string]Event, error) {
	return executeThreads(threads, false)
}

func ExecuteThreadsIgnoreError(threads []Thread) (map[string]Event, error) {
	return executeThreads(threads, true)
}
