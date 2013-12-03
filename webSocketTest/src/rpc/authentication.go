/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 11/30/13
 * Time: 2:30 PM
 * To change this template use File | Settings | File Templates.
 */
package rpc

func Login(account string) bool {
	//if login fail
	if account == "fail" {
		return false
	}


	//notify other user he had login in
	// handle sync or asyc
	if account == "ok" {

//		//get this user online friend id list
//		ids := []string{12,45,56}
//
//		channel := &broadcast_channel.Channel{ids}
//		arg := map[string]interface {}{"Id":account}
//		channel.Push("auth.LoginSuccess", arg)

	}

	return true
}



