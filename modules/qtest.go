package modules

import (
	"fmt"

	lazlo "github.com/djosephsen/lazlo/lib"
)

var QuestionTest = &lazlo.Module{
	Name: `QuestionTest`,
	Usage: `"%BOTNAME% askme foo" : replies with the question: foo?
%BOTNAME% qtest : runs an automated question test`,
	Run: func(b *lazlo.Broker) {
		cb1 := b.MessageCallback(`(?i)(ask *me) (.*)`, true)
		cb2 := b.MessageCallback(`(?i)(qtest)`, true)
		for {
			select {
			case newReq := <-cb1.Chan:
				go newQuestion(b, newReq) // an example of dynamic question/response
			case newReq := <-cb2.Chan:
				go runTest(b, newReq) // an example of scripted question/response
			}
		}
	},
}

func newQuestion(b *lazlo.Broker, req lazlo.PatternMatch) {
	lazlo.Logger.Info("new question")
	qcb := b.QuestionCallback(req.Event.User, req.Match[2])
	answer := <-qcb.Answer
	response := fmt.Sprintf("You answered: '%s'", answer)
	b.Say(response, qcb.DMChan)
}

func runTest(b *lazlo.Broker, req lazlo.PatternMatch) {
	dmChan := b.GetDM(req.Event.User)
	user := b.SlackMeta.GetUserName(req.Event.User)
	b.Say(fmt.Sprintf(`hi %s! I'm going to ask you a few questions.`, user), dmChan)
	qcb := b.QuestionCallback(req.Event.User, `what is your name?`)
	name := <-qcb.Answer
	qcb = b.QuestionCallback(req.Event.User, `what is your quest?`)
	quest := <-qcb.Answer
	qcb = b.QuestionCallback(req.Event.User, `what is your favorite color?`)
	color := <-qcb.Answer
	b.Say(fmt.Sprintf(`awesome. you said your name is %s, your quest is %s and your favorite color is %s`, name, quest, color), dmChan)
}
