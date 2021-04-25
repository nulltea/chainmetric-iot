package shared

import (
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	compositeKeySeparatorRune = utf8.RuneSelf
	compositeKeySeparator = string(compositeKeySeparatorRune)
)

var LevelDB  *leveldb.DB

func initLevelDB() {
	var (
		path = viper.GetString("device.local_cache_path")
		err error
	)

	if len(path) == 0 {
		Logger.Fatal("failed to initialise LevelDB: local path not provided")
	}

	if LevelDB, err = leveldb.OpenFile(path, nil); err != nil {
		Logger.Fatal(errors.Wrap(err, "failed to initialise LevelDB"))
	}
}

// FormCompositeKey provides composite key for specified `objectType` by combining the given `attributes`.
func FormCompositeKey(objectType string, attributes ...string) string {
	ck := objectType + compositeKeySeparator
	for _, attr := range attributes {
		ck += attr + compositeKeySeparator
	}
	return ck
}

// SplitCompositeKey retrieves object type and attributes from `compositeKey`.
func SplitCompositeKey(compositeKey string) (string, []string) {
	componentIndex := 1
	components := []string{}
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == compositeKeySeparatorRune {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}
	return components[0], components[1:]
}
