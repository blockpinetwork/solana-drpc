package utils

func MapContains(src map[interface{}]interface{}, key interface{}) bool {
	_,ok := src[key]
	return ok
}
