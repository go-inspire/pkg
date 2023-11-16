/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"log/slog"
	"os"
	"testing"
)

func Test_Slog(t *testing.T) {
	slog.Info("debug")
	//slog.Debug("debug", "k1")
	slog.Error("debug", "k1", "v1")
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	l.Info("debug", "k1", "v1")

}
