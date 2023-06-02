package utils

import (
	"github.com/fatih/color"
	"time"
)

type PromptLevel int

const (
	Info PromptLevel = iota
	Warning
	Fatal
	Normal
)

type (
	Prompt struct {
		msg       string
		level     PromptLevel
		timestamp time.Time
	}

	Prompter struct {
		prompts []Prompt
		colors  map[color.Attribute]*color.Color
	}
)

func NewPrompter() *Prompter {
	return &Prompter{
		prompts: make([]Prompt, 16),
		colors:  make(map[color.Attribute]*color.Color),
	}
}

func (p *Prompter) printColor(attribute color.Attribute) *color.Color {
	v, ok := p.colors[attribute]
	if ok {
		return v
	} else {
		nc := color.New(attribute)
		p.colors[attribute] = nc
		return nc
	}
}

func (p *Prompter) Clear() {
	p.prompts = make([]Prompt, 16) // Clear the slice
}

func (p *Prompter) LogFatal(msg string) {
	p.prompts = append(p.prompts, Prompt{msg: msg, level: Fatal, timestamp: time.Now()})
	p.Fatal(msg)
}

func (p *Prompter) Fatal(msg ...string) {
	hr, r := p.printColor(color.FgHiRed), p.printColor(color.FgRed)
	_, _ = r.Print("[Fatal]", " ")
	for _, v := range msg {
		_, _ = hr.Print(v)
	}
	_, _ = hr.Print("\n")
}

func (p *Prompter) LogInfo(msg string) {
	p.prompts = append(p.prompts, Prompt{msg: msg, level: Info, timestamp: time.Now()})
	p.Fatal(msg)
}

func (p *Prompter) Info(msg ...string) {
	hr, r := p.printColor(color.FgHiCyan), p.printColor(color.FgCyan)
	_, _ = r.Print("[Info]", " ")
	for _, v := range msg {
		_, _ = hr.Print(v)
	}
	_, _ = hr.Print("\n")
}

func (p *Prompter) LogWarning(msg string) {
	p.prompts = append(p.prompts, Prompt{msg: msg, level: Warning, timestamp: time.Now()})
	p.Fatal(msg)
}

func (p *Prompter) Warning(msg ...string) {
	hr, r := p.printColor(color.FgHiYellow), p.printColor(color.FgYellow)
	_, _ = r.Print("[Warning]", " ")
	for _, v := range msg {
		_, _ = hr.Print(v)
	}
	_, _ = hr.Print("\n")
}

func (p *Prompter) LogNormal(msg string) string {
	p.prompts = append(p.prompts, Prompt{msg: msg, level: Normal, timestamp: time.Now()})
	return msg
}

func (p *Prompter) Error(err error) {
	if err == nil || err.Error() == "" {
		return
	}
	p.Fatal(err.Error())
}

func (p *Prompter) HandleError(handler DelayedErrorHandler) {
	p.Error(handler.Error())
}

type DelayedErrorHandler interface {
	Catch(err error)
	Error() error
}
