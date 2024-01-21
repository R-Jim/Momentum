package spawner

type SpawnTypeCode int

type SpawnRunnerData struct {
	HealthValue int
}

func getSpawnTypeData(code SpawnTypeCode) interface{} {
	switch code {
	case 1:
		return SpawnRunnerData{
			HealthValue: 2,
		}
	}
	return nil
}
