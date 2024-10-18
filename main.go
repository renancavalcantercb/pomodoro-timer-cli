package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	timerActive bool
	startTime   time.Time
	duration    = 25 * time.Minute
	timer       *time.Timer
	mu          sync.Mutex
	stateFile   = "/tmp/pomodoro_state.json"
)

type TimerState struct {
	Active    bool
	StartTime time.Time
}

func saveState() {
	state := TimerState{
		Active:    timerActive,
		StartTime: startTime,
	}

	file, err := os.Create(stateFile)
	if err != nil {
		fmt.Println("Error saving state:", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(state)
	if err != nil {
		fmt.Println("Error encoding state:", err)
	}
}

func loadState() {
	file, err := os.Open(stateFile)
	if err != nil {
		return
	}
	defer file.Close()

	var state TimerState
	err = json.NewDecoder(file).Decode(&state)
	if err != nil {
		fmt.Println("Error decoding state:", err)
		return
	}

	timerActive = state.Active
	startTime = state.StartTime

	if timerActive {
		elapsed := time.Since(startTime)
		remaining := duration - elapsed
		if remaining > 0 {
			timer = time.NewTimer(remaining)
			go func() {
				<-timer.C
				mu.Lock()
				timerActive = false
				saveState()
				mu.Unlock()
				fmt.Println("\nPomodoro finished!")
			}()
		} else {
			timerActive = false
			saveState()
			fmt.Println("\nPomodoro already finished!")
		}
	}
}

func main() {
	loadState()

	rootCmd := &cobra.Command{Use: "pomodoro"}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new 25-minute Pomodoro",
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			loadState()

			if timerActive {
				fmt.Println("Pomodoro already running!")
				return
			}
			timerActive = true
			startTime = time.Now()

			saveState()

			timer = time.NewTimer(duration)
			fmt.Println("Starting a new Pomodoro: 25 minutes")

			go func() {
				<-timer.C
				mu.Lock()
				timerActive = false
				saveState()
				mu.Unlock()
				fmt.Println("\nPomodoro finished!")
			}()
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the current Pomodoro",
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			loadState()

			if !timerActive {
				fmt.Println("No active Pomodoro to stop.")
				return
			}

			if timer != nil {
				timer.Stop()
				timerActive = false
				saveState()
				fmt.Println("Pomodoro stopped.")
			} else {
				fmt.Println("Error: Timer is nil.")
			}
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check the status of the current Pomodoro",
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			loadState()

			if !timerActive {
				fmt.Println("No active Pomodoro.")
				return
			}
			elapsed := time.Since(startTime)
			remaining := duration - elapsed
			fmt.Printf("Pomodoro running: %.2f minutes remaining\n", remaining.Minutes())
		},
	}

	rootCmd.AddCommand(startCmd, stopCmd, statusCmd)
	rootCmd.Execute()
}
