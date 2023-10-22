package joystick

import (
	"machine"
)

// Joystick アナログジョイスティック制御用 構造体
type Joystick struct {
	xStick machine.ADC
	yStick machine.ADC
	sw     machine.Pin
}

// New 初期化
func New(xStick machine.Pin, yStick machine.Pin, sw machine.Pin) (Joystick, error) {
	var err error
	joystick := Joystick{}
	joystick.xStick = machine.ADC{Pin: xStick}
	joystick.yStick = machine.ADC{Pin: yStick}
	joystick.sw = sw

	err = joystick.xStick.Configure(machine.ADCConfig{})
	if err != nil {
		return Joystick{}, err
	}

	err = joystick.yStick.Configure(machine.ADCConfig{})
	if err != nil {
		return Joystick{}, err
	}

	joystick.sw.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	return joystick, nil
}

// GetDataX X軸のアナログ値を取得
func (j Joystick) GetDataX() uint16 {
	return j.xStick.Get()
}

// GetDataY Y軸のアナログ値を取得
func (j Joystick) GetDataY() uint16 {
	return j.yStick.Get()
}

// GetDataSw スイッチの状態を取得
func (j Joystick) GetDataSw() bool {
	return j.sw.Get()
}
