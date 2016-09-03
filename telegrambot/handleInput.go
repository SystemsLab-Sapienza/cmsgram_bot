package main

func handleInput(msg string, user int) error {
	var rm ResponseMessage

	lcmd, err := getLastCommand(user)
	if err != nil {
		return err
	}
	if len(lcmd) == 0 {
		rm.Send("Comando non riconosciuto.", user)
		return nil
	}

	switch lcmd {
	case "/cerca":
		res, err := getInfo(msg)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			rm.Send("Nessun docente trovato.", user)
			return nil
		}

		rm.Send(res, user)
	default:
		rm.Send("Comando non riconosciuto.", user)
		return nil
	}

	// Reset last command
	err = setLastCommand("", user)
	if err != nil {
		return err
	}

	return nil
}
