package exconfig

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"

	"gitlab.weibo.cn/gdp/gdp/encoding/toml"
	"gitlab.weibo.cn/gdp/libs/set"
)

func Int(reply *api.KVPair, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(reply.Value), 10, 0)
	return int(n), err
}

func Int64(reply *api.KVPair, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(reply.Value), 10, 64)
}

func String(reply *api.KVPair, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return string(reply.Value), nil
}

func Bytes(reply *api.KVPair, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	return reply.Value, nil
}

func Bool(reply *api.KVPair, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(reply.Value))
}

func Strings(reply *api.KVPair, err error, separator string) ([]string, error) {
	var result []string
	if err != nil {
		return result, err
	}
	if value := strings.TrimSpace(string(reply.Value)); len(value) > 0 {
		for _, val := range strings.Split(value, separator) {
			result = append(result, val)
		}
	}
	return result, err
}

func ByteSlices(reply *api.KVPair, err error, separator string) ([][]byte, error) {
	var result [][]byte
	if err != nil {
		return result, err
	}
	if value := strings.TrimSpace(string(reply.Value)); len(value) > 0 {
		for _, val := range strings.Split(value, separator) {
			result = append(result, []byte(val))
		}
	}
	return result, err
}

func Sets(reply *api.KVPair, err error, separator string) (set.Set, error) {
	sets := set.NewSet()
	if err != nil {
		return sets, err
	}
	if value := strings.TrimSpace(string(reply.Value)); len(value) > 0 {
		for _, val := range strings.Split(value, separator) {
			sets.Add(val)
		}
	}
	return sets, nil
}

func Json(reply *api.KVPair, err error, v interface{}) error {
	if err != nil {
		return err
	}
	return json.Unmarshal(reply.Value, v)
}

func Toml(reply *api.KVPair, err error, v interface{}) error {
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(reply.Value), v)
	return err
}

func Yaml(reply *api.KVPair, err error, v interface{}) error {
	if err != nil {
		return err
	}
	return yaml.Unmarshal(reply.Value, v)
}
