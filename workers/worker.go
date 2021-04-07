/*
Package workers is an internal package. Package realizes the IManager interface.
*/
package workers

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/testobject"
)

const (
	workerStopChLength = 2
	hoursInYear        = 24 * 356 * 100
)

// Worker is an implementation of faces.IWorker Interface .
type Worker struct {
	sync.RWMutex

	// global storage for items
	workBench faces.IWorkBench

	isLast        bool
	isStarted     bool
	globalStop    bool
	id            string
	activeWorkers *int32

	stopCh chan struct{}
	wg     *sync.WaitGroup

	name    faces.Name
	in, out faces.IChan
	errCh   faces.IChan

	typ             faces.ManagerType
	nextManagerName faces.Name

	giveBirth faces.GiveBirth
	handler   faces.IHandler
	tracer    faces.ITrace

	// need to use in test mode
	testObject faces.ITestObject
}

// NewWorker is constructor.
func NewWorker(id string, name faces.Name, wb faces.IWorkBench, in, out, errCh faces.IChan, giveBirth faces.GiveBirth,
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

		stopCh: make(chan struct{}, workerStopChLength),
		tracer: tr,

		testObject: testobject.Empty(),

		workBench: wb,
	}, nil
}

// SetBorderCond is a setter. It sets up the condition to stop, start the worker and log the info about next manager.
func (w *Worker) SetBorderCond(typ faces.ManagerType, isLast bool, nextManagerName faces.Name) {
	w.Lock()
	defer w.Unlock()

	w.isLast = isLast
	w.nextManagerName = nextManagerName
	w.typ = typ
}

// SetTestMode is a simple setter. It attaches the testObject.
func (w *Worker) SetTestMode(testObject faces.ITestObject) {
	w.Lock()
	defer w.Unlock()

	if testObject == nil {
		testObject = testobject.Empty()
	}

	w.testObject = testObject
}

// GetBorderCond is a getter. It returns the condition to stop, start the worker and log the info about next manager.
func (w *Worker) GetBorderCond() (faces.Name, faces.ManagerType, bool) {
	return w.name, w.typ, w.isLast
}

// IsLast is a simple getter. It returns flag "is Manager last".
func (w *Worker) IsLast() bool {
	return w.isLast
}

// Name is a simple getter. It returns Manager name.
func (w *Worker) Name() faces.Name {
	return w.name
}

// ID is a simple getter. It returns Worker ID.
func (w *Worker) ID() string {
	return w.id
}

// Stop stops the worker.
func (w *Worker) Stop() {
	if w.isStarted {
		w.stopCh <- struct{}{}
	}
}

func (w *Worker) logf(format string, a ...interface{}) {
	if w.tracer != nil {
		w.tracer.LazyPrintf(format, a...)
	}
}

func (w *Worker) startHandler(ctx context.Context) error {
	if !w.testObject.IsTestMode() {
		// simple start if not it's test mode
		return w.handler.Start(ctx)
	}

	var err error

	values := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(w.testObject.TestObject())}
	st := reflect.TypeOf(w.handler)

	if _, ok := st.MethodByName(faces.StartTestHandlerPrefix + w.testObject.Suffix()); ok {
		res := reflect.ValueOf(w.handler).MethodByName(faces.StartTestHandlerPrefix + w.testObject.Suffix()).Call(values)

		if len(res) > 0 && !res[0].IsNil() {
			err = res[0].Interface().(error)
		}

		return err
	}

	if _, ok := st.MethodByName(faces.StartTestHandlerPrefix); ok {
		res := reflect.ValueOf(w.handler).MethodByName(faces.StartTestHandlerPrefix).Call(values)

		if len(res) > 0 && !res[0].IsNil() {
			err = res[0].Interface().(error)
		}

		return err
	}

	// simple start if not it's test mode
	return nil
}

// Start runs the worker.
func (w *Worker) Start(ctx context.Context) error {
	if err := w.startHandler(ctx); err != nil {
		return err
	}

	w.job(ctx)

	return nil
}

func (w *Worker) stopHandler(ctx context.Context) {
	if w.testObject.IsTestMode() {
		values := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(w.testObject.TestObject())}
		st := reflect.TypeOf(w.handler)

		if _, ok := st.MethodByName(faces.StopTestHandlerPrefix + w.testObject.Suffix()); ok {
			reflect.ValueOf(w.handler).MethodByName(faces.StopTestHandlerPrefix + w.testObject.Suffix()).Call(values)

			return
		}

		if _, ok := st.MethodByName(faces.StopTestHandlerPrefix); ok {
			reflect.ValueOf(w.handler).MethodByName(faces.StopTestHandlerPrefix).Call(values)

			return
		}
	}

	// simple start if not it's test mode
	w.handler.Stop(ctx)
}

