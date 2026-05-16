// Package buildinfo holds build-time injected version metadata.
package buildinfo

// 本文件承载构建期注入的版本元数据。

var (
	// Version 表示当前可执行文件的版本号，默认值会在构建时被覆盖。
	Version = "dev"
	// Commit 表示当前构建对应的 git commit，默认值会在构建时被覆盖。
	Commit = "unknown"
	// BuildTime 表示当前构建时间，默认值会在构建时被覆盖。
	BuildTime = "unknown"
)

// Info 表示当前二进制的完整构建元数据。
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
}

// Current 返回当前构建元数据快照。
func Current() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}
}
