package load_balancer


var (
	methodMap = make(map[string]*Method)
)

type Method struct {
	Name string
	Type int
	Group string
}

func AddMethodInfo(methodName string, methodType int, methodGroup string) {
	group := "default"
	if len(methodGroup) > 0 {
		group = methodGroup
	}
	methodMap["getProgramAccounts"] = &Method{Name:methodName, Type: methodType, Group: group}
}

func GetMethodInfo(methodName string) *Method {
	return methodMap[methodName]
}