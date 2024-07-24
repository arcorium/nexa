package forward

import (
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
  "maps"
)

func NewConfig(allowEmpty bool, opts ...ConfigOption) Config {
  conf := Config{
    allowEmpty: allowEmpty,
  }

  for _, opt := range opts {
    opt(&conf)
  }
  return conf
}

// WithAdditionalData append new data into outgoing context
func WithAdditionalData(data map[string]string) ConfigOption {
  return func(conf *Config) {
    conf.additional = data
  }
}

// WithExcludeKeys only allow metadata with those keys to be forwarded, when used with WithIncludeKeys
// The last called function will override the last one
func WithExcludeKeys(keys ...string) ConfigOption {
  return func(conf *Config) {
    conf.isInclude = false
    conf.keys = keys
  }
}

// WithIncludeKeys only allow metadata with those keys to be forwarded, when used with WithExcludeKeys
// The last called function will override the last one
func WithIncludeKeys(keys ...string) ConfigOption {
  return func(conf *Config) {
    conf.isInclude = true
    conf.keys = keys
  }
}

// WithExpectedKeys will make sure the metadata has this keys
func WithExpectedKeys(keys ...string) ConfigOption {
  return func(conf *Config) {
    conf.expectedKeys = keys
  }
}

type ConfigOption = func(config *Config)

type Config struct {
  allowEmpty   bool
  additional   map[string]string
  expectedKeys []string // keys on metadata that should present on metadata, when the keys aren't present it will return error
  isInclude    bool     // Determine if keys is include or exclude
  keys         []string // Empty means to either forward or remove all the metadata
}

func (c *Config) HasKey() bool {
  return len(c.keys) > 0
}

func (c *Config) HasAdditionalData() bool {
  return c.additional != nil
}

func (c *Config) determineMetadata(old metadata.MD) metadata.MD {
  var newMd = make(metadata.MD)
  if c.HasKey() {
    // Expected keys only
    if c.isInclude {
      // Only forward expected key
      for _, v := range c.keys {
        get := old.Get(v)
        newMd.Set(v, get...)
      }
    } else {
      // Copy and delete except keys
      maps.Copy(newMd, old)
      for _, v := range c.keys {
        delete(newMd, v)
      }
    }
    return newMd
  }
  maps.Copy(newMd, old)
  return newMd
}

func (c *Config) getNewMetadata(old metadata.MD) metadata.MD {
  newMD := c.determineMetadata(old)
  if c.HasAdditionalData() {
    c.addAdditionalData(newMD)
  }
  return newMD
}

func (c *Config) addAdditionalData(md metadata.MD) {
  for k, v := range c.additional {
    md.Set(k, v)
  }
}

func (c *Config) validate(md metadata.MD) error {
  // Check expected keys
  for _, key := range c.expectedKeys {
    _, ok := md[key]
    if !ok {
      return status.Errorf(codes.InvalidArgument, "expected metadata with key: %s", key)
    }
  }
  return nil
}
