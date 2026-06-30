package broker

import (
	"errors"
	"fmt"
	"ndm_broker/queue"
	"sync"
	"time"
)

type Queue struct {
	Queue       *queue.Queue[string]
	Subscribers *queue.Queue[chan string]
}

func New(initQueueCapacity int) (*Broker, error) {
	if initQueueCapacity < 1 {
		return nil, errors.New("broker.New(): initQueueCapacity < 1")
	}

	return &Broker{
		queues:            make(map[string]*Queue),
		initQueueCapacity: initQueueCapacity,
	}, nil
}

type Broker struct {
	mu     sync.Mutex
	queues map[string]*Queue

	initQueueCapacity int
}

func (o *Broker) Put(queueName string, message string) error {
	q, err := o.getQueue(queueName)
	if err != nil {
		return fmt.Errorf("Broker.Put(): %w", err)
	}

	q.Queue.Put(message)

	return nil
}

func (o *Broker) Subscribe(queueName string) (<-chan string, error) {
	q, err := o.getQueue(queueName)
	if err != nil {
		return nil, fmt.Errorf("Broker.Subscribe(): %w", err)
	}

	ch := make(chan string)
	q.Subscribers.Put(ch)

	return ch, nil
}

func (o *Broker) getQueue(queueName string) (*Queue, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if q, ok := o.queues[queueName]; ok {
		return q, nil
	}

	q, err := queue.New[string](o.initQueueCapacity)
	if err != nil {
		return nil, fmt.Errorf("Broker.getQueue(): can't create msgs queue: %w", err)
	}

	s, err := queue.New[chan string](o.initQueueCapacity)
	if err != nil {
		return nil, fmt.Errorf("Broker.getQueue(): can't create subs queue: %w", err)
	}

	r := &Queue{q, s}
	o.queues[queueName] = r

	go o.queueHandler(q, s)

	return r, nil
}

func (o *Broker) queueHandler(msgs *queue.Queue[string], subs *queue.Queue[chan string]) {
	delay := 10 * time.Millisecond

	for {
		// Читаем сообщение из очереди
		msg, ok := msgs.Get()
		// Если очередь пустая, спим
		if !ok {
			time.Sleep(delay)
			continue
		}

		for {
			// Получаем подписчика из очереди
			sub, ok := subs.Get()
			// Если подписчиков нет, спим
			if !ok {
				time.Sleep(delay)
				continue
			}

			written := false

			// Каналы синхронные, поэтому может точно узнать, получил ли подписчик сообщение
			select {
			case sub <- msg:
				written = true
			default:
			}

			// Сообщение доставлено, переходим к чтению следующего сообщения из очереди
			if written {
				break
			}

			// Сообщение не доставлено, значит подписчик "отвалился" по таймауту.
			// Переходим к следующему подписчику, чтобы отправить ему текущее сообщение.
		}
	}
}
