// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"sort"

	"golang.org/x/exp/maps"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgListCommands implements `listCommands` command.
//
// The passed context is canceled when the client connection is closed.
func (h *Handler) MsgListCommands(connCtx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	cmdList := must.NotFail(types.NewDocument())
	names := maps.Keys(h.Commands())
	sort.Strings(names)

	for _, name := range names {
		cmd := h.Commands()[name]
		if cmd.Help == "" {
			continue
		}

		cmdList.Set(name, must.NotFail(types.NewDocument(
			"help", cmd.Help,
		)))
	}

	var reply wire.OpMsg
	must.NoError(reply.SetSections(wire.MakeOpMsgSection(
		must.NotFail(types.NewDocument(
			"commands", cmdList,
			"ok", float64(1),
		)),
	)))

	return &reply, nil
}
