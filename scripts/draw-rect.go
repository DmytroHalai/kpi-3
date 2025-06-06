//go:build ignoretest
// +build ignoretest

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func sendPostRequest(data string) error {
	url := "http://localhost:17000/"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func moveAlongSegment(x1, y1, x2, y2 float64, steps int, delay time.Duration) error {
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t

		cmd := fmt.Sprintf("move %.4f %.4f", x, y)
		if err := sendPostRequest(cmd); err != nil {
			return fmt.Errorf("move error: %w", err)
		}
		if err := sendPostRequest("update"); err != nil {
			return fmt.Errorf("update error: %w", err)
		}
		time.Sleep(delay)
	}
	return nil
}

func main() {
	commands := []string{
		"green",
		"bgrect 0.25 0.25 0.75 0.75",
		"figure 0.5 0.5",
		"update",
	}

	for _, cmd := range commands {
		if err := sendPostRequest(cmd); err != nil {
			fmt.Printf("Error sending command '%s': %v\n", cmd, err)
			return
		}
	}

	points := [][2]float64{
		{0.25, 0.25}, // A
		{0.75, 0.25}, // B
		{0.75, 0.75}, // C
		{0.25, 0.75}, // D
	}

	stepsPerEdge := 50
	delay := 50 * time.Millisecond

	for {
		for i := 0; i < len(points); i++ {
			start := points[i]
			end := points[(i+1)%len(points)]
			if err := moveAlongSegment(start[0], start[1], end[0], end[1], stepsPerEdge, delay); err != nil {
				fmt.Printf("Error moving along segment: %v\n", err)
				return
			}
		}
	}
}
