package astolform

type compiledRule struct {
	fn       ValidatorFunc
	param    string
	ruleName string
}

type fieldMetaData struct {
	index     int
	fieldName string
	rules     []compiledRule
}

type structCache struct {
	fields []fieldMetaData
}
