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

package registry

import (
	"github.com/FerretDB/FerretDB/internal/backends/mysql"
	"github.com/FerretDB/FerretDB/internal/handler"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// init registers "mysql" handler.
func init() {
	registry["mysql"] = func(opts *NewHandlerOpts) (*handler.Handler, CloseBackendFunc, error) {
		b, err := mysql.NewBackend(&mysql.NewBackendParams{
			URI: opts.MySQLURL,
			L:   opts.Logger.Named("mysql"),
			P:   opts.StateProvider,
		})
		if err != nil {
			return nil, nil, lazyerrors.Error(err)
		}

		handlerOpts := &handler.NewOpts{
			Backend:     b,
			TCPHost:     opts.TCPHost,
			ReplSetName: opts.ReplSetName,

			SetupDatabase: opts.SetupDatabase,
			SetupUsername: opts.SetupUsername,
			SetupPassword: opts.SetupPassword,
			SetupTimeout:  opts.SetupTimeout,

			L:             opts.Logger.Named("mysql"),
			ConnMetrics:   opts.ConnMetrics,
			StateProvider: opts.StateProvider,

			DisablePushdown:         opts.DisablePushdown,
			EnableNestedPushdown:    opts.EnableNestedPushdown,
			CappedCleanupPercentage: opts.CappedCleanupPercentage,
			CappedCleanupInterval:   opts.CappedCleanupInterval,
			EnableNewAuth:           opts.EnableNewAuth,
			BatchSize:               opts.BatchSize,
			MaxBsonObjectSizeBytes:  opts.MaxBsonObjectSizeBytes,
		}

		h, err := handler.New(handlerOpts)
		if err != nil {
			return nil, b.Close, lazyerrors.Error(err)
		}

		return h, b.Close, nil
	}
}
