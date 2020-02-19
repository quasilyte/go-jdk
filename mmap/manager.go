package mmap

import (
	"fmt"
)

type Manager struct {
	mapped []Descriptor
}

func (m *Manager) Close() error {
	for _, d := range m.mapped {
		if err := Free(d); err != nil {
			return fmt.Errorf("free mapped memory: %v", err)
		}
	}
	return nil
}

func (m *Manager) AllocateExecutable(length int) ([]byte, error) {
	// TODO(quasilyte): re-use mapped regions.
	d, buf, err := Executable(length)
	if err != nil {
		return nil, err
	}
	m.mapped = append(m.mapped, d)
	return buf, nil
}
