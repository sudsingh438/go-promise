package promise

import (
	"sync"
)

type Promise struct {
	pending bool
	executor func(resolve func(interface{}), reject func(error))
	result interface{}
	err error
	mutex sync.Mutex
	wg sync.WaitGroup
}

func New(executor func(resolve func(interface{}), reject func(error))) *Promise {
	var promise = &Promise{
		pending:  true,
		executor: executor,
		result:   nil,
		err:      nil,
		mutex:    sync.Mutex{},
		wg:       sync.WaitGroup{},
	}

	promise.wg.Add(1)

	go func() {
		promise.executor(promise.resolve, promise.reject)
	}()

	return promise
}

func (promise *Promise) resolve(resolution interface{}) {
	promise.mutex.Lock()

	if !promise.pending {
		promise.mutex.Unlock()
		return
	}

	switch result := resolution.(type) {
	case *Promise:
		flattenedResult, err := result.Await()
		if err != nil {
			promise.mutex.Unlock()
			promise.reject(err)
			return
		}
		promise.result = flattenedResult
	default:
		promise.result = result
	}
	promise.pending = false

	promise.wg.Done()
	promise.mutex.Unlock()
}

func (promise *Promise) reject(err error) {
	promise.mutex.Lock()
	defer promise.mutex.Unlock()

	if !promise.pending {
		return
	}

	promise.err = err
	promise.pending = false

	promise.wg.Done()
}

func (promise *Promise) Then(fulfillment func(data interface{}) interface{}) *Promise {
	return New(func(resolve func(interface{}), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(err)
			return
		}
		resolve(fulfillment(result))
	})
}

func (promise *Promise) Catch(rejection func(err error) error) *Promise {
	return New(func(resolve func(interface{}), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(rejection(err))
			return
		}
		resolve(result)
	})
}

func (promise *Promise) Await() (interface{}, error) {
	promise.wg.Wait()
	return promise.result, promise.err
}

func Resolve(resolution interface{}) *Promise {
	return New(func(resolve func(interface{}), reject func(error)) {
		resolve(resolution)
	})
}

func Reject(err error) *Promise {
	return New(func(resolve func(interface{}), reject func(error)) {
		reject(err)
	})
}
