package slackerlib

import(
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
)


type ApiRequest struct{
	URL		string
	Values	url.Values
	Broker	*Broker
}

//base function for communicating with the slack api
func MakeAPIReq(req ApiRequest)(*ApiResponse, error){
	resp:=new(ApiResponse)
	req.Values.Set(`token`, req.Bot.Config.Token)

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
func (b *Broker) getMeASocket() error {
   var req = ApiRequest{
      URL: `https://slack.com/api/rtm.start`,
		Values: make(url.Values),
      Broker: b,
   }
   authResp,err := MakeAPIReq(req)
   if err != nil{
      return err
   }

   if authResp.URL == ""{
      return fmt.Errorf("Auth failure")
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
      return fmt.Errorf("no dice dialing that websocket: %v", err)
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
   replyText:=fmt.Sprintf(`%s: %s`, event.Sbot.Meta.GetUserName(event.User), s)
   return event.Sbot.Send(&Event{
      Type:    event.Type,
      Channel: event.Channel,
      Text:    replyText,
      })
}

// convinience function to respond to a message event
func (event *Event) Respond(s string) chan map[string]interface{}{
   return event.Sbot.Send(&Event{
      Type:    event.Type,
      Channel: event.Channel,
      Text:    s,
      })
}

// convinience function to join a channel
// bots aren't actually allowed to use this command (I should probably delete this)
func (channel *Channel) Join(bot *Sbot) (*ApiResponse, error){
   var req = ApiRequest{
      URL: `https://slack.com/api/channels.join`,
		Values: url.Values{`name`: {channel.Name}},
      Bot: bot,
   }
	resp, err := MakeAPIReq(req)
	return resp, err	
}
