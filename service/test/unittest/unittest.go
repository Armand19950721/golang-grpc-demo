package unittest

import (
	"service/utils"
	"time"

	"github.com/google/uuid"
)

func Flow() {
	testIsTimeBetween()
	testRemane()
	testEmailLimit()
	// TestDbMasterSlave()
}

func testEmailLimit() {
	id := uuid.New()

	for i := 1; i < utils.UserEmailLimit+4; i++ {
		valid, err := utils.CheckUserEmailLimit(id)

		if err != nil {
			panic("CheckUserEmailLimit err:" + err.Error())
		}

		if i >= utils.UserEmailLimit+1 && valid {
			panic("CheckUserEmailLimit test fail")
		}

		if !valid {
			utils.PrintObj("testEmailLimit !valid return")
			return
		}
	}
}

func testIsTimeBetween() {
	now := time.Now()
	past := now.Add(time.Second * -341)
	future := now.Add(time.Second * 341)

	res := utils.IsTimeBetween(now, past, future) // success
	if !res {
		panic("utils.IsTimeBetween test fail")
	}
	res = utils.IsTimeBetween(now, now, future) // success
	if !res {
		panic("utils.IsTimeBetween test fail")
	}
	res = utils.IsTimeBetween(now, past, now) // success
	if !res {
		panic("utils.IsTimeBetween test fail")
	}
	res = utils.IsTimeBetween(past, now, future) // fail
	if res {
		panic("utils.IsTimeBetween test fail")
	}
	res = utils.IsTimeBetween(past, future, now) // fail
	if res {
		panic("utils.IsTimeBetween test fail")
	}
}

func testRemane() {
	res := utils.AutoRename("abc")
	if res != "abc(1)" {
		panic("utils.AutoRename fail")
	}

	res = utils.AutoRename("abc(1)")
	if res != "abc(2)" {
		panic("utils.AutoRename fail")
	}

	res = utils.AutoRename("abc(12345)")
	if res != "abc(12346)" {
		panic("utils.AutoRename fail")
	}

	res = utils.AutoRename("abc(1234333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333335)")
	if res != "abc(1234333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333335)(1)" {
		panic("utils.AutoRename fail")
	}

	res = utils.AutoRename("abc(1234333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333335)(1)")
	if res != "abc(1234333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333335)(2)" {
		panic("utils.autoRename fail")
	}
}
