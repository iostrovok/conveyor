/*
	Internal package. Package realizes the IWorker interface.
*/
package workers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/iostrovok/conveyor/faces"
)

/*
	....
*/

type Worker struct {
	sync.RWMutex

	name    faces.Name
	id      string
	in, out faces.IChan
	errCh   faces.IChan

	isStarted  bool
	globalStop bool
	stopCh     chan struct{}
	wg         *sync.WaitGroup

	isLast bool
	typ    faces.ManagerType

	giveBirth faces.GiveBirth
	handler   faces.IHandler
	tracer    faces.ITrace

	activeWorkers *int32
}

func NewWorker(id string, name faces.Name, in, out, errCh faces.IChan, giveBirth faces.GiveBirth,
	wg *sync.WaitGroup, tr faces.ITrace, activeWorkers *int32) (faces.IWorker, error) {

	handler, err := giveBirth(name)
	if err != nil {
		return nil, err
	}

	return &Worker{
		id:            id,
		name:          name,
		in:            in,
		out:           out,
		errCh:         errCh,
		giveBirth:     giveBirth,
		handler:       handler,
		wg:            wg,
		activeWorkers: activeWorkers,

		stopCh: make(chan struct{}, 2),
		tracer: tr,
	}, nil
}

func (w *Worker) SetBorderCond(typ faces.ManagerType, isLast bool) {
	w.Lock()
	defer w.Unlock()

	w.isLast = isLast
	w.typ = typ
}

func (w *Worker) IsLast() bool {
	w.RLock()
	defer w.RUnlock()

	return w.isLast
}

func (w *Worker) Name() faces.Name {
	return w.name
}

func (w *Worker) ID() string {
	return w.id
}

func (w *Worker) Stop() {
	if w.isStarted {
		w.stopCh <- struct{}{}
	}
}

func (w *Worker) logTrace(format string, a ...interface{}) {
	if w.tracer != nil {
		w.tracer.LazyPrintf(format, a...)
	}
}

func (w *Worker) Start(ctx context.Context) error {

	if err := w.handler.Start(ctx); err != nil {
		return err
	}

	w.wg.Add(1)
	go func() {
		w.isStarted = true
		defer func() {
			w.handler.Stop()
			w.isStarted = false
			w.wg.Done()
		}()

		w.logTrace("%s is started", w.id)

		for {
			if w.globalStop {
				return
			}

			select {
			case <-ctx.Done(): // global context
				w.logTrace("%s is stopped by global context", w.id)
				return
			case <-w.stopCh:
				w.logTrace("%s is stopped by message", w.id)
				return
			case item, ok := <-w.in.ChanOut():
				if !ok {
					w.logTrace("%s channel in is close", w.id)
					return
				}
				item.LogTraceFinishTime("[%s] time in chan", w.name)

				w.process(ctx, item)
			}
		}
	}()

	return nil
}

func (w *Worker) process(ctx context.Context, item faces.IItem) {

	// set up which handler was last.
	item.SetLastHandler(w.name)

	// check that item should be skipped

	if needToSkip(w.typ, w.name, item) {
		// needed handler is not  found
		if w.isLast && w.typ == faces.WorkerManagerType {
			// no more handlers after that. Fix error.
			err := errors.New("it is the last handler: skipped name [" + string(item.GetSkipToName()) + "] can't be processed")
			logError(w.name, err, item)
			pushToNotNilChan(w.errCh, item)
		} else {
			pushToNotNilChan(w.out, item)
		}

		return
	}

	// main action
	err := w.run(ctx, item)
	logError(w.name, err, item)
	w.debriefingOfFlight(err, item)
}

func doit(internalErr chan error, handler faces.IHandler, item faces.IItem) {
	defer func() {
		if e := recover(); e != nil {
			internalErr <- fmt.Errorf("%+v", e)
		}
	}()

	internalErr <- handler.Run(item)
}

func (w *Worker) run(ctx context.Context, item faces.IItem) error {

	atomic.AddInt32(w.activeWorkers, 1)
	defer atomic.AddInt32(w.activeWorkers, -1)

	internalErr := make(chan error, 1)
	go doit(internalErr, w.handler, item)

	select {
	case <-ctx.Done():
		w.globalStop = true
		return errors.New(w.id + " processing is stopped by global context")
	case <-item.GetContext().Done():
		return errors.New(w.id + " processing is stopped by context")
	case err := <-internalErr:
		return err
	}
}

// logError fix log/debug/trace
func logError(name faces.Name, err error, item faces.IItem) {
	if err == nil {
		item.LogTraceFinishTime("[%s] success", name)
		return
	}

	item.LogTraceFinishTime("[%s] has an error", name)
	item.SetHandlerError(name)
	item.AddError(err)
}

func (w *Worker) debriefingOfFlight(err error, item faces.IItem) {

	switch w.typ {
	case faces.FinalManagerType:
		if !w.isLast {
			pushToNotNilChan(w.out, item)
		}
	case faces.ErrorManagerType:
		pushToNotNilChan(w.out, item)
	default: //  faces.WorkerManagerType
		if err == nil {
			pushToNotNilChan(w.out, item)
		} else {
			pushToNotNilChan(w.errCh, item)
		}
	}
}

func pushToNotNilChan(ch faces.IChan, item faces.IItem) {
	if ch != nil {
		ch.Push(item)
	}
}

func needToSkip(typ faces.ManagerType, name faces.Name, item faces.IItem) bool {

	// never skip system handlers
	if typ != faces.WorkerManagerType {
		item.CleanSkipToName()
		return false
	}

	skipName := item.GetSkipToName()
	if skipName == faces.EmptySkipName {
		// skipping has not requested
		return false
	}

	if skipName == name {
		// Attention, last station! processing...
		item.CleanSkipToName()
		return false
	}

	return true
}
