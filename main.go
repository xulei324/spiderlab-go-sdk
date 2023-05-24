package crawlab

import (
	"github.com/xulei324/spiderlab-go-sdk/entity"
	"github.com/xulei324/spiderlab-go-sdk/utils"
)

func SaveItem(item entity.Item) (err error) {
	if err := utils.SaveItem(item); err != nil {
		return err
	}
	return nil
}
