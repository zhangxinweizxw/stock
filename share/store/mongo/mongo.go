package mongo
//
//import (
//    "io"
//
//    "stock/share/logging"
//)
//
//var (
//	MgoDb *mgo.Database
//)
//
//func Init(source string, dbName string) error {
//	if MgoDb == nil {
//		session, err := GetMongo(source)
//		if err != nil {
//			return err
//		}
//		news := session.Copy()
//		MgoDb = news.DB(dbName)
//	}
//	return nil
//}
//
//func GetMongo(source string) (*mgo.Session, error) {
//	session, err := mgo.Dial(source)
//	if err != nil {
//		logging.Fatal(err)
//	}
//
//	//设置会话模式
//	session.SetMode(mgo.Monotonic, true)
//	return session, err
//}
//
//// query
//func Query(tableName string, query map[string]interface{}, result interface{}, limit int, sort ...string) error {
//	var err error
//	if limit == 0 {
//		err = MgoDb.C(tableName).Find(query).Sort(sort...).All(result)
//	} else {
//		err = MgoDb.C(tableName).Find(query).Sort(sort...).Limit(limit).All(result)
//	}
//	//when happend relication change, it will refresh session and try do once.
//	if err == io.EOF {
//		MgoDb.Session.Refresh()
//		err = MgoDb.C(tableName).Find(query).Sort(sort...).Limit(limit).All(result)
//	}
//	return err
//}
//
////single row
//func QueryRow(tableName string, query map[string]interface{}, result interface{}) error {
//	err := MgoDb.C(tableName).Find(query).One(result)
//	if err == io.EOF {
//		MgoDb.Session.Refresh()
//		err = MgoDb.C(tableName).Find(query).One(result)
//	}
//	return err
//}
//
//func Pipe(tableName string, pipeline interface{}, result interface{}) error {
//	err := MgoDb.C(tableName).Pipe(pipeline).All(result)
//	if err == io.EOF {
//		MgoDb.Session.Refresh()
//		err = MgoDb.C(tableName).Pipe(pipeline).All(result)
//	}
//	return err
//}
//
//// insert
//func Insert(tableName string, args ...interface{}) error {
//	err := MgoDb.C(tableName).Insert(args...)
//	if err == io.EOF {
//		MgoDb.Session.Refresh()
//		err = MgoDb.C(tableName).Insert(args...)
//	}
//	return err
//
//}
//
//// update
//func Update(tableName string, query map[string]interface{}, update interface{}) error {
//	_, err := MgoDb.C(tableName).UpdateAll(query, update)
//	if err == io.EOF {
//		MgoDb.Session.Refresh()
//		_, err = MgoDb.C(tableName).UpdateAll(query, update)
//	}
//	return err
//}
//
//// delete
//func Delete(tableName string, query map[string]interface{}) error {
//	_, err := MgoDb.C(tableName).RemoveAll(query)
//	if err == io.EOF {
//		MgoDb.Session.Refresh()
//		_, err = MgoDb.C(tableName).RemoveAll(query)
//	}
//
//	return err
//}
//
//// close mongo pool
//func Close() {
//	if MgoDb != nil && MgoDb.Session != nil {
//		MgoDb.Session.Close()
//	}
//}
