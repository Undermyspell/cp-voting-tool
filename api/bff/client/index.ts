
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
        usersOnlineCount: 0,
        addQuestion(question) {
            this.questions.push(question)
            setTimeout(() => {
                //@ts-ignore
                htmx.process(document.getElementById('question-list'));
            }, 0);
        },
        updateQuestion(question) {
            this.questions = this.questions.map((q) => (q.Id === question.Id ? Object.assign({}, q, { ...question }) : q))
            setTimeout(() => {
                //@ts-ignore
                htmx.process(document.getElementById('question-list'));
            }, 0);
        },
        deleteQuestion(question: any): void {
            this.questions = this.questions.filter(q => q.Id !== question.Id);
        },
        updateUserOnlineCount(usersOnlineCount: number) {
            this.usersOnlineCount = usersOnlineCount
        },
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
                //@ts-ignore
                htmx.ajax('GET', '/q/s/page/true', {target:'body', swap:'innerHTML'})
                break
            case "stop_session":
                //@ts-ignore
                htmx.ajax('GET', '/q/s/page/false', {target:'body', swap:'innerHTML'})
                break
            case "user_connected":
            case "user_disconnected":
                Alpine.store('questionData').updateUserOnlineCount(data.UserCount)
                break
            case "new_question":
                Alpine.store('questionData').addQuestion(data)
                break
            case "upvote_question":
            case "undo_upvote_question":
            case "update_question":
                Alpine.store('questionData').updateQuestion(data)
                break
            case "delete_question":
                Alpine.store('questionData').deleteQuestion(data)
                break
            case "answer_question":
                
                break
        }
    }).connect();
    }
})();

