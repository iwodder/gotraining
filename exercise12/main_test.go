package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func Test_renameFiles(t *testing.T) {
	type args struct {
		newPattern string
	}
	tests := []struct {
		name     string
		args     args
		setup    func(string) []string
		expected map[string]struct{}
	}{
		{
			name: "Renames one file",
			args: args{
				newPattern: "birthday_%d_of_%d.txt",
			},
			setup: func(dir string) []string {
				ret := make([]string, 0)
				for _, v := range []string{"birthday_001.txt"} {
					f, _ := os.Create(filepath.Join(dir, v))
					ret = append(ret, f.Name())
				}
				return ret
			},
			expected: map[string]struct{}{
				"birthday_1_of_1.txt": {},
			},
		},
		{
			name: "Renames two files",
			args: args{
				newPattern: "birthday_%d_of_%d.txt",
			},
			setup: func(dir string) []string {
				ret := make([]string, 0)
				for _, v := range []string{"birthday_001.txt", "birthday_002.txt"} {
					f, _ := os.Create(filepath.Join(dir, v))
					ret = append(ret, f.Name())
				}
				return ret
			},
			expected: map[string]struct{}{
				"birthday_1_of_2.txt": {},
				"birthday_2_of_2.txt": {},
			},
		},
		{
			name: "Renames many files",
			args: args{
				newPattern: "birthday_%d_of_%d.txt",
			},
			setup: func(dir string) []string {
				ret := make([]string, 0)
				for _, v := range []string{"birthday_001.txt", "birthday_002.txt", "birthday_004.txt", "birthday_007.txt"} {
					f, _ := os.Create(filepath.Join(dir, v))
					ret = append(ret, f.Name())
				}
				return ret
			},
			expected: map[string]struct{}{
				"birthday_1_of_4.txt": {},
				"birthday_2_of_4.txt": {},
				"birthday_3_of_4.txt": {},
				"birthday_4_of_4.txt": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			files := tt.setup(dir)
			renameFiles(files, tt.args.newPattern)

			_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					info, _ := d.Info()
					_, ok := tt.expected[info.Name()]
					if ok {
						delete(tt.expected, info.Name())
					}
				}
				return nil
			})

			if len(tt.expected) != 0 {
				t.Errorf("Unable to find the following expected files: %v", tt.expected)
			}
		})
	}
}
