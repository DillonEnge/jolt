package messagequeue

import (
	"context"
	"fmt"
	"sync"

	"github.com/DillonEnge/jolt/database"
)

type Store struct {
	topics map[string]topic
	mu     sync.RWMutex
}

type topic struct {
	cs []TopicChan
	ms []*database.Message
}

type TopicChan chan *database.Message

func NewStore() *Store {
	return &Store{
		topics: make(map[string]topic),
	}
}

func (s *Store) AddTopic(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.topics[name]; !ok {
		s.topics[name] = topic{
			cs: make([]TopicChan, 0),
			ms: make([]*database.Message, 0),
		}
	}
}

func (s *Store) RemoveTopic(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.topics[name]; ok {
		for _, v := range t.cs {
			close(v)
		}
		delete(s.topics, name)
	}
}

func (s *Store) Subscribe(ctx context.Context, topicName string) (TopicChan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.topics[topicName]; ok {
		c := make(TopicChan)
		t.cs = append(t.cs, c)
		s.topics[topicName] = t

		go func() {
			for _, v := range t.ms {
				select {
				case <-ctx.Done():
					break
				}

				c <- v
			}
		}()
		return c, nil
	}

	return nil, fmt.Errorf("topic not found: %s", topicName)
}

func (s *Store) Unsubscribe(topicName string, tc TopicChan) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.topics[topicName]; ok {
		for i, c := range t.cs {
			if c == tc {
				drainChan(c)
				t.cs = append(t.cs[:i], t.cs[i+1:]...)
				s.topics[topicName] = t
			}
		}
	}
}

func (s *Store) Publish(topicName string, message *database.Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if t, ok := s.topics[topicName]; ok {
		for _, c := range t.cs {
			go func() {
				c <- message
			}()
		}
	}
}

func drainChan(c chan *database.Message) {
	m := &database.Message{}
	for m != nil {
		m = <-c
	}
	close(c)
}
