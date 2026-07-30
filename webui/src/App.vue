<script setup>
	import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
	export default {}
</script>

<template>
	<header v-if="$route.name !== 'login'" class="top-navbar">
		<div class="container-fluid d-flex justify-content-between align-items-center px-4">
			<a class="brand-logo" href="#/home">WasaText</a>
			<div class="nav-links">
				<RouterLink to="/home" class="nav-link-custom">Chat</RouterLink>
				<RouterLink to="/profile" class="nav-link-custom">Profilo</RouterLink>
				<RouterLink to="/login" class="nav-link-custom logout-link">Esci</RouterLink>
			</div>
		</div>
	</header>

	<div class="container-fluid p-0">
		<RouterView />
	</div>
</template>

<style>
.top-navbar {
	background-color: var(--color-bg-white);
	border-bottom: 1px solid var(--color-border);
	height: 60px;
	display: flex;
	align-items: center;
}
.brand-logo {
	color: var(--color-primary);
	font-weight: 600;
	font-size: 1.4rem;
	text-decoration: none;
	letter-spacing: -0.5px;
}
.nav-links {
	display: flex;
	gap: 24px;
}
.nav-link-custom {
	color: var(--color-text-main);
	text-decoration: none;
	font-weight: 500;
	font-size: 0.95rem;
	transition: color 0.2s ease;
}
.nav-link-custom:hover, .nav-link-custom.router-link-active {
	color: var(--color-primary);
}
.logout-link {
	color: var(--color-text-muted);
}
</style>
