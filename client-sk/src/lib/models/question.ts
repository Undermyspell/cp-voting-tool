export interface Question {
	Text: string
	Id: string
	Votes: number
	Voted: boolean
	Owned: boolean
	Answered: boolean
	Anonymous: boolean
	Creator?: string
}
