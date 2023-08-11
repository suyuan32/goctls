package importschema

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/suyuan32/goctls/rpc/execx"
	"path/filepath"
	"strings"
)

// FormatFile formats the file to adjust simple admin
func FormatFile(ctx *GenContext) error {
	files, err := fileutil.ListFileNames(ctx.OutputDir)
	if err != nil {
		return err
	}

	for _, v := range files {
		filePath := filepath.Join(ctx.OutputDir, v)

		fileStr, err := fileutil.ReadFileToString(filePath)
		if err != nil {
			return err
		}

		if !strings.Contains(fileStr, "),\n") {
			fileStr = strings.ReplaceAll(fileStr, "),", "),\n\t\t")
			fileStr = strings.ReplaceAll(fileStr, "ent.Field{field", "ent.Field{\n\t\tfield")
			fileStr = strings.ReplaceAll(fileStr, "ent.Index{index", "ent.Index{\n\t\tindex")
		}

		if !strings.Contains(fileStr, "WithComments") && strings.Contains(fileStr, "Comment") {
			fileStr = strings.ReplaceAll(fileStr, ".Comment",
				".\n\t\tAnnotations(entsql.WithComments(true)).\n\t\tComment")

			if !strings.Contains(fileStr, "dialect/entsql") {
				importIndex := strings.Index(fileStr, "import (")
				fileStr = fileStr[:importIndex+8] + "\n\t\"entgo.io/ent/dialect/entsql\"" + fileStr[importIndex+8:]
			}
		}

		if strings.Contains(fileStr, "uuid \"github.com/satori/go.uuid\"") {
			fileStr = strings.ReplaceAll(fileStr, "uuid \"github.com/satori/go.uuid\"",
				"\"github.com/gofrs/uuid/v5\"")
		}

		if err := fileutil.WriteStringToFile(filePath, fileStr, false); err != nil {
			return err
		}

		_, err = execx.Run(fmt.Sprintf("gofmt -s -w %s", filePath), ctx.OutputDir)
		if err != nil {
			return errors.Join(err, errors.New("failed to format the files, please install gofmt"))
		}

	}

	return nil
}
