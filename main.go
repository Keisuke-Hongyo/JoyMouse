package main

import (
	"JoyMouse/RKJXV122400R"
	"machine"
	"machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"
)

// スティックの値格納用構造体
type switchState struct {
	xValue  uint16
	yValue  uint16
	swState bool
}

func getData(ch chan<- switchState) {
	var s switchState

	// ジョイスティックで使用するのポート初期化
	j, err := joystick.New(machine.A0, machine.A1, machine.D2)
	if err != nil {
		panic(err)
	}
	for {
		// 取得データから移動量を計算
		s.xValue = j.GetDataX() >> 12 // 16bit -> 4bit 変換
		s.yValue = j.GetDataY() >> 12 // 16bit -> 4bit 変換

		// スイッチの状態を取得
		if j.GetDataSw() == false {
			time.Sleep(100 * time.Millisecond)
			if j.GetDataSw() == false {
				s.swState = false
			} else {
				s.swState = true
			}
		} else {
			s.swState = true
		}
		ch <- s
		time.Sleep(5 * time.Millisecond)
	}
}

// LED点灯制御用
func ledCtrl(ch chan<- bool) {

	ledR := machine.LED_RED
	ledG := machine.LED_GREEN
	ledB := machine.LED_BLUE

	ledR.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledG.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledB.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledG.High()
	ledB.High()

	for {
		ledR.High()
		time.Sleep(100 * time.Millisecond)
		ledR.Low()
		time.Sleep(100 * time.Millisecond)
		ch <- true
	}
}

// 初期化ルーチン
func init() {
	// ADC初期化
	machine.InitADC()
}

// メインルーチン
func main() {
	var s switchState

	ch1 := make(chan switchState, 1)
	ch2 := make(chan bool, 1)

	m := mouse.Port()
	kb := keyboard.Port()
	go getData(ch1)
	go ledCtrl(ch2)

	for {
		// Receive channel data from Goroutine
		select {
		// Control of ADC
		case s = <-ch1:
			m.Move(int(s.yValue)-8, int(s.xValue)-8)
			if s.swState == false {
				_, err := kb.Write([]byte("Hello !!"))
				if err != nil {
					panic(err)
				}
				time.Sleep(200 * time.Millisecond)
			}
		// Control of LED
		case <-ch2:
		}
	}
}
