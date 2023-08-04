package parser

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shadiestgoat/goSankey/common"
)

func Parse(file string) (*common.Chart, error) {
	errMgr := MultiError{Errors: []error{}}

	out, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			errMgr.Errors = append(errMgr.Errors, fmt.Errorf("file '%s' could not be found", file))
		} else {
			errMgr.Errors = append(errMgr.Errors, err)
		}

		return nil, errMgr
	}
	
	var sections = map[string][]string{}

	curSec := ""

	for _, l := range strings.Split(string(out), "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "[[") && strings.HasSuffix(l, "]]") {
			curSec = strings.ToLower(strings.TrimSpace(l[2:len(l)-2]))
			continue
		}
		if strings.HasPrefix(l, ";") || strings.HasPrefix(l, "#") {
			continue
		}
		if curSec == "" {
			continue
		}
		if l == "" {
			continue
		}
		if sections[curSec] == nil {
			sections[curSec] = []string{}
		}
		sections[curSec] = append(sections[curSec], l)
	}

	c := &common.Chart{}

	c.Config = config(sections["config"])
	nodes, steps, err := nodes(sections["nodes"], c.Config)
	errMgr.Errors = append(errMgr.Errors, err)

	if len(nodes) == 0 {
		errMgr.Errors = append(errMgr.Errors, fmt.Errorf("no nodes were loaded"))
		return nil, errMgr
	}

	c.Steps = steps
	c.Connections, err = connections(sections["connections"], nodes)
	errMgr.Errors = append(errMgr.Errors, err)

	return c, errMgr.Optional()
}


func parseOpt(l string) (string, string) {
	out := strings.SplitN(l, "=", 2)
	if len(out) != 2 {
		return "", ""
	}
	return strings.ToLower(strings.TrimSpace(out[0])), strings.TrimSpace(out[1])
}

func parseColor(raw string) *common.Color {
	if raw == "" {
		return nil
	}
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "#")
	raw = strings.TrimPrefix(raw, "0x")

	switch len(raw) {
	case 8:
		return parseColorPure(raw[:6])
	case 3:
		return parseColorPure(string([]byte{raw[0], raw[0], raw[1], raw[1], raw[2], raw[2]}))
	case 6:
		return parseColorPure(raw)
	}
	return nil
}

func parseColorPure(v string) *common.Color {
	out1, err := strconv.ParseUint(v[:2], 16, 8)
	if err != nil {
		return nil
	}
	out2, err := strconv.ParseUint(v[2:4], 16, 8)
	if err != nil {
		return nil
	}
	out3, err := strconv.ParseUint(v[4:], 16, 8)
	if err != nil {
		return nil
	}
	return &common.Color{
		R: uint8(out1),
		G: uint8(out2),
		B: uint8(out3),
	}
}

func parseBool(raw string) *bool {
	raw = strings.ToLower(strings.TrimSpace(raw))
	
	v := false

	switch raw {
	case "yes", "true", "1":
		v = true
	case "no", "false", "0":
		v = false
	default:
		return nil
	}

	return &v
}