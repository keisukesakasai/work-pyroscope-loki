package profile

import (
	"context"
	"pyroscope-loki-app/internal/log"
	"pyroscope-loki-app/internal/utils"
	"runtime"

	"github.com/grafana/pyroscope-go"
)

const (
	PyroscopeEndpointURLEnv = "PYROSCOPE_ENDPOINT_URL"
)

func Start(serviceAddress string) {
	ctx := context.Background()
	logger := log.GetLoggerFromCtx(ctx)

	appVersion := utils.GetEnv(utils.AppVersionEnv, "unknown")
	serviceName := utils.GetEnv(utils.ServiceNameEnv, "unknown")

	// profiling 設定
	/* Mutex Profile 設定（オプション）
	mutexProfileRate は Mutex Profile の収集される頻度です。
	mutexProfileRate = 1 のとき全ての Mutex Event が収集されます。
	mutexProfileRate > 1 のとき mutexProfileRate 回のうち 1 回 Mutex Profile が収集されます。
	*/
	mutexProfileRate := 1
	runtime.SetMutexProfileFraction(mutexProfileRate) // ・・・(1)
	/* Block Profile 設定（オプション）
	blockProfileRate は Block Profile をサンプルする際の Block 時間（ns）です。
	blockProfileRate = 0 のとき Block Profile が無効になります。
	blockProfileRate > 0 のとき blockProfileRate n秒単位で Block ごとに Block Profile が収集されます。
	*/
	blockProfileRate := 1
	runtime.SetBlockProfileRate(blockProfileRate) // ・・・(2)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: serviceName,
		// Pyroscope のエンドポイントを設定
		ServerAddress: serviceAddress,
		Logger:        pyroscope.StandardLogger,

		// タグを設定することで、タグ指定でのプロファイル表示や、タグ間のプロファイル比較ができ便利です
		Tags: map[string]string{
			utils.AppVersionKey:  appVersion,
			utils.ServiceNameKey: serviceName,
		},

		ProfileTypes: []pyroscope.ProfileType{
			// デフォルトで取得するプロファイル
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// オプショナルで取得するプロファイル
			pyroscope.ProfileGoroutines,
			// ・・・(1) の設定が必要
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			// ・・・(2) の設定が必要
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	logger.Info("Start Profile...")
}
