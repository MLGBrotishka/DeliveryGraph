package main

import (
	"fmt"
	"os"
	"os/exec"
	
)
func removeFolder(path string) error {
	cmd := exec.Command("rm", "-rf", path)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func main() {
	err := deleteFilesAndFolder()
	if err != nil {
		fmt.Println("Ошибка при удалении файлов и папки:", err)
		return
	}

	fmt.Println("Файлы и папка успешно удалены.")
}

func deleteFilesAndFolder() error {
	// Удаление файлов
	err := os.Remove("..pypgr")
	if err != nil {
		return err
	}

	err = os.Remove("..pypgr_names")
	if err != nil {
		return err
	}

	// Удаление папки и её содержимого
	err = removeFolder("./chunks")
	if err != nil {
		return err
	}

	// Переход на папку выше
	err = os.Chdir("..")
	if err != nil {
		return err
	}

	// Удаление файла
	err = os.Remove("map")
	if err != nil {
		return err
	}

	return nil
}
