package seeds

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func SeedAll(app *pocketbase.PocketBase) {
	// FileSystem
	fs, err := app.NewFilesystem()
	if err != nil {
		log.Fatal("cant initialize filesystem: ", err)
	}
	defer fs.Close()

	// Seeds
	TemplatesSeed(app, fs)
}

func JoinSeedsDataPath(elem ...string) string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении пути к исполняемому файлу: %v", err)
	}

	return filepath.Join(append([]string{pwd, "seeds_data"}, elem...)...)
}

func NewFileFromSeedsData(elem ...string) *filesystem.File {
	path := JoinSeedsDataPath(elem...)

	file, err := filesystem.NewFileFromPath(path)
	if err != nil {
		log.Fatal("cant get contacts file: ", path, err)
	}

	return file
}
