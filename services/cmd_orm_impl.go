package services

import (
	"context"
	"fmt"

	"github.com/geekymedic/neon-cli/templates"
	"github.com/iancoleman/strcase"
)

func (s *GenerateServer) ORM(ctx context.CancelFunc, reply *GenServerORMArg) (*EmptyReply, error) {
	table := reply.Table
	tableName := table.TableNames()[0]
	var snakeProperty []templates.ORMProperty
	var camelProperty []templates.ORMProperty
	for _, item := range table.Column(tableName) {
		snakeProperty = append(snakeProperty, templates.ORMProperty{
			Name: strcase.ToSnake(item.ColumnName),
			Type: item.Type,
		})
		camelProperty = append(camelProperty, templates.ORMProperty{
			Name: strcase.ToCamel(item.ColumnName),
			Type: item.Type,
		})
	}
	var data = templates.ORMTplArg{
		ShortTable:    strcase.ToLowerCamel(tableName),
		FullTable:     strcase.ToCamel(tableName),
		TableName:     tableName,
		SnakeProperty: snakeProperty,
		CamelProperty: camelProperty,
	}
	txt, err := templates.ParseTemplate(templates.ORMTemplate, data)
	if err != nil {
		return &EmptyReply{}, err
	}

	fmt.Println(txt)
	return &EmptyReply{}, nil
}
