package manager

import (
	"context"
	"fmt"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log-old"
	manifest "github.com/consensysquorum/quorum-key-manager/src/manifests/types"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

const ManagerID = "LocalManager"

type Config struct {
	Path string
}

type LocalManager struct {
	path   string
	isDir  bool
	isLive bool

	msgs []Message

	loaded chan struct{}
	err    error
	logger log.Logger
}

func NewLocalManager(cfg *Config, logger log.Logger) (*LocalManager, error) {
	fs, err := os.Stat(cfg.Path)
	if err == nil {
		return &LocalManager{
			path:   cfg.Path,
			loaded: make(chan struct{}),
			isDir:  fs.IsDir(),
			logger: logger,
		}, nil
	}

	if os.IsNotExist(err) {
		return nil, errors.InvalidParameterError("folder or file does not exists. %s", cfg.Path)
	}

	return nil, err
}

type subscription struct {
	kinds    map[manifest.Kind]struct{}
	messages chan<- []Message
	errors   chan error
	stop     chan struct{}
	done     chan struct{}
	logger   *log_old.Logger
}

func (sub *subscription) Unsubscribe() error {
	close(sub.stop)
	<-sub.done
	close(sub.errors)
	return nil
}

func (sub *subscription) Error() <-chan error { return sub.errors }

func (sub *subscription) inbox(msgs []Message) {
	var submsgs []Message
	for _, msg := range msgs {
		if msg.Err != nil {
			sub.logger.WithError(msg.Err).Error("failed to load manifest")
			continue
		}

		if sub.kinds == nil {
			submsgs = append(submsgs, msg)
			continue
		}

		if _, ok := sub.kinds[msg.Manifest.Kind]; ok {
			submsgs = append(submsgs, msg)
		}
	}

	sub.messages <- submsgs
}

func (ll *LocalManager) Subscribe(kinds []manifest.Kind, messages chan<- []Message) (Subscription, error) {
	sub := &subscription{
		messages: messages,
		errors:   make(chan error, 1),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
		logger:   ll.logger,
	}

	if kinds != nil {
		sub.kinds = make(map[manifest.Kind]struct{})
		for _, kind := range kinds {
			sub.kinds[kind] = struct{}{}
		}
	}

	go ll.processSub(sub)

	return sub, nil
}

func (ll *LocalManager) processSub(sub *subscription) {
	defer close(sub.done)

	select {
	case <-ll.loaded:
		if ll.err != nil {
			sub.errors <- ll.err
		} else {
			sub.inbox(ll.msgs)
		}
	case <-sub.stop:
	}
}

func (ll *LocalManager) load() error {
	ll.logger.WithField("path", ll.path).WithField("isDir", ll.isDir).Debug("reading manifest items")

	if !ll.isDir {
		ll.msgs = append(ll.msgs, ll.buildMessages(ll.path)...)
		return nil
	}

	return filepath.Walk(ll.path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(fp) == ".yml" || filepath.Ext(fp) == ".yaml" {
			ll.msgs = append(ll.msgs, ll.buildMessages(fp)...)
		}

		return nil
	})
}

func (ll *LocalManager) Start(context.Context) error {
	defer func() {
		close(ll.loaded)
		ll.isLive = true
	}()

	ll.err = ll.load()
	return ll.err
}

func (ll *LocalManager) buildMessages(fp string) []Message {
	val := validator.New()
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return []Message{newCreateActionMsg(nil, err)}
	}

	mnf := &manifest.Manifest{}
	if err = yaml.Unmarshal(data, mnf); err == nil {
		if err2 := val.Struct(mnf); err2 != nil {
			return []Message{newCreateActionMsg(nil, err2)}
		}

		return []Message{newCreateActionMsg(mnf, nil)}
	}

	mnfs := []*manifest.Manifest{}
	if err = yaml.Unmarshal(data, &mnfs); err != nil {
		return []Message{newCreateActionMsg(nil, err)}
	}

	msgs := []Message{}
	for _, mnf := range mnfs {
		if err := val.Struct(mnf); err != nil {
			msgs = append(msgs, newCreateActionMsg(nil, err))
		} else {
			msgs = append(msgs, newCreateActionMsg(mnf, nil))
		}
	}

	return msgs
}

func newCreateActionMsg(mnf *manifest.Manifest, err error) Message {
	return Message{
		Loader:   ManagerID,
		Action:   CreateAction,
		Manifest: mnf,
		Err:      err,
	}
}

func (ll *LocalManager) Stop(context.Context) error {
	ll.isLive = false
	return nil
}

func (ll *LocalManager) Error() error { return ll.err }
func (ll *LocalManager) Close() error { return nil }

func (ll *LocalManager) ID() string { return ManagerID }
func (ll *LocalManager) CheckLiveness() error {
	if ll.isLive {
		return nil
	}
	return fmt.Errorf("service %s is not live", ll.ID())
}

func (ll *LocalManager) CheckReadiness() error {
	for _, msg := range ll.msgs {
		if msg.Err != nil {
			return msg.Err
		}
	}

	return ll.Error()
}
