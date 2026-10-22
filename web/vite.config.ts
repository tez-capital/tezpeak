import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { demoApi } from './demo-api';

export default defineConfig(({ mode }) => ({
	plugins: [sveltekit(), mode === 'demo' && demoApi()],
	server: {
		host: "0.0.0.0",
		// demo mode answers /api itself
		proxy: mode === 'demo' ? undefined : {
			"/api": {
				target: "http://127.0.0.1:8733/",
				changeOrigin: true,
				secure: false,
			}
		}
	},
	build: {
		target: "es2020",
		commonjsOptions: {
			transformMixedEsModules: true,
		},
	},
}));
