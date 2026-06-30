package queue

import (
	"errors"
	"sync"
)

func New[T any](initCapacity int) (*Queue[T], error) {
	if initCapacity <= 0 {
		return nil, errors.New("queue.New(): initCapacity <= 0")
	}

	return &Queue[T]{
		items: make([]T, initCapacity),
	}, nil
}

// Queue кольцевая растущая очередь
type Queue[T any] struct {
	items      []T
	readIndex  int
	writeIndex int

	mu       sync.Mutex
	defaultV T
}

func (o *Queue[T]) Put(v T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	switch {
	case o.writeIndex == len(o.items): // Дошли до конца массива
		// Если в начале есть место для нового эл-та и еще одного для "зазора" между writeIndex & readIndex
		if o.readIndex > 1 {
			// Пишем в начало массива
			o.writeIndex = 0
		} else {
			// Иначе, расширяем массив (гошный алгоритм расширения)
			o.items = append(o.items, o.defaultV)
			// "Заполняем" массив полностью
			o.items = o.items[:cap(o.items)]
		}

	// Запись слева "догнала" чтение
	// Между индексами записи и чтения оставляем одну пустую ячейку
	case o.writeIndex+1 == o.readIndex:
		// Расширяем массив и копируем данные в правильном порядке

		oldSlice := o.items
		o.items = append(o.items, o.defaultV)

		// Копируем в начало массива эл-ты, которые должны быть прочитаны первыми
		rightChunk := oldSlice[o.readIndex:]
		copy(o.items[0:], rightChunk)

		// Копируем эл-ты, которые писали последними в начало массива
		leftChunk := oldSlice[:o.readIndex-1] // Пустую ячейку не копируем
		copy(o.items[len(rightChunk):], leftChunk)

		// "Заполняем" массив полностью
		o.items = o.items[:cap(o.items)]

		o.readIndex = 0
		o.writeIndex = len(leftChunk) + len(rightChunk)
	}

	o.items[o.writeIndex] = v
	o.writeIndex++
}

func (o *Queue[T]) Get() (T, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Читать нечего, возвращаем пустой результат
	if o.readIndex == o.writeIndex {
		return o.defaultV, false
	}

	if o.writeIndex < o.readIndex {
		if o.readIndex == len(o.items) {
			o.readIndex = 0
		}
	}

	v := o.items[o.readIndex]
	o.readIndex++

	return v, true
}
