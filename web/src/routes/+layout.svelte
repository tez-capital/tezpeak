<script lang="ts">
	import '@src/nodespecific.ts';
	import '@src/styles/default.sass';

	import {
		APP_ID,
		APP_STATUS_LEVEL,
		APP_CONNECTION_STATUS,
	} from '@app/state';

	$: animateStatusBar = $APP_STATUS_LEVEL === 'warning' || $APP_STATUS_LEVEL === 'error';
	$: centerStatusBar = $APP_STATUS_LEVEL === 'ok';
</script>

<div class="layout-grid">
	<div class="background"></div>
	<header>
		<!-- <slot name="header" /> -->
		<div class="title-wrap">
			<a class="unstyle-link" href="/about">
				<h4>{$APP_ID}</h4>
			</a>
			<div class="connection-status">
				<div class="connection-status-sign" class:connected={$APP_CONNECTION_STATUS === "connected"} class:reconnecting={$APP_CONNECTION_STATUS === "reconnecting"}></div>
				{$APP_CONNECTION_STATUS}
			</div>
		</div>
		<div
			class="status-bar"
			class:ok={$APP_STATUS_LEVEL === "ok"}
			class:warning={$APP_STATUS_LEVEL === "warning"}
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

<!-- <div id="menu-layer"></div> -->
<style lang="sass">
:root
	--menu-gap: var(--spacing)

.background
	position: fixed
	height: 100vh
	width: 100vw
	pointer-events: none
	background-image: url('/assets/images/svg/bg.svg')
	background-size: cover
	background-position: center
	z-index: -2

.layout-grid
	position: relative
	display: grid
	width: 100vw
	height: 100vh
	grid-template-columns: minmax(100px, 1fr)
	grid-template-rows: var(--header-height) 1fr
	grid-template-areas: "header" "main" "footer"
	grid-column-gap: var(--menu-gap)
	color: var(--text-color)
	
	header
		position: fixed
		height: var(--header-height)
		width: 100vw
		display: flex
		justify-content: left
		align-items: center
		
		z-index: 100
		background-color: var(--background-color)

		.title-wrap
			display: grid
			width: 100vw
			grid-template-columns: auto 1fr auto
			align-items: center
			padding: 0 var(--spacing)
			text-transform: uppercase
		
		.connection-status
			grid-column: 3
		
			.connection-status-sign
				display: inline-block
				width: 12px
				height: 12px
				border-radius: 50%
				margin-right: var(--spacing-f2)
				background-color: var(--error-color)

				&.connected
					background-color: var(--success-color) !important
				
				&.reconnecting
					background-color: var(--warning-color) !important

		.status-bar
			position: absolute
			left: 0
			bottom: 0
			width: 100vw
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
		position: relative
		grid-area: main
		overflow-x: hidden
		//min-height: 100vh

	footer
		position: fixed
		display: flex
		width: 100%
		grid-area: footer
		justify-content: center
		color: var(--text-color)
		padding-bottom: var(--spacing-f2)
		bottom: 0
		z-index: -1

	
@keyframes slide
	0%
		background-position: left
	50%
		background-position: right
	100%
		background-position: left

</style>