func (w *Worker) job(ctx context.Context) {
	w.wg.Add(1)

	dur := w.handler.TickerDuration()
	if dur == time.Duration(0) {
		dur = time.Hour * hoursInYear // 100 years by default
	}

	go func(ticker *time.Ticker) {
		defer func() {
			ticker.Stop()
			w.stopHandler(ctx)
			w.isStarted = false
			w.wg.Done()
		}()

		w.isStarted = true
		w.logf("%s is started", w.id)

		for {
			if w.globalStop {
				return
			}

			select {
			case <-ctx.Done(): // global context
				w.logf("%s is stopped by global context", w.id)

				return
			case <-w.stopCh:
				w.logf("%s is stopped by message", w.id)

				return
			case <-ticker.C:
				w.logf("%s ticker is running", w.id)
				w.handler.TickerRun(ctx)
			case i, ok := <-w.in.ChanOut():
				if !ok {
					w.logf("%s channel in is close", w.id)

					return
				}

				if item, err := w.workBench.Get(i); err == nil {
					item.ReceivedFromChannel()
					item.BeforeProcess(w.name)
					item.LogTraceFinishTimef("[%s] time in chan", w.name)
					nextCh, nextName := w.process(ctx, i, item)
					if nextName != "" {
						item.PushedToChannel(nextName)
					}

					// send to next manager or return index to workBench
					if nextCh != nil {
						nextCh.Push(i)
					} else {
						w.workBench.Clean(i)
					}
				} else {
					w.logf("%s gets error for %d workBench.Get: %s", w.id, i, err.Error())
				}
			}
		}
	}(time.NewTicker(dur))
}

func (w *Worker) process(ctx context.Context, index int, item faces.IItem) (faces.IChan, faces.Name) {
	// set up which handler was last.
	item.SetLastHandler(w.name)

	// check that item should be skipped
	find, err := item.NeedToSkip(w)
	if err != nil {
		// no more handlers after that. Fix error.
		logError(w.name, err, item)
		item.AfterProcess(w.name, err)

		if !w.isLast && w.typ != faces.FinalManagerType {
			//item.PushedToChannel(faces.ErrorName)
			return w.errCh, faces.ErrorName
		}

		//pushToNotNilChan(w.errCh, index)

		return w.errCh, ""
	}

	if find {
		item.AfterProcess(w.name, err)

		// needed handler is not found
		if !w.isLast && w.typ != faces.FinalManagerType {
			//item.PushedToChannel(w.nextManagerName)

			return w.out, w.nextManagerName
		}

		//pushToNotNilChan(w.out, index)

		return w.out, ""
	}

	// main action
	err = w.run(ctx, item)
	logError(w.name, err, item)
	item.AfterProcess(w.name, err)
	return w.checkDebriefingOfFlight(err)

	//w.debriefingOfFlight(err, item)
}

func doit(internalErr chan error, handler faces.IHandler, item faces.IItem) {
	defer func() {
		if e := recover(); e != nil {
			internalErr <- errors.Errorf("%+v", e)
		}
	}()

	internalErr <- handler.Run(item)
}

func doitWithTest(internalErr chan error, handler faces.IHandler, item faces.IItem) {
	defer func() {
		if e := recover(); e != nil {
			internalErr <- errors.Errorf("%+v", e)
		}
	}()

	var err error

	values := []reflect.Value{reflect.ValueOf(item), reflect.ValueOf(item.GetTestObject().TestObject())}

	suffix := item.GetTestObject().Suffix()
	st := reflect.TypeOf(handler)

	if _, ok := st.MethodByName(faces.RunTestHandlerPrefix + suffix); ok {
		res := reflect.ValueOf(handler).MethodByName(faces.RunTestHandlerPrefix + suffix).Call(values)
		if len(res) > 0 && !res[0].IsNil() {
			err = res[0].Interface().(error)
		}
		internalErr <- err

		return
	}

	if _, ok := st.MethodByName(faces.RunTestHandlerPrefix); ok {
		res := reflect.ValueOf(handler).MethodByName(faces.RunTestHandlerPrefix).Call(values)
		if len(res) > 0 && !res[0].IsNil() {
			err = res[0].Interface().(error)
		}
		internalErr <- err

		return
	}

	internalErr <- handler.(faces.IHandler).Run(item)
}

func (w *Worker) run(ctx context.Context, item faces.IItem) error {
	atomic.AddInt32(w.activeWorkers, 1)
	defer atomic.AddInt32(w.activeWorkers, -1)

	internalErr := make(chan error, 1)

	if item.GetTestObject() == nil || !item.GetTestObject().IsTestMode() {
		go doit(internalErr, w.handler, item)
	} else {
		go doitWithTest(internalErr, w.handler, item)
	}

	var err error
	select {
	case <-ctx.Done():
		err = errors.New(w.id + " processing is stopped by global context")

		// whole process is stopped. Factor doesn't work more.
		w.globalStop = true

		// stops current item on conveyor
		item.Cancel()
	case <-item.GetContext().Done():
		// it's not obviously, customer handler should check it.
		err = errors.New(w.id + " processing is stopped by item context")
	case err = <-internalErr: /* does nothing */
	}

	return err
}

// logError fix log/debug/trace.
func logError(name faces.Name, err error, item faces.IItem) {
	if err == nil {
		item.LogTraceFinishTimef("[%s] success", name)

		return
	}

	item.LogTraceFinishTimef("[%s] has an error", name)
	item.SetHandlerError(name)
	item.AddError(err)
}

func (w *Worker) checkDebriefingOfFlight(err error) (faces.IChan, faces.Name) {
	switch w.typ {
	case faces.FinalManagerType:
		if !w.isLast {
			return w.out, w.nextManagerName
		}
	case faces.ErrorManagerType:
		return w.out, faces.ErrorName
	case faces.WorkerManagerType:
		if err == nil {
			return w.out, w.nextManagerName
		} else {
			return w.errCh, faces.ErrorName
		}
	}

	return nil, ""
}
