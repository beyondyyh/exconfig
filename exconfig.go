package exconfig

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-hclog"
)

const (
	defaultRetryTimes    = 10               // service discovery fail retry times
	defaultDiscoverySpan = 60 * time.Second // service discovery time span
)

var (
	ErrNil       = errors.New("consul api nil returned")
	ErrNotExists = errors.New("key not exists")
)

// Config...
type Config struct {
	ConsulServerAddr string
	Datacenter       string
	KeyPrefix        string
	DiscoverySpan    time.Duration
	Logger           hclog.Logger
}

// defaultLogger used to generate a Config instance with deault config
func DefaultConfig() *Config {
	cfg := &Config{
		ConsulServerAddr: "http://127.0.0.1:8500",
		Datacenter:       "dc1",
		KeyPrefix:        "hello",
		DiscoverySpan:    defaultDiscoverySpan,
		Logger:           defaultLogger(),
	}
	return cfg
}

// defaultLogger used to generate a logger with deault config
func defaultLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name: "exconfig",
	})
}

// Manifest...
type Manifest struct {
	cfg       *Config
	logger    hclog.Logger
	apiConfig *api.Config
	apiClient *api.Client
	waitIndex uint64
	muLock    sync.RWMutex
	stop      chan struct{}
	warehouse map[string]*api.KVPair
}

// New returns instance of Manifest
func New(cfg *Config) (manifest *Manifest, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic info: %v", r)
		}
	}()

	// bootstrap the config
	defConfig := DefaultConfig()

	if cfg.ConsulServerAddr == "" {
		cfg.ConsulServerAddr = defConfig.ConsulServerAddr
	}
	if cfg.Datacenter == "" {
		cfg.Datacenter = defConfig.Datacenter
	}
	if cfg.KeyPrefix == "" {
		cfg.KeyPrefix = defConfig.KeyPrefix
	}
	if cfg.DiscoverySpan == 0 {
		cfg.DiscoverySpan = defConfig.DiscoverySpan
	}
	if cfg.Logger == nil {
		cfg.Logger = defConfig.Logger
	}

	apiConfig := &api.Config{
		Address:    cfg.ConsulServerAddr,
		Datacenter: cfg.Datacenter,
	}
	apiClient, err := api.NewClient(apiConfig)
	if err != nil {
		return
	}

	manifest = &Manifest{
		cfg:       cfg,
		logger:    cfg.Logger,
		apiConfig: apiConfig,
		apiClient: apiClient,
		muLock:    sync.RWMutex{},
		stop:      make(chan struct{}), // blocking
		warehouse: make(map[string]*api.KVPair, 0),
	}
	cfg.Logger = nil

	go manifest.discovery()
	runtime.Gosched()

	return manifest, nil
}

func (m *Manifest) Close() {
	m.stop <- struct{}{}
}

// GenerateEnv used to print consul-api env
func (m *Manifest) GenerateEnv() []string {
	env := m.apiConfig.GenerateEnv()
	env = append(env,
		fmt.Sprintf("Datacenter=%s", m.apiConfig.Datacenter),
		fmt.Sprintf("KeyPrefix=%s", m.cfg.KeyPrefix),
	)
	return env
}

// discovery service discover
func (m *Manifest) discovery() {
	var retryTimes int
	for {
		select {
		case <-m.stop:
			m.logger.Info("service stop", "keyPrefix", m.cfg.KeyPrefix)
			m.warehouse = nil
			close(m.stop)
			return
		default:
			options := &api.QueryOptions{
				WaitIndex: m.waitIndex,
				WaitTime:  m.cfg.DiscoverySpan,
			}
			kvPairs, meta, err := m.apiClient.KV().List(m.cfg.KeyPrefix, options)
			if err != nil {
				m.logger.Warn("service update fail", "keyPrefix", m.cfg.KeyPrefix, "error", err)
				if retryTimes < defaultRetryTimes {
					retryTimes++
				}
				time.Sleep(time.Second * time.Duration(retryTimes))
				continue
			}

			retryTimes = 0
			if meta.LastIndex != m.waitIndex {
				var topic string
				if m.waitIndex == 0 {
					topic = "service load success"
				} else {
					topic = "service update success"
				}
				m.logger.Info(topic, "keyPrefix", m.cfg.KeyPrefix)
				m.waitIndex = meta.LastIndex
				m.setWarehouse(kvPairs)
			}
		}
	}
}

// setWarehouse used to set kv pairs to warehouse
func (m *Manifest) setWarehouse(kvPairs api.KVPairs) {
	m.muLock.Lock()
	defer m.muLock.Unlock()

	for _, pair := range kvPairs {
		m.warehouse[pair.Key] = pair
	}
}

// Acquire used to acquire the value by key
func (m *Manifest) Acquire(key string) (*api.KVPair, error) {
	m.muLock.RLock()
	defer m.muLock.RUnlock()

	key = m.cfg.KeyPrefix + "/" + strings.TrimLeft(key, "/")
	reply, ok := m.warehouse[key]
	if !ok {
		return nil, ErrNotExists
	}

	if reply == nil {
		return nil, ErrNil
	}

	return reply, nil
}
