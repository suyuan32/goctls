package importschema

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/strutil"

	"github.com/suyuan32/goctls/rpc/execx"
)

// FormatFile formats the file to adjust simple admin
func FormatFile(ctx *GenContext) error {
	files, err := fileutil.ListFileNames(ctx.OutputDir)
	if err != nil {
		return err
	}

	pluralSuffix := ""

	if ctx.PluralTable {
		pluralSuffix = "s"
	}

	for _, v := range files {
		filePath := filepath.Join(ctx.OutputDir, v)

		fileStr, err := fileutil.ReadFileToString(filePath)
		if err != nil {
			return err
		}

		if strings.Contains(fileStr, "entsql.WithComments(true)") {
			continue
		}

		if !strings.Contains(fileStr, "),\n") {
			fileStr = strings.ReplaceAll(fileStr, "),", "),\n\t\t")
			fileStr = strings.ReplaceAll(fileStr, "ent.Field{field", "ent.Field{\n\t\tfield")
			fileStr = strings.ReplaceAll(fileStr, "ent.Index{index", "ent.Index{\n\t\tindex")
		}

		if !strings.Contains(fileStr, "WithComments") && strings.Contains(fileStr, "Comment") {
			if !strings.Contains(fileStr, "dialect/entsql") {
				importIndex := strings.Index(fileStr, "import (")
				fileStr = fileStr[:importIndex+8] + "\n\t\"entgo.io/ent/dialect/entsql\"" + fileStr[importIndex+8:]
			}
		}

		if strings.Contains(fileStr, "uuid \"github.com/satori/go.uuid\"") {
			fileStr = strings.ReplaceAll(fileStr, "uuid \"github.com/satori/go.uuid\"",
				"\"github.com/gofrs/uuid/v5\"")
		}

		fileStr = removeOldAnnotation(fileStr, getSchemaName(fileStr))

		fileStr += fmt.Sprintf(`

func (%s) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "%s%s"},
	}
}`, getSchemaName(fileStr), strings.ToLower(strutil.SnakeCase(getSchemaName(fileStr))), pluralSuffix)

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

func getSchemaName(data string) string {
	typeIndex := strings.Index(data, "type ")
	structIndex := strings.Index(data, " struct")

	return strings.Trim(data[typeIndex+4:structIndex], " ")
}

func removeOldAnnotation(data, schemaName string) string {
	lastFuncIndex := strings.LastIndex(data, fmt.Sprintf("func (%s) Annotations()", schemaName))
	if lastFuncIndex == -1 {
		return data
	}
	funcEndIndex := 0
	count := 0
	for i := lastFuncIndex; i < len(data); i++ {
		if data[i] == '{' {
			count++
		} else if data[i] == '}' {
			count--

			if count == 0 {
				funcEndIndex = i + 1
				break
			}
		}
	}
	return data[:lastFuncIndex] + data[funcEndIndex:]
}
