// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/library/hggen/internal/consts"
	"hotgo/internal/library/hggen/internal/utility/mlog"
	"hotgo/internal/library/hggen/internal/utility/utils"
)

func generateEntity(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathEntity = gfile.Join(in.Path, in.EntityPath)
	in.genItems.AppendDirPath(dirPathEntity)
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}

		var (
			newTableName                    = in.NewTableNames[i]
			entityFilePath                  = filepath.FromSlash(gfile.Join(dirPathEntity, gstr.CaseSnake(newTableName)+".go"))
			structDefinition, appendImports = generateStructDefinition(ctx, generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
				StructName:           formatFieldName(newTableName, FieldNameCaseCamel),
				FieldMap:             fieldMap,
				IsDo:                 false,
			})
			entityContent = generateEntityContent(
				ctx,
				in,
				newTableName,
				formatFieldName(newTableName, FieldNameCaseCamel),
				structDefinition,
				appendImports,
			)
		)
		in.genItems.AppendGeneratedFilePath(entityFilePath)
		err = gfile.PutContents(entityFilePath, strings.TrimSpace(entityContent))
		if err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", entityFilePath, err)
		} else {
			utils.GoFmt(entityFilePath)
			mlog.Print("generated:", gfile.RealPath(entityFilePath))
		}
	}
}

func generateEntityContent(
	ctx context.Context, in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string, appendImports []string,
) string {
	entityContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoEntityPath, consts.TemplateGenDaoEntityContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(ctx, structDefine, false, appendImports),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
			tplVarPackageName:        filepath.Base(in.EntityPath),
		},
	)
	entityContent = replaceDefaultVar(in, entityContent)
	return entityContent
}
