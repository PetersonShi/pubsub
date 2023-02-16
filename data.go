// License statement
// ---------------------------------------------------------------------------------
//
// Licensed under The MIT License
// For full copyright and license information, please see the MIT-LICENSE.txt
// Redistributions of files must retain the above copyright notice.
//
// @author    shitou<rookiefun@sina.com>
// @copyright shitou<rookiefun@sina.com>
// @license   http://www.opensource.org/licenses/mit-license.php MIT License
// ---------------------------------------------------------------------------------

package pubsub


type msgData struct {
	values map[string]interface{}
}

//Set
func(p *msgData)SetValue(key string,value interface{}) {p.values[key] = value}
func(p *msgData)SetString(key string,value string)     {p.SetValue(key,value)}

//Get
func(p *msgData)GetValue(key string)interface{}{
	return p.values[key]
}
func(p *msgData)GetString(key string)string{
	if v,ok := p.GetValue(key).(string); ok{
		return v
	}
	return ""
}

func newMsgData()*msgData {
	return &msgData{values: make(map[string]interface{})}
}