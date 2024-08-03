import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { join } from 'path'

export default defineConfig({
	plugins: [sveltekit()],
	resolve: {
		alias: {
			'@src': join(__dirname, 'src'),
			'@components': join(__dirname, 'src/components'),
			'@starlight': join(__dirname, 'src/components/starlight'),
			'@la': join(__dirname, 'src/components/la'),
			'@app': join(__dirname, 'src/app'),
		}
	},
	server: {
		host: "0.0.0.0",
		proxy: {
			"/api": {
				target: "http://localhost:8733",
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
});
