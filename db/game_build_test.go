package db

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestVerifyGameBuild(t *testing.T) {
	_ = config.Load("../config.json")
	InitializeSQL()

	err := VerifyGameBuild(GameBuild{
		QuaverAPIDll:          "4323f6613895445c4a36ad58eb0a2d73",
		QuaverServerClientDll: "c56c6e5127676992522ec0dfa21a3d84",
		QuaverServerCommonDll: "6405048e715c1e94d2063e22e41cce09",
		QuaverSharedDll:       "c8631772af10471bce5ad100d11839fa",
	})

	if err != nil {
		t.Fatal(err)
	}

	CloseSQLConnection()
}
