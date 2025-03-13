package database

import (
	"strings"

	"gorm.io/gorm"
)

// ApplyConditions 动态多条件查询函数，根据给定的查询多条件动态构建查询语句。
// db: gorm.DB 对象，用于执行查询的数据库连接。
// conditions: 一个 map，包含查询的条件，key 是字段名，value 是查询的值。
// 其中，字符串类型的值将使用模糊查询（LIKE），其他类型的值使用精确查询（=）。
func ApplyConditions(db *gorm.DB, conditions map[string]interface{}) *gorm.DB {
	// 如果没有提供任何条件，直接返回原始的 db 对象，不进行任何修改
	if len(conditions) == 0 {
		return db
	}

	var conditionBuilt []string // 用于存放每个条件的构建部分（字符串）
	var values []interface{}    // 用于存放每个条件的值（参数）

	// 遍历传入的条件，构建查询字符串
	for key, value := range conditions {
		// 如果值是非空的字符串类型，执行模糊查询（LIKE）
		if strVal, ok := value.(string); ok && strVal != "" {
			conditionBuilt = append(conditionBuilt, key+" LIKE ?") // 将字段名和 LIKE 查询条件拼接
			values = append(values, "%"+strVal+"%")                // 添加参数，并用 % 进行模糊匹配
		} else {
			// 否则，执行精确查询（=）
			conditionBuilt = append(conditionBuilt, key+" = ?") // 将字段名和 = 查询条件拼接
			values = append(values, value)                      // 添加精确查询的值
		}
	}

	// 将所有的查询条件使用 OR 连接，并执行查询
	// 将构建的查询条件和对应的值传递给 db.Where
	return db.Where(strings.Join(conditionBuilt, " OR "), values...)
}
