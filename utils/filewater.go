package utils

import (
	"fmt"
	"log"

	"gopkg.in/fsnotify.v1"
)

// FileWatcher 结构体用于封装文件监听功能
type FileWatcher struct {
	watcher  *fsnotify.Watcher
	filePath string
	onChange func(string)
}

// New 创建一个新的 FileWatcher 实例
func New(filePath string, onChange func(string)) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating fsnotify watcher: %v", err)
	}

	// 添加监控的文件
	err = watcher.Add(filePath)
	if err != nil {
		return nil, fmt.Errorf("error adding file to watcher: %v", err)
	}

	return &FileWatcher{
		watcher:  watcher,
		filePath: filePath,
		onChange: onChange,
	}, nil
}

// Start 开始监听文件变化
func (fw *FileWatcher) Start() {
	fmt.Printf("Watching %s for changes...\n", fw.filePath)

	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
					fw.onChange(event.Name)
				}
			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("error watching file: %v", err)
			}
		}
	}()
}

// Close 关闭文件监听
func (fw *FileWatcher) Close() {
	fw.watcher.Close()
}
