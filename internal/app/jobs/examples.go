// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"fmt"
	"time"
)

func exampleJob() {
	fmt.Printf("Every seconds, %s\n", time.Now().Format("15:04:05"))
}
