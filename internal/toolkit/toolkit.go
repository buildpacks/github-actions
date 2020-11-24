/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package toolkit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

//go:generate mockery --all  --inpackage --case=underscore

type Toolkit interface {
	AddPath(paths ...string) error
	ExportVariable(name string, value string) error

	GetInput(name string) (string, bool)
	SetOutput(name string, value string)
	GetState(name string) (string, bool)
	SetState(name string, value string)
	AddMask(mask string)

	StartGroup(title string)
	EndGroup()

	IsDebug() bool
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Warning(a ...interface{})
	Warningc(context MessageContext)
	Warningf(format string, a ...interface{})
	Error(a ...interface{})
	Errorc(context MessageContext)
	Errorf(format string, a ...interface{})
}

type MessageContext struct {
	File    string
	Line    string
	Column  string
	Message string
}

func (m *MessageContext) String() string {
	var s []string
	if m.File != "" {
		s = append(s, fmt.Sprintf("file=%s", m.File))
	}
	if m.Line != "" {
		s = append(s, fmt.Sprintf("line=%s", m.Line))
	}
	if m.Column != "" {
		s = append(s, fmt.Sprintf("col=%s", m.Column))
	}

	return fmt.Sprintf("%s::%s", strings.Join(s, ","), escape(m.Message))
}

func FailedError(a ...interface{}) error {
	return errors.New(errorString(a...))
}

func FailedErrorc(context MessageContext) error {
	return errors.New(errorStringc(context))
}

func FailedErrorf(format string, a ...interface{}) error {
	return errors.New(errorStringf(format, a...))
}

type DefaultToolkit struct {
	once sync.Once

	Environment map[string]string
	Writer      io.Writer
}

func (d *DefaultToolkit) AddPath(paths ...string) error {
	d.once.Do(d.init)

	path, ok := d.Environment["GITHUB_PATH"]
	if !ok {
		return fmt.Errorf("$GITHUB_PATH must be set")
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open %s\n%w", path, err)
	}
	defer f.Close()

	for _, p := range paths {
		_, _ = fmt.Fprintln(f, p)
	}

	return nil
}

func (d *DefaultToolkit) ExportVariable(name string, value string) error {
	d.once.Do(d.init)

	path, ok := d.Environment["GITHUB_ENV"]
	if !ok {
		return fmt.Errorf("$GITHUB_ENV must be set")
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open %s\n%w", path, err)
	}
	defer f.Close()

	if strings.ContainsRune(value, '\n') {
		if _, err := fmt.Fprintln(f, fmt.Sprintf("%s<<EOF\n%s\nEOF", name, value)); err != nil {
			return fmt.Errorf("unable to write variable")
		}
	} else {
		_, _ = fmt.Fprintln(f, fmt.Sprintf("%s=%s", name, value))
	}

	return nil
}

func (d *DefaultToolkit) GetInput(name string) (string, bool) {
	d.once.Do(d.init)
	s, ok := d.Environment[fmt.Sprintf("INPUT_%s", strings.ToUpper(name))]
	return s, ok
}

func (d *DefaultToolkit) SetOutput(name string, value string) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::set-output name=%s::%s\n", name, escape(value))
}

func (d *DefaultToolkit) GetState(name string) (string, bool) {
	d.once.Do(d.init)
	s, ok := d.Environment[fmt.Sprintf("STATE_%s", strings.ToUpper(name))]
	return s, ok
}

func (d *DefaultToolkit) SetState(name string, value string) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::save-state name=%s::%s\n", name, escape(value))
}

func (d *DefaultToolkit) AddMask(mask string) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::add-mask::%s\n", escape(mask))
}

func (d *DefaultToolkit) StartGroup(title string) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::group::%s\n", title)
}

func (d *DefaultToolkit) EndGroup() {
	d.once.Do(d.init)
	_, _ = fmt.Fprintln(d.Writer, "::endgroup::")
}

func (d *DefaultToolkit) IsDebug() bool {
	d.once.Do(d.init)

	t, err := strconv.ParseBool(d.Environment["RUNNER_DEBUG"])
	if err != nil {
		return false
	}

	return t
}

func (d *DefaultToolkit) Debug(a ...interface{}) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::debug::%s\n", escape(fmt.Sprint(a...)))
}

func (d *DefaultToolkit) Debugf(format string, a ...interface{}) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::debug::%s\n", escape(fmt.Sprintf(format, a...)))
}

func (d *DefaultToolkit) Error(a ...interface{}) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintln(d.Writer, errorString(a...))
}

func (d *DefaultToolkit) Errorc(context MessageContext) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintln(d.Writer, errorStringc(context))
}

func (d *DefaultToolkit) Errorf(format string, a ...interface{}) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintln(d.Writer, errorStringf(format, a...))
}

func (d *DefaultToolkit) Warning(a ...interface{}) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::warning ::%s\n", escape(fmt.Sprint(a...)))
}

func (d *DefaultToolkit) Warningc(context MessageContext) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintln(d.Writer, "::warning", escape(context.String()))
}

func (d *DefaultToolkit) Warningf(format string, a ...interface{}) {
	d.once.Do(d.init)
	_, _ = fmt.Fprintf(d.Writer, "::warning ::%s\n", escape(fmt.Sprintf(format, a...)))
}

func (d *DefaultToolkit) init() {
	if d.Environment == nil {
		d.Environment = make(map[string]string)

		for _, s := range os.Environ() {
			t := strings.SplitN(s, "=", 2)
			d.Environment[t[0]] = t[1]
		}
	}

	if d.Writer == nil {
		d.Writer = os.Stdout
	}
}

func errorString(a ...interface{}) string {
	return fmt.Sprintf("::error ::%s", escape(fmt.Sprint(a...)))
}

func errorStringc(context MessageContext) string {
	return fmt.Sprintf("::error %s", context.String())
}

func errorStringf(format string, a ...interface{}) string {
	return fmt.Sprintf("::error ::%s", escape(fmt.Errorf(format, a...).Error()))
}

func escape(s string) string {
	return strings.ReplaceAll(s, "\n", "%0A")
}
