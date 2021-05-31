package manifest

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/types"
	"gopkg.in/yaml.v2"
)

const ManagerID = "LocalManager"

type Config struct {
	Path string
}

type LocalManager struct {
	path          string
	isDir         bool
	subscriptions []*subscription
}

func NewLocalManager(cfg *Config) (*LocalManager, error) {
	fs, err := os.Stat(cfg.Path)
	if err == nil {
		return &LocalManager{
			path:  cfg.Path,
			isDir: fs.IsDir(),
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
}

func (sub *subscription) Unsubscribe() error {
	close(sub.errors)
	return nil
}

func (sub *subscription) Error() <-chan error { return sub.errors }

func (ll *LocalManager) Subscribe(kinds []manifest.Kind, messages chan<- []Message) (Subscription, error) {
	sub := &subscription{
		messages: messages,
		errors:   make(chan error),
	}
	if kinds != nil {
		sub.kinds = make(map[manifest.Kind]struct{})
		for _, kind := range kinds {
			sub.kinds[kind] = struct{}{}
		}
	}

	ll.subscriptions = append(ll.subscriptions, sub)

	return sub, nil
}

func (ll *LocalManager) Start(context.Context) error {
	msgs := []Message{}
	if ll.isDir {
		err := filepath.Walk(ll.path, func(fp string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(fp) == ".yml" || filepath.Ext(fp) == ".yaml" {
				msgs = append(msgs, ll.buildMessages(fp)...)
			}

			return nil
		})

		if err != nil {
			return err
		}
	} else {
		msgs = append(msgs, ll.buildMessages(ll.path)...)
	}

	for _, sub := range ll.subscriptions {
		go func(sub *subscription) {
			var submsgs []Message
			for _, msg := range msgs {
				if sub.kinds == nil {
					submsgs = append(submsgs, msg)
					continue
				}

				if _, ok := sub.kinds[msg.Manifest.Kind]; ok {
					submsgs = append(submsgs, msg)
				}
			}
			sub.messages <- submsgs
		}(sub)
	}

	return nil
}

func (ll *LocalManager) buildMessages(fp string) []Message {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return []Message{{
			Loader: ManagerID,
			Action: CreateAction,
			Err:    err,
		}}
	}

	mnf := &manifest.Manifest{}
	if err = yaml.Unmarshal(data, mnf); err == nil {
		return []Message{{
			Loader:   ManagerID,
			Action:   CreateAction,
			Manifest: mnf,
		}}
	}

	mnfs := []*manifest.Manifest{}
	if err = yaml.Unmarshal(data, &mnfs); err == nil {
		msgs := []Message{}
		for _, mnf := range mnfs {
			msgs = append(msgs, Message{
				Loader:   ManagerID,
				Action:   CreateAction,
				Manifest: mnf,
			})
		}
		return msgs
	}

	return []Message{{
		Loader: ManagerID,
		Action: CreateAction,
		Err:    errors.InvalidFormatError("cannot read manifest file %s", fp),
	}}
}

func (ll *LocalManager) Stop(context.Context) error { return nil }
func (ll *LocalManager) Error() error               { return nil }
func (ll *LocalManager) Close() error               { return nil }
