
import Alpine from "alpinejs"
import { Centrifuge } from 'centrifuge';


(function(){
    //@ts-ignore
    window.Alpine = Alpine

    Alpine.store("questionData", {
        async init() {
            this.user = await (await fetch('/user')).json()
            initCentrifugo(this.user) 
        },
        questions: JSON.parse(document.getElementById('questions')!.textContent!),
        get sortedQuestions() {
            return this.questions.sort((a, b) => b.Votes - a.Votes)
        },
        user: null,
        addQuestion(question) {
            this.questions.push(question)
        }
    })

    Alpine.start() 

    const initCentrifugo = async(user) => {
    console.log("init centrifugo")

    const centrifuge = new Centrifuge("ws://localhost:3333/api/v1/connection/websocket", {
        token: user.Token
    });
    centrifuge.on('connecting', function (ctx) {
        console.log(`connecting: ${ctx.code}, ${ctx.reason}`);
    }).on('connected', function (ctx) {
        console.log(`connected over ${ctx.transport}`);
    }).on('disconnected', function (ctx) {
        console.log(`disconnected: ${ctx.code}, ${ctx.reason}`);
    }).on('message', function (msg) {
        console.log(`message: ${JSON.stringify(msg)}`);

        const data = JSON.parse(msg.data.Payload)
        const eventType = msg.data.EventType

        switch(eventType){
            case "start_session":
                
                break
            case "stop_session":
                
                break
            case "user_connected":
            case "user_disconnected":
                
                break
            case "new_question":
                Alpine.store('questionData').addQuestion(data)
            case "undo_upvote_question":
                
                break
            case "update_question":
                
                break
            case "delete_question":
                
                break
            case "answer_question":
                
                break
        }
    }).connect();
    }
})();

