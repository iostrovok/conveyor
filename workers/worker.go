package workers

/*
	Internal package. Package realizes the IWorker interface.
*/

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

type Worker struct {
	sync.RWMutex

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

		stopCh: make(chan struct{}, workerStopChLength),
		tracer: tr,

		testObject: testobject.Empty(),
	}, nil
}

func (w *Worker) SetBorderCond(typ faces.ManagerType, isLast bool, nextManagerName faces.Name) {
	w.Lock()
	defer w.Unlock()

	w.isLast = isLast
	w.nextManagerName = nextManagerName
	w.typ = typ
}

func (w *Worker) SetTestMode(testObject faces.ITestObject) {
	w.Lock()
	defer w.Unlock()

	if testObject == nil {
		testObject = testobject.Empty()
	}

	w.testObject = testObject
}

func (w *Worker) GetBorderCond() (faces.Name, faces.ManagerType, bool) {
	return w.name, w.typ, w.isLast
}

func (w *Worker) IsLast() bool {
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

func (w *Worker) logf(format string, a ...interface{}) {
	if w.tracer != nil {
		w.tracer.LazyPrintf(format, a...)
	}
}

func (w *Worker) startHandler(ctx context.Context) error {
	var err error

	if w.testObject.IsTestMode() {
		values := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(w.testObject.TestObject())}
		st := reflect.TypeOf(w.handler)

		if _, ok := st.MethodByName("StartTest_" + w.testObject.Suffix()); ok {
			values := reflect.ValueOf(w.handler).MethodByName("StartTest_" + w.testObject.Suffix()).Call(values)

			if !values[0].IsNil() {
				err = values[0].Interface().(error)
			}

			return err
		}

		if _, ok := st.MethodByName("StartTest"); ok {
			values := reflect.ValueOf(w.handler).MethodByName("StartTest").Call(values)

			if !values[0].IsNil() {
				err = values[0].Interface().(error)
			}

			return err
		}
	}

	// simple start if not it's test mode
	return w.handler.Start(ctx)
}

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

		if _, ok := st.MethodByName("StopTest_" + w.testObject.Suffix()); ok {
			reflect.ValueOf(w.handler).MethodByName("StopTest_" + w.testObject.Suffix()).Call(values)

			return
		}

		if _, ok := st.MethodByName("StopTest"); ok {
			reflect.ValueOf(w.handler).MethodByName("StopTest").Call(values)

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
			case item, ok := <-w.in.ChanOut():
				if !ok {
					w.logf("%s channel in is close", w.id)

					return
				}

				item.ReceivedFromChannel()
				item.BeforeProcess(w.name)
				item.LogTraceFinishTimef("[%s] time in chan", w.name)
				w.process(ctx, item)
			}
		}
	}(time.NewTicker(dur))
}

func (w *Worker) process(ctx context.Context, item faces.IItem) {
	// set up which handler was last.
	item.SetLastHandler(w.name)

	// check that item should be skipped
	find, err := item.NeedToSkip(w)
	if err != nil {
		// no more handlers after that. Fix error.
		logError(w.name, err, item)
		item.AfterProcess(w.name, err)

		if !w.isLast && w.typ != faces.FinalManagerType {
			item.PushedToChannel(faces.ErrorName)
		}

		pushToNotNilChan(w.errCh, item)

		return
	}

	if find {
		item.AfterProcess(w.name, err)

		// needed handler is not found
		if !w.isLast && w.typ != faces.FinalManagerType {
			item.PushedToChannel(w.nextManagerName)
		}

		pushToNotNilChan(w.out, item)

		return
	}

	// main action
	err = w.run(ctx, item)
	logError(w.name, err, item)
	item.AfterProcess(w.name, err)
	w.debriefingOfFlight(err, item)
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

	if _, ok := st.MethodByName("RunTest_" + suffix); ok {
		values := reflect.ValueOf(handler).MethodByName("RunTest_" + suffix).Call(values)

		if !values[0].IsNil() {
			err = values[0].Interface().(error)
		}
		internalErr <- err

		return
	}

	if _, ok := st.MethodByName("RunTest"); ok {
		values := reflect.ValueOf(handler).MethodByName("RunTest").Call(values)
		if !values[0].IsNil() {
			err = values[0].Interface().(error)
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

func (w *Worker) debriefingOfFlight(err error, item faces.IItem) {
	switch w.typ {
	case faces.FinalManagerType:
		if !w.isLast {
			item.PushedToChannel(w.nextManagerName)
			pushToNotNilChan(w.out, item)
		}
	case faces.ErrorManagerType:
		item.PushedToChannel(faces.ErrorName)
		pushToNotNilChan(w.out, item)
	case faces.WorkerManagerType:
		if err == nil {
			item.PushedToChannel(w.nextManagerName)
			pushToNotNilChan(w.out, item)
		} else {
			item.PushedToChannel(faces.ErrorName)
			pushToNotNilChan(w.errCh, item)
		}
	}
}

func pushToNotNilChan(ch faces.IChan, item faces.IItem) {
	if ch != nil {
		ch.Push(item)
	}
}
