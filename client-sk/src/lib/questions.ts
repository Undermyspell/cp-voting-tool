import { derived, get, writable, type Writable } from 'svelte/store';
import type { Question } from './models/question';
import { deleteRequest, getRequest, postRequest, putRequest } from './api';
// import { eventSource } from './eventsource';
import { activeSessison, userOnline } from './session';
import { centrifuge } from './centrifuge';
import type { MessageContext } from 'centrifuge';
import type { VotingEvent } from './models/voting-event';

const allQuestions: Writable<Question[]> = writable([]);
export const questions = derived(allQuestions, ($allQuestions: Question[]) =>
	$allQuestions.filter((q) => q.Answered === false)
);
export const answeredQuestions = derived(allQuestions, ($allQuestions: Question[]) =>
	$allQuestions.filter((q) => q.Answered === true)
);
export const sessionActive = writable(false);
export const isAutosortActive = writable(false);
export const unAnsweredCount = derived(questions, ($questions) => $questions.length);
export const answeredCount = derived(
	answeredQuestions,
	($answeredQuestions) => $answeredQuestions.length
);


centrifuge.subscribe((centrifuge) => {
	if(centrifuge) {
		centrifuge.on("message", (msg: MessageContext) => {
			console.log("Received event: ", msg)
			const event: VotingEvent = msg.data as VotingEvent

			const data = JSON.parse(event.Payload)

			switch(event.EventType){
				case "start_session":
					console.log('start listener');
					activeSessison.set(true);
					break
				case "stop_session":
					activeSessison.set(false);
					clearQuestions();
					break
				case "user_connected":
				case "user_disconnected":
					userOnline.set(data.UserCount);
					break
				case "new_question":
					questionAdded(data);
					break
				case "upvote_question":
				case "undo_upvote_question":
					updateVote(data);
					break
				case "update_question":
					questionEdited(data);
					break
				case "delete_question":
					questionDeleted(data);
					break
				case "answer_question":
					questionAnswered(data);
					break
			}
			
		})
	}
})

// eventSource.subscribe((eventSource) => {
// 	if (eventSource) {
// 		eventSource.addEventListener('new_question', (event) => {
// 			const data = JSON.parse(event.data);
// 			questionAdded(data);
// 		});
// 		eventSource.addEventListener('upvote_question', (event) => {
// 			const data = JSON.parse(event.data);
// 			updateVote(data);
// 		});
// 		eventSource.addEventListener('undo_upvote_question', (event) => {
// 			const data = JSON.parse(event.data);
// 			updateVote(data);
// 		});
// 		eventSource.addEventListener('update_question', (event) => {
// 			const data = JSON.parse(event.data);
// 			questionEdited(data);
// 		});
// 		eventSource.addEventListener('delete_question', (event) => {
// 			const data = JSON.parse(event.data);
// 			questionDeleted(data);
// 		});
// 		eventSource.addEventListener('answer_question', (event) => {
// 			const data = JSON.parse(event.data);
// 			questionAnswered(data);
// 		});
// 	}
// });

export const getQuestions = async () => {
	try {
		const repsonse = await getRequest({ path: '/question/session' });
		if (repsonse.ok) {
			activeSessison.set(true);
			const data = (await repsonse.json()) as Question[];
			allQuestions.set(data);
			updateAutosort(true);
			sortAndUpdateQuestions();
			updateAutosort(false);
		}
	} catch (error) {
		console.log(error);
	}
};

export function getQuestion(id: string) {
	return get(allQuestions).find((q) => q.Id === id);
}

export function clearQuestions() {
	allQuestions.set([]);
}
export const updateAutosort = (value: boolean) => {
	isAutosortActive.set(value);
	if (value === true) {
		sortAndUpdateQuestions();
	}
};

export const postQuestion = async (questionText: string, anonymous: boolean) => {
	await postRequest({
		path: '/question/new',
		body: JSON.stringify({ anonymous: anonymous, text: questionText })
	});
};

export const updateQuestion = async (payload: { Id: string; Anonymous: boolean; Text: string }) => {
	await putRequest({
		path: '/question/update',
		body: JSON.stringify({ id: payload.Id, text: payload.Text, anonymous: payload.Anonymous })
	});
};

export const voteQuestion = async (questionId: string) => {
	await putRequest({ path: `/question/upvote/${questionId}` });
};

export const undoVoteQuestion = async (questionId: string) => {
	await putRequest({ path: `/question/undovote/${questionId}` });
};

export const answerQuestion = async (questionId: string) => {
	await putRequest({ path: `/question/answer/${questionId}` });
};

export const deleteQuestion = async (questionId: string) => {
	await deleteRequest({ path: `/question/delete/${questionId}` });
};

const questionAdded = (question: Question) => {
	allQuestions.set([...get(allQuestions), question]);
	sortAndUpdateQuestions();
};

const questionDeleted = (payload: { Id: string }) => {
	allQuestions.set([...get(allQuestions).filter((q) => q.Id !== payload.Id)]);
	sortAndUpdateQuestions();
};

const updateVote = (payload: { Id: string; Votes: number; Voted: boolean }) => {
	allQuestions.set(
		get(allQuestions).map((q) => (q.Id === payload.Id ? Object.assign({}, q, { ...payload }) : q))
	);
	sortAndUpdateQuestions();
};

const questionEdited = (payload: {
	Id: string;
	Text: string;
	Creator: string;
	Anonymous: boolean;
}) => {
	allQuestions.set(
		get(allQuestions).map((q) => (q.Id === payload.Id ? Object.assign({}, q, { ...payload }) : q))
	);
	sortAndUpdateQuestions();
};

const questionAnswered = (payload: { Id: string }) => {
	allQuestions.set(
		get(allQuestions).map((q) =>
			q.Id === payload.Id ? Object.assign({}, q, { Answered: true }) : q
		)
	);
	sortAndUpdateQuestions();
};

const sortAndUpdateQuestions = () => {
	if (get(isAutosortActive) === true) {
		const sortedArray = [...get(allQuestions)].sort((a, b) => b.Votes - a.Votes);
		allQuestions.set(sortedArray);
	}
};
