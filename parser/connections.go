package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shadiestgoat/goSankey/common"
)

func connections(inp []string, nodes map[string]*common.Node) ([]*common.Connection, error) {
	conns := map[string]*common.Connection{}
	connIds := []string{}
	errMgr := &MultiError{}

	for _, l := range inp {
		out := strings.SplitN(l, "->", 2)
		lineErr := fmt.Errorf("line '%s' is not formatted correctly", l)

		if len(out) != 2 {
			errMgr.Errors = append(errMgr.Errors, lineErr)
			continue
		}
		fromID := strings.TrimSpace(out[0])
		out2 := strings.SplitN(out[1], ":", 2)
		if len(out2) != 2 {
			errMgr.Errors = append(errMgr.Errors, lineErr)
			continue
		}
		destID := strings.TrimSpace(out2[0])
		amt, err := strconv.Atoi(strings.TrimSpace(out2[1]))
		if err != nil || amt == 0 {
			errMgr.Errors = append(errMgr.Errors, lineErr)
			continue
		}

		from := nodes[fromID]
		dest := nodes[destID]

		if from == nil || dest == nil {
			errMgr.Errors = append(errMgr.Errors, fmt.Errorf("id not recognized on line '%s'", l))
			continue
		}

		if from.Step >= dest.Step {
			errMgr.Errors = append(errMgr.Errors, fmt.Errorf("backwards connection on line '%s'", l))
			continue
		}

		from.TotalOut += amt
		dest.TotalIn += amt
		
		cID := fromID + "-" + destID
		if conns[cID] == nil {
			conns[cID] = &common.Connection{
				Origin: from,
				Dest:   dest,
				Amount: amt,
			}
			connIds = append(connIds, cID)
		} else {
			conns[cID].Amount += amt
		}

	}

	outConns := []*common.Connection{}

	for _, c := range connIds {
		outConns = append(outConns, conns[c])
	}

	return outConns, errMgr.Optional()
}