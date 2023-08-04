package parser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/shadiestgoat/colorutils"
	"github.com/shadiestgoat/goSankey/common"
)

// Returns a map of id -> node & a step array for the nodes
func nodes(inp []string, conf *common.Config) (map[string]*common.Node, [][]*common.Node, error) {
	errors := MultiError{}
	nodeIds := []string{}
	nodeIndex := map[string]int{}
	nodeMap := map[string][]string{}
	curNode := ""

	for i, l := range inp {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			curNode = strings.TrimSpace(l[1 : len(l)-1])
			continue
		}
		if curNode == "" {
			continue
		}
		if l == "" {
			continue
		}

		if nodeMap[curNode] == nil {
			nodeIds = append(nodeIds, curNode)
			nodeIndex[curNode] = i
			nodeMap[curNode] = []string{}
		}
		nodeMap[curNode] = append(nodeMap[curNode], l)
	}

	nodes := map[string]*common.Node{}

	stepValues := []int{}
	oldToNewStep := map[int]int{}

	for _, nID := range nodeIds {
		n := node(nID, nodeMap[nID], conf)
		if n == nil {
			errors.Errors = append(errors.Errors, fmt.Errorf("node '%s' didn't have a STEP included", nID))
			continue
		}

		nodes[nID] = n
		if _, ok := oldToNewStep[n.Step]; ok {
			continue
		}
		oldToNewStep[n.Step] = 0
		stepValues = append(stepValues, n.Step)
	}

	sort.IntSlice(stepValues).Sort()

	for i, v := range stepValues {
		oldToNewStep[v] = i
	}

	steps := make([][]*common.Node, len(stepValues))

	for _, nID := range nodeIds {
		n := nodes[nID]
		if n == nil {
			continue
		}
		
		n.Step = oldToNewStep[n.Step]

		if len(steps[n.Step]) == 0 {
			steps[n.Step] = []*common.Node{n}
		} else {
			steps[n.Step] = append(steps[n.Step], n)
		}
	}
	
	return nodes, steps, errors.Optional()
}

func node(id string, lines []string, conf *common.Config) *common.Node {
	n := &common.Node{
		ID: id,
	}

	stepSet := false

	for _, l := range lines {
		k, v := parseOpt(l)
		if k == "" {
			continue
		}

		switch k {
		case "title":
			n.Title = v
		case "step":
			out, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			stepSet = true
			n.Step = int(out)
		case "color":
			out := parseColor(v)
			if out == nil {
				continue
			}
			n.Color = out
		}
	}

	if !stepSet {
		return nil
	}

	if n.Title == "" {
		n.Title = n.ID
	}

	if n.Color == nil {
		r, g, b := colorutils.NewContrastColor(7.5, conf.BackgroundIsLight, colorutils.RelativeLuminosity(conf.Background.R, conf.Background.G, conf.Background.B))
		n.Color = &common.Color{
			R: r,
			G: g,
			B: b,
		}
	}

	return n
}
