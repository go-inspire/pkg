/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package server

import "context"

// 服务接口，又 #app 统一管理
type Server interface {
	// Start  启动服务
	Start(ctx context.Context) error

	// Stop 停止服务
	Stop(ctx context.Context) error
}
