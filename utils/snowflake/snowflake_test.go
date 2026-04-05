package snowflake

import (
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/yitter/idgenerator-go/regworkerid"
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	registerId := regworkerid.RegisterOne(regworkerid.RegisterConf{
		Address:         "192.168.2.159:6379",
		Password:        "123456",
		DB:              0,
		LifeTimeSeconds: 60,
		MaxWorkerId:     8,
	})
	var options = idgen.NewIdGeneratorOptions(uint16(registerId))
	//options.WorkerIdBitLength = 10
	idgen.SetIdGenerator(options)

	var newId = idgen.NextId()
	t.Log(newId)
	time.Sleep(time.Hour)

}
