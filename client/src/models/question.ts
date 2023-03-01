export interface Question {
	Text: string
	Id: string
	Votes: number
	Owned: boolean
	Answered: boolean
	Anonymous: boolean
	Creator?: string
}
