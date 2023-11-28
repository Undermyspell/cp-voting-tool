import { json2csv } from 'json-2-csv';
import { questions } from './questions';
import { get } from 'svelte/store';

export const export2csv = async () => {
	const currentQuestions = get(questions);
	const data = await json2csv(currentQuestions, {
		excelBOM: true,
		delimiter: { field: ';' },
		keys: ['Text', 'Votes']
	});
	const csvBlob = new Blob([data], { type: 'text/csv' });
	const url = window.URL.createObjectURL(csvBlob);
	const tmpLink = document.createElement('a');
	tmpLink.setAttribute('href', url);
	tmpLink.setAttribute('download', 'questions.csv');
	tmpLink.click();
};
