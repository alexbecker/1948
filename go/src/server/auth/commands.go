package auth

func Manage(command string, args ...string) {
	var err error
	switch command {
	case "add":
		err = AddUser(args[0], args[1])
	case "grant":
		err = AddRole(args[0], args[1])
	case "makeparent":
		err = MakeParent(args[0], args[1])
	case "print":
		err = PrintUsers()
	default:
		panic("Command must be one of \"add\", \"grant\", \"makeparent\", \"print\"")
	}

	if err != nil {
		panic(err)
	}
}
