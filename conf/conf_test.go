package conf

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

var (
	folderName = "test_conf"
	fileName   = "test_token.json"
)

func TestConfig_Save(t *testing.T) {

	defer func() {
		os.RemoveAll(folderName)
	}()

	tests := []struct {
		conf           config
		fileName, want string
	}{
		{
			config{
				"www.baidu.com",
				"12345678",
				"2018-01-01 12:00:00",
			},
			"test1_w_token.json",
			`{"host":"www.baidu.com","token":"12345678","expired_at":"2018-01-01 12:00:00"}`,
		},
		{
			config{
				"www.youku.com",
				"1234512345",
				"2019-01-01 00:00:00",
			},
			"test2_w_token.json",
			`{"host":"www.youku.com","token":"1234512345","expired_at":"2019-01-01 00:00:00"}`,
		},
	}
	for _, test := range tests {

		if err := test.conf.Save(path.Join(folderName, test.fileName)); err != nil {
			t.Errorf("conf save meet error %v", err)
		}

		byteBuffer, _ := ioutil.ReadFile(path.Join(folderName, test.fileName))
		if string(byteBuffer) != test.want {
			t.Errorf("%s want: %s get: %s", test.fileName, test.want, byteBuffer)
		}
	}

}

func TestConfig_IsValid(t *testing.T) {
	tests := []struct {
		conf config
		host string
		now  time.Time
		want bool
	}{
		{
			config{
				"www.youku.com",
				"1234512345",
				"2018-01-01 01:00:00",
			},
			"www.youku.com",
			time.Date(2018, 1, 1, 0, 55, 55, 0, time.UTC),
			true,
		},
		{
			config{
				"www.youku1.com",
				"1234512345",
				"2018-01-01 01:00:00",
			},
			"www.youku.com",
			time.Date(2018, 1, 1, 0, 55, 55, 0, time.UTC),
			false,
		},
		{
			config{
				"www.youku.com",
				"1234512345",
				"2018-01-01 01:00:00",
			},
			"www.youku.com",
			time.Date(2018, 1, 1, 1, 5, 0, 0, time.UTC),
			false,
		},
	}

	for _, test := range tests {
		if got := test.conf.IsValid(test.host, test.now); got != test.want {
			t.Errorf("conf %v, host: %s, now: %s, want: %t, got: %t", test.conf, test.host, test.now, test.want, got)
		}
	}
}

func TestConfig_Update(t *testing.T) {

	tests := []struct {
		jsonStr string
		want    *config
	}{
		{
			`{"host":"www.baidu.com","token":"12345678","expired_at":"2018-01-01 12:00:00"}`,
			&config{"www.baidu.com", "12345678", "2018-01-01 12:00:00"},
		},
		{
			`{"host":"www.baidu.com","token":"12345678","expired_at":""}`,
			&config{"www.baidu.com", "12345678", ""},
		},
	}
	for _, test := range tests {
		got := new(config)
		r := []byte(test.jsonStr)
		got.Update(r)

		if *got != *test.want {
			t.Errorf("want: %+v got: %+v", test.want, got)
		}
	}
}

func TestCheckTimeValid(t *testing.T) {
	tests := []struct {
		timeStr string
		now     time.Time
		want    bool
	}{
		{
			"2018-01-01 01:00:00",
			time.Date(2018, 1, 1, 0, 55, 55, 0, time.UTC),
			true,
		},
		{
			"2018-01-01 01:00:00",
			time.Date(2018, 1, 1, 2, 0, 0, 0, time.UTC),
			false,
		},
	}

	for _, test := range tests {
		got, err := checkTimeValid(test.timeStr, test.now)
		if err != nil {
			t.Errorf("checkTimeValid got %s got error %s", test.timeStr, err)
			continue
		}
		if got != test.want {
			t.Errorf("check %s want=%t, got=%t", test.timeStr, test.want, got)
		}
	}
}
