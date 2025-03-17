package jsonconfig

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"slices"
	"sync"

	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log/logattr"
	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

type configPreset struct {
	Key      string   `json:"key"`
	Services []string `json:"services"`
}

type filePreset struct {
	Presets []configPreset `json:"presets"`
}

func (fp filePreset) replace(p configPreset) {
	idx := slices.IndexFunc(fp.Presets, func(item configPreset) bool {
		return item.Key == p.Key
	})

	if idx == -1 {
		fp.Presets = append(fp.Presets, p)
		return
	}

	fp.Presets[idx] = p
}

type jsonconfig struct {
	m sync.Mutex
}

func (repo *jsonconfig) Save(ctx context.Context, p domain.ConfigPreset) (domain.ConfigPreset, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	l := log.Logger().With(logattr.Caller("jsonconfig: save"))

	l.Debug("preparing to read file")
	f, err := openFile()
	if err != nil {
		l.Error("error opening file", logattr.Error(err))

		return domain.ConfigPreset{}, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		l.Error("error reading buffer", logattr.Error(err))

		return domain.ConfigPreset{}, err
	}
	l.Debug("buf created", slog.String("buf", string(buf)))

	var fp filePreset
	err = json.Unmarshal(buf, &fp)
	if err != nil {
		l.Error("error unmarshalling", logattr.Error(err))

		return domain.ConfigPreset{}, err
	}

	fp.replace(configPreset{
		Key:      p.Key,
		Services: p.Services,
	})

	buf, err = json.MarshalIndent(fp, "", "  ")
	if err != nil {
		l.Error("error marshalling JSON", logattr.Error(err))

		return domain.ConfigPreset{}, err
	}

	err = f.Truncate(0)
	if err != nil {
		l.Error("failed to truncate file", logattr.Error(err))

		return domain.ConfigPreset{}, err
	}

	_, err = f.WriteAt(buf, io.SeekStart)
	if err != nil {
		l.Error("failed to write to file", logattr.Error(err))

		return domain.ConfigPreset{}, err
	}

	return p, nil
}

func New() *jsonconfig {
	return &jsonconfig{}
}
