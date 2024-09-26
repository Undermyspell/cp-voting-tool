export type Question ={
	Text: string;
	Id: string;
	Votes: number;
	Voted: boolean;
	Owned: boolean;
	Answered: boolean;
	Anonymous: boolean;
	Creator?: string;
}

export type QuestionDataStore = {
	init: () => Promise<void>,
	questions: Question[]
	readonly sortedAnsweredQuestions: Question[]
	readonly  sortedUnansweredQuestions: Question[],
	user: any,
	usersOnlineCount: number,
	addQuestion(this: QuestionDataStore, question: Question): () => void,
	updateQuestion(this: QuestionDataStore, question: Question): () => void,
	deleteQuestion(this: QuestionDataStore, question: Question): () => void,
	answerQuestion(this: QuestionDataStore, question: Question): () => void,
	updateUserOnlineCount(this: QuestionDataStore, usersOnlineCount: number): () => void,
}