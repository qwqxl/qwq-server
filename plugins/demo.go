package plugins

// ---------- 插件实现 ----------
//func init() {
//	// 注册插件到系统
//	plugin.RegisterBuiltin("greeter", func(config string) plugin.PluginInterface {
//		return &GreeterPlugin{prefix: config}
//	})
//}
//
//type GreeterPlugin struct {
//	prefix string
//}
//
//func (gp *GreeterPlugin) Execute(params map[string]interface{}) (string, error) {
//	name, ok := params["name"].(string)
//	if !ok {
//		name = "Guest"
//	}
//	return fmt.Sprintf("%s %s!", gp.prefix, name), nil
//}
