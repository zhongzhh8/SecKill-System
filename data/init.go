package data

import "SecKill/conf"

func init()  {
	config, err := conf.GetAppConfig()
	if err != nil {
		panic("failed to load data config: " + err.Error())
	}

	initMysql(config)
	initRedisConnection(config)
}

func Close()  {
	err := client.Close()
	if err != nil {
		print("Error on closing redisService client.")
	}
	err = Db.Close()
	if err != nil {
		print("Error on closing dbService client.")
	}
}
