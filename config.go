package app

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	UnkownConfigFormatError = errors.New("UnkownConfigFormat")
)

type Config interface {
	IsSet() bool

	StringVal() string
	Int64Val() int64
	Float64Val() float64
	StringSliceVal() []string
	DurationVal() time.Duration
	TimeVal() time.Time
	BoolVal() bool
	InterfaceVal() interface{}
}

type ConfigSource interface {
	Load() error
	Get(k string) Config
	Has(k string) bool
}

type AppOptions struct {
	tw table.Writer

	Print   bool
	Source  ConfigSource
	FlagSet *flag.FlagSet
	Prefix  string
	Args    []string
}

func (ao *AppOptions) getFieldNameAndType(typeVal reflect.StructField) (string, string) {

	vals := strings.Split(typeVal.Tag.Get("config"), ",")

	if len(vals) == 0 {
		return "", ""
	}

	key := vals[0]

	if len(key) == 0 {
		return "", ""
	}

	if len(ao.Prefix) > 0 {
		key = strings.Join([]string{ao.Prefix, key}, ".")
	}

	if len(vals) > 1 {
		return key, vals[1]
	}

	return key, "string"

}

func (ao *AppOptions) getFieldUsage(typeVal reflect.StructField) string {
	return typeVal.Tag.Get("usage")
}

func (ao *AppOptions) getFieldDefaultValue(typeVal reflect.StructField) string {
	return typeVal.Tag.Get("default")
}

func (ao *AppOptions) ApplyFlags(keySet any) {

	t := reflect.TypeOf(keySet)

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("Can only fill flags for struct pointers")
	}

	e := t.Elem()

	configType := reflect.TypeOf((*Config)(nil)).Elem()

	for i := 0; i < e.NumField(); i++ {

		typeVal := e.Field(i)

		key, flagType := ao.getFieldNameAndType(typeVal)

		if len(key) == 0 {
			continue
		}

		if typeVal.Type.Implements(configType) {

			key := ConvertKeyCase(key, KebabCase)

			usage := strings.Join([]string{
				"[%s]\n\t",
				ao.getFieldUsage(typeVal),
				"\n",
			}, "")

			defVal := ao.getFieldDefaultValue(typeVal)

			switch strings.ToLower(flagType) {

			case "bool", "boolean":

				val := false

				switch strings.ToLower(defVal) {

				case "t", "true", "y", "yes":
					val = true

				case "f", "false", "n", "no", "not", "":
					break

				default:
					panic(fmt.Sprintf("Invalid boolean default value '%s' (error: %s)", key, defVal))

				}

				ao.FlagSet.Bool(key, val, fmt.Sprintf(usage, "boolean"))

			case "int64", "int", "integer":

				val := int64(0)

				if len(defVal) > 0 {

					intVal, err := strconv.ParseInt(defVal, 10, 64)

					if err != nil {
						panic(fmt.Errorf("Invalid int64 default value '%s' (error: %s)", key, err))
					}

					val = intVal
				}

				ao.FlagSet.Int64(key, val, fmt.Sprintf(usage, "int64"))

			case "float64", "float", "double":

				val := float64(0)

				if len(defVal) > 0 {

					floatVal, err := strconv.ParseFloat(defVal, 64)

					if err != nil {
						panic(fmt.Errorf("Invalid float64 default value '%s' (error: %s)", key, err))
					}

					val = floatVal
				}

				ao.FlagSet.Float64(key, val, fmt.Sprintf(usage, "float64"))

			case "duration":

				val := time.Duration(0)

				if len(defVal) > 0 {

					durationVal, err := time.ParseDuration(defVal)

					if err != nil {
						panic(fmt.Errorf("Invalid duration default value '%s' (error: %s)", key, err))
					}

					val = durationVal

				}

				ao.FlagSet.Duration(key, val, fmt.Sprintf(usage, "duration"))

			case "str", "string":
				ao.FlagSet.String(key, defVal, fmt.Sprintf(usage, "string"))

			case "time", "datetime", "date":
				ao.FlagSet.String(key, defVal, fmt.Sprintf(usage, "datetime"))

			default:
				panic(fmt.Errorf("Unable to define the type for flag '%s' (type: %s)", key, flagType))
			}

			continue
		}

		if typeVal.Type.Kind() == reflect.Ptr && typeVal.Type.Elem().Kind() == reflect.Struct {
			typeCb := &AppOptions{
				Source:  ao.Source,
				FlagSet: ao.FlagSet,
				Prefix:  key,
			}
			typeCb.ApplyFlags(reflect.New(typeVal.Type.Elem()).Interface())
		}

	}

}

func (ao *AppOptions) ApplyConfigs(keySet any) {

	v := reflect.ValueOf(keySet)

	if ao.Print {
		ao.tw = table.NewWriter()
		ao.tw.AppendHeader(table.Row{"Configuration key", "Type", "Value"})
	}

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("Can only fill configs of struct pointers")
	}

	e := v.Elem()
	elType := e.Type()

	configType := reflect.TypeOf((*Config)(nil)).Elem()

	for i := 0; i < e.NumField(); i++ {

		fieldVal := e.Field(i)

		if !fieldVal.IsValid() {
			continue
		}

		fieldType := fieldVal.Type()

		typeVal := elType.Field(i)

		if !typeVal.IsExported() {
			continue
		}

		key, _ := ao.getFieldNameAndType(typeVal)

		if fieldType.Implements(configType) {

			cfg := ao.Source.Get(key)

			if ao.tw != nil {
				ao.tw.AppendRow(table.Row{key, fmt.Sprintf("%T", cfg), fmt.Sprintf("%+v", cfg)})
			}

			fieldVal.Set(reflect.ValueOf(cfg))
			continue
		}

		if fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().Kind() == reflect.Struct {
			fieldCb := &AppOptions{
				Source:  ao.Source,
				FlagSet: ao.FlagSet,
				Prefix:  key,
				Print:   false,
				tw:      ao.tw,
			}

			fieldCb.ApplyConfigs(fieldVal.Interface())
		}

	}

	if ao.Print {
		fmt.Println("Configuration table:")
		fmt.Println(ao.tw.Render())
	}

}
