package lib

import(
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/fatih/structs"
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"strconv"
)


type ApiRequest struct{
	URL		string
	Values	url.Values
	Broker	*Broker
}

//base function for communicating with the slack api
func MakeAPIReq(req ApiRequest)(*ApiResponse, error){
	resp:=new(ApiResponse)
	req.Values.Set(`token`, req.Broker.Config.Token)

	reply, err := http.PostForm(req.URL, req.Values)
   if err != nil{
      return resp, err
   }
   defer reply.Body.Close()

	dec := json.NewDecoder(reply.Body)
   err = dec.Decode(resp)
	if err != nil {
		return resp, fmt.Errorf("Couldn't decode json. ERR: %v", err)
	}
	return resp, nil
}

// Go forth and get a websocket for RTM and all the Slack Team Metadata
func (b *Broker) getASocket() (*websocket.Conn, *ApiResponse, error) {
   var req = ApiRequest{
      URL: `https://slack.com/api/rtm.start`,
		Values: make(url.Values),
      Broker: b,
   }
   authResp,err := MakeAPIReq(req)
   if err != nil{
      return nil, nil, err
   }

   if authResp.URL == ""{
      return nil, nil, fmt.Errorf("Auth failure")
   }
   wsURL := strings.Split(authResp.URL, "/")
   wsURL[2] = wsURL[2] + ":443"
   authResp.URL = strings.Join(wsURL, "/")
   Logger.Debug(`Team Wesocket URL: `, authResp.URL)

   var Dialer websocket.Dialer
   header := make(http.Header)
   header.Add("Origin", "http://localhost/")

   ws, _, err := Dialer.Dial(authResp.URL, header)
   if err != nil{
      return nil, nil, fmt.Errorf("no dice dialing that websocket: %v", err)
   }

   //yay we're websocketing
   return ws, authResp, nil
}

// parses sBot.Meta to return a user's Name field given its ID
func (meta *ApiResponse) GetUserName(id string) string{
   for _,user := range meta.Users{
      if user.ID == id{
         return user.Name
      }
   }
   return ``
}

// parses sBot.Meta to return a pointer to a user object given its ID
func (meta *ApiResponse) GetUser(id string) *User{
   for _,user := range meta.Users{
      if user.ID == id{
         return &user
      }
   }
   return nil
}

// parses sBot.Meta to return a pointer to a user object given its Name
func (meta *ApiResponse) GetUserByName(name string) *User{
   for _,user := range meta.Users{
      if user.Name == name{
         return &user
      }
   }
   return nil
}

// parses sBot.Meta to return a pointer to a channel object given its ID
func (meta *ApiResponse) GetChannel(id string) *Channel{
   for _,channel := range meta.Channels{
      if channel.ID == id{
         return &channel
      }
   }
   return nil
}

// parses sBot.Meta to return a pointer to a channel object given its Name
func (meta *ApiResponse) GetChannelByName(name string) *Channel{
   for _,channel := range meta.Channels{
      if channel.Name == name{
         return &channel
      }
   }
   return nil
}

// convinience function to reply to a message event
func (event *Event) Reply(s string) chan map[string]interface{}{
   replyText:=fmt.Sprintf(`%s: %s`, event.Broker.SlackMeta.GetUserName(event.User), s)
   return event.Broker.Send(&Event{
      Type:    event.Type,
      Channel: event.Channel,
      Text:    replyText,
      })
}

// convinience function to respond to a message event
func (event *Event) Respond(s string) chan map[string]interface{}{
   return event.Broker.Send(&Event{
      Type:    event.Type,
      Channel: event.Channel,
      Text:    s,
      })
}

// this is a confusing hack that I'm using because slack's RTM websocket
// doesn't seem to support their own markup syntax. So anything that looks
// like it has markup in it is sent into this function by the write thread
// instead of into the websocket where it belongs. 
func apiPostMessage(e Event) {
   var req = ApiRequest{
      URL: `https://slack.com/api/chat.postMessage`,
      Values: make(url.Values),
		Broker: e.Broker,
   }
   req.Values.Set(`channel`, e.Channel)
   req.Values.Set(`text`, e.Text)
   req.Values.Set(`id`, strconv.Itoa(int(e.ID)))
   req.Values.Set(`as_user`, e.Broker.Config.Name)
   req.Values.Set(`pretty`, `1`)
   authResp,_ := MakeAPIReq(req)
	s:=structs.New(authResp) // convert this to a map[string]interface{} why not? hax. 
	resp := s.Map()
	if replyVal,isReply:= resp[`reply_to`]; isReply{
		if replyVal != nil{
			e.Broker.handleApiReply(resp)
		}
	}
}
