package main

func (a *app) handleCommand(input string) (bool, error) {
	// FIXME: uncomment in the event commands are needed at the CLI level.
	// if strings.HasPrefix(input, cmdPrefix) {
	// 	input = strings.TrimFunc(input, func(c rune) bool {
	// 		// Trim off any unexpected input
	// 		return unicode.IsSpace(c) || !unicode.IsLetter(c) && !unicode.IsNumber(c) && !unicode.IsSymbol(c) && c != '*'
	// 	})

	// 	if strings.Compare(input, "help") == 0 {
	// 		// Currently no commands, so no help

	// 		// Return false cause we still want the backend process to evaluate
	// 		// this command.
	// 		return false, nil
	// 	}
	// }

	return false, nil
}
