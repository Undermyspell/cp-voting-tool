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
	user: User | null,
	usersOnlineCount: number,
	autoSortQuestions: boolean,
	addQuestion(this: QuestionDataStore, question: Question): () => void,
	updateQuestion(this: QuestionDataStore, question: Question): () => void,
	deleteQuestion(this: QuestionDataStore, question: Question): () => void,
	answerQuestion(this: QuestionDataStore, question: Question): () => void,
	updateUserOnlineCount(this: QuestionDataStore, usersOnlineCount: number): () => void,
}

export type ThemeData = {
	init: () => void,
	isDarkMode: boolean,
	toggleTheme: () => void,
	setTheme:(theme: 'dark' | 'light' | 'system') => void
}

export type User ={
	Email: string, 
	Name: string, 
	Token: string,
 	Role: 'admin' | 'session_admin' | 'contributor'
}