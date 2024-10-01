
import Alpine from "alpinejs"
import { Centrifuge } from 'centrifuge';
import type { Question, QuestionDataStore, ThemeData, User } from './types';

(function(){
    //@ts-ignore
    window.Alpine = Alpine

    Alpine.store("questionData", {
        async init(this: QuestionDataStore) {
            this.user = await (await fetch('/user')).json() as User
            initCentrifugo(this.user) 
        },
        questions: JSON.parse(document.getElementById('questions')!.textContent!),
        get sortedAnsweredQuestions(): Question[] {
            return (this as QuestionDataStore).questions.filter(q => q.Answered).sort((a, b) => b.Votes - a.Votes)
        },
        get sortedUnansweredQuestions(): Question[]  {
            return (this as QuestionDataStore).questions.filter(q => !q.Answered).sort((a, b) => b.Votes - a.Votes)
        },
        user: null,
        usersOnlineCount: 0,
        addQuestion(this: QuestionDataStore, question: Question) {
            this.questions.push(question)
        },
        updateQuestion(this: QuestionDataStore, question: Question) {
            this.questions = this.questions.map((q) => (q.Id === question.Id ? Object.assign({}, q, { ...question }) : q))
        },
        deleteQuestion(this: QuestionDataStore, question: Question): void {
            this.questions = this.questions.filter(q => q.Id !== question.Id);
        },
        answerQuestion(this: QuestionDataStore, question: Question): void {
            this.questions = this.questions.map((q) => (q.Id === question.Id ? Object.assign({}, q, { Answered: true }) : q))
        },
        updateUserOnlineCount(this: QuestionDataStore, usersOnlineCount: number) {
            this.usersOnlineCount = usersOnlineCount
        },
    } as QuestionDataStore)


    Alpine.store("theme",{
        isDarkMode: false,
        init() {
            console.log("INIT")
            const savedTheme = localStorage.getItem('theme') || 'system';
            if (savedTheme === 'dark') {
              this.isDarkMode = true;
            } else if (savedTheme === 'light') {
              this.isDarkMode = false;
            } else {
              this.isDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
            }
  
            window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (event) => {
                if (localStorage.getItem('theme') === 'system') {
                  this.isDarkMode = event.matches;
                }
            });
        },
        toggleTheme() {
            if (localStorage.getItem('theme') === 'light') {
              this.setTheme('dark');
            } else if (localStorage.getItem('theme') === 'dark') {
              this.setTheme('system');
            } else {
              this.setTheme('light');
            }
          },
        setTheme(theme) {
            if (theme === 'dark') {
              this.isDarkMode = true;
              localStorage.setItem('theme', 'dark');
            } else if (theme === 'light') {
              this.isDarkMode = false;
              localStorage.setItem('theme', 'light');
            } else {
              this.isDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
              localStorage.setItem('theme', 'system');
            }
          }
    } as ThemeData)

    Alpine.start() 

    const initCentrifugo = async(user) => {
    console.log("init centrifugo")

    const protocol = window.location.protocol === "https:" ? "wss" : "ws"
    const centrifuge = new Centrifuge(`${protocol}://${window.location.host}/api/v1/connection/websocket`, {
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
                (Alpine.store('questionData') as QuestionDataStore).updateUserOnlineCount(data.UserCount)
                break
            case "new_question":
                (Alpine.store('questionData') as QuestionDataStore).addQuestion(data)
                break
            case "upvote_question":
            case "undo_upvote_question":
            case "update_question":
                (Alpine.store('questionData') as QuestionDataStore).updateQuestion(data)
                break
            case "delete_question":
                (Alpine.store('questionData') as QuestionDataStore).deleteQuestion(data)
                break
            case "answer_question":
                (Alpine.store('questionData') as QuestionDataStore).answerQuestion(data)
                break
        }
    }).connect();
    }
})();

