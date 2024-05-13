/*
Copyright 2022 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPickDevice(t *testing.T) {
	mc := jdownloader.NewMockClient()
	_, err := pickDevice(mc)
	assert.Error(t, err)
	mc.SetDevices(&[]jdownloader.DeviceInfo{{
		Name: "mock",
	}})
	dev, err := pickDevice(mc)
	assert.NoError(t, err)
	assert.Equal(t, "mock", dev)
}

func pfloat64(c float64) *float64 {
	return &c
}

func pint64(c int64) *int64 {
	return &c
}

func TestFormatSpeed(t *testing.T) {
	assert.Equal(t, "1000 B/s", formatSpeed(pfloat64(1000.0)))
	assert.Equal(t, "9.8 KiB/s", formatSpeed(pfloat64(10000.0)))
	assert.Equal(t, "976.6 KiB/s", formatSpeed(pfloat64(1000000.0)))
	assert.Equal(t, "95.4 MiB/s", formatSpeed(pfloat64(100000000.0)))
	assert.Equal(t, "N/A", formatSpeed(nil))
}

func TestFormatEta(t *testing.T) {
	assert.Equal(t, "00:16:40", formatEta(pint64(1000)))
	assert.Equal(t, "8 days 21:32:34", formatEta(pint64(768754)))
	assert.Equal(t, "N/A", formatEta(nil))
}
