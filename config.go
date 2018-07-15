package config

import (
	"sync"
	"io"
	"fmt"
	"errors"
)

// package version
const Version = "1.0.1"

// supported config format
const (
	Json = "json"
	Yml  = "yml"
	Yaml = "yaml"
	Toml = "toml"
)

type stringArr []string
type stringMap map[string]string

// Decoder for decode yml,json,toml defFormat content
type Decoder func(blob []byte, v interface{}) (err error)
type Encoder func(v interface{}) (out []byte, err error)

// Options config options
type Options struct {
	// parse env value. like: "${EnvName}" "${EnvName|default}"
	ParseEnv bool
	// config is readonly
	Readonly bool
	// default write format
	DumpFormat string
	// default input format
	ReadFormat string
}

// Config
type Config struct {
	// config instance name
	name string
	lock sync.RWMutex

	// config options
	opts *Options
	// all config data
	data map[string]interface{}

	// loaded config files
	loadedFiles []string

	// decoders["toml"] = func(blob []byte, v interface{}) (err error){}
	// decoders["yaml"] = func(blob []byte, v interface{}) (err error){}
	decoders map[string]Decoder
	encoders map[string]Encoder

	// cache got config data
	intCaches map[string]int
	strCaches map[string]string
	arrCaches map[string]stringArr
	mapCaches map[string]stringMap
}

// New
func New(name string) *Config {
	return &Config{
		name: name,
		data: make(map[string]interface{}),

		// init options
		opts: &Options{DumpFormat: Json, ReadFormat: Json},

		// default add json driver
		encoders: map[string]Encoder{Json: JsonEncoder},
		decoders: map[string]Decoder{Json: JsonDecoder},
	}
}

/*************************************************************
 * config setting
 *************************************************************/

// SetOptions
func (c *Config) SetOptions(opts *Options) {
	c.opts = opts

	if c.opts.DumpFormat == "" {
		c.opts.DumpFormat = Json
	}

	if c.opts.ReadFormat == "" {
		c.opts.ReadFormat = Json
	}
}

// Readonly
func (c *Config) Readonly(readonly bool) {
	c.opts.Readonly = readonly
}

// Name get config name
func (c *Config) Name() string {
	return c.name
}

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
}

// SetDriver set a decoder and encoder for a format.
func (c *Config) SetDriver(format string, decoder Decoder, encoder Encoder)  {
	c.SetDecoder(format, decoder)
	c.SetEncoder(format, encoder)
}

// HasDecoder
func (c *Config) HasDecoder(format string) bool {
	if format == Yml {
		format = Yaml
	}

	_, ok := c.decoders[format]
	return ok
}

// SetDecoder
func (c *Config) SetDecoder(format string, decoder Decoder) {
	if format == Yml {
		format = Yaml
	}

	c.decoders[format] = decoder
}

// SetDecoders
func (c *Config) SetDecoders(decoders map[string]Decoder) {
	for format, decoder := range decoders {
		c.SetDecoder(format, decoder)
	}
}

// SetEncoder
func (c *Config) SetEncoder(format string, encoder Encoder) {
	if format == Yml {
		format = Yaml
	}

	c.encoders[format] = encoder
}

// SetEncoders
func (c *Config) SetEncoders(encoders map[string]Encoder) {
	for format, encoder := range encoders {
		c.SetEncoder(format, encoder)
	}
}

// HasEncoder
func (c *Config) HasEncoder(format string) bool {
	if format == Yml {
		format = Yaml
	}

	_, ok := c.encoders[format]
	return ok
}

/*************************************************************
 * helper methods
 *************************************************************/

// WriteTo Write out config data representing the current state to a writer.
func (c *Config) WriteTo(out io.Writer) (n int64, err error) {
	return c.DumpTo(out, c.opts.DumpFormat)
}

// DumpTo use the format(json,yaml,toml) dump config data to a writer
func (c *Config) DumpTo(out io.Writer, format string) (n int64, err error) {
	var ok bool
	var encoder Encoder

	if format == Yml {
		format = Yaml
	}

	if encoder, ok = c.encoders[format]; !ok {
		err = errors.New("no exists or no register encoder for the format: " + format)
		return
	}

	// encode data to string
	encoded, err := encoder(&c.data)
	if err != nil {
		return
	}

	// write content to out
	num, err := fmt.Fprintln(out, string(encoded))
	if err != nil {
		return
	}

	return int64(num), nil
}

// ClearAll
func (c *Config) ClearAll() {
	c.ClearData()
	c.ClearCaches()

	c.loadedFiles = []string{}
}

// ClearData
func (c *Config) ClearData() {
	c.data = make(map[string]interface{})
}

// ClearCaches
func (c *Config) ClearCaches() {
	c.intCaches = nil
	c.strCaches = nil
	c.mapCaches = nil
	c.arrCaches = nil
}

// initCaches
func (c *Config) initCaches() {
	c.intCaches = map[string]int{}
	c.strCaches = map[string]string{}
	c.arrCaches = map[string]stringArr{}
	c.mapCaches = map[string]stringMap{}
}
