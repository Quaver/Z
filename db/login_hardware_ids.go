package db

import "time"

type LoginHardwareId struct {
	Id                    int    `db:"id"`
	UserId                int    `db:"user_id"`
	CpuId                 string `db:"cpu_id"`
	DiskId                string `db:"disk_id"`
	CpuDiskId             string `db:"cpu_disk_id"`
	QuaverDll             string `db:"quaver_dll"`
	QuaverAPIDll          string `db:"quaver_api_dll"`
	QuaverServerClientDll string `db:"quaver_server_client_dll"`
	QuaverServerCommonDll string `db:"quaver_server_common_dll"`
	QuaverSharedDll       string `db:"quaver_shared_dll"`
	Occurrences           int    `db:"occurrences"`
	Timestamp             int64  `db:"timestamp"`
}

// InsertLoginHardwareId Logs the hardware ids and client build signatures used by a user during login.
func InsertLoginHardwareId(userId int, cpuId string, diskId string, cpuDiskId string, build GameBuild) error {
	timestamp := time.Now().UnixMilli()

	updateQuery := "UPDATE login_hardware_ids SET " +
		"quaver_dll = ?, quaver_api_dll = ?, quaver_server_client_dll = ?, quaver_server_common_dll = ?, quaver_shared_dll = ?, occurrences = occurrences + 1, timestamp = ? " +
		"WHERE user_id = ? AND cpu_id = ? AND disk_id = ? AND cpu_disk_id = ?"

	result, err := SQL.Exec(updateQuery, build.QuaverDll, build.QuaverAPIDll, build.QuaverServerClientDll,
		build.QuaverServerCommonDll, build.QuaverSharedDll, timestamp, userId, cpuId, diskId, cpuDiskId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		return nil
	}

	query := "INSERT INTO login_hardware_ids " +
		"(user_id, cpu_id, disk_id, cpu_disk_id, quaver_dll, quaver_api_dll, quaver_server_client_dll, quaver_server_common_dll, quaver_shared_dll, occurrences, timestamp) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = SQL.Exec(query, userId, cpuId, diskId, cpuDiskId, build.QuaverDll, build.QuaverAPIDll,
		build.QuaverServerClientDll, build.QuaverServerCommonDll, build.QuaverSharedDll, 1, timestamp)

	if err != nil {
		return err
	}

	return nil
}
