package system

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Executable 获取可执行文件绝对路径
func Executable() string {
	bin, _ := os.Executable()
	return filepath.Dir(bin)
}

// Getwd 获取工作目录
func Getwd() string {
	dir, _ := os.Getwd()
	return dir
}

// GetDirSize 获取目录大小，单位 Bit
func GetDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

type FileInfo struct {
	os.FileInfo
	Path string
}

// GlobFiles 按照升序排列所有文件
func GlobFiles(path string) ([]FileInfo, error) {
	files := make([]FileInfo, 0, 8)
	err := filepath.Walk(path, func(ppath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, FileInfo{
			FileInfo: info,
			Path:     ppath,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 按照文件的修改时间升序排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	return files, nil
}

func CleanOldFiles(files []FileInfo, count int) (int, error) {
	logs := make([]string, 0, count)

	// 删除文件
	num := count
	for _, file := range files {
		if err := os.Remove(file.Path); err != nil {
			slog.Error("文件删除失败", "err", err)
			continue
		}
		num--
		logs = append(logs, file.Path)
		if num <= 0 {
			break
		}
	}

	if len(logs) > 0 {
		slog.Info("删除旧文件", "logs", logs)
	}

	return count - num, nil
}

// RemoveEmptyDirs 删除空目录，性能优化版，增加时间范围过滤
func RemoveEmptyDirs(ctx context.Context, rootDir string, start, end time.Time) error {
	// 使用 map 统计每个目录的文件数量
	dirFileCount := make(map[string]int)
	cleanRootDir := filepath.Clean(rootDir)
	dirFileCount[filepath.Dir(cleanRootDir)]++
	err := filepath.Walk(cleanRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != cleanRootDir {
			// 只处理在时间范围内的目录
			mod := info.ModTime()
			if mod.Before(start) || mod.After(end) {
				return filepath.SkipDir
			}
			dirFileCount[path] = 0
		}
		parentDir := filepath.Dir(path)
		if parentDir != path { // 确保不是根目录
			dirFileCount[parentDir]++
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 遍历 map，将 == 0 的数量的目录删除
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			changed := false
			for dir, count := range dirFileCount {
				if count != 0 {
					continue
				}
				if err := os.RemoveAll(dir); err != nil {
					return err
				}
				delete(dirFileCount, dir)
				parentDir := filepath.Dir(dir)
				if parentDir != dir {
					dirFileCount[parentDir]--
					changed = true
				}

			}
			if !changed {
				return nil
			}
		}
	}
}

// Abs 获取绝对目录
// 与 filepath.Abs 的区别是，这个以可执行文件目录为工作目录
func Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	bin, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(bin), path), nil
}
