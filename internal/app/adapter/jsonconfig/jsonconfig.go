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

	fp.Presets = replacePreset(fp, configPreset{Key: p.Key, Services: p.Services})
	l.Debug("modify preset", logattr.Any("presets", fp.Presets))

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

func replacePreset(fp filePreset, p configPreset) []configPreset {
	presets := make([]configPreset, len(fp.Presets))
	copy(presets, fp.Presets)

	idx := slices.IndexFunc(presets, func(item configPreset) bool {
		return item.Key == p.Key
	})

	if idx != -1 {
		presets[idx] = p

		return presets
	}

	presets = append(presets, p)

	return presets
}

func (repo *jsonconfig) Fetch(ctx context.Context) ([]domain.ConfigPreset, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	l := log.Logger().With(logattr.Caller("jsonconfig: fetch"))

	l.Debug("preparing to read file")
	f, err := openFile()
	if err != nil {
		l.Error("error opening file", logattr.Error(err))

		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		l.Error("error reading buffer", logattr.Error(err))

		return nil, err
	}
	l.Debug("buf created", slog.String("buf", string(buf)))

	var fp filePreset
	err = json.Unmarshal(buf, &fp)
	if err != nil {
		l.Error("error unmarshalling", logattr.Error(err))

		return nil, err
	}

	res := make([]domain.ConfigPreset, len(fp.Presets))
	for i, p := range fp.Presets {
		res[i] = domain.NewConfigPreset(p.Key, p.Services)
	}

	return res, nil
}

func New() *jsonconfig {
	return &jsonconfig{}
}
