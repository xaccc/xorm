package xorm

import (
	"database/sql"
	"strings"
)

type SqlMapsExecutor struct {
	session *Session
	sqlkeys interface{}
	parmas  interface{}
	err     error
}

func (sqlMapsExecutor *SqlMapsExecutor) Execute() ([][]map[string]interface{}, map[string][]map[string]interface{}, error) {
	defer sqlMapsExecutor.session.resetStatement()
	if sqlMapsExecutor.session.IsAutoClose {
		defer sqlMapsExecutor.session.Close()
	}

	if sqlMapsExecutor.err != nil {
		return nil, nil, sqlMapsExecutor.err
	}

	var model_1_results ResultMap
	var model_2_results sql.Result
	var err error

	sqlModel := 1

	if sqlMapsExecutor.session.IsSqlFuc == true {
		err := sqlMapsExecutor.session.Begin()
		if err != nil {
			return nil, nil, err
		}
	}

	switch sqlMapsExecutor.sqlkeys.(type) {
	case string:
		sqlkey := sqlMapsExecutor.sqlkeys.(string)
		sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkey)
		sqlStr = strings.TrimSpace(sqlStr)
		sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

		if sqlMapsExecutor.parmas == nil {
			switch sqlCmd {
			case "select", "desc":
				model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey).Query()
				sqlModel = 1
			case "insert", "delete", "update", "create", "drop":
				model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey).Execute()
				sqlModel = 2
			default:
				sqlModel = 3
			}
		} else {
			switch sqlMapsExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, ok := sqlMapsExecutor.parmas.([]map[string]interface{})
				if !ok {
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err1 := sqlMapsExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, ErrParamsType
				}
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap[0]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap[0]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}

			case map[string]interface{}:
				parmaMap, ok := sqlMapsExecutor.parmas.(map[string]interface{})
				if !ok {
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err1 := sqlMapsExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, ErrParamsType
				}

				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}

			default:
				if sqlMapsExecutor.session.IsSqlFuc == true {
					err1 := sqlMapsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}
		}

		resultSlice := make([][]map[string]interface{}, 1)

		if sqlModel == 1 {
			if model_1_results.Error != nil {
				if sqlMapsExecutor.session.IsSqlFuc == true {
					err1 := sqlMapsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, model_1_results.Error
			}

			resultSlice[0] = make([]map[string]interface{}, len(model_1_results.Results))
			resultSlice[0] = model_1_results.Results
			return resultSlice, nil, nil
		} else if sqlModel == 2 {
			if err != nil {
				if sqlMapsExecutor.session.IsSqlFuc == true {
					err1 := sqlMapsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, err
			}

			resultMap := make([]map[string]interface{}, 1)
			resultMap[0] = make(map[string]interface{})

			//todo all database support LastInsertId
			LastInsertId, _ := model_2_results.LastInsertId()

			resultMap[0]["LastInsertId"] = LastInsertId
			RowsAffected, err := model_2_results.RowsAffected()
			if err != nil {
				return nil, nil, err
			}
			resultMap[0]["RowsAffected"] = RowsAffected
			resultSlice[0] = resultMap
			return resultSlice, nil, nil
		} else {
			resultSlice[0] = nil
		}
	case []string:

		sqlkeysSlice := sqlMapsExecutor.sqlkeys.([]string)
		n := len(sqlkeysSlice)
		resultSlice := make([][]map[string]interface{}, n)
		parmaSlice := make([]map[string]interface{}, n)

		if sqlMapsExecutor.parmas == nil {
			for i, _ := range sqlkeysSlice {
				sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkeysSlice[i])
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkeysSlice[i]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkeysSlice[i]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}

				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultSlice[i] = make([]map[string]interface{}, len(model_1_results.Results))
					resultSlice[i] = model_1_results.Results

				} else if sqlModel == 2 {
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultSlice[i] = make([]map[string]interface{}, 1)
					resultSlice[i] = resultMap

				} else {
					resultSlice[i] = nil
				}

			}

		} else {
			switch sqlMapsExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaSlice = sqlMapsExecutor.parmas.([]map[string]interface{})

			default:
				if sqlMapsExecutor.session.IsSqlFuc == true {
					err1 := sqlMapsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}

			for i, _ := range sqlkeysSlice {
				sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkeysSlice[i])
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

				if parmaSlice[i] == nil {
					switch sqlCmd {
					case "select", "desc":
						model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkeysSlice[i]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkeysSlice[i]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				} else {
					switch sqlCmd {
					case "select", "desc":
						model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkeysSlice[i], &parmaSlice[i]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkeysSlice[i], &parmaSlice[i]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				}

				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultSlice[i] = make([]map[string]interface{}, len(model_1_results.Results))
					resultSlice[i] = model_1_results.Results

				} else if sqlModel == 2 {
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultSlice[i] = make([]map[string]interface{}, 1)
					resultSlice[i] = resultMap

				} else {
					resultSlice[i] = nil
				}
			}

		}

		if sqlMapsExecutor.session.IsSqlFuc == true {
			err1 := sqlMapsExecutor.session.Commit()
			if err1 != nil {
				return nil, nil, err1
			}
		}
		return resultSlice, nil, nil

	case map[string]string:

		sqlkeysMap := sqlMapsExecutor.sqlkeys.(map[string]string)
		n := len(sqlkeysMap)
		resultsMap := make(map[string][]map[string]interface{}, n)
		parmasMap := make(map[string]map[string]interface{}, n)

		if sqlMapsExecutor.parmas == nil {
			for k, _ := range sqlkeysMap {
				sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkeysMap[k])
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkeysMap[k]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkeysMap[k]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}

				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultsMap[k] = make([]map[string]interface{}, len(model_1_results.Results))
					resultsMap[k] = model_1_results.Results

				} else if sqlModel == 2 {
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultsMap[k] = make([]map[string]interface{}, 1)
					resultsMap[k] = resultMap

				} else {
					resultsMap[k] = nil
				}

			}

		} else {
			switch sqlMapsExecutor.parmas.(type) {
			case map[string]map[string]interface{}:
				parmasMap = sqlMapsExecutor.parmas.(map[string]map[string]interface{})

			default:
				if sqlMapsExecutor.session.IsSqlFuc == true {
					err1 := sqlMapsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}

			for k, _ := range sqlkeysMap {
				sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkeysMap[k])
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				if parmasMap[k] == nil {
					switch sqlCmd {
					case "select", "desc":
						model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkeysMap[k]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkeysMap[k]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				} else {
					parmaMap := parmasMap[k]
					switch sqlCmd {
					case "select", "desc":
						model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkeysMap[k], &parmaMap).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkeysMap[k], &parmaMap).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				}

				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultsMap[k] = make([]map[string]interface{}, len(model_1_results.Results))
					resultsMap[k] = model_1_results.Results

				} else if sqlModel == 2 {
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlMapsExecutor.session.IsSqlFuc == true {
							err1 := sqlMapsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultsMap[k] = make([]map[string]interface{}, 1)
					resultsMap[k] = resultMap

				} else {
					resultsMap[k] = nil
				}
			}
		}

		if sqlMapsExecutor.session.IsSqlFuc == true {
			err1 := sqlMapsExecutor.session.Commit()
			if err1 != nil {
				return nil, nil, err1
			}
		}
		return nil, resultsMap, nil

	}

	return nil, nil, nil
}
