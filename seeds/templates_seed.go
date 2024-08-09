package seeds

import (
	"log"

	"lamas/models"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TemplatesSeed(app *pocketbase.PocketBase, fs *filesystem.System) {
	dao := app.Dao()

	// Collection
	collection, _ := dao.FindCollectionByNameOrId(models.TemplateCollectionId)
	baseFilesPath := collection.BaseFilesPath()

	// Update
	TemplateUpdate(dao, fs, baseFilesPath, "start:guest", "lifebuoy.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "start:user", "fish.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "start:admin", "admin.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "rules", "gates.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "settings", "sea.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "settings:username", "alert.png")
	TemplateUpdate(dao, fs, baseFilesPath, "vpn:default", "otters.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "vpn:help", "routers.jpg")
	TemplateUpdate(dao, fs, baseFilesPath, "watch", "cinema.jpg")
}

func TemplateUpdate(dao *daos.Dao, fs *filesystem.System, baseFilesPath string, slug string, fileName string) {
	item, err := models.GetTemplate(dao, slug)
	if err != nil {
		log.Println("cant get Template item:", slug, err)
		return
	}

	// Already seeded
	if item.Banner != "" {
		return
	}

	// File
	file := NewFileFromSeedsData("templates", fileName)
	item.Banner = file.Name
	fileKey := baseFilesPath + "/" + item.Id + "/" + file.Name

	// Upload
	if err := fs.UploadFile(file, fileKey); err != nil {
		log.Fatal("cant upload Template file: ", err)
	}

	// Persist the changes
	if err := dao.Save(item); err != nil {
		fs.Delete(fileKey)
		log.Fatal("cant save Template item: ", err)
	}
}
