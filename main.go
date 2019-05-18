package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func detectFrontend() string {
	if _, err := exec.LookPath("ludo"); err == nil {
		return "ludo"
	}

	if _, err := exec.LookPath("retroarch"); err == nil {
		return "retroarch"
	}

	return ""
}

func main() {

	frontend := detectFrontend()
	if frontend == "" {
		log.Fatal("No libretro frontend enabled in systemd")
	}

	raspiAdaptor := raspi.NewAdaptor()
	raspiAdaptor.SetName(frontend)

	resetBtn := gpio.NewButtonDriver(raspiAdaptor, "3")
	resetBtn.SetName("resetBtn")

	powerBtn := gpio.NewButtonDriver(raspiAdaptor, "5")
	powerBtn.SetName("powerBtn")

	power := gpio.NewLedDriver(raspiAdaptor, "7")
	power.SetName("power")

	led := gpio.NewLedDriver(raspiAdaptor, "8")
	led.SetName("led")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	robot := gobot.NewRobot("superpi-shutdown-libretro",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{resetBtn, powerBtn, power, led},
		func() {
			power.On()
			led.On()

			powerBtn.On(gpio.ButtonRelease, func(data interface{}) {
				log.Print("Shutting down...")

				gobot.Every(100*time.Millisecond, func() {
					led.Toggle()
				})

				if err := exec.Command("systemctl", "stop", frontend).Run(); err != nil {
					log.Print(err)
					return
				}

				time.Sleep(time.Second)

				if err := exec.Command("shutdown", "-h", "now").Run(); err != nil {
					log.Print(err)
				}
			})

			var (
				resetting bool         = false
				ticker    *time.Ticker = nil
			)

			resetBtn.On(gpio.ButtonRelease, func(data interface{}) {

				ticker = gobot.Every(100*time.Millisecond, func() {
					led.Toggle()
				})

				time.AfterFunc(2*time.Second, func() {
					if !resetBtn.Active {
						log.Print("Resetting...")

						resetting = true

						if err := exec.Command("systemctl", "stop", frontend).Run(); err != nil {
							log.Print(err)
							return
						}

						time.Sleep(time.Second)

						if err := exec.Command("shutdown", "-r", "now").Run(); err != nil {
							log.Print(err)
						}
					}
				})
			})

			resetBtn.On(gpio.ButtonPush, func(data interface{}) {
				if ticker != nil && !resetting {
					log.Printf("Restarting %s...", frontend)

					ticker.Stop()
					ticker = nil

					led.On()

					if err := exec.Command("systemctl", "restart", frontend).Run(); err != nil {
						log.Print(err)
					}
				}
			})
		},
	)

	go func() {
		if err := robot.Start(); err != nil {
			log.Print(err)
		}
	}()

	<-sigs
	power.Off()
}
