package gormutils

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetModelCommonQuery(databaseSession *gorm.DB, model interface{}, filter interface{}) *gorm.DB {
	switch filter.(type) {
	case map[string]interface{}:
		query := filter.(map[string]interface{})
		return databaseSession.Where(query).Model(model)
	case *map[string]interface{}:
		query := filter.(*map[string]interface{})
		return databaseSession.Where(*query).Model(model)
	case []clause.Expression:
		query := filter.([]clause.Expression)
		return databaseSession.Clauses(query...).Model(model)
	case *[]clause.Expression:
		query := filter.(*[]clause.Expression)
		return databaseSession.Clauses(*query...).Model(model)
	default:
		return databaseSession.Where(filter).Model(model)
	}
}

func GetQueryWithPagination(commonStatement *gorm.DB, field string, order string, pageLimit, pageIndex int) *gorm.DB {
	return commonStatement.Order(field + " " + order).Offset(pageLimit * pageIndex).Limit(pageLimit)
}

func GetQueryWhereIn(databaseSession *gorm.DB, fieldName string, values []string) clause.IN {
	valuesInterface := make([]interface{}, len(values))
	for i, value := range values {
		valuesInterface[i] = value
	}

	return GetQueryWhereInInterfaceValues(databaseSession, fieldName, valuesInterface)
}

func GetQueryWhereInInterfaceValues(databaseSession *gorm.DB, fieldName string, values []interface{}) clause.IN {
	return clause.IN{Column: fieldName, Values: values}
}

func GetQueryWithSelectSubItemCount(commonStatement *gorm.DB, currentModelTableName string, itemModelTableName string, itemMappingCurrentModelIDField string, currentModelIDField string, countFieldName string) *gorm.DB {
	return commonStatement.Select(
		GetQueryWithSelectSubItemCountSQLString(
			currentModelTableName, itemModelTableName,
			itemMappingCurrentModelIDField, currentModelIDField,
			countFieldName,
			false,
		),
	)
}

func GetQueryWithSelectSubItemCountSQLString(currentModelTableName string, itemModelTableName string, itemMappingCurrentModelIDField string, currentModelIDField string, countFieldName string, skipSelectCurrentModelAllFields bool) string {
	currentModelTableName = strconv.Quote(currentModelTableName)
	itemModelTableName = strconv.Quote(itemModelTableName)

	result := fmt.Sprintf("(select count(*) AS %s FROM %s WHERE %s.%s = %s.%s)",
		countFieldName, itemModelTableName,
		itemModelTableName, itemMappingCurrentModelIDField, currentModelTableName, currentModelIDField,
	)
	if skipSelectCurrentModelAllFields {
		return result
	}

	return fmt.Sprintf("%s.*, %s",
		currentModelTableName,
		result,
	)
}

// SetField ...
func SetField(databaseSession *gorm.DB, model interface{}, fieldName string, value interface{}, UUIDFieldName string, ids ...string) (int64, error) {

	commonStatement := GetModelCommonQuery(databaseSession, model, GetQueryWhereIn(databaseSession, UUIDFieldName, ids))
	var queryResult *gorm.DB

	queryResult = commonStatement.Updates(
		map[string]interface{}{
			fieldName: value,
		},
	)
	if queryResult.Error != nil {
		return 0, queryResult.Error
	}

	return queryResult.RowsAffected, nil
}
