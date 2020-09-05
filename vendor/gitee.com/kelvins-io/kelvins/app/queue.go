package app

import (
	"context"
	"flag"
	"fmt"
	"gitee.com/kelvins-io/common/convert"
	"gitee.com/kelvins-io/common/log"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/internal/config"
	"gitee.com/kelvins-io/kelvins/internal/logging"

	"gitee.com/kelvins-io/kelvins/setup"
	queue_log "github.com/RichardKnop/machinery/v1/log"
	"time"
)

// RunQueueApplication runs queue application.
func RunQueueApplication(application *kelvins.QueueApplication) {
	if application.Name == "" {
		logging.Fatal("Application name can't not be empty")
	}

	flag.Parse()
	application.LoggerRootPath = *loggerPath
	application.Type = kelvins.AppTypeQueue
	err := runQueue(application)
	if err != nil {
		logging.Fatalf("RunQueueApplication err: %v", err)
	}
}

// runQueue runs queue application.
func runQueue(queueApp *kelvins.QueueApplication) error {

	// 1. load config
	err := config.LoadDefaultConfig(queueApp.Application)
	if err != nil {
		return err
	}
	if queueApp.LoadConfig != nil {
		err = queueApp.LoadConfig()
		if err != nil {
			return err
		}
	}

	// 2. init application
	err = initApplication(queueApp.Application)
	if err != nil {
		return err
	}

	// 3. setup vars
	err = setupQueueVars(queueApp)
	if err != nil {
		return err
	}
	if queueApp.SetupVars != nil {
		err = queueApp.SetupVars()
		if err != nil {
			return err
		}
	}

	// 4. apollo hot update listen
	//config.TriggerApolloHotUpdateListen(queueApp.Application)

	// 5. start server
	errorsChan := make(chan error)

	logging.Infof("Start queue consume")
	concurrency := len(queueApp.GetNamedTaskFuncs())
	if kelvins.QueueServerSetting != nil {
		concurrency = kelvins.QueueServerSetting.WorkerConcurrency
	}
	logging.Infof("Count of worker goroutine: %d", concurrency)
	consumerTag := queueApp.Application.Name + convert.Int64ToStr(time.Now().Local().UnixNano())

	var queueList = []string{""}
	queueList = append(queueList, kelvins.QueueServerSetting.CustomQueueList...)

	for _, customQueue := range queueList {
		cTag := consumerTag
		if len(customQueue) > 0 {
			cTag = customQueue + "-" + consumerTag
		}

		logging.Infof("Consumer Tag: %s", cTag)
		worker := queueApp.QueueServer.TaskServer.NewCustomQueueWorker(cTag, concurrency, customQueue)
		worker.LaunchAsync(errorsChan)
	}

	return <-errorsChan
}

// setupQueueVars ...
func setupQueueVars(queueApp *kelvins.QueueApplication) error {
	err := setupCommonVars(queueApp.Application)
	if err != nil {
		return err
	}

	queueApp.QueueLogger, err = log.GetBusinessLogger("queue.consume")
	if err != nil {
		return fmt.Errorf("kelvinslog.GetBusinessLogger: %v", err)
	}
	queue_log.Set(&queueLogger{
		logger: queueApp.QueueLogger,
	})

	if queueApp.GetNamedTaskFuncs == nil && queueApp.RegisterEventHandler == nil {
		return fmt.Errorf("Lack of implement GetNamedTaskFuncs And RegisterEventHandler")
	}
	if kelvins.QueueRedisSetting != nil && kelvins.QueueRedisSetting.Broker != "" {
		queueApp.QueueServer, err = setup.NewRedisQueue(kelvins.QueueRedisSetting, queueApp.GetNamedTaskFuncs())
		if err != nil {
			return err
		}
		return nil
	}
	if kelvins.QueueAliAMQPSetting != nil && kelvins.QueueAliAMQPSetting.VHost != "" {
		queueApp.QueueServer, err = setup.NewAliAMQPQueue(kelvins.QueueAliAMQPSetting, queueApp.GetNamedTaskFuncs())
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Lack of kelvinsQueue* section config")
}

// queueLogger implements machinery log interface.
type queueLogger struct {
	logger *log.LoggerContext
}

// Print uses logger to log info msg.
func (q *queueLogger) Print(a ...interface{}) {
	q.logger.Info(context.Background(), fmt.Sprint(a...))
}

// Printf uses logger to log info msg.
func (q *queueLogger) Printf(format string, a ...interface{}) {
	q.logger.Infof(context.Background(), format, a...)
}

// Println uses logger to log info msg.
func (q *queueLogger) Println(a ...interface{}) {
	q.logger.Info(context.Background(), fmt.Sprint(a...))
}

// Fatal uses kelvins.ErrLogger to log err msg.
func (q *queueLogger) Fatal(a ...interface{}) {
	kelvins.ErrLogger.Error(context.Background(), fmt.Sprint(a...))
}

// Fatalf uses kelvins.ErrLogger to log err msg.
func (q *queueLogger) Fatalf(format string, a ...interface{}) {
	kelvins.ErrLogger.Errorf(context.Background(), format, a...)
}

// Fatalln uses kelvins.ErrLogger to log err msg.
func (q *queueLogger) Fatalln(a ...interface{}) {
	kelvins.ErrLogger.Error(context.Background(), fmt.Sprint(a...))
}

// Panic uses kelvins.ErrLogger to log err msg.
func (q *queueLogger) Panic(a ...interface{}) {
	kelvins.ErrLogger.Error(context.Background(), fmt.Sprint(a...))
}

// Panicf uses kelvins.ErrLogger to log err msg.
func (q *queueLogger) Panicf(format string, a ...interface{}) {
	kelvins.ErrLogger.Errorf(context.Background(), format, a)
}

// Panicln uses kelvins.ErrLogger to log err msg.
func (q *queueLogger) Panicln(a ...interface{}) {
	kelvins.ErrLogger.Error(context.Background(), fmt.Sprint(a...))
}
