package mongo
//
//type Model struct {
//	TableName string `json:"-"`
//	TeamID    int64  `json:"-"`
//}
//
//// 获取单条数据
//func (this *Model) GetSingle(exps map[string]interface{}, result interface{}) error {
//	return QueryRow(this.TableName, exps, result)
//}
//
//// 获取多条数据
//func (this *Model) GetMulti(exps map[string]interface{}, result interface{}, limit int, sort ...string) error {
//	return Query(this.TableName, exps, result, limit, sort...)
//}
//
//// 删除数据
//func (this *Model) Delete(exps map[string]interface{}) error {
//	return Delete(this.TableName, exps)
//}
//
//// 插入数据
//func (this *Model) Insert(args ...interface{}) error {
//	return Insert(this.TableName, args...)
//}
//
//// 更新数据
//func (this *Model) Update(exps map[string]interface{}, params interface{}) error {
//	return Update(this.TableName, exps, params)
//}
//
//func (this *Model) Pipe(pipeline interface{}, result interface{}) error {
//	return Pipe(this.TableName, pipeline, result)
//}
