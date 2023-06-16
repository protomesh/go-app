package app

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/list"
)

type App interface {
	Log() Logger
	Close()
}

type Dependency interface {
	Attach(app App, dep interface{})
}

type Injector[D any] struct {
	app App
	dep D
}

func (a *Injector[D]) Attach(app App, dep any) {
	a.app = app
	a.dep = dep.(D)
}

func (a *Injector[D]) Dependency() D {
	return a.dep
}

func (a *Injector[Dependency]) Log() Logger {
	return a.app.Log()
}

func Inject[D any](app App, dep D) {
	inject(app, dep, false, nil)
}

func InjectAndPrint[D any](app App, dep D) {
	inject(app, dep, true, nil)
}

func inject[D any](app App, dep D, print bool, lw list.Writer) {

	depVal := reflect.ValueOf(dep)
	appDep := reflect.TypeOf((*Dependency)(nil)).Elem()

	if print {
		lw = list.NewWriter()
		lw.SetStyle(list.StyleBulletSquare)
	}

	if depVal.Kind() == reflect.Ptr && depVal.Elem().Kind() == reflect.Struct {

		depEl := depVal.Elem()
		depType := reflect.TypeOf(dep)

		for i := 0; i < depEl.NumField(); i++ {

			fieldVal := depEl.Field(i)

			if fieldVal.Type().Implements(appDep) && fieldVal.Kind() == reflect.Ptr {

				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}

				fieldEl := fieldVal.Elem()

				appInj := fieldEl.FieldByName("Injector")
				if !fieldEl.IsValid() || appInj.Kind() != reflect.Ptr {
					continue
				}

				if appInj.IsZero() {
					appInj.Set(reflect.New(appInj.Type().Elem()))
				}

				if lw != nil {
					lw.AppendItem(fmt.Sprintf("%s\n[%s ---> %s]\n", depType.Elem().Field(i).Name, depType.String(), fieldVal.Type().String()))
					lw.Indent()
				}

				appInst := appInj.Interface()

				appDep := appInst.(Dependency)

				appDep.Attach(app, dep)

				fieldInst := fieldVal.Interface()

				inject(fieldInst.(App), fieldInst, false, lw)

				if lw != nil {
					lw.UnIndent()
				}

			}

		}

	}

	if print {
		fmt.Println("Dependency hierarchy:")
		fmt.Println(lw.Render())
	}

}

const (
	configFile_cfg = "file"
)

func init() {

	flag.String(ConvertKeyCase(configFile_cfg, KebabCase), "", "[string]\n\tPath to config file (JSON, TOML or YAML)\n")

}

type app struct {
	cfg        ConfigSource
	logBuilder interface {
		Sync() error
	}
	log Logger
}

func NewApp[D Dependency](deps D, opts *AppOptions) App {

	logBuilder := &loggerBuilder[D]{}

	if opts.FlagSet != nil {

		args := opts.Args
		if args != nil {
			args = os.Args[1:]
		}

		opts.ApplyFlags(deps)
		opts.ApplyFlags(logBuilder)

		opts.FlagSet.Parse(args)

	}

	cfg := NewCompositeSource(
		NewFlagSource(opts.KeyCase, opts.FlagSet),
		NewEnvSource(opts.KeyCase),
	)

	err := cfg.Load()
	if err != nil {
		panic(err)
	}

	configFile := cfg.Get(configFile_cfg)

	if configFile.IsSet() {

		cfg = NewCompositeSource(
			NewFlagSource(opts.KeyCase, opts.FlagSet),
			NewEnvSource(opts.KeyCase),
			NewFileSource(configFile.StringVal()),
		)

		err := cfg.Load()
		if err != nil {
			panic(err)
		}

	}

	opts.Source = cfg

	opts.ApplyConfigs(logBuilder)

	log := logBuilder.build()

	appInstance := &app{
		cfg:        cfg,
		logBuilder: logBuilder,
		log:        log,
	}

	opts.Source = cfg

	if opts.Print {
		InjectAndPrint(appInstance, deps)
	} else {
		Inject(appInstance, deps)
	}

	opts.ApplyConfigs(deps)

	return appInstance

}

func (a *app) Log() Logger {
	return a.log
}

func (a *app) Close() {
	a.logBuilder.Sync()
}
