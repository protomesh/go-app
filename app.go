package app

import (
	"fmt"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/list"
)

var (
	_ Dependency = &Injector[any]{}
)

type App interface {
	Log() Logger
}

type Dependency interface {
	Attach(app any, dep any)
}

type Injector[D any] struct {
	app App
	dep D
}

func (a *Injector[D]) Attach(app any, dep any) {
	a.app = app.(App)
	a.dep = dep.(D)
}

func (a *Injector[D]) Dependency() D {
	return a.dep
}

func (a *Injector[Dependency]) Log() Logger {
	return a.app.Log()
}

func Inject[D any](app App, dep D) {
	InjectAny(app, dep, false, nil)
}

func InjectAndPrint[D any](app App, dep D) {
	InjectAny(app, dep, true, nil)
}

func InjectAny[D any](app any, dep D, print bool, lw list.Writer) {

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

				InjectAny(fieldInst, fieldInst, false, lw)

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

type AppWithClose interface {
	App
	Close()
}

type app struct {
	log interface {
		Logger
		Sync() error
	}

	ConfigFile Config `config:"config.file,str" usage:"Path to config file (JSON, TOML or YAML)"`
}

func NewApp[D Dependency](deps D, opts *AppOptions) AppWithClose {

	logBuilder := &loggerBuilder[D]{}
	appInstance := &app{}

	if opts.FlagSet != nil {

		args := opts.Args
		if args != nil {
			args = os.Args[1:]
		}

		opts.ApplyFlags(appInstance)
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

	opts.ApplyConfigs(appInstance)

	if appInstance.ConfigFile.IsSet() {

		cfg = NewCompositeSource(
			NewFlagSource(opts.KeyCase, opts.FlagSet),
			NewEnvSource(opts.KeyCase),
			NewFileSource(appInstance.ConfigFile.StringVal()),
		)

		err := cfg.Load()
		if err != nil {
			panic(err)
		}

	}

	opts.Source = cfg

	opts.ApplyConfigs(logBuilder)

	appInstance.log = logBuilder.build()

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
	a.log.Sync()
}
