package db

type GameBuild struct {
	Id                    int    `db:"id"`
	Version               string `db:"version"`
	QuaverDll             string `db:"quaver_dll"` // Not used during check intentionally due to differences across client platforms.
	QuaverAPIDll          string `db:"quaver_api_dll"`
	QuaverServerClientDll string `db:"quaver_server_client_dll"`
	QuaverServerCommonDll string `db:"quaver_server_common_dll"`
	QuaverSharedDll       string `db:"quaver_shared_dll"`
	Allowed               bool   `db:"allowed"`
	Timestamp             int64  `db:"timestamp"`
}

// VerifyGameBuild Checks to see if the provided game build is valid. Returns an error if invalid
func VerifyGameBuild(build GameBuild) error {
	var result GameBuild

	const query string = "SELECT * FROM game_builds WHERE quaver_api_dll = ? AND quaver_server_client_dll = ? AND quaver_server_common_dll = ? AND quaver_shared_dll = ? AND allowed = 1 LIMIT 1"
	err := SQL.Get(&result, query, build.QuaverAPIDll, build.QuaverServerClientDll, build.QuaverServerCommonDll, build.QuaverSharedDll)

	if err != nil {
		return err
	}

	return nil
}
