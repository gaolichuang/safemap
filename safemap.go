package safemap

type SafeMap interface {
	Insert(string, interface{})
	Delete(string)
	Find(string) (interface{}, bool)
	Len() int
	Update(string, UpdateFunc)
	Close() map[string]interface{}
}

func New() SafeMap {
	cmdC := make(cmdChannel)
	go cmdC.run()
	return cmdC
}

type UpdateFunc func(interface{}, bool) interface{}

type cmdChannel chan CommandData

type CommandData struct {
	action  commandAction
	key     string
	value   interface{}
	result  chan<- interface{}
	data    chan<- map[string]interface{}
	updater UpdateFunc
}

type commandAction int

const (
	remove commandAction = iota
	end
	find
	insert
	length
	update
)

func (cmdC cmdChannel) Insert(key string, value interface{}) {
	cmdC <- CommandData{action: insert, key: key, value: value}
}

func (cmdC cmdChannel) Delete(key string) {
	cmdC <- CommandData{action: remove, key: key}
}

type findResult struct {
	value interface{}
	found bool
}

func (cmdC cmdChannel) Find(key string) (interface{}, bool) {
	reply := make(chan interface{})
	cmdC <- CommandData{action: find, key: key, result: reply}
	result := (<-reply).(findResult)
	return result.value, result.found
}

func (cmdC cmdChannel) Len() int {
	reply := make(chan interface{})
	cmdC <- CommandData{action: length, result: reply}
	return (<-reply).(int)
}

func (cmdC cmdChannel) Update(key string, updater UpdateFunc) {
	cmdC <- CommandData{action: update, key: key, updater: updater}
}

func (cmdC cmdChannel) Close() map[string]interface{} {
	reply := make(chan map[string]interface{})
	cmdC <- CommandData{action: end, data: reply}
	return <-reply
}

func (cmdC cmdChannel) run() {
	store := make(map[string]interface{})
	for cmd := range cmdC {
		switch cmd.action {
		case insert:
			store[cmd.key] = cmd.value
		case remove:
			delete(store, cmd.key)
		case find:
			value, found := store[cmd.key]
			cmd.result <- findResult{value, found}
		case length:
			cmd.result <- len(store)
		case update:
			value, found := store[cmd.key]
			store[cmd.key] = cmd.updater(value, found)
		case end:
			close(cmdC)
			cmd.data <- store
		}
	}

}
