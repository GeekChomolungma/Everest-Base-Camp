package handlerfactory

var FactoryInstanceMap = make(map[string]BaseFactory)

func Init() {
	FactoryInstanceMap["HuoBi"] = &HuoBiFactory{}
}
