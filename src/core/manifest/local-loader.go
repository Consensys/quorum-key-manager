package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"gopkg.in/yaml.v2"
)

const LoaderID = "LocalLoader"

type LocalLoader struct {
	path          string
	isDir         bool
	subscriptions []chan<- []Message
}

func NewLocalLoader(path string) (*LocalLoader, error) {
	fs, err := os.Stat(path)
	if err == nil {
		return &LocalLoader{
			path:  path,
			isDir: fs.IsDir(),
		}, nil
	}

	if os.IsNotExist(err) {
		return nil, errors.InvalidParameterError("folder or file does not exists. %s", path)
	}

	return nil, err
}

func (ll *LocalLoader) Subscribe(mnfsts chan<- []Message) (Subscription, error) {
	ll.subscriptions = append(ll.subscriptions, mnfsts)
	// @TODO Implemented unsubscribe and error methods
	return nil, nil
}

func (ll *LocalLoader) Start() error {
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

	for _, s := range ll.subscriptions {
		go func(sb chan<- []Message) {
			sb <- msgs
		}(s)
	}

	return nil
}

func (ll *LocalLoader) buildMessages(fp string) []Message {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return []Message{{
			Loader: LoaderID,
			Action: CreateAction,
			Err:    err,
		}}
	}

	mnf := &Manifest{}
	if err = yaml.Unmarshal(data, mnf); err == nil {
		return []Message{{
			Loader:   LoaderID,
			Action:   CreateAction,
			Manifest: mnf,
		}}
	}

	mnfs := []*Manifest{}
	if err = yaml.Unmarshal(data, &mnfs); err == nil {
		msgs := []Message{}
		for _, mnf := range mnfs {
			msgs = append(msgs, Message{
				Loader:   LoaderID,
				Action:   CreateAction,
				Manifest: mnf,
			})
		}
		return msgs
	}

	return []Message{{
		Loader: LoaderID,
		Action: CreateAction,
		Err:    errors.InvalidFormatError("cannot read manifest file %s", fp),
	}}
}
