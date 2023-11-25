import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		fs: {
			allow: ['styled-system'],
		},
	},
	build: {
		target: 'esnext'
	},
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	}
});
