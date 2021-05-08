package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"gopkg.in/yaml.v2"
)

type LocalLoader struct {
	path          string
	isDir         bool
	subscriptions []chan<- []*Message
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

func (ll *LocalLoader) Subscribe(mnfsts chan<- []*Message) (Subscription, error) {
	ll.subscriptions = append(ll.subscriptions, mnfsts)
	// @TODO Implemented unsubscribe and error methods 
	return nil, nil
}

func (ll *LocalLoader) Start() error {
	msgs := []*Message{}
	if ll.isDir {
		err := filepath.Walk(ll.path, func(fp string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}
			
			//@TODO Read only *.yml and *.yaml

			msgs = append(msgs, ll.buildMessage(fp))
			return nil
		})

		if err != nil {
			return err
		}
	} else {
		msgs = append(msgs, ll.buildMessage(ll.path))
	}
	
	// Send messages to subcripted
	for _, s := range(ll.subscriptions) {
		go func(sb chan<- []*Message) {
			sb <- msgs
		}(s)
	} 

	return nil
}

func (ll *LocalLoader) buildMessage(fp string) *Message {
	msg := &Message{
		Loader: "LocalLoader",
		Action: CreateAction,
	}

	data, err := ioutil.ReadFile(fp)
	if err != nil {
		msg.Err = err
		return msg
	}

	err = yaml.Unmarshal(data, &msg.Manifest)
	if err != nil {
		msg.Err = err
		return msg
	}

	return msg
}

