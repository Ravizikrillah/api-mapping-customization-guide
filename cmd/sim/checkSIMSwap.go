package main

// SimSwapPlugin is a plugin for checking SIM swap.
type SimSwapPlugin struct{}

// Execute checks SIM swap with the arguments maxAge and score.
func (p SimSwapPlugin) Execute(args ...interface{}) interface{} {
	if len(args) != 2 {
		return "leng kurang"
	}

	maxAge, ok1 := args[0].(float64)
	score, ok2 := args[1].(float64)

	if !ok1 || !ok2 {
		return "Error"
	}

	if maxAge < 24 && score == 1 {
		return true
	} else if maxAge >= 24 && maxAge <= 48 && score == 2 {
		return true
	} else if maxAge > 48 && maxAge <= 72 && score == 3 {
		return true
	} else if maxAge > 72 && score == 4 {
		return true
	} else {
		return false
	}
}

// SimSwapPluginInstance is a variable that stores the plugin instance.
var SimSwapPluginInstance SimSwapPlugin

// go build -buildmode=plugin -o plugins/sim_swap_plugin.so sim/checkSIMSwap.go
