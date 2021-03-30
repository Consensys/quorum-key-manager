package node

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/proxy"
// 	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
// )

// var NodeKind manifest.Kind = "Node"

// type NodeSpecs struct {
// 	Addr  string        `json:"addr"`
// 	Proxy *proxy.Config `json:"proxy,omitempty"`
// }

// type manager struct {
// 	mux   sync.RWMutex
// 	nodes map[string]*nodeBundle
// }

// type nodeBundle struct {
// 	manifest *manifest.Manifest

// 	transport http.RoundTripper
// 	proxy     http.Handler
// 	req       *http.Request

// 	err error
// }

// func (node *nodeBundle) Request() *http.Request {
// 	return node.req
// }

// func (node *nodeBundle) Transport() http.RoundTripper {
// 	return node.transport
// }

// func (node *nodeBundle) Proxy() http.Handler {
// 	return node.proxy
// }

// func New() Manager {
// 	return &manager{
// 		mux:   sync.RWMutex{},
// 		nodes: make(map[string]*nodeBundle),
// 	}
// }

// func (m *manager) Load(ctx context.Context, mnfsts ...*manifest.Manifest) error {
// 	m.mux.Lock()
// 	defer m.mux.Unlock()
// 	for _, mnf := range mnfsts {
// 		if err := m.load(ctx, mnf); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (m *manager) GetNode(_ context.Context, name string) (Node, error) {
// 	m.mux.RLock()
// 	defer m.mux.RUnlock()
// 	if nodeBundle, ok := m.nodes[name]; ok {
// 		return nodeBundle, nil
// 	}

// 	return nil, fmt.Errorf("node not found")
// }

// func (m *manager) List(_ context.Context) ([]string, error) {
// 	m.mux.RLock()
// 	defer m.mux.RUnlock()

// 	nodeNames := []string{}
// 	for name := range m.nodes {
// 		nodeNames = append(nodeNames, name)
// 	}

// 	return nodeNames, nil
// }

// func (m *manager) load(_ context.Context, mnf *manifest.Manifest) error {
// 	if _, ok := m.nodes[mnf.Name]; ok {
// 		return fmt.Errorf("node %q already exist", mnf.Name)
// 	}

// 	switch mnf.Kind {
// 	case NodeKind:
// 		node := new(nodeBundle)
// 		m.nodes[mnf.Name] = node

// 		specs := new(NodeSpecs)
// 		if err := mnf.UnmarshalSpecs(specs); err != nil {
// 			node.err = err
// 			return err
// 		}

// 		prxy, err := proxy.New(specs.Proxy, nil, nil)
// 		if err != nil {
// 			node.err = err
// 			return err
// 		}

// 		req, err := http.NewRequest(http.MethodPost, specs.Addr, nil)
// 		if err != nil {
// 			node.err = err
// 			return err
// 		}

// 		node.manifest = mnf
// 		node.proxy = prxy
// 		node.transport = prxy.Transport
// 		node.req = req
// 	default:
// 		return fmt.Errorf("invalid manifest kind %s", mnf.Kind)
// 	}

// 	return nil
// }
