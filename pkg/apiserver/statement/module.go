// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package statement

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newService),
	fx.Invoke(registerRouter),
)
