<svelte:options immutable />

<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '@components/starlight/components/Button.svelte';
	import { TEZPEAK_VERSION, TEZPEAK_CODENAME } from '@src/common/constants';
	import { DiscordIcon, TwitterIcon, TelegramIcon } from '@src/components/app/icons';
</script>

<div class="about-wrap">
	<div class="about">
		<div class="logo not-selectable">
			<svg
				width="100%"
				height="100%"
				viewBox="0 0 315 235"
				version="1.1"
				xmlns="http://www.w3.org/2000/svg"
				xmlns:xlink="http://www.w3.org/1999/xlink"
				xml:space="preserve"
				style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:2;"
			>
				<g transform="matrix(1,0,0,1,-60.6666,-529.167)">
					<g id="White" transform="matrix(0.787415,0,0,0.636664,60.6666,529.167)">
						<clipPath id="_clip1">
							<rect x="0" y="0" width="399.374" height="368.057" />
						</clipPath>
						<g clip-path="url(#_clip1)">
							<g transform="matrix(2.11886,0,0,2.62057,657.769,-740.562)">
								<g transform="matrix(1,0,0,1,0,186.389)">
									<path
										d="M-242.435,225.841L-217.584,194.204L-192.829,225.841L-122.481,225.841L-162.243,158.839L-170.941,175.757L-217.584,95.66L-263.941,174.419L-273.308,158.839L-310.966,225.841L-242.435,225.841ZM-278.431,204.401L-268.093,185.682L-273.308,175.757L-298.254,220.202L-278.431,204.401ZM-201.877,158.513C-202.136,166.616 -179.445,208.658 -177.996,207.692L-172.417,197.98L-172.418,197.978L-171.196,195.653L-155.389,210.677L-161.751,189.634C-173.889,178.573 -187.388,168.268 -201.875,158.513L-201.877,158.51L-201.877,158.511L-201.879,158.51L-201.877,158.513ZM-161.8,175.757L-161.751,189.634L-143.644,207.692L-161.8,175.757ZM-266.869,197.98L-229.126,149.53L-223.499,158.644L-218.498,152.895L-206.29,170.274L-207.774,154.158L-217.621,140.452L-210.58,133.63L-198.835,145.101L-217.932,112.704L-266.869,197.98Z"
										style="fill:#fff;"
									/>
								</g>
								<g transform="matrix(1,0,0,1,0,186.389)">
									<g opacity="0.75">
										<path
											d="M-273.308,175.757L-273.308,190.025L-298.254,220.202L-278.431,204.401L-268.093,185.682L-273.308,175.757Z"
											style="fill:#fff;"
										/>
									</g>
								</g>
								<g transform="matrix(1,0,0,1,0,186.389)">
									<g opacity="0.75">
										<path
											d="M-198.835,145.101L-217.932,112.704L-217.621,140.452L-210.58,133.63L-198.835,145.101"
											style="fill:#fff;"
										/>
									</g>
								</g>
								<g transform="matrix(1,0,0,1,0,186.389)">
									<g opacity="0.75">
										<path
											d="M-177.998,207.692C-179.448,208.658 -202.143,166.608 -201.879,158.51L-172.419,197.98L-177.998,207.692Z"
											style="fill:#fff;"
										/>
									</g>
								</g>
								<g transform="matrix(1,0,0,1,-0.00227542,186.389)">
									<g opacity="0.55">
										<path
											d="M-155.389,210.677L-161.751,189.634C-173.891,178.572 -187.39,168.266 -201.879,158.51L-172.419,197.98L-171.196,195.653L-155.389,210.677Z"
											style="fill:#fff;"
										/>
									</g>
								</g>
								<g transform="matrix(1,0,0,1,0,186.389)">
									<g opacity="0.75">
										<path
											d="M-161.8,175.757L-161.751,189.634L-143.644,207.692L-161.8,175.757Z"
											style="fill:#fff;"
										/>
									</g>
								</g>
								<g transform="matrix(1.02742,0,0,1.02742,6.17497,176.954)">
									<circle cx="-216.95" cy="228.768" r="8.566" style="fill:#fff;" />
								</g>
							</g>
							<g transform="matrix(1.26998,0,0,1.57069,-671.943,-228.049)">
								<g transform="matrix(14,0,0,14,606.325,351.459)">
									<path
										d="M0.489,-0.699L0.489,-0.641L0.295,-0.641L0.295,-0L0.225,-0L0.225,-0.641L0.03,-0.641L0.03,-0.699L0.489,-0.699Z"
										style="fill-rule:nonzero;"
									/>
								</g>
								<g transform="matrix(14,0,0,14,617.441,351.459)">
									<path
										d="M0.125,-0.062L0.482,-0.062L0.482,-0L0.043,-0L0.043,-0.056L0.395,-0.637L0.049,-0.637L0.049,-0.699L0.477,-0.699L0.477,-0.643L0.125,-0.062Z"
										style="fill-rule:nonzero;"
									/>
								</g>
								<g transform="matrix(14,0,0,14,628.627,351.459)">
									<path
										d="M0.048,-0.35C0.048,-0.419 0.063,-0.48 0.093,-0.535C0.123,-0.589 0.164,-0.631 0.216,-0.662C0.268,-0.692 0.326,-0.707 0.39,-0.707C0.468,-0.707 0.535,-0.689 0.591,-0.652C0.646,-0.615 0.687,-0.562 0.712,-0.495L0.63,-0.495C0.61,-0.542 0.58,-0.579 0.539,-0.606C0.498,-0.633 0.448,-0.646 0.39,-0.646C0.339,-0.646 0.293,-0.634 0.252,-0.61C0.211,-0.586 0.178,-0.552 0.155,-0.506C0.132,-0.462 0.12,-0.409 0.12,-0.35C0.12,-0.291 0.132,-0.239 0.155,-0.194C0.178,-0.149 0.211,-0.115 0.252,-0.091C0.293,-0.067 0.339,-0.055 0.39,-0.055C0.448,-0.055 0.498,-0.068 0.539,-0.095C0.58,-0.121 0.61,-0.157 0.63,-0.204L0.712,-0.204C0.687,-0.137 0.646,-0.086 0.59,-0.049C0.534,-0.012 0.467,0.006 0.39,0.006C0.326,0.006 0.268,-0.009 0.216,-0.04C0.164,-0.07 0.123,-0.112 0.093,-0.166C0.063,-0.22 0.048,-0.281 0.048,-0.35Z"
										style="fill-rule:nonzero;"
									/>
								</g>
							</g>
						</g>
					</g>
				</g>
			</svg>
		</div>
		<h2 class="title">TEZPEAK</h2>
		<h5 class="version not-selectable">
			{TEZPEAK_VERSION}
			{TEZPEAK_CODENAME}
		</h5>
		<div class="social">
			<a class="svg-icon" href="https://discord.gg/2PnCr62YgF" target="_blank" rel="noreferrer">
				<DiscordIcon />
			</a>
			<a class="svg-icon" href="https://twitter.com/TezCapital" target="_blank" rel="noreferrer">
				<TwitterIcon />
			</a>
			<a class="svg-icon" href="https://t.me/tezcapital" target="_blank" rel="noreferrer">
				<TelegramIcon />
			</a>
		</div>
		<div class="navigation">
			<div class="home">
				<Button label="Home" on:click={() => goto("/")} />
			</div>
		</div>
	</div>
</div>

<style lang="sass">
	.about-wrap
		height: 100%
		margin: 0px !important
		min-height: 300px !important
		overflow: hidden
		--footer-height: 30px
		color: var(--text-color)
		
		.about
			display: flex
			flex-direction: column
			justify-content: center
			align-items: center
			height: calc(100% - var(--footer-height))
			margin:0 !important
			
			.logo
				width: 350px
				max-width: 40vw
				max-height: 20vh
			.title
				margin-top: min(var(--spacing-x2), 3vh)
				margin-bottom: min(var(--spacing), 3vh)
			.version
				margin-top: 0
				margin-bottom: min(var(--spacing-x2), 3vh)
				font-size: 1.5rem
				font-weight: 500
			.social
				display: grid
				grid-template-columns: auto auto auto
				grid-gap: var(--spacing)

				a
					margin: 0 0.5rem
					color: var(--text-color)
					transition: transform 0.1s ease-in-out
					border-radius: var(--border-radius)
					text-decoration: none
					width: 50px
					
					&:hover
						transform: scale(1.1)

.navigation
	display: grid
	padding: var(--spacing-x3)
	grid-template-columns: 1fr 200px 1fr

	.home
		grid-column: 2
</style>
