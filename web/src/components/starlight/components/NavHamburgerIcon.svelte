<svelte:options immutable />

<script lang="ts">
	export let expanded = false
</script>

<div class="nav-icon-wrap">
	<div id="nav_expand_btn_icon" class:expanded>
		<div />
		<div />
		<div />
	</div>
</div>

<style lang="sass">
	.nav-icon-wrap
		--icon-line-height: 2px

		display: flex
		align-items: center
		justify-content: center
		margin-top: calc(var(--spacing) * -1)
		margin-bottom: calc(var(--spacing) * -1)
		max-width: var(--nav-menu-size)

		#nav_expand_btn_icon 
			display: block
			position: relative
			outline: transparent
			width: calc(var(--nav-menu-size) / 2)
			height: calc(var(--nav-menu-size) / 2.5)
						
			div
				position: absolute
				left: 0
				width: var(--icon-line-height)
				height: var(--icon-line-height)
				filter: invert(100%)
				background-color: var(--nav-menu-background-color)
				border-radius: var(--icon-line-height)
				transition: 0.2s ease-in-out

			div:first-child 
				top: 0
				width: calc(var(--nav-menu-size) / 3)

			div:nth-child(2)
				top: calc(50%  - var(--icon-line-height) / 2)
				width: calc(var(--nav-menu-size) / 5)

			div:last-child
				top: calc(100%  - var(--icon-line-height))
				width: calc(var(--nav-menu-size) / 2.5)

			&.expanded 
				div:first-child
					top: calc(50%  - var(--icon-line-height) / 2)
					width: calc(var(--nav-menu-size) / 2)
					transform: rotateZ(45deg)

				div:nth-child(2)
					width: 0
					margin-left: calc(var(--nav-menu-size) / 4)

				div:last-child
					top: calc(50%  - var(--icon-line-height) / 2)
					width: calc(var(--nav-menu-size) / 2)
					transform: rotateZ(-45deg)
</style>