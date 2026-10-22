<script lang="ts">
	import '@src/nodespecific.ts';
	import '@alis-is/starlight/style.sass';
	import '@src/styles/default.sass';

	import { APP_STATUS_LEVEL } from '@app/state/status';
	import { APP_ID, APP_CONNECTION_STATUS } from '@app/state/index';
	import { theme, cycleTheme } from '@app/theme';
	import { SunIcon, MoonIcon } from '@alis-is/starlight/la';

	$: animateStatusBar = $APP_STATUS_LEVEL === 'warning' || $APP_STATUS_LEVEL === 'error';
	$: centerStatusBar = $APP_STATUS_LEVEL === 'ok';
</script>

<div class="layout">
	<header>
		<div class="title-wrap">
			<a class="unstyle-link title" href="/about">
				<h4>{$APP_ID}</h4>
			</a>
			<div class="connection-status" title="connection status">
				<div
					class="connection-status-sign"
					class:connected={$APP_CONNECTION_STATUS === 'connected'}
					class:reconnecting={$APP_CONNECTION_STATUS === 'reconnecting'}
					class:paused={$APP_CONNECTION_STATUS === 'paused'}
				></div>
				<span class="connection-status-label">{$APP_CONNECTION_STATUS}</span>
			</div>
			<button
				class="unstyle-button theme-switch"
				title="Switch to {$theme === 'dark' ? 'light' : 'dark'} theme"
				aria-label="Switch to {$theme === 'dark' ? 'light' : 'dark'} theme"
				onclick={cycleTheme}
			>
				{#if $theme === 'dark'}<MoonIcon />{:else}<SunIcon />{/if}
				<span class="theme-label">{$theme}</span>
			</button>
		</div>
		<div
			class="status-bar"
			class:ok={$APP_STATUS_LEVEL === 'ok'}
			class:warning={$APP_STATUS_LEVEL === 'warning'}
			class:animate={animateStatusBar}
			class:center={centerStatusBar}
		></div>
	</header>
	<main>
		<slot />
	</main>
	<footer>
		© {new Date().getFullYear()} tez.capital
	</footer>
</div>

<style lang="sass">
.layout
	display: flex
	flex-direction: column
	min-height: 100vh
	min-height: 100dvh
	color: var(--text-color)

header
	position: sticky
	top: 0
	height: var(--header-height)
	display: flex
	align-items: center
	z-index: 40
	background-color: var(--nav-menu-background-color, var(--card-background-color))
	backdrop-filter: var(--blur-filter)
	border-bottom: 1px solid var(--frost-border-color)

	.title-wrap
		display: flex
		align-items: center
		gap: var(--spacing)
		box-sizing: border-box
		width: 100%
		max-width: var(--content-max-width)
		margin: 0 auto
		padding: 0 var(--spacing)
		text-transform: uppercase

	.title
		flex: 1
		min-width: 0

		h4
			margin: 0
			white-space: nowrap
			overflow: hidden
			text-overflow: ellipsis

	.connection-status
		display: flex
		align-items: center
		gap: var(--spacing-f2)
		font-size: 0.9rem

		.connection-status-sign
			width: 10px
			height: 10px
			border-radius: 50%
			background-color: var(--error-color)

			&.connected
				background-color: var(--success-color)

			&.reconnecting, &.paused
				background-color: var(--warning-color)

	.theme-switch
		display: flex
		align-items: center
		gap: var(--spacing-f2)
		padding: var(--spacing-f2) var(--spacing)
		border-radius: 999px
		border: 1px solid var(--frost-border-color)
		font-size: 0.8rem
		text-transform: capitalize
		transition: background-color 0.15s

		&:hover
			background-color: var(--accent-soft)

		:global(svg)
			width: 18px
			height: 18px
			fill: var(--text-color)

	.status-bar
		position: absolute
		left: 0
		bottom: -1px
		width: 100%
		height: 2px
		background-size: 200% 100% !important
		background: linear-gradient(to right, transparent, var(--error-color), transparent)

		&.ok
			background: linear-gradient(to right, transparent, var(--success-color), transparent)

		&.warning
			background: linear-gradient(to right, transparent, var(--warning-color), transparent)

		&.animate
			animation: slide 5s linear infinite

		&.center
			background-position: center

main
	flex: 1
	box-sizing: border-box
	width: 100%
	max-width: var(--content-max-width)
	margin: 0 auto
	padding: var(--spacing)

footer
	display: flex
	justify-content: center
	padding: var(--spacing)
	color: var(--hint-color)
	font-size: 0.85rem

@media (max-width: 600px)
	header
		.connection-status-label, .theme-label
			display: none

@keyframes slide
	0%
		background-position: left
	50%
		background-position: right
	100%
		background-position: left
</style>
